package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/yanzay/tbot/v2"
)

// isAdmin проверяет является ли пользователь одним из админов бота
func (b *zrutyBot) isAdmin(id string) bool {
	if _, ok := b.Admins[id]; ok {
		return true
	}
	return false
}

// isUser проверяет есть ли пользователь среди отслеживаемых
func (b *zrutyBot) isUser(id string) bool {
	if _, ok := b.Users[id]; ok {
		return true
	}
	return false
}

// isInGroup позволяет проверить состоит ли пользователь в группе
func (b *zrutyBot) isInGroup(groupID string, userID int) bool {
	gcm, err := b.Client.GetChatMember(groupID, userID)
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
	var (
		users          = m.NewChatMembers
		usersAdded int = 0
	)
	for _, u := range users {
		uid := strconv.Itoa(u.ID)
		if !zruty.isUser(uid) {
			b.Users[uid] = &User{
				ID:          uid,
				FirstName:   u.FirstName,
				LastName:    u.LastName,
				Username:    u.Username,
				Groups:      make(map[string]string),
				IsBot:       u.IsBot,
				FirstSeen:   time.Now(),
				CheckPassed: false,
			}
			b.Users[uid].Groups[m.Chat.ID] = m.Chat.Title
			log.Printf("Добавлен новый пользователь %s %s(%s)",
				u.FirstName,
				u.LastName,
				u.Username,
			)
			usersAdded++
		} else {
			log.Printf("Пользователь %s %s(%s) уже отслеживается…",
				u.FirstName,
				u.LastName,
				u.Username,
			)
			if _, ok := b.Users[uid].Groups[m.Chat.ID]; !ok {
				b.Users[uid].Groups[m.Chat.ID] = m.Chat.Title
			}
		}
	}
	log.Printf("Добавлено пользователей: %v / %v", usersAdded, len(users))
}

// welcomeUsers отправляет новым пользователям приветственное сообщение
func (b *zrutyBot) welcomeUsers(m *tbot.Message) {
	log.Printf("Chat ID is: %v", m.Chat.ID)
	var users = m.NewChatMembers
	for _, u := range users {
		_, err := b.Client.SendMessage(
			m.Chat.ID,
			fmt.Sprintf(
				b.WelcomeMessage,
				u.FirstName,
			),
		)
		if err != nil {
			log.Print(err)
		}
	}
}

// delUser удаляет пользователя из списка отслеживаемых
func (b *zrutyBot) delUser(id string) {
	delete(b.Users, id)
	log.Print("Пользователь больше не отслеживается…")
}

// checkUsers проверяет зарегистрированных пользователей и банит лишних
func (b *zrutyBot) checkUsers() {
	// Проходим циклом по всем зарегистрированным пользователям
	for id, u := range b.Users {
		// Для каждого пользователя обходим его группы и проверяем,
		// состоит ли пользователь в группе
		uid, err := strconv.Atoi(id)
		if err != nil {
			log.Print(err)
		}
		for gid, gTitle := range u.Groups {
			if !b.isInGroup(gid, uid) {
				log.Printf(
					"Пользователь @%s уже не в группе %s",
					u.Username,
					gTitle,
				)
				b.delUser(id)
			} else if !u.CheckPassed {
				if int(time.Since(u.FirstSeen).Hours()) > b.BanAfter {
					log.Printf(
						"Кикаем пользователя @%s",
						u.Username)
					_, err := b.Client.SendMessage(
						gid,
						fmt.Sprintf(
							kickMessage,
							u.FirstName,
							deathCauses[rand.Intn(len(deathCauses))],
						),
					)
					if err != nil {
						log.Print(err)
					}
					err = b.Client.KickChatMember(gid, uid)
					if err != nil {
						log.Print(err)
					}
					b.delUser(id)
					return
				}
			} else if u.CheckPassed {
				log.Print("Пользователь уже отправлял сообщения в чат")
				b.delUser(id)
			}
		}
	}
}

func (b *zrutyBot) makeBackup() {
	backup, err := json.MarshalIndent(b, "", " ")
	if err != nil {
		log.Print(err)
	}
	err = ioutil.WriteFile("config.json", backup, 0644)
	if err != nil {
		log.Print(err)
	}
	log.Print("Настройки сохранены.")
}

func (b *zrutyBot) restoreBackup() {
	backup, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Print(err)
	}
	log.Print("Чтение файла настроек успешно")
	err = json.Unmarshal(backup, b)
	if err != nil {
		log.Print(err)
	}
}

func (b *zrutyBot) shutdown() {
	b.makeBackup()
	for id := range b.Admins {
		_, err := b.Client.SendMessage(
			id,
			"Бот остановлен…",
		)
		if err != nil {
			log.Print(err)
		}
	}
	os.Exit(0)
}
