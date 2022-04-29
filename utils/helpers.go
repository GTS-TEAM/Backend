package utils

import (
	"encoding/json"
	"log"
)

func LogError(in string, err error) {
	log.Println("\n[ERROR] " + "[" + in + "] " + err.Error())
}

func LogInfo(in string, message string) {
	log.Println("\n[INFO] " + "[" + in + "] " + message)
}

func BindStruct(obj interface{}, to interface{}) error {
	marshal, err := json.Marshal(obj)

	if err != nil {
		return err
	}
	return json.Unmarshal(marshal, to)
}
