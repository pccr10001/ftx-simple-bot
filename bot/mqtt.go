package bot

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pccr10001/jrb-bot/app"
	"io/ioutil"
	"log"
)

var MQTT mqtt.Client

func InitMQTT(handler mqtt.MessageHandler) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:%d", app.AppConfig.Mqtt.Server, app.AppConfig.Mqtt.Port))
	opts.SetClientID(app.AppConfig.Mqtt.ClientId)
	opts.SetUsername(app.AppConfig.Mqtt.Username)
	opts.SetPassword(app.AppConfig.Mqtt.Password)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetTLSConfig(NewTlsConfig())
	MQTT = mqtt.NewClient(opts)
	if token := MQTT.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	log.Println("MQTT Connected")

	if token := MQTT.Subscribe(app.AppConfig.Mqtt.Topic, 2, handler); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Println("MQTT Subscribed")
}

func NewTlsConfig() *tls.Config {
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(app.AppConfig.Mqtt.CAFile)
	if err != nil {
		log.Fatalln(err.Error())
	}
	certpool.AppendCertsFromPEM(ca)
	return &tls.Config{
		RootCAs: certpool,
	}
}
