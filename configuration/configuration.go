package configuration

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type ConfigurationElement struct {
	Hostname         string
	Port             string
	Username         string
	Password         string
	ExcludedEntities []string
	DatabasesSuffix  string
}

type Configuration struct {
	Db1 ConfigurationElement
	Db2 ConfigurationElement
}

func LoadConfiguration(path string) Configuration {
	file, e := ioutil.ReadFile(path)

	if e != nil {
		log.Fatalf("Configuration file not found at path %s", path)
	}

	var jsontype Configuration

	if e = json.Unmarshal(file, &jsontype); e != nil {
		log.Fatalf("Invalid configuration file due to %s", e.Error())
	}

	return jsontype
}
