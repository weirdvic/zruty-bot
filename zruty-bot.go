package main

import (
	"database/sql"
	_ "embed"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/yanzay/tbot/v2"
)

type user struct {
	userID        int
	firstName     string
	lastName      string
	username      string
	groupID       string
	groupTitle    string
	isBot         bool
	firstSeenAt   sql.NullTime
	checkPassedAt sql.NullTime
}

type zrutyBot struct {
	// Токен Telegram бота
	token string
	// Объект Client бота для удобства
	client tbot.Client
	// Объект базы данных
	db *sql.DB
}

var (
	zruty      zrutyBot
	logUpdates bool
)

// main инициализирует и запускает бота, выполняя следующие действия:
// 1. Открывает подключение к базе данных и применяет миграции.
// 2. Инициализирует объект бота, загружая конфигурацию.
// 3. Настраивает обработку системных сигналов для корректного завершения работы.
// 4. Создаёт и настраивает объект Telegram бота, регистрируя команды и обработчики сообщений.
// 5. Запускает сервер бота и уведомляет администраторов о начале работы.
// 6. Периодически проверяет пользователей, применяя правила бана.
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logUpdates = strings.ToLower(os.Getenv("LOG_UPDATES")) == "true"
	db := openDB("zruty.sqlite3")
	defer db.Close()
	if err := applyMigrations(db); err != nil {
		log.Fatalf("❌ Ошибка применения миграций: %v", err)
	}
	zruty.db = db
	if err := zruty.init(); err != nil {
		log.Fatalf("❌ Ошибка инициализации бота: %v", err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(b *zrutyBot) {
		<-c
		b.shutdown()
	}(&zruty)
	bot := tbot.New(zruty.token)
	log.Print("✅ Бот создан…")
	zruty.client = *bot.Client()

	bot.Use(updatesHandler)
	bot.HandleMessage(`^/start.*`, startHandler)
	bot.HandleMessage(`^/report.*`, reportHandler)
	bot.HandleMessage(`^/underAttackSwitch.*`, underAttackSwitchHandler)
	bot.HandleMessage(``, defaultHandler)
	bot.HandleCallback(zruty.callbackHandler)

	go func(bot *tbot.Server) {
		err := bot.Start()
		if err != nil {
			log.Fatal(err)
		}
	}(bot)
	log.Print("🚀 Бот запущен…")
	zruty.notifyAdmins("😎 Бот начал работу…")

	// Функция для запуска проверки пользователей с определённым интервалом
	func(b *zrutyBot) {
		for {
			time.Sleep(1 * time.Minute)
			b.checkUsers()
		}
	}(&zruty)
}
