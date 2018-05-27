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

func putDoc(payload []byte) {
	elasticUrl := "elasticsearch:9200/sensors/weather"
	timestamp := time.Now().Format(time.RFC822Z)

	type PayloadMessage struct {
		Title 	 string  `json:"title"`
		Temp     float64 `json:"temp"`
		Humidity float64 `json:"humidity"`
	}

	type ElasticPayload struct {
		PayloadMessage
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
		Timestamp:      fmt.Sprintf("%s", timestamp),
	}

	fmt.Printf("Elastic Struct: %+v\n", elasticPayload)

	jsonPayload, err := json.Marshal(elasticPayload)

	if err != nil {
		log.Fatal("Error converting elasticsearch payload to json.")
	}

	fmt.Printf("Elastic Doc: %s\n", jsonPayload)

	req, err := http.NewRequest("POST", elasticUrl, bytes.NewBuffer(jsonPayload))

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
	opts := mqtt.NewClientOptions().AddBroker("mosquitto:1883").SetClientID("mini_server")
	opts.SetKeepAlive(15 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(10 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe("sensors/house", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	i := 0
	for i < 1 {
		time.Sleep(30 * time.Second)
	}

}
