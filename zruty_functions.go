package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/yanzay/tbot/v2"
)

// init инициализирует бота, читая из переменной окружения BOT_TOKEN,
// если она задана, или из БД, если переменная не задана.
// Если токен из переменной окружения отличается от токена в БД, то
// обновляет его в БД.
func (b *zrutyBot) init() error {
	botTokenEnv := os.Getenv("BOT_TOKEN")
	if botTokenEnv != "" {
		// Если переменная задана, проверим и, если нужно, обновим БД
		storedToken, err := getSetting(b.db, "botToken")
		switch {
		case err == sql.ErrNoRows:
			_, err = b.db.Exec(`INSERT INTO settings (key, value) VALUES ('botToken', ?)`, botTokenEnv)
			if err != nil {
				return fmt.Errorf("не удалось записать токен бота в БД: %w", err)
			}
			log.Println("✅ Токен бота записан в БД из переменной окружения")
		case err != nil:
			return fmt.Errorf("ошибка чтения токена бота из БД: %w", err)
		case storedToken != botTokenEnv:
			_, err = b.db.Exec(`UPDATE settings SET value = ? WHERE key = 'botToken'`, botTokenEnv)
			if err != nil {
				return fmt.Errorf("не удалось обновить токен бота в БД: %w", err)
			}
			log.Println("ℹ️ Токен бота в БД обновлён на значение из переменной окружения")
		default:
			log.Println("✅ Токен бота в БД совпадает с переменной окружения")
		}
		b.token = botTokenEnv
	} else {
		// Если переменная не задана, читаем из БД
		storedToken, err := getSetting(b.db, "botToken")
		if err == sql.ErrNoRows || storedToken == "" {
			return fmt.Errorf("не задан токен: переменная BOT_TOKEN пуста и отсутствует в settings")
		} else if err != nil {
			return fmt.Errorf("ошибка чтения токена из БД: %w", err)
		}
		b.token = storedToken
		log.Println("✅ Токен успешно получен из БД")
	}

	envAdminID := os.Getenv("BOT_ADMIN_ID")
	if envAdminID != "" {
		var exists bool
		err := b.db.QueryRow(`SELECT EXISTS (SELECT 1 FROM admins WHERE user_id = ?)`, envAdminID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("ошибка проверки администратора: %w", err)
		}

		if !exists {
			// Добавим его в users и admins
			now := time.Now().UTC()
			_, err = b.db.Exec(`
				INSERT INTO users (id, check_passed_at)
				VALUES (?, ?)
			`, envAdminID, now)
			if err != nil {
				return fmt.Errorf("не удалось создать пользователя: %w", err)
			}

			_, err = b.db.Exec(`INSERT INTO admins (user_id) VALUES (?)`, envAdminID)
			if err != nil {
				return fmt.Errorf("не удалось назначить администратора: %w", err)
			}
			log.Printf("✅ Админ %s добавлен в базу", envAdminID)
		} else {
			log.Printf("✅ Админ %s уже есть в базе", envAdminID)
		}
	} else {
		// Если переменная не задана — проверим, что в БД есть хоть один админ
		var count int
		err := b.db.QueryRow(`SELECT COUNT(*) FROM admins`).Scan(&count)
		if err != nil {
			return fmt.Errorf("не удалось посчитать админов: %w", err)
		}
		if count == 0 {
			return fmt.Errorf("не задан BOT_ADMIN_ID и в таблице admins нет записей")
		}
		log.Println("✅ Найден хотя бы один админ в базе")
	}

	// Проверяем переменную UNDER_ATTACK и обновляем значение в БД, если нужно
	underAttackEnv := os.Getenv("UNDER_ATTACK")
	switch {
	case underAttackEnv == "0" || strings.ToLower(underAttackEnv) == "false":
		// Отключаем режим "Под атакой" в БД
		_, err := b.db.Exec(`UPDATE settings SET value = 'false' WHERE key = 'underAttack'`)
		if err != nil {
			log.Printf("❌ Не удалось отключить режим \"Под атакой\": %v", err)
		}
	case underAttackEnv == "1" || strings.ToLower(underAttackEnv) == "true":
		// Включаем режим "Под атакой" в БД
		_, err := b.db.Exec(`UPDATE settings SET value = 'true' WHERE key = 'underAttack'`)
		if err != nil {
			log.Printf("❌ Не удалось включить режим \"Под атакой\": %v", err)
		}
	default:
		log.Printf("Значение UNDER_ATTACK не распознано: %s", underAttackEnv)
	}

	return nil
}

