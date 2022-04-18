package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func callStove() HspStove {
	response, err := http.Get(fmt.Sprintf("http://%s/status.cgi", os.Getenv("HSP_STOVE_IP")))
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var stoveResponse HspStove
	json.Unmarshal(body, &stoveResponse)

	return stoveResponse
}

func calculatePin(nonce string, pin string) [16]byte {
	bPin := md5.Sum([]byte(pin))
	return md5.Sum([]byte(nonce + hex.EncodeToString(bPin[:])))
}

func command(targetTemp *int, start *bool, weekProgramStart *bool) {
	command := HspCommand{targetTemp, start, weekProgramStart}
	commandJson, _ := json.Marshal(command)
	stove := callStove()
	calculatedPin := calculatePin(stove.Meta.Nonce, os.Getenv("HSP_STOVE_PIN"))
	client := &http.Client{}
	commandBuffer := bytes.NewBuffer(commandJson)
	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("http://%s/status.cgi", os.Getenv("HSP_STOVE_IP")),
		commandBuffer)
	headers := http.Header{
		"X-HS-PIN":     []string{hex.EncodeToString(calculatedPin[:])},
		"X-BACKEND-IP": []string{"https://app.hsp.com"},
	}

	req.Header = headers
	response, _ := client.Do(req)
	defer response.Body.Close()
	_, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	pollValue()
}

func clearStoveError(errCode int) {
	var err = []HspError{
		{ErrorCode: errCode},
	}
	command := HspCleanError{SeenError: err}
	commandJson, _ := json.Marshal(command)
	stove := callStove()
	calculatedPin := calculatePin(stove.Meta.Nonce, os.Getenv("HSP_STOVE_PIN"))
	client := &http.Client{}
	commandBuffer := bytes.NewBuffer(commandJson)
	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("http://%s/status.cgi", os.Getenv("HSP_STOVE_IP")),
		commandBuffer)
	headers := http.Header{
		"X-HS-PIN":     []string{hex.EncodeToString(calculatedPin[:])},
		"X-BACKEND-IP": []string{"https://app.hsp.com"},
	}

	req.Header = headers
	response, _ := client.Do(req)
	defer response.Body.Close()

	pollValue()
}
