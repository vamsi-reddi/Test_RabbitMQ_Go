package config

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing_RabbitMQ/types"

	"github.com/go-playground/validator"
)

var CfgObj *types.Configurations

func LoadConfig(secretData string) bool {
	jsonParser := json.NewDecoder(strings.NewReader(secretData))

	err := jsonParser.Decode(&CfgObj)

	if err != nil {
		fmt.Println("error parsing config: ", err)
		return false
	}

	validate := validator.New()

	err = validate.Struct(CfgObj)

	if err != nil {
		fmt.Println("error validating config: ", err)
		return false
	}

	return true
}