// isAdmin проверяет, является ли пользователь с указанным идентификатором администратором бота.
// Возвращает true, если пользователь найден в таблице администраторов, иначе false.
func (b *zrutyBot) isAdmin(id string) bool {
	var exists bool
	uid, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("❌ Ошибка при проверке администратора: %v", err)
		return false
	}
	err = b.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM admins WHERE user_id = ?)`, uid).Scan(&exists)
	if err != nil {
		log.Printf("❌ Ошибка при проверке администратора: %v", err)
		return false
	}
	return exists
}

// isUser проверяет, зарегистрирован ли пользователь с указанным идентификатором.
// Возвращает true, если пользователь найден в таблице users, иначе false.
func (b *zrutyBot) isUser(id string) bool {
	var exists bool
	uid, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("❌ Ошибка при проверке isUser(%s): %v", id, err)
		return false
	}
	err = b.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`, uid).Scan(&exists)
	if err != nil {
		log.Printf("❌ Ошибка при проверке isUser(%s): %v", id, err)
		return false
	}
	return exists
}

// isValidGroup проверяет, существует ли группа с указанным идентификатором в базе данных.
// Возвращает true, если группа найдена, иначе false.
func (b *zrutyBot) isValidGroup(id string) bool {
	var valid bool
	err := b.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM chats WHERE chat_id = ?)`, id).Scan(&valid)
	if err != nil {
		log.Printf("❌ Ошибка при проверке isValidGroup(%s): %v", id, err)
		return false
	}
	return valid
}

// isInGroup проверяет, является ли пользователь с указанным userID
// участником группы groupID.
// Возвращает true, если пользователь является администратором, владельцем,
// участником или ограниченным участником, иначе false.
func (b *zrutyBot) isInGroup(groupID string, userID int) bool {
	gcm, err := b.client.GetChatMember(groupID, userID)
	if err != nil {
		log.Print(err)
	}
	switch gcm.Status {
	case "owner":
		return true
	case "administrator":
		return true
	case "member":
		return true
	case "restricted":
		return true
	case "kicked":
		return false
	default:
		return false
	}
}

// addUsers регистрирует новых пользователей
func (b *zrutyBot) addUsers(m *tbot.Message) {
	users := m.NewChatMembers
	usersAdded := 0

	for _, u := range users {
		uid := strconv.Itoa(u.ID)
		groupID, err := strconv.Atoi(m.Chat.ID)
		if err != nil {
			log.Println("❌ Неверный идентификатор группы:", m.Chat.ID)
		}

		// Проверка: есть ли пользователь
		if !b.isUser(uid) {
			// Добавляем в таблицу users
			_, err := b.db.Exec(`
				INSERT INTO users (id, first_name, last_name, username, is_bot, first_seen_at)
				VALUES (?, ?, ?, ?, ?, ?)
			`, u.ID, u.FirstName, u.LastName, u.Username, u.IsBot, time.Now())
			if err != nil {
				log.Printf("❌ Ошибка добавления пользователя %d: %v", u.ID, err)
				continue
			}
			log.Printf("✅ Добавлен новый пользователь %s %s (@%s)", u.FirstName, u.LastName, u.Username)
		} else {
			log.Printf("🔄 Пользователь %s %s (@%s) уже отслеживается", u.FirstName, u.LastName, u.Username)
		}

		// Убедимся, что есть связь пользователь-группа
		_, err = b.db.Exec(`
			INSERT OR IGNORE INTO user_chats (user_id, chat_id)
			VALUES (?, ?)
		`, uid, groupID)
		if err != nil {
			log.Printf("❌ Ошибка записи связи user->chat: %v", err)
			continue
		}
		usersAdded++
	}

	log.Printf("✅ Добавлено новых пользователей: %d / %d", usersAdded, len(users))
}

// restrictUser ограничивает возможности пользователя в чате на заданное количество секунд.
// Параметры:
// - chatID: идентификатор чата, в котором нужно замутить пользователя.
// - userID: идентификатор пользователя, которому ограничиваются возможности в чате.
func (b *zrutyBot) restrictUser(chatID string, userID int) {
	err := b.client.RestrictChatMember(
		chatID,
		userID,
		&tbot.ChatPermissions{
			CanSendMessages:       false,
			CanSendMediaMessages:  false,
			CanSendPolls:          false,
			CanSendOtherMessages:  false,
			CanAddWebPagePreviews: false,
		},
	)
	if err != nil {
		log.Printf("❌ Не удалось замутить пользователя %d: %v", userID, err)
	}
	log.Printf("🚫 Пользователю %d запрещено писать сообщения в чат %s", userID, chatID)
}

// unrestrictUser восстанавливает возможности пользователя в чате.
// Параметры:
// - chatID: идентификатор чата, в котором нужно восстановить права пользователя.
// - userID: идентификатор пользователя, которому восстанавливаются права в чате.
func (b *zrutyBot) unrestrictUser(chatID string, userID int) {
	err := b.client.RestrictChatMember(
		chatID,
		userID,
		&tbot.ChatPermissions{
			CanSendMessages:       true,
			CanSendMediaMessages:  true,
			CanSendPolls:          true,
			CanSendOtherMessages:  true,
			CanAddWebPagePreviews: true,
		},
	)
	if err != nil {
		log.Printf("❌ Не удалось вернуть права пользователя %d: %v", userID, err)
	}
	log.Printf("🚫 Пользователю %d снова разрешено писать сообщения в чат %s", userID, chatID)
}

// sendChallengeMsg отправляет сообщение с кнопкой для верификации пользователя.
// Параметры:
// - chatID: идентификатор чата, в который отправляется сообщение.
// - button: кнопка, предоставляющая пользователю возможность подтвердить, что он не робот.
// Возвращает идентификатор отправленного сообщения.
func (b *zrutyBot) sendChallengeMsg(chatID string, button tbot.InlineKeyboardButton) (messageID int) {
	challengeMessage, err := getSetting(b.db, "challengeMessage")
	if err != nil {
		log.Printf("❌ Не удалось прочитать challengeMessage: %v", err)
		challengeMessage = "Подтвердите, что вы не робот"
	}
	msg, err := b.client.SendMessage(chatID,
		challengeMessage,
		tbot.OptInlineKeyboardMarkup(&tbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]tbot.InlineKeyboardButton{
				{button},
			},
		}))
	if err != nil {
		log.Printf("❌ Не удалось отправить сообщение для верификации пользователя: %v", err)
		return 0
	}
	return msg.MessageID
}

// verifyUser - функция, которая ждет 30 секунд, а затем проверяет, может ли пользователь
// отправлять сообщения в чате. Если пользователь не может отправлять сообщения,
// то он не прошел верификацию. Функция также удаляет сообщение для
// верификации пользователя.
func (b *zrutyBot) verifyUser(chatID string, userID int, challengeMsgID int) {
	challengeTime, err := getSetting(b.db, "challengeTime")
	if err != nil {
		log.Printf("❌ Ошибка при получении challengeTime: %v", err)
		challengeTime = "30"
	}
	defaultSleep := 30
	sleep := defaultSleep
	if val, err := strconv.Atoi(challengeTime); err == nil {
		sleep = val
	} else {
		log.Printf("❌ Ошибка при преобразовании challengeTime (%q) в число: %v. Используется значение по умолчанию: %d", challengeTime, err, defaultSleep)
	}
	time.Sleep(time.Duration(sleep) * time.Second)
	chatMember, err := b.client.GetChatMember(chatID, userID)
	if err != nil {
		log.Printf("❌ Не удалось получить информацию о пользователе: %v", err)
		return
	}
	if !chatMember.CanSendMessages {
		log.Printf("✅ Пользователь %d не прошел верификацию в чате %s", userID, chatID)
		// TODO: решить что делаем с пользователями, которые не прошли верификацию
	}
	err = b.client.DeleteMessage(chatID, challengeMsgID)
	if err != nil {
		log.Printf("❌ Не удалось удалить сообщение для верификации пользователя: %v", err)
	}
}

// welcomeUsers отправляет новым пользователям приветственное сообщение
func (b *zrutyBot) welcomeUsers(m *tbot.Message) {
	var (
		users            = m.NewChatMembers
		muteDuration int = 0
	)
	welcomeMessage, err := getSetting(b.db, "welcomeMessage")
	if err != nil {
		log.Printf("❌ Не удалось прочитать welcomeMessage: %v", err)
		welcomeMessage = `Добро пожаловать, <a href="tg://user?id=%d">%s</a>!`
	}
	underAttack, err := isSettingEnabled(b.db, "underAttack")
	if err != nil {
		log.Printf("❌ Не удалось прочитать underAttack: %v", err)
		underAttack = false
	}
	err = b.db.QueryRow(`SELECT CAST(value AS INTEGER) FROM settings WHERE key = 'muteDuration'`).Scan(&muteDuration)
	if err != nil {
		log.Printf("❌ Не удалось прочитать muteDuration: %v", err)
		muteDuration = 0
	}
	for _, u := range users {
		// Отправляем приветственное сообщение
		_, err := b.client.SendMessage(
			m.Chat.ID,
			fmt.Sprintf(
				welcomeMessage,
				u.ID,
				u.FirstName,
			),
			tbot.OptParseModeHTML,
		)
		if err != nil {
			log.Printf("❌ Не удалось отправить приветственное сообщение: %v", err)
		}
		// Если включен режим "Под атакой", то запускаем механизм проверки
		if underAttack {
			b.restrictUser(m.Chat.ID, u.ID)
			challengeButton := tbot.InlineKeyboardButton{
				Text:         "Я не робот!",
				CallbackData: "verify_" + strconv.Itoa(u.ID),
			}
			challengeMsgID := b.sendChallengeMsg(m.Chat.ID, challengeButton)
			go b.verifyUser(m.Chat.ID, u.ID, challengeMsgID)
		}
	}
}

// delUser удаляет пользователя из базы данных.
// Эта функция выполняет следующие действия:
// 1. Проверяет, зарегистрирован ли пользователь.
// 2. Если пользователь является администратором, удаляет его из таблицы администраторов.
// 3. Удаляет все ассоциации пользователь-группа.
// 4. Удаляет пользователя из таблицы пользователей.
// Все операции выполняются в рамках транзакции.
func (b *zrutyBot) delUser(id string) {
	log.Print("Пользователь больше не отслеживается…")
	tx, err := b.db.Begin()
	if err != nil {
		log.Println("❌ Ошибка начала транзакции: %w", err)
		return
	}
	if !b.isUser(id) {
		log.Printf("❌ Пользователь %s не зарегистрирован", id)
		return
	}
	if b.isAdmin(id) {
		_, err = tx.Exec(`DELETE FROM admins WHERE user_id = ?`, id)
		if err != nil {
			log.Println("❌ Ошибка удаления администратора: %w", err)
		} else {
			log.Println("✅ Администратор %w удалён", id)
		}
	}
	_, err = tx.Exec(`DELETE FROM user_chats WHERE chat_id = ?`, id)
	if err != nil {
		log.Println("❌ Ошибка удаления связи user->chat: %w", err)
	} else {
		log.Println("✅ Привязка чатов для пользователя %w удалена", id)
	}
	_, err = tx.Exec(`DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		log.Println("❌ Ошибка удаления пользователя: %w", err)
	} else {
		log.Println("✅ Пользователь %w удалён", id)
	}
	err = tx.Commit()
	if err != nil {
		log.Println("❌ Ошибка завершения транзакции: %w", err)
	}
}

