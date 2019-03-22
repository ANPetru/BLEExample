package main

import (
	"fmt"

	"github.com/jaracil/ei"

	"github.com/gopherjs/gopherjs/js"
	"github.com/jaracil/goco/ble"
	"github.com/jaracil/psgo"
	_ "github.com/jaracil/psgo/psjs"
)

var document = js.Global.Get("document")
var devices [][]string

func main() {
	subscriber := psgo.NewSubscriber(msgSubscriber)
	subscriber.Subscribe("web.scan.bluetooth.devices")
	subscriber2 := psgo.NewSubscriber(msgSubscriber)
	subscriber2.Subscribe("web.stop.bluetooth.devices")
	subscriber3 := psgo.NewSubscriber(msgSubscriber)
	subscriber3.Subscribe("web.connect.bluetooth.devices")
}

func msgSubscriber(msg *psgo.Msg) {
	subscribers := map[string]func(msg *psgo.Msg){
		"web.scan.bluetooth.devices":    scanDevices,
		"web.stop.bluetooth.devices":    stopScaning,
		"web.connect.bluetooth.devices": connectDevice,
	}
	subscribers[msg.To](msg)
}

func connectDevice(msg *psgo.Msg) {
	devID := msg.Dat
	fmt.Println("connect main")
	connectedDev, err := ble.Connect(ei.N(devID).StringZ(), onDisconnect)
	if err == nil {
		msg.Answer(connectedDev.Characteristics(), nil)
		fmt.Println(connectedDev.Name())
	} else {
		fmt.Println("Couldn't connect")
	}
}

func onDisconnect(p *ble.Peripheral) {
	fmt.Println("Disconnected")
}

func scanDevices(msg *psgo.Msg) {
	devices = nil
	opt := []string{}
	ble.StartScan(opt, onFoundDevice, false)
}

func stopScaning(msg *psgo.Msg) {
	err := ble.StopScan()
	if err != nil {
		devices = nil
		devices = append(devices, []string{"Error", "Error"})
	}
	msg.Answer(devices, nil)

}

func onFoundDevice(p *ble.Peripheral) {
	if p.Name() == "undefined" {
		return
	}
	devices = append(devices, []string{p.Name(), p.ID()})
}
