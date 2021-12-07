package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/yanzay/tbot/v2"
)

func startHandler(m *tbot.Message) {
	if m.Chat.Type == "private" && zruty.isAdmin(m.Chat.ID) {
		_, err := zruty.Client.SendMessage(m.Chat.ID, greetAdmin)
		if err != nil {
			log.Print(err)
		}
	} else {
		_, err := zruty.Client.SendMessage(m.Chat.ID, notAdmin)
		if err != nil {
			log.Print(err)
		}
	}
}

func reportHandler(m *tbot.Message) {
	var (
		report     string
		usersCount int = 0
	)
	if zruty.isAdmin(m.Chat.ID) && m.Chat.Type == "private" {
		if len(zruty.Users) != 0 {
			report += "```\nЕсть отслеживаемые пользователи:\n\n"
			for _, u := range zruty.Users {
				var userChats []string
				usersCount++
				for _, title := range u.Groups {
					userChats = append(userChats, title)
				}
				report += fmt.Sprintf(
					"%d.\t%s %s @%s %s назад\nВ чатах: %s\n",
					usersCount,
					u.FirstName,
					u.LastName,
					u.Username,
					time.Since(u.FirstSeen),
					strings.Join(userChats, ", "),
				)
			}
			report += fmt.Sprintf("\nВсего пользователей %d\n```", usersCount)
		} else {
			report += "`Нет отслеживаемых пользователей`"
		}
	}
	_, err := zruty.Client.SendMessage(
		m.Chat.ID,
		report,
		tbot.OptParseModeMarkdown,
	)
	if err != nil {
		log.Print(err)
	}
}

func backupHandler(m *tbot.Message) {
	if m.Chat.Type == "private" && zruty.isAdmin(m.Chat.ID) {
		zruty.makeBackup()
		_, err := zruty.Client.SendMessage(
			fmt.Sprint(m.From.ID),
			"Core dumped",
		)
		if err != nil {
			log.Print(err)
		}
	}
}

func defaultHandler(m *tbot.Message) {
	if zruty.isValidGroup(m.Chat.ID) &&
		(m.Chat.Type == "supergroup" || m.Chat.Type == "group") {
		if len(m.NewChatMembers) > 0 {
			zruty.addUsers(m)
			zruty.welcomeUsers(m)
			return
		} else if u := strconv.Itoa(m.From.ID); zruty.isUser(u) {
			zruty.Users[u].CheckPassed = true
			zruty.notifyAdmins(
				fmt.Sprintf(
					"пользователь @%s прошёл проверку",
					zruty.Users[u].Username,
				),
			)
			zruty.delUser(u)
			log.Printf("Пользователь %s %s(@%s) написал сообщение в чат!",
				zruty.Users[u].FirstName,
				zruty.Users[u].LastName,
				zruty.Users[u].Username,
			)
		}
	}
}
