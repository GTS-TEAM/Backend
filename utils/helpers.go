package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

func LogError(err error) {
	log.Println("[ERROR] " + err.Error())
}

func LogInfo(message string) {
	log.Println("[INFO] " + message)
}

func BindStruct(obj interface{}, to interface{}) error {
	marshal, err := json.Marshal(obj)

	fmt.Printf("Marshal: %s\n", string(marshal))

	if err != nil {
		return err
	}
	return json.Unmarshal(marshal, to)
}
