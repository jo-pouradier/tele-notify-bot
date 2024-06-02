package master

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/grpc"
)

type Master interface {
	Serve()
	RegisterCommands([]any)
	AddAdmin(int)
	AddAdmins([]int)
}

type MasterImpl struct {
	conns  map[string]grpc.ClientConn
	botApi *tgbotapi.BotAPI
	admins []int
	// all commands and there actions
	commands map[string]func(tgbotapi.Update) string
}

func NewMaster(debug bool) *MasterImpl {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOCKEN"))
	if err != nil {
		log.Panic("Telegram key not found")
	}
	id, err := strconv.Atoi(os.Getenv("ADMIN_USER"))
	log.Printf("admin: %d", id)

	if err != nil {
		log.Fatal("Admin user Not found, add Env variale : ADMIN_USER")
	}
	if err != nil {
		log.Panic("Error with ADMIN env ID")
	}

	bot.Debug = debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	commands := map[string]func(tgbotapi.Update) string{
		"metrics": ProcessMetrics,
	}

	return &MasterImpl{
		commands: commands,
		botApi:   bot,
		admins:   []int{id},
	}

}

func (m *MasterImpl) Serve() {
	// m.setCommands()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := m.botApi.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			str := m.processMessage(update)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, str)
			m.botApi.Send(msg)
		}
	}
}

// func (m *MasterImpl) setCommands() {
// 	var cmds = []string{}
// 	for k := range m.commands {
// 		cmds = append(cmds, k)
// 	}
// 	m.botApi
// }

func (m *MasterImpl) processMessage(u tgbotapi.Update) string {
	// check user ID
	fromID := int(u.SentFrom().ID)
	for _, id := range m.admins {
		if fromID == id {
			continue
		} else {
			log.Printf("Got message from unauthorized %d", fromID)
			return "Your are not authorized"
		}
	}

	// check if command
	if !u.Message.IsCommand() {
		return "Only commands are accepted"
	}

	return m.commands[u.Message.Command()](u)

}

func ProcessMetrics(u tgbotapi.Update) string {
	// m := pb.NewMetricsServiceClient(m.conns)
	// metricsCtx, metricsCancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer metricsCancel()
	// metrics, _ := m.Metrics(metricsCtx, &pb.Empty{})
	// log.Printf("get metrics: %s", metrics)
	return "metrics..."
}
