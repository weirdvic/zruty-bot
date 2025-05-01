package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/yanzay/tbot/v2"
)

func startHandler(m *tbot.Message) {
	if m.Chat.Type == "private" && zruty.isAdmin(m.Chat.ID) {
		_, err := zruty.client.SendMessage(m.Chat.ID, greetAdmin)
		if err != nil {
			log.Print(err)
		}
	} else {
		_, err := zruty.client.SendMessage(m.Chat.ID, notAdmin)
		if err != nil {
			log.Print(err)
		}
	}
}

func reportHandler(m *tbot.Message) {
	// Проверяем, что сообщение от админа и это приватный чат
	if !zruty.isAdmin(m.Chat.ID) || m.Chat.Type != "private" {
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
			reportBuilder.WriteString(fmt.Sprintf(
				"%d.\t%s %s @%s %s назад\nВ чатах: %s\n",
				i,
				u.firstName,
				u.lastName,
				u.username,
				time.Since(u.firstSeenAt).Round(time.Second),
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

func defaultHandler(m *tbot.Message) {
	if !zruty.isValidGroup(m.Chat.ID) ||
		(m.Chat.Type != "supergroup" && m.Chat.Type != "group") {
		return
	}
	// Если в чате появились новые участники
	if len(m.NewChatMembers) > 0 {
		zruty.addUsers(m)
		zruty.welcomeUsers(m)
		return
	}
	// Проверка — отправил ли сообщение один из отслеживаемых пользователей
	uid := m.From.ID
	var (
		username   string
		firstName  string
		lastName   string
		userExists bool
	)

	err := zruty.db.QueryRow(`
		SELECT username, first_name, last_name
		FROM users
		WHERE id = ? AND check_passed_at != NULL
	`, uid).Scan(&username, &firstName, &lastName)

	switch {
	case err == sql.ErrNoRows:
		// Пользователя нет или он уже прошёл проверку
		return
	case err != nil:
		log.Printf("❌ Ошибка при проверке пользователя %d @%s: %v", uid, username, err)
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
			log.Printf("❌ Не удалось обновить check_passed для пользователя %d @%s: %v", uid, username, err)
			return
		}

		// Уведомляем админов
		zruty.notifyAdmins(fmt.Sprintf(
			"✅ Пользователь %d @%s прошёл проверку",
			uid,
			username,
		))

		log.Printf("✅ Пользователь %s %s(@%s) написал сообщение в чат!",
			firstName,
			lastName,
			username,
		)
	}
}
