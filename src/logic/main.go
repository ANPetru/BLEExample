package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/jaracil/ei"
	"github.com/jaracil/goco/ble"
	"github.com/jaracil/psgo"
	_ "github.com/jaracil/psgo/psjs"
)

var document = js.Global.Get("document")
var devices [][]string
var connectedDev *ble.Peripheral

type Characteristic struct {
	Service  string `json:"service"`
	CharacId string `json:"characId"`
}

type UniversalDTO struct {
	Data interface{} `json:"data"`
	// more fields with important meta-data about the message...
}

func main() {
	subscriber := psgo.NewSubscriber(msgSubscriber)
	subscriber.Subscribe("web.scan.bluetooth.devices")
	subscriber.Subscribe("web.stop.bluetooth.devices")
	subscriber.Subscribe("web.connect.bluetooth.devices")
	subscriber.Subscribe("web.read.bluetooth.devices")
	subscriber.Subscribe("web.disconnect.bluetooth.devices")
	subscriber.Subscribe("web.get.connected.bluetooth.devices")

}

func msgSubscriber(msg *psgo.Msg) {
	subscribers := map[string]func(msg *psgo.Msg){
		"web.scan.bluetooth.devices":          scanDevices,
		"web.stop.bluetooth.devices":          stopScaning,
		"web.connect.bluetooth.devices":       connectDevice,
		"web.read.bluetooth.devices":          readCharacteristic,
		"web.disconnect.bluetooth.devices":    disconnectDevice,
		"web.get.connected.bluetooth.devices": getConnectedDevice,
	}
	fmt.Println(msg.To)
	subscribers[msg.To](msg)
}

func getConnectedDevice(msg *psgo.Msg) {
	var result [2]string
	result[0] = connectedDev.Name()
	result[1] = connectedDev.ID()
	msg.Answer(result, nil)
}

func disconnectDevice(msg *psgo.Msg) {
	ble.Disconnect(connectedDev.ID())
	connectedDev = nil

}

func readCharacteristic(msg *psgo.Msg) {
	dtoToSend := UniversalDTO{msg.Dat}
	byteData, _ := json.Marshal(dtoToSend)

	charac := &Characteristic{}
	recivedDTO := UniversalDTO{Data: charac}
	json.Unmarshal(byteData, &recivedDTO)
	response, _ := ble.Read(connectedDev.ID(), charac.Service, charac.CharacId)
	msg.Answer(getStringFromBA(response), nil)
}

func getStringFromBA(arr []byte) []string {
	var result []string
	if strings.Index(string(arr), ",") > -1 {
		ind := 0
		str := ""
		for i := 1; i < len(arr)-1; i++ {

			if arr[i] == 44 {
				fmt.Println(str)
				str = strings.Replace(str, "\"", "", -1)
				fmt.Println(str)

				result = append(result, str)
				ind++
				i++
				str = ""
			}
			if arr[i] < 127 && arr[i] > 28 {
				str += string(arr[i])
			}
		}
	} else {
		str := ""
		for i := range arr {
			if arr[i] < 127 && arr[i] > 28 {
				str += string(arr[i])
			}
		}
		result = append(result, str)
	}
	return result

}

func connectDevice(msg *psgo.Msg) {
	if connectedDev == nil {
		cD, err := ble.Connect(ei.N(msg.Dat).StringZ(), nil)
		connectedDev = cD
		if err == nil {
			msg.Answer(connectedDev.Characteristics(), nil)
		} else {
			msg.Answer("Error", nil)
		}
	} else {
		msg.Answer(connectedDev.Characteristics(), nil)
	}
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
