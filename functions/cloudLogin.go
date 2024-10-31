package functions

import (
	"bytes"
	"encoding/json"
	"goback-client/data"
	"log"
	"net/http"
)

var URL = "http://" + data.Config.Cloud + ":" + data.Config.Port

func Login() bool {
	var loginURL = URL + "/api/v1/user/login"
	jsonData := []byte(`{"username":"` + data.Config.Username + `","password":"` + data.Config.Password + `"}`)
	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var loginResp LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		log.Println("Error decoding login response:", err)
		return false
	}

	if loginResp.Status == 200 {
		data.Config.JWTTOKEN = loginResp.Data.AccessToken
		if !AuthedPing() {
			log.Println("AuthedPing failed")
			return false
		}
		return true
	}
	return false
}

func AuthedPing() bool {
	var pingURL = URL + "/api/v1/authed/ping"
	req, _ := http.NewRequest("GET", pingURL, nil)
	req.Header.Set("Authorization", data.Config.JWTTOKEN)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending ping request:", err)
		return false
	}
	defer resp.Body.Close()

	var pingResp PingResponse
	err = json.NewDecoder(resp.Body).Decode(&pingResp)
	if err != nil {
		return false
	}

	return pingResp.Msg == "authed pong!"
}
