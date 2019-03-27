package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jaracil/ei"
	"github.com/jaracil/goco/ble"
	"github.com/jaracil/psgo"
	_ "github.com/jaracil/psgo/psjs"
)

var devices [][]string
var connectedDev *ble.Peripheral

type CharacteristicToRead struct {
	Service  string `json:"service"`
	CharacID string `json:"characId"`
}

type CharacteristicToWrite struct {
	Service  string `json:"service"`
	CharacID string `json:"characId"`
	Message  string `json:"message"`
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
	subscriber.Subscribe("web.write.bluetooth.devices")
	subscriber.Subscribe("web.disconnect.bluetooth.devices")
	subscriber.Subscribe("web.get.connected.bluetooth.devices")
	subscriber.Subscribe("web.rssi.get.bluetooth.devices")

}

func msgSubscriber(msg *psgo.Msg) {
	subscribers := map[string]func(msg *psgo.Msg){
		"web.scan.bluetooth.devices":          scanDevices,
		"web.stop.bluetooth.devices":          stopScaning,
		"web.connect.bluetooth.devices":       connectDevice,
		"web.read.bluetooth.devices":          readCharacteristic,
		"web.disconnect.bluetooth.devices":    disconnectDevice,
		"web.get.connected.bluetooth.devices": getConnectedDevice,
		"web.write.bluetooth.devices":         writeToCharac,
		"web.rssi.get.bluetooth.devices":      getRSSI,
	}
	fmt.Println(msg.To)
	subscribers[msg.To](msg)
}

func getRSSI(msg *psgo.Msg) {
	rssi, err := ble.ReadRSSI(connectedDev.ID())
	if err != nil {
		msg.Answer("Error", nil)
	} else {
		msg.Answer(fmt.Sprintf("%d", rssi), nil)
	}
}

func writeToCharac(msg *psgo.Msg) {
	if ble.IsConnected(connectedDev.ID()) {
		dtoToSend := UniversalDTO{msg.Dat}
		byteData, _ := json.Marshal(dtoToSend)
		charac := &CharacteristicToWrite{}
		recivedDTO := UniversalDTO{Data: charac}
		json.Unmarshal(byteData, &recivedDTO)
		fmt.Println("Writing message: " + charac.Message + " to " + connectedDev.ID() + " - " + charac.Service + " - " + charac.CharacID)
		err := ble.Write(connectedDev.ID(), charac.Service, charac.CharacID, []byte(charac.Message))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Message s1ent successfully")
		}
	} else {
		//device is disconnected
	}
}

func getConnectedDevice(msg *psgo.Msg) {
	var result [2]string
	result[0] = connectedDev.Name()
	result[1] = connectedDev.ID()
	msg.Answer(result, nil)
}

func disconnectDevice(msg *psgo.Msg) {
	if connectedDev != nil {
		err := ble.Disconnect(connectedDev.ID())
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("Removing connected device")
		connectedDev = nil
	}
}

func readCharacteristic(msg *psgo.Msg) {
	if ble.IsConnected(connectedDev.ID()) {
		dtoToSend := UniversalDTO{msg.Dat}
		byteData, _ := json.Marshal(dtoToSend)
		charac := &CharacteristicToRead{}
		recivedDTO := UniversalDTO{Data: charac}
		json.Unmarshal(byteData, &recivedDTO)
		response, _ := ble.Read(connectedDev.ID(), charac.Service, charac.CharacID)
		fmt.Println(response)
		msg.Answer(getStringFromBA(response), nil)
	} else {
		msg.Answer("disconnected", nil)
	}

}

func getStringFromBA(arr []byte) []string {
	var result []string
	if strings.Index(string(arr), ",") > -1 {
		ind := 0
		str := ""
		for i := 1; i < len(arr)-1; i++ {

			if arr[i] == 44 {
				str = strings.Replace(str, "\"", "", -1)
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
		answered := false
		time.AfterFunc(5*time.Second, func() {
			if !answered {
				answered = true
				msg.Answer("Error", nil)
			}
		})
		cD, err := ble.Connect(ei.N(msg.Dat).StringZ(), nil)
		if err == nil && !answered {
			connectedDev = cD
			msg.Answer(connectedDev.Characteristics(), nil)
		} else {
			msg.Answer("Error", nil)
		}
		answered = true
	} else {
		msg.Answer(connectedDev.Characteristics(), nil)
	}
}

func scanDevices(msg *psgo.Msg) {
	devices = nil
	ble.StartScan([]string{}, onFoundDevice, false)
}

func stopScaning(msg *psgo.Msg) {
	err := ble.StopScan()
	if err != nil {
		devices = nil
		devices = append(devices, []string{"Error", "Error"})
		devices = append(devices, []string{"Error", "Error"})
		devices = append(devices, []string{"Error", "Error"})
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