// checkUsers проверяет зарегистрированных пользователей и банит лишних
func (b *zrutyBot) checkUsers() {
	var banAfter int
	err := b.db.QueryRow(`SELECT CAST(value AS INTEGER) FROM settings WHERE key = 'banAfter'`).Scan(&banAfter)
	if err != nil {
		log.Println("❌ Ошибка чтения токена из БД: %w", err)
		return
	}
	rows, err := b.db.Query(`
	SELECT
		u.id,
		u.first_name,
		u.last_name,
		u.username,
		uc.chat_id,
		c.title,
		u.is_bot,
		u.first_seen_at,
		u.check_passed_at  
	FROM users u
	LEFT JOIN user_chats uc ON u.id = uc.user_id
	LEFT JOIN chats c ON uc.chat_id = c.chat_id
	WHERE u.check_passed_at IS NULL
	`)
	if err != nil {
		log.Printf("❌ Ошибка получения пользователей: %v", err)
		return
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var u user
		err := rows.Scan(&u.userID,
			&u.firstName,
			&u.lastName,
			&u.username,
			&u.groupID,
			&u.groupTitle,
			&u.isBot,
			&u.firstSeenAt,
			&u.checkPassedAt,
		)
		if err != nil {
			log.Printf("❌ Ошибка чтения строки: %v", err)
			continue
		}
		users = append(users, u)
	}

	for _, u := range users {
		uid := strconv.Itoa(u.userID)
		inGroup := b.isInGroup(u.groupID, u.userID)
		if !inGroup {
			log.Printf("✅ Пользователь @%s больше не в группе %s", u.username, u.groupTitle)
			b.delUser(uid)
			continue
		}

		// Если пользователь находится в группе, но не прошёл проверку
		if u.firstSeenAt.Valid && time.Since(u.firstSeenAt.Time).Hours() > float64(banAfter) {
			log.Printf("Кикаем пользователя @%s", u.username)

			err = b.client.KickChatMember(u.groupID, u.userID)
			if err != nil {
				log.Printf("❌ Ошибка при удалении пользователя %d @%s: %v", u.userID, u.username, err)
				b.notifyAdmins(fmt.Sprintf(
					`❌ Не удалось удалить пользователя <a href="tg://user?id=%d">%s</a> из группы %s : %v`,
					u.userID, u.firstName, u.groupTitle, err,
				))
			} else {
				b.notifyAdmins(fmt.Sprintf(
					`✅ Пользователь <a href="tg://user?id=%d">%s</a> был удалён из группы %s`,
					u.userID, u.firstName, u.groupTitle,
				))
				kickMessage, err := getSetting(b.db, "kickMessage")
				if err != nil {
					log.Printf("❌ Не удалось прочитать kickMessage: %v", err)
					kickMessage = `Пользователь %s покидает чат: %s`
				}
				_, err = b.client.SendMessage(
					u.groupID,
					fmt.Sprintf(
						kickMessage,
						u.firstName,
						deathCauses[rand.Intn(len(deathCauses))],
					),
				)
				if err != nil {
					log.Printf("❌ Ошибка отправки сообщения в чат: %v", err)
				}
			}
			b.delUser(uid)
			continue
		}
	}
}

// notifyAdmins отправляет системное уведомление всем администраторам бота.
// Получает список администраторов из базы данных и отправляет им сообщение
// через Telegram API. В случае ошибки при получении списка администраторов,
// чтении идентификатора администратора или отправке сообщения, ошибка логируется.
func (b *zrutyBot) notifyAdmins(message string) {
	rows, err := b.db.Query(`SELECT user_id FROM admins`)
	if err != nil {
		log.Printf("❌ notifyAdmins: ошибка при получении списка админов: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Printf("❌ notifyAdmins: ошибка чтения admin user_id %v: %v", id, err)
			continue
		}

		_, err := b.client.SendMessage(
			strconv.Itoa(id),
			fmt.Sprintf("Системное уведомление: %s", message),
			tbot.OptParseModeHTML,
		)
		if err != nil {
			log.Printf("❌ notifyAdmins: ошибка отправки сообщения %d: %v", id, err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("❌ notifyAdmins: ошибка итерации: %v", err)
	}
}

// shutdown завершает работу бота, уведомляя администраторов о его остановке.
// После отправки уведомления, приложение завершает выполнение, используя os.Exit(0).
func (b *zrutyBot) shutdown() {
	b.notifyAdmins("бот остановлен…")
	os.Exit(0)
}
