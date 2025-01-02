package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jasonlvhit/gocron"
)

var mqttClient mqtt.Client

func main() {
	checkForEnv()
	mqttClient = initMqtt()
	autodiscovery()
	pollValue()
	go executeJob()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	subscribeMqtt()
	<-c

}

func autodiscovery() {
	stove := callStove()
	hspDevice := HspDevice{
		Ids:             []string{stove.Meta.SerialNumber},
		Name:            "HSP Pallazza",
		Manufacturer:    "Haas+Sohn",
		Model:           stove.Meta.StoveType,
		SoftwareVersion: stove.Meta.SoftwareVersion,
	}

	var targetTemp = HspSensorDiscovery{
		Name:              fmt.Sprintf("HSP %s Target Temperature", stove.Meta.SerialNumber),
		UniqueId:          fmt.Sprintf("hsp-%s-target_temp", stove.Meta.SerialNumber),
		UnitOfMeasurement: "°C",
		DeviceClass:       "temperature",
		ForceUpdate:       true,
		StateTopic:        fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber),
		ValueTemplate:     "{{ value_json['sp_temp'] }}",
		Device:            hspDevice,
	}

	jsonValue, _ := json.Marshal(targetTemp)
	mqttClient.Publish(fmt.Sprintf("homeassistant/sensor/hsp_%s_target_temp/config", stove.Meta.SerialNumber), 1, true, jsonValue)

	var currentTemp = HspSensorDiscovery{
		Name:              fmt.Sprintf("HSP %s Current Temperature", stove.Meta.SerialNumber),
		UniqueId:          fmt.Sprintf("hsp-%s-is_temp", stove.Meta.SerialNumber),
		UnitOfMeasurement: "°C",
		DeviceClass:       "temperature",
		ForceUpdate:       true,
		StateTopic:        fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber),
		ValueTemplate:     "{{ value_json['is_temp'] }}",
		Device:            hspDevice,
	}
	jsonValue, _ = json.Marshal(currentTemp)
	mqttClient.Publish(fmt.Sprintf("homeassistant/sensor/hsp_%s_current_temp/config", stove.Meta.SerialNumber), 1, true, jsonValue)

	var tvlTemp = HspSensorDiscovery{
		Name:              fmt.Sprintf("HSP %s Forward Flow Temperature", stove.Meta.SerialNumber),
		UniqueId:          fmt.Sprintf("hsp-%s-tvl_temp", stove.Meta.SerialNumber),
		UnitOfMeasurement: "°C",
		DeviceClass:       "temperature",
		ForceUpdate:       true,
		StateTopic:        fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber),
		ValueTemplate:     "{{ value_json['tvl_temp'] }}",
		Device:            hspDevice,
	}
	jsonValue, _ = json.Marshal(tvlTemp)
	mqttClient.Publish(fmt.Sprintf("homeassistant/sensor/hsp_%s_forward_flow_temp/config", stove.Meta.SerialNumber), 1, true, jsonValue)

	var heatCurve = HspSensorDiscovery{
		Name:          fmt.Sprintf("HSP %s Heating Curve", stove.Meta.SerialNumber),
		UniqueId:      fmt.Sprintf("hsp-%s-ht_char", stove.Meta.SerialNumber),
		ForceUpdate:   true,
		StateTopic:    fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber),
		ValueTemplate: "{{ value_json['ht_char'] }}",
		Device:        hspDevice,
	}
	jsonValue, _ = json.Marshal(heatCurve)
	mqttClient.Publish(fmt.Sprintf("homeassistant/sensor/hsp_%s_heating_curve/config", stove.Meta.SerialNumber), 1, true, jsonValue)

	var ignitions = HspSensorDiscovery{
		Name:          fmt.Sprintf("HSP %s Ignitions", stove.Meta.SerialNumber),
		UniqueId:      fmt.Sprintf("hsp-%s-ignitions", stove.Meta.SerialNumber),
		ForceUpdate:   true,
		StateTopic:    fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber),
		ValueTemplate: "{{ value_json['ignitions'] }}",
		Device:        hspDevice,
	}
	jsonValue, _ = json.Marshal(ignitions)
	mqttClient.Publish(fmt.Sprintf("homeassistant/sensor/hsp_%s_ignitions/config", stove.Meta.SerialNumber), 1, true, jsonValue)
	
	var cleaningIn = HspSensorDiscovery{
		Name:              fmt.Sprintf("HSP %s Cleaning In", stove.Meta.SerialNumber),
		UniqueId:          fmt.Sprintf("hsp-%s-cleaning_in", stove.Meta.SerialNumber),
		ForceUpdate:       true,
		UnitOfMeasurement: "h",
		StateTopic:        fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber),
		ValueTemplate:     "{{ (value_json['cleaning_in']/60) | round(1) }}",
		Device:            hspDevice,
	}
	jsonValue, _ = json.Marshal(cleaningIn)
	mqttClient.Publish(fmt.Sprintf("homeassistant/sensor/hsp_%s_cleaning_in/config", stove.Meta.SerialNumber), 1, true, jsonValue)

	var maintenanceIn = HspSensorDiscovery{
		Name:              fmt.Sprintf("HSP %s Maintenance In", stove.Meta.SerialNumber),
		UniqueId:          fmt.Sprintf("hsp-%s-maintenance_in", stove.Meta.SerialNumber),
		ForceUpdate:       true,
		UnitOfMeasurement: "kg",
		StateTopic:        fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber),
		ValueTemplate:     "{{ value_json['maintenance_in'] }}",
		Device:            hspDevice,
	}
	jsonValue, _ = json.Marshal(maintenanceIn)
	mqttClient.Publish(fmt.Sprintf("homeassistant/sensor/hsp_%s_maintenance_in/config", stove.Meta.SerialNumber), 1, true, jsonValue)

	var operatingHours = HspSensorDiscovery{
		Name:              fmt.Sprintf("HSP %s Operating Hours", stove.Meta.SerialNumber),
		UniqueId:          fmt.Sprintf("hsp-%s-operating_hours", stove.Meta.SerialNumber),
		ForceUpdate:       true,
		UnitOfMeasurement: "h",
		StateTopic:        fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber),
		ValueTemplate:     "{{ value_json['on_time'] }}",
		Device:            hspDevice,
	}
	jsonValue, _ = json.Marshal(operatingHours)
	mqttClient.Publish(fmt.Sprintf("homeassistant/sensor/hsp_%s_operating_hours/config", stove.Meta.SerialNumber), 1, true, jsonValue)

	var err = HspSensorDiscovery{
		Name:          fmt.Sprintf("HSP %s Error", stove.Meta.SerialNumber),
		UniqueId:      fmt.Sprintf("hsp-%s-error", stove.Meta.SerialNumber),
		ForceUpdate:   true,
		StateTopic:    fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber),
		ValueTemplate: "{{ value_json['error'] }}",
		Device:        hspDevice,
	}
	jsonValue, _ = json.Marshal(err)
	mqttClient.Publish(fmt.Sprintf("homeassistant/sensor/hsp_%s_err/config", stove.Meta.SerialNumber), 1, true, jsonValue)

	var consumption = HspSensorDiscovery{
		Name:              fmt.Sprintf("HSP %s Consumption", stove.Meta.SerialNumber),
		UniqueId:          fmt.Sprintf("hsp-%s-consumption", stove.Meta.SerialNumber),
		ForceUpdate:       true,
		UnitOfMeasurement: "kg",
		StateTopic:        fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber),
		ValueTemplate:     "{{ value_json.consumption }}",
		Device:            hspDevice,
	}
	jsonValue, _ = json.Marshal(consumption)
	mqttClient.Publish(fmt.Sprintf("homeassistant/sensor/hsp_%s_consumption/config", stove.Meta.SerialNumber), 1, true, jsonValue)

	var powerSwitch = HspSwitchDiscovery{
		Name:          fmt.Sprintf("HSP %s Power", stove.Meta.SerialNumber),
		UniqueId:      fmt.Sprintf("hsp-%s-prg", stove.Meta.SerialNumber),
		ForceUpdate:   true,
		StateTopic:    fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber),
		ValueTemplate: "{{ value_json['prg'] }}",
		Device:        hspDevice,
		CommandTopic:  fmt.Sprintf("hsp-%s/command/power", stove.Meta.SerialNumber),
		StateOff:      false,
		StateOn:       true,
	}
	jsonValue, _ = json.Marshal(powerSwitch)
	mqttClient.Publish(fmt.Sprintf("homeassistant/switch/hsp_%s_power/config", stove.Meta.SerialNumber), 1, true, jsonValue)

	var weekProgramSwitch = HspSwitchDiscovery{
		Name:          fmt.Sprintf("HSP %s Weekprogram", stove.Meta.SerialNumber),
		UniqueId:      fmt.Sprintf("hsp-%s-wprg", stove.Meta.SerialNumber),
		ForceUpdate:   true,
		StateTopic:    fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber),
		ValueTemplate: "{{ value_json['wprg'] }}",
		Device:        hspDevice,
		CommandTopic:  fmt.Sprintf("hsp-%s/command/weekProgram", stove.Meta.SerialNumber),
		StateOff:      false,
		StateOn:       true,
	}
	jsonValue, _ = json.Marshal(weekProgramSwitch)
	mqttClient.Publish(fmt.Sprintf("homeassistant/switch/hsp_%s_weekprogram/config", stove.Meta.SerialNumber), 1, true, jsonValue)

	var roomModeSwitch = HspSwitchDiscovery{
		Name:          fmt.Sprintf("HSP %s Room Mode", stove.Meta.SerialNumber),
		UniqueId:      fmt.Sprintf("hsp-%s-room_mode", stove.Meta.SerialNumber),
		ForceUpdate:   true,
		StateTopic:    fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber),
		ValueTemplate: "{{ value_json['room_mode'] }}",
		Device:        hspDevice,
		CommandTopic:  fmt.Sprintf("hsp-%s/command/roomMode", stove.Meta.SerialNumber),
		StateOff:      false,
		StateOn:       true,
	}
	jsonValue, _ = json.Marshal(roomModeSwitch)
	mqttClient.Publish(fmt.Sprintf("homeassistant/switch/hsp_%s_room_mode/config", stove.Meta.SerialNumber), 1, true, jsonValue)

	var ecoModeSwitch = HspSwitchDiscovery{
		Name:          fmt.Sprintf("HSP %s Eco Mode", stove.Meta.SerialNumber),
		UniqueId:      fmt.Sprintf("hsp-%s-eco_mode", stove.Meta.SerialNumber),
		ForceUpdate:   true,
		StateTopic:    fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber),
		ValueTemplate: "{{ value_json['eco_mode'] }}",
		Device:        hspDevice,
		CommandTopic:  fmt.Sprintf("hsp-%s/command/ecoMode", stove.Meta.SerialNumber),
		StateOff:      false,
		StateOn:       true,
	}
	jsonValue, _ = json.Marshal(ecoModeSwitch)
	mqttClient.Publish(fmt.Sprintf("homeassistant/switch/hsp_%s_eco_mode/config", stove.Meta.SerialNumber), 1, true, jsonValue)

	var heatCurveInputNumber = HspInputNumberDiscovery{
		Name:          fmt.Sprintf("HSP %s Heating Curve", stove.Meta.SerialNumber),
		UniqueId:      fmt.Sprintf("hsp-%s-heating_curve", stove.Meta.SerialNumber),
		ForceUpdate:   true,
		StateTopic:    fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber),
		ValueTemplate: "{{ value_json['ht_char'] }}",
		Device:        hspDevice,
		CommandTopic:  fmt.Sprintf("hsp-%s/command/heatingCurve", stove.Meta.SerialNumber),
		NumberMin:     1,
		NumberMax:     4,
	}
	jsonValue, _ = json.Marshal(heatCurveInputNumber)
	mqttClient.Publish(fmt.Sprintf("homeassistant/input_number/hsp_%s_heating_curve/config", stove.Meta.SerialNumber), 1, true, jsonValue)

	var hspClimate = HspClimateDiscovery{
		Device:                   hspDevice,
		Name:                     fmt.Sprintf("HSP %s Temperature", stove.Meta.SerialNumber),
		UniqueId:                 fmt.Sprintf("hsp-%s-climate", stove.Meta.SerialNumber),
		ForceUpdate:              true,
		TempCommandTopic:         fmt.Sprintf("hsp-%s/command/target_temperature", stove.Meta.SerialNumber),
		TempStateTopic:           fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber),
		TempStateTemplate:        "{{ value_json['sp_temp'] }}",
		CurrentTempStateTopic:    fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber),
		CurrentTempStateTemplate: "{{ value_json['is_temp'] }}",
		MinTemp:                  "10",
		MaxTemp:                  "30",
		TempStep:                 "1",
		ModeStateTopic:           fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber),
		ModeStateTemplate:        "{{ 'heat' }}",
		Modes:                    []string{"heat"},
	}
	jsonValue, _ = json.Marshal(hspClimate)
	mqttClient.Publish(fmt.Sprintf("homeassistant/climate/hsp_%s/config", stove.Meta.SerialNumber), 1, true, jsonValue)

	var hspClearErrorBtn = HspButtonDiscovery{
		Device:        hspDevice,
		Name:          "Clear Error",
		UniqueId:      fmt.Sprintf("hsp-%s-clear-err-btn", stove.Meta.SerialNumber),
		CommandTopic:  fmt.Sprintf("hsp-%s/command/clean_error", stove.Meta.SerialNumber),
		ValueTemplate: "{{ 'non' }}",
		ForceUpdate:   false,
	}
	jsonValue, _ = json.Marshal(hspClearErrorBtn)
	mqttClient.Publish(fmt.Sprintf("homeassistant/button/hsp_%s_clear_error/config", stove.Meta.SerialNumber), 1, true, jsonValue)
}
func executeJob() {
	interval, _ := strconv.ParseInt(os.Getenv("HSP_POLL_INTERVAL"), 10, 64)
	gocron.Every(uint64(interval)).Second().Do(pollValue)
	<-gocron.Start()
}
func pollValue() {
	stove := callStove()
	log.Println("Publish Stove Values")
	json, _ := json.Marshal(stove)
	mqttClient.Publish(fmt.Sprintf("hsp-%s/result", stove.Meta.SerialNumber), 1, false, json)
}
func subscribeMqtt() {
	stove := callStove()

	topics := make(map[string]byte)
	topics[fmt.Sprintf("hsp-%s/command/power", stove.Meta.SerialNumber)] = 1
	topics[fmt.Sprintf("hsp-%s/command/weekProgram", stove.Meta.SerialNumber)] = 1
	topics[fmt.Sprintf("hsp-%s/command/roomMode", stove.Meta.SerialNumber)] = 1
	topics[fmt.Sprintf("hsp-%s/command/ecoMode", stove.Meta.SerialNumber)] = 1
	topics[fmt.Sprintf("hsp-%s/command/target_temperature", stove.Meta.SerialNumber)] = 1
	topics[fmt.Sprintf("hsp-%s/command/forward_flow_temperature", stove.Meta.SerialNumber)] = 1
	topics[fmt.Sprintf("hsp-%s/command/heating_curve", stove.Meta.SerialNumber)] = 1
	topics[fmt.Sprintf("hsp-%s/command/clean_error", stove.Meta.SerialNumber)] = 1

	token := mqttClient.SubscribeMultiple(topics, func(client mqtt.Client, message mqtt.Message) {
		log.Printf("Topic %s was published, Value: %s \r\n", message.Topic(), string(message.Payload()))
		if message.Topic() == fmt.Sprintf("hsp-%s/command/power", stove.Meta.SerialNumber) {
			payload, _ := strconv.ParseBool(string(message.Payload()))
			command(nil, nil, nil, BoolPointer(payload), nil, nil, nil)
		}
		if message.Topic() == fmt.Sprintf("hsp-%s/command/weekProgram", stove.Meta.SerialNumber) {
			payload, _ := strconv.ParseBool(string(message.Payload()))
			command(nil, nil, nil, nil, BoolPointer(payload), nil, nil)
		}
		if message.Topic() == fmt.Sprintf("hsp-%s/command/roomMode", stove.Meta.SerialNumber) {
			payload, _ := strconv.ParseBool(string(message.Payload()))
			command(nil, nil, nil, nil, nil, BoolPointer(payload), nil)
		}
		if message.Topic() == fmt.Sprintf("hsp-%s/command/ecoMode", stove.Meta.SerialNumber) {
			payload, _ := strconv.ParseBool(string(message.Payload()))
			command(nil, nil, nil, nil, nil, nil, BoolPointer(payload))
		}
		if message.Topic() == fmt.Sprintf("hsp-%s/command/target_temperature", stove.Meta.SerialNumber) {
			payload, _ := strconv.ParseFloat(string(message.Payload()), 0)
			var p int = int(payload)
			command(IntPointer(p), nil, nil, nil, nil, nil, nil)
		}
		if message.Topic() == fmt.Sprintf("hsp-%s/command/forward_flow_temperature", stove.Meta.SerialNumber) {
			payload, _ := strconv.ParseFloat(string(message.Payload()), 0)
			var p int = int(payload)
			command(nil, IntPointer(p), nil, nil, nil, nil, nil)
		}
		if message.Topic() == fmt.Sprintf("hsp-%s/command/heating_curve", stove.Meta.SerialNumber) {
			payload, _ := strconv.ParseFloat(string(message.Payload()), 0)
			var p float = float(payload)
			command(nil, nil, FloatPointer(p), nil, nil, nil, nil)
		}
		if message.Topic() == fmt.Sprintf("hsp-%s/command/clean_error", stove.Meta.SerialNumber) {
			currentError := callStove()
			if len(currentError.Error) > 0 {
				errCode := currentError.Error[0].ErrorCode
				clearStoveError(errCode)
			}

		}
	})

	token.Wait()
}
func checkForEnv() {
	envVars := [...]string{"HSP_STOVE_IP", "HSP_STOVE_PIN", "MQTT_PORT", "MQTT_IP", "HSP_POLL_INTERVAL"}
	for _, element := range envVars {
		env := os.Getenv(element)
		if len(env) == 0 {
			log.Printf("Environment variable %s must be set!\r\n", element)
			os.Exit(0)
		}
	}
}
