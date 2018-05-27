package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

func putDoc(docTitle string, payload []byte) {
	url := "http://192.168.1.13:9200/sensors/weather"

	timestamp := time.Now().Format(time.RFC822Z)

	type PayloadMessage struct {
		Temp     float64 `json:"temp"`
		Humidity float64 `json:"humidity"`
	}

	type ElasticPayload struct {
		PayloadMessage
		Title     string `json:"title"`
		Timestamp string `json:"timestamp"`
	}

	fmt.Printf("Received Payload: %s\n", payload)
	var payloadMessage PayloadMessage
	err := json.Unmarshal(payload, &payloadMessage)

	if err != nil {
		log.Fatal("Unable to parse message from sensors.")
	}

	elasticPayload := ElasticPayload{
		PayloadMessage: payloadMessage,
		Title:          docTitle,
		Timestamp:      fmt.Sprintf("%s", timestamp),
	}

	fmt.Printf("Elastic Struct: %+v\n", elasticPayload)

	jsonPayload, err := json.Marshal(elasticPayload)

	if err != nil {
		log.Fatal("Error converting elasticsearch payload to json.")
	}

	fmt.Printf("Elastic Doc: %s\n", jsonPayload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))

	if err != nil {
		log.Fatal("Error posting to elasticsearch")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
	putDoc("testing", msg.Payload())
}

func main() {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("192.168.1.13:1883").SetClientID("mini_server")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe("sensors/testing", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	i := 0
	for i < 1 {
		time.Sleep(30 * time.Second)
	}

}
