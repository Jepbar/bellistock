package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	DbConnect      string `json:"db_connection"`
	ListenAndServe string `json:"listen_and_serve"`
	SecretOfJwt    string `json:"secret_of_jwt"`
	Validation     int64  `json:"valid"`
}

func ReadJsonFile() Config {
	plan, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	var usr Config
	err = json.Unmarshal(plan, &usr)
	return usr
}
