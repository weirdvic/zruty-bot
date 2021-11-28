package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/yanzay/tbot/v2"
)

func startHandler(m *tbot.Message) {
	if m.Chat.Type == "private" && WB.isAdmin(m.Chat.ID) {
		_, err := WB.Client.SendMessage(m.Chat.ID, greetAdmin)
		if err != nil {
			log.Print(err)
		}
	} else {
		_, err := WB.Client.SendMessage(m.Chat.ID, notAdmin)
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
	if WB.isAdmin(m.Chat.ID) && m.Chat.Type == "private" {
		if len(WB.Users) != 0 {
			report += "```\nЕсть отслеживаемые пользователи:\n\n"
			for _, u := range WB.Users {
				usersCount++
				report += fmt.Sprintf(
					"%d.\t%s %s @%s %.0f мин. назад\n",
					usersCount,
					u.FirstName,
					u.LastName,
					u.Username,
					time.Since(u.FirstSeen).Minutes())
			}
			report += fmt.Sprintf("\nВсего пользователей %d\n```", usersCount)
		} else {
			report += "`Нет отслеживаемых пользователей`"
		}
	}
	_, err := WB.Client.SendMessage(
		m.Chat.ID,
		report,
		tbot.OptParseModeMarkdown,
	)
	if err != nil {
		log.Print(err)
	}
}

func defaultHandler(m *tbot.Message) {
	if m.Chat.Type == "supergroup" || m.Chat.Type == "group" {
		if len(m.NewChatMembers) > 0 {
			WB.addUsers(m)
			WB.welcomeUsers(m)
			return
		} else if u := strconv.Itoa(m.From.ID); WB.isUser(u) {
			WB.Users[u].CheckPassed = true
			log.Printf("Пользователь %s %s(@%s) написал сообщение в чат!",
				WB.Users[u].FirstName,
				WB.Users[u].LastName,
				WB.Users[u].Username,
			)
		}
	}
}
