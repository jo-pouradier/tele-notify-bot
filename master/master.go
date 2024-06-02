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
		"metrics": ProcessMetrics,
	}

	for k, v := range commands {
		b.RegisterHandler(bot.HandlerTypeMessageText, fmt.Sprintf("/%s", k), bot.MatchTypeExact, v)
	}

	return &MasterImpl{
		commands: commands,
		b:        b,
		admins:   []int{id},
	}

}

func handlerDef(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}

func (m *MasterImpl) Serve() {
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
