package nxrm

import "fmt"

var config = map[string]interface{}{
	"username":  "testymctestface",
	"password":  "mybigsecret",
	"url":       "http://localhost:8081/repository/tf-backend",
	"subpath":   "this/here",
	"stateName": "demo.tfstate",
	"timeout":   30,
}

func mismatchError(configName string, got interface{}) string {
	return fmt.Sprintf("%s: expected: %s got: %s", configName, config[configName], got.(string))
}
