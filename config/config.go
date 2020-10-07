package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	DbConnect      string `json:"db_connection"`
	ListenAndServe string `json:"listen_and_serve"`
	SecretOfJwt    string `json:"secret_of_jwt"`
	Validation     int64  `json:"valid"`
}

var Conf Config

func ReadJsonFile() {
	var plan []byte
	plan, err1 := ioutil.ReadFile("configuration.json")
	if err1 != nil {
		panic(err1)
	}
	err := json.Unmarshal(plan, &Conf)
	if err != nil {
		fmt.Println(err)
	}
}
