package bot

import (
	"github.com/pccr10001/jrb-bot/app"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

var Bot *tb.Bot

func InitBot() {
	var err error
	Bot, err = tb.NewBot(tb.Settings{
		Token:  app.AppConfig.Telegram.Token,
		Poller: &tb.LongPoller{Timeout: 5 * time.Second},
	})
	if err != nil {
		log.Fatalln(err)
	}

	Bot.Handle(tb.OnText, func(m *tb.Message) {
		log.Printf("From %s (%d) :%s", m.Sender.Username, m.Sender.ID, m.Text)
		//Bot.Send(m.Sender, "Hello World!")
	})

	Bot.Start()
}

func SendMessage(msg string) {
	log.Printf("TG Sent: %s\n", msg)
	Bot.Send(&tb.User{ID: 622343903}, msg)
}
