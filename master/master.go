package master

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"google.golang.org/grpc"
)

type Master interface {
	Serve()
	RegisterCommands([]any)
	AddAdmin(int)
	AddAdmins([]int)
}

type MasterImpl struct {
	conns map[string]grpc.ClientConn
	// bot *tgbotapi.BotAPI
	b      *bot.Bot
	admins []int
	// all commands and there actions
	commands map[string]bot.HandlerFunc
}

func NewMaster(debug bool) *MasterImpl {
	opts := []bot.Option{
		bot.WithDefaultHandler(handlerDef),
		bot.WithCallbackQueryDataHandler("TESTING", bot.MatchTypeExact, handlerCallback),
	}
	if debug {
		opts = append(opts, bot.WithDebug())
	}
	b, err := bot.New(os.Getenv("TELEGRAM_TOCKEN"), opts...)

	if err != nil {
		log.Panicf("Error while connecting to bot: %v", err)
	}

	id, err := strconv.Atoi(os.Getenv("ADMIN_USER"))
	if err != nil {
		log.Fatal("Admin user Not found, add Env variale : ADMIN_USER")
	}

	ctx := context.Background()
	name, _ := b.GetMyName(ctx, &bot.GetMyNameParams{})
	log.Printf("Authorized on account %s", name)

	commands := map[string]bot.HandlerFunc{
		"metrics":  ProcessMetrics,
		"test":     TestCommand1,
		"keyboard": TestKeyboardButton,
	}

	for k, v := range commands {
		b.RegisterHandler(bot.HandlerTypeMessageText, fmt.Sprintf("/%s", k), bot.MatchTypeExact, v)
	}
	b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{Command: "metrics", Description: "get server metrics"},
			{Command: "test", Description: "just a test"},
			{Command: "keyboard", Description: "get keyboard buttons"},
		},
	})

	return &MasterImpl{
		commands: commands,
		b:        b,
		admins:   []int{id},
	}

}

func handlerDef(ctx context.Context, b *bot.Bot, update *models.Update) {
	log.Printf("Update: %v", update)
	mess := update.Message
	if mess == nil {
		log.Printf("No message: %+v", update.Message)
		return
	}
	id := update.Message.Chat.ID
	if id == 0 {
		log.Printf("Get nil chat id: %v", id)
		id = 2
	}
	msg := update.Message.Text
	if len(msg) == 0 {
		msg = "did not send a message"
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   msg,
	})
}

func handlerCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            "ok",
		// ShowAlert:       true,
	})
}

func (m *MasterImpl) Serve() {
	m.SetMenuButtons()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	m.b.Start(ctx)
}

func ProcessMetrics(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		// Text:      "*" + bot.EscapeMarkdown("metrics...*ok*") + "*",
		Text:      "metrics ok",
		ParseMode: models.ParseModeMarkdown,
	})
}

func TestCommand1(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "here are buttons",
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text:         "testing",
						CallbackData: "TESTING",
					},
				},
			},
		},
	})
}

func TestKeyboardButton(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "here are buttons",
		ReplyMarkup: models.ReplyKeyboardMarkup{
			OneTimeKeyboard: true,
			Keyboard: [][]models.KeyboardButton{
				{
					{Text: "first keyboard button"},
					{Text: "second keyboard button"},
				},
				{
					{Text: "third"},
				},
			},
		},
	})
}

func (m *MasterImpl) SetMenuButtons() {
	ctx := context.Background()
	params := &bot.SetChatMenuButtonParams{
		MenuButton: models.MenuButtonCommands{Type: models.MenuButtonTypeCommands}, // Example: Using default menu button
	}
	ok, err := m.b.SetChatMenuButton(ctx, params)

	if err != nil {
		log.Fatalf("Error setting menu button: %v", err)
		return
	}
	log.Printf("Set Menu Button: %v", ok)

}
