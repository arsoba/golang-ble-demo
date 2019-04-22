package main

import (
	"ble/mqttclient"
	"ble/thigy"
	"crypto/rand"
	"log"
	"runtime"
	"time"

	"github.com/tarm/serial"
)

var done = make(chan struct{})

var mqqtClient mqttclient.MqttClient
var topic string

func main() {
	thigyInstance := thigy.NewThigy()
	go thigyInstance.Init()

	go listenSerial(thigyInstance)

	go func() {
		bytes := make([]byte, 4)
		bytes[0] = 1
		bytes[1] = 255
		bytes[2] = 0
		bytes[3] = 0
		thigyInstance.Msg <- thigy.BleMsg{
			thigy.UisUUID,
			thigy.UisLedUUID,
			bytes,
		}
	}()

	mqqtClient = mqttclient.NewClient()
	mqqtClient.Connect()
	topic = "temp"
	go publishState(thigyInstance)

	<-done
}

func publishState(thigy *thigy.Thigy) {
	for {
		runtime.Gosched()
		if thigy.Connected {
			mqqtClient.Publish(topic, thigy.ThigyState)
			time.Sleep(time.Duration(2000) * time.Millisecond)
		}
	}
}

func listenSerial(thigyInstance *thigy.Thigy) {
	c := &serial.Config{Name: "/dev/ttyACM0", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	for {
		buf := make([]byte, 5)
		_, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		bytes := make([]byte, 4)
		rand.Read(bytes)
		bytes[0] = 1

		thigyInstance.Msg <- thigy.BleMsg{
			thigy.UisUUID,
			thigy.UisLedUUID,
			bytes,
		}
	}
}
