package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/yanzay/tbot/v2"
)

func startHandler(m *tbot.Message) {
	if m.Chat.Type == "private" && zruty.isAdmin(m.Chat.ID) {
		greetAdminMessage, err := getSetting(zruty.db, "greetAdminMessage")
		if err != nil {
			log.Printf("❌ Не удалось прочитать greetAdminMessage: %v", err)
			greetAdminMessage = `Приветствую!`
		}
		_, err = zruty.client.SendMessage(m.Chat.ID, greetAdminMessage)
		if err != nil {
			log.Printf("❌ Не удалось отправить сообщение: %v", err)
		}
	} else {
		notAdminMessage, err := getSetting(zruty.db, "notAdminMessage")
		if err != nil {
			log.Printf("❌ Не удалось прочитать notAdminMessage: %v", err)
			notAdminMessage = `Вы не являетесь администратором этого бота.`
		}
		_, err = zruty.client.SendMessage(m.Chat.ID, notAdminMessage)
		if err != nil {
			log.Printf("❌ Не удалось отправить сообщение: %v", err)
		}
	}
}

func reportHandler(m *tbot.Message) {
	// Проверяем, что сообщение от админа и это приватный чат
	if !zruty.isAdmin(m.Chat.ID) || m.Chat.Type != "private" {
		notAdminMessage, err := getSetting(zruty.db, "notAdminMessage")
		if err != nil {
			log.Printf("❌ Не удалось прочитать notAdminMessage: %v", err)
			notAdminMessage = `Вы не являетесь администратором этого бота.`
		}
		_, err = zruty.client.SendMessage(m.Chat.ID, notAdminMessage)
		if err != nil {
			log.Printf("❌ Не удалось отправить сообщение: %v", err)
		}
		return
	}
	rows, err := zruty.db.Query(`
		SELECT
			u.id,
			u.first_name,
			u.last_name,
			u.username,
			u.first_seen_at,
			c.title
		FROM users u
		LEFT JOIN user_chats uc ON u.id = uc.user_id
		LEFT JOIN chats c ON uc.chat_id = c.chat_id
		WHERE u.check_passed_at IS NULL
		ORDER BY u.first_seen_at ASC
	`)
	if err != nil {
		log.Printf("❌ Ошибка получения пользователей: %v", err)
		return
	}
	defer rows.Close()

	var (
		reportBuilder strings.Builder
		users         = make(map[int]*user)
		userGroups    = make(map[int][]string)
	)

	for rows.Next() {
		var u user
		err := rows.Scan(&u.userID, &u.firstName, &u.lastName, &u.username, &u.firstSeenAt, &u.groupTitle)
		if err != nil {
			log.Printf("❌ Ошибка чтения строки: %v", err)
			continue
		}
		if _, exists := users[u.userID]; !exists {
			users[u.userID] = &u
		}
		userGroups[u.userID] = append(userGroups[u.userID], u.groupTitle)
	}

	if len(users) == 0 {
		reportBuilder.WriteString("`Нет отслеживаемых пользователей`")
	} else {
		reportBuilder.WriteString("```\nЕсть отслеживаемые пользователи:\n\n")
		i := 1
		for _, u := range users {
			groupTitles := strings.Join(userGroups[u.userID], ", ")
			var timeSinceStr string
			if duration := durationSince(u.firstSeenAt); duration != nil {
				timeSinceStr = duration.String()
			} else {
				timeSinceStr = "unknown"
			}
			reportBuilder.WriteString(fmt.Sprintf(
				"%d.\t%s %s @%s %s назад\nВ чатах: %s\n",
				i,
				u.firstName,
				u.lastName,
				u.username,
				timeSinceStr,
				groupTitles,
			))
			i++
		}
		reportBuilder.WriteString(fmt.Sprintf("\nВсего пользователей %d\n```", len(users)))
	}

	_, err = zruty.client.SendMessage(
		m.Chat.ID,
		reportBuilder.String(),
		tbot.OptParseModeMarkdown,
	)
	if err != nil {
		log.Printf("❌ Ошибка отправки отчёта: %v", err)
	}
}

