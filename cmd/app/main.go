package main

import (
	"context"
	"fmt"

	"main/internal/announcement"
	"main/internal/callback"
	chatstatus "main/internal/chatStatus"
	"main/internal/commands"
	"main/internal/config"
	"main/internal/core"
	"main/internal/messages"
	"main/internal/queries"
	"main/internal/users"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {

	err := config.Init()
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, config.AppConfig.DatabaseDNS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	if err := dbpool.Ping(ctx); err != nil {
		panic(err)
	}
	queriesСlient := queries.New(dbpool)

	chat := chatstatus.NewChatStatus(queriesСlient)

	announcement := announcement.NewAnnouncement(queriesСlient)

	// Create bot and enable debugging info
	bot, err := telego.NewBot(config.AppConfig.BotToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	usersStorage := users.NewUsersClient(queriesСlient, config.AppConfig.LoginAdmin, bot)

	core := core.NewCore(usersStorage, chat, announcement, bot)
	commands := commands.NewCommands(bot, usersStorage, chat, core)
	callback := callback.NewCallBack(queriesСlient, announcement, usersStorage, bot, core, chat, commands)

	messages := messages.NewMessages(announcement, usersStorage, chat, commands, core, bot)

	// Call method getMe (https://core.telegram.org/bots/api#getme)
	botUser, err := bot.GetMe()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Print Bot information
	fmt.Printf("Bot user: %+v\n", botUser)

	// Get updates channel
	updates, _ := bot.UpdatesViaLongPolling(nil)

	// Stop reviving updates from updates channel
	defer bot.StopLongPolling()

	// Loop through all updates when they came
	for update := range updates {

		// press button
		if update.CallbackQuery != nil {
			if update.CallbackQuery.From.Username == config.AppConfig.LoginAdmin {
				core.SaveAdminChatID(update.CallbackQuery.From.ID)
			}
			go callback.CallBack(ctx, update.CallbackQuery)
			continue
		}

		// Check if update contains message
		if update.Message != nil {
			if update.Message.From.Username == config.AppConfig.LoginAdmin {
				core.SaveAdminChatID(update.Message.From.ID)
			}

			// Get chat ID from message
			chatID := tu.ID(update.Message.Chat.ID)

			if update.Message.Text == "/start" {
				go commands.Start(ctx, chatID, update.Message.From.ID, config.AppConfig.Rules)
			} else if update.Message.Text == "/add" {
				go commands.Add(ctx, chatID, update.Message.From.ID)
			} else if update.Message.Text == "/delete" {
				go commands.Delete(ctx, chatID, update.Message.From.ID)

			} else { // some text

				status := chat.Get(ctx, update.Message.From.ID)

				switch status {

				case 1:
					go messages.Text(ctx, update.Message) // text input

				case 2:
					go messages.Photo(ctx, update.Message)

				case 3:
					go messages.Number(ctx, update.Message)

				default:
					core.SendDefaultMessage(chatID)

				}
			}
		}
	}

}
