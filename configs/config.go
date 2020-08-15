package configs

import (
	"io/ioutil"
	"path/filepath"

	"github.com/mrzack99s/mrzack-dns-server/structs"

	"gopkg.in/yaml.v3"
)

var SystemConfig structs.SystemConfig

func ParseSystemConfig() {
	filename, _ := filepath.Abs("./config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)

	err = yaml.Unmarshal(yamlFile, &SystemConfig)
	if err != nil {
		panic(err)
	}

}
