package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pccr10001/jrb-bot/app"
	"github.com/pccr10001/jrb-bot/bot"
	"github.com/pccr10001/jrb-bot/ftx"
	"github.com/pccr10001/jrb-bot/model"
	"log"
	"strings"
)

func main() {

	app.ParseConfig()
	ftx.Init()
	bot.InitMQTT(SignalHandler)
	bot.InitBot()
}

func jsonMarshal(val interface{}) string {
	b, _ := json.Marshal(val)
	return string(b)
}

func SignalHandler(client mqtt.Client, message mqtt.Message) {
	var signal model.Signal

	err := json.Unmarshal(message.Payload(), &signal)
	if err != nil {
		log.Printf("failed to decode signal: %s\n", string(message.Payload()))
		return
	}

	if ftx.CoinMap[signal.Symbol] == "" {
		return
	}

	bot.SendMessage(fmt.Sprintf("Signal Received: %s %s", signal.Symbol, signal.Direction))
	position, err := ftx.GetPosition(signal.Symbol)
	if err != nil {
		log.Println(err)
		bot.SendMessage(fmt.Sprintf("Failed to receive position: %s", signal.Symbol))
		return
	}

	if position != nil {
		bot.SendMessage(fmt.Sprintf("Current Position: %s %s %s @ %s",
			betterFormat(position.Size),
			position.Future,
			position.Side,
			betterFormat(position.RecentAverageOpenPrice),
		))

		if (position.Side == "buy" && signal.Direction == "up") ||
			(position.Side == "sell" && signal.Direction == "down") {
			return
		}

		order, err := ftx.FillPosition(*position)
		if err != nil {
			bot.SendMessage(fmt.Sprintf("Failed to fill position: %s", position.Future))
			log.Println(err)
			return
		}
		earn := 0.0
		if position.Side == ftx.SIDE_BUY {
			earn = (order.AvgFillPrice - position.RecentAverageOpenPrice) * order.Size
		} else {
			earn = (position.RecentAverageOpenPrice - order.AvgFillPrice) * order.Size
		}
		bot.SendMessage(fmt.Sprintf("Position filled: %s %s %s @ %s EARN %s",
			betterFormat(position.Size),
			position.Future,
			position.Side,
			betterFormat(order.AvgFillPrice),
			betterFormat(earn),
		))
	}

	side := ftx.SIDE_SELL
	if signal.Direction == "up" {
		side = ftx.SIDE_BUY
	}
	order, err := ftx.PlaceOrder(signal.Symbol, 100, side)
	if err != nil {
		log.Println(err)
		bot.SendMessage(fmt.Sprintf("Failed to place position: %s", position.Future))
		return
	}

	bot.SendMessage(fmt.Sprintf("Position placed: %s %s %s @ %s",
		betterFormat(order.Size),
		order.Future,
		order.Side,
		betterFormat(order.AvgFillPrice),
	))
}

func betterFormat(num float64) string {
	s := fmt.Sprintf("%.6f", num)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}
