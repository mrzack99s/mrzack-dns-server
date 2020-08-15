package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func ParseRecords() map[string]string {

	jsonFile, err := os.Open("./dnsRecords.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]string
	json.Unmarshal([]byte(byteValue), &result)
	return result
}
