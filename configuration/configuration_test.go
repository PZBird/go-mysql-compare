package configuration

import (
	"reflect"
	"testing"
)

func TestLoadConfiguration(t *testing.T) {
	Db1 := ConfigurationElement{
		Hostname:         "localhost",
		Port:             "3306",
		Username:         "user",
		Password:         "password",
		ExcludedEntities: make([]string, 0),
		DatabasesSuffix:  "_stage",
	}
	Db2 := ConfigurationElement{
		Hostname:         "localhost",
		Port:             "3307",
		Username:         "user",
		Password:         "password",
		ExcludedEntities: make([]string, 0),
		DatabasesSuffix:  "",
	}
	mock := Configuration{
		Db1,
		Db2,
	}

	config := LoadConfiguration("../config.example.json")

	equal := reflect.DeepEqual(mock, config)

	if equal == false {
		t.Fatalf("Configuration load failed. Example didn't match for a mock.")
	}
}
