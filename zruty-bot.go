package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/yanzay/tbot/v2"
)

type User struct {
	ID          string            `json:"id"`
	FirstName   string            `json:"first_name"`
	LastName    string            `json:"last_name"`
	Username    string            `json:"username"`
	Groups      map[string]string `json:"groups"`
	IsBot       bool
	FirstSeen   time.Time
	CheckPassed bool
}

type botConfig struct {
	// Токен Telegram бота
	Token string `json:"token"`
	// Объект Client бота для удобства
	Client tbot.Client
	// Количество часов, после которых аккаунт нужно банить
	BanAfter int `json:"ban_after"`
	// Список администраторов бота
	Admins map[string]*User `json:"admins"`
	// Список отслеживаемых пользователей
	Users map[string]*User
}

var (
	// Встраиваем файл конфигурации в бинарник при компиляции
	//go:embed config.json
	configFile []byte
	WB         botConfig
)

func init() {
	err := json.Unmarshal(configFile, &WB)
	if err != nil {
		log.Fatal(err)
	}
	if WB.Admins == nil {
		WB.Admins = make(map[string]*User)
	}
	if WB.Users == nil {
		WB.Users = make(map[string]*User)
	}
	rand.Seed(time.Now().UnixNano())
}

func main() {
	bot := tbot.New(WB.Token)
	log.Print("Bot created…")
	WB.Client = *bot.Client()
	bot.HandleMessage(`.*start.*`, startHandler)
	bot.HandleMessage(`^/report.*`, reportHandler)
	bot.HandleMessage(``, defaultHandler)

	go func(bot *tbot.Server) {
		err := bot.Start()
		if err != nil {
			log.Fatal(err)
		}
	}(bot)
	log.Print("Bot started…")

	// Функция для запуска проверки пользователей с определённым интервалом
	func(b *botConfig) {
		for {
			time.Sleep(1 * time.Hour)
			b.checkUsers()
		}
	}(&WB)
}
