package main

import (
	_ "embed"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yanzay/tbot/v2"
)

type User struct {
	ID          string            `json:"id"`
	FirstName   string            `json:"first_name"`
	LastName    string            `json:"last_name"`
	Username    string            `json:"username"`
	Groups      map[string]string `json:"groups"`
	IsBot       bool              `json:"is_bot"`
	FirstSeen   time.Time         `json:"first_seen"`
	CheckPassed bool              `json:"check_passed"`
}

type zrutyBot struct {
	// Токен Telegram бота
	Token string `json:"token"`
	// Объект Client бота для удобства
	Client tbot.Client
	// Количество часов, после которых аккаунт нужно банить
	BanAfter int `json:"ban_after"`
	// Список администраторов бота
	Admins map[string]*User `json:"admins"`
	// Список отслеживаемых пользователей
	Users map[string]*User `json:"users"`
	// Приветственное сообщение
	WelcomeMessage string `json:"welcome_message"`
}

var (
	zruty zrutyBot
)

func init() {
	if zruty.Admins == nil {
		zruty.Admins = make(map[string]*User)
	}
	if zruty.Users == nil {
		zruty.Users = make(map[string]*User)
	}
	if _, err := os.Stat("config.json"); err == nil {
		zruty.restoreBackup()
	} else {
		log.Fatalf("Недоступен файл конфигурации: %v", err)
	}
	// Инициализация ГПСЧ
	rand.Seed(time.Now().UnixNano())

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(b *zrutyBot) {
		<-c
		b.shutdown()
	}(&zruty)
}

func main() {
	bot := tbot.New(zruty.Token)
	log.Print("Бот создан…")
	zruty.Client = *bot.Client()

	if _, err := os.Stat("config.json"); err == nil {
		log.Print("config.json обнаружен…")
		zruty.restoreBackup()
	}

	bot.HandleMessage(`.*start.*`, startHandler)
	bot.HandleMessage(`^/report.*`, reportHandler)
	bot.HandleMessage(`^/backup.*`, backupHandler)
	bot.HandleMessage(``, defaultHandler)

	go func(bot *tbot.Server) {
		err := bot.Start()
		if err != nil {
			log.Fatal(err)
		}
	}(bot)
	log.Print("Бот запущен…")

	// Функция для запуска проверки пользователей с определённым интервалом
	func(b *zrutyBot) {
		for {
			time.Sleep(1 * time.Minute)
			b.checkUsers()
		}
	}(&zruty)
}