// defaultHandler обрабатывает любые сообщения, отправленные в чате.
// Он не реагирует на сообщения, отправленные ботами.
//
// Если в чате появились новые участники, он добавляет их в базу данных
// и отправляет им приветственное сообщение.
//
// Если отправитель сообщения уже является отслеживаемым,
// обновляется отметка о прохождении проверки
// и уведомляются администраторы.
func defaultHandler(m *tbot.Message) {
	if !zruty.isValidGroup(m.Chat.ID) ||
		(m.Chat.Type != "supergroup" && m.Chat.Type != "group") {
		log.Printf("❌ Неизвестный тип чата: %s", m.Chat.Type)
		return
	}
	// Если в чате появились новые участники
	if len(m.NewChatMembers) > 0 {
		log.Printf("👥 Новые участники: %v", m.NewChatMembers)
		zruty.addUsers(m)
		zruty.welcomeUsers(m)
		return
	}
	// Проверка — отправил ли сообщение один из отслеживаемых пользователей
	uid := m.From.ID
	var (
		username   sql.NullString
		firstName  sql.NullString
		lastName   sql.NullString
		userExists bool
	)

	err := zruty.db.QueryRow(`
		SELECT username, first_name, last_name
		FROM users
		WHERE id = ? AND check_passed_at IS NULL
	`, uid).Scan(&username, &firstName, &lastName)

	switch {
	case err == sql.ErrNoRows:
		// Пользователя нет или он уже прошёл проверку
		return
	case err != nil:
		log.Printf("❌ Ошибка при проверке пользователя %d @%v: %v", uid, username, err)
		return
	default:
		userExists = true
	}

	if userExists {
		// Обновляем check_passed
		_, err := zruty.db.Exec(`
			UPDATE users SET check_passed_at = CURRENT_TIMESTAMP WHERE id = ?
		`, uid)
		if err != nil {
			log.Printf("❌ Не удалось обновить check_passed для пользователя %d @%v: %v", uid, username, err)
			return
		}

		// Уведомляем админов
		zruty.notifyAdmins(fmt.Sprintf(
			"✅ Пользователь %d @%v прошёл проверку",
			uid,
			username,
		))

		log.Printf("✅ Пользователь %v %v(@%v) написал сообщение в чат!",
			firstName,
			lastName,
			username,
		)
	}
}

// underAttackSwitchHandler - переключает режим underAttack. Если сообщение пришло
// от админа в приватном чате, то переключает режим и отправляет подтверждение.
func underAttackSwitchHandler(m *tbot.Message) {
	// Проверяем, что сообщение от админа и это приватный чат
	if !zruty.isAdmin(m.Chat.ID) || m.Chat.Type != "private" {
		notAdminMessage, err := getSetting(zruty.db, "notAdminMessage")
		if err != nil {
			log.Printf("❌ Не удалось прочитать notAdminMessage: %v", err)
			notAdminMessage = `Вы не являетесь администратором этого бота.`
		}
		_, err = zruty.client.SendMessage(m.Chat.ID, notAdminMessage)
		if err != nil {
			log.Printf("❌ Не удалось отправить сообщение: %v", err)
		}
		return
	}
	err := flipSetting(zruty.db, "underAttack")
	if err != nil {
		log.Printf("❌ Не удалось изменить underAttack: %v", err)
	}
	underAttack, err := getSetting(zruty.db, "underAttack")
	if err != nil {
		log.Printf("❌ Не удалось получить underAttack: %v", err)
	}
	_, err = zruty.client.SendMessage(m.Chat.ID, fmt.Sprintf("Значение underAttack изменено на: %s", underAttack))
	if err != nil {
		log.Printf("❌ Ошибка отправки сообщения: %v", err)
	}
}

// callbackHandler - обработчик callback запросов
// Он обрабатывает callback'и, связанные с кнопкой верификации,
// и в ответ на них разрешает пользователю писать сообщения в чат.
func (b *zrutyBot) callbackHandler(cq *tbot.CallbackQuery) {
	// Случай нажатия на кнопку верификации
	if cq.Data != "" && len(cq.Data) > 7 && cq.Data[:7] == "verify_" {
		challengeUserID, err := strconv.Atoi(strings.TrimPrefix(cq.Data, "verify_"))
		if err != nil {
			log.Printf("❌ Не удалось извлечь ID пользователя из %s: %v", cq.Data, err)
			_ = b.client.AnswerCallbackQuery(cq.ID, tbot.OptText("❌ Не удалось извлечь ID пользователя из сообщения"))
			return
		}

		if cq.From.ID != challengeUserID {
			_ = b.client.AnswerCallbackQuery(cq.ID, tbot.OptText("❌ Вы не можете подтвердить другого пользователя"))
			return
		}
		b.unrestrictUser(cq.Message.Chat.ID, challengeUserID)
		log.Printf("✅ Пользователю %d разрешено писать сообщения в чат %s", challengeUserID, cq.Message.Chat.ID)
	}
}
