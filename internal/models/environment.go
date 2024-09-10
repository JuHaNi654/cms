package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var enviromentVariables = map[string]string{
	"Env":          "ENVIRONMENT",
	"Port":         "PORT",
	"DatabaseAddr": "DATABASE_ADDR",
}

var Environment = &environment{}

type environment struct {
	Root         string
	Env          string
	Port         string
	DatabaseAddr string
}

func (e environment) IsProduction() bool {
	return e.Env == "production"
}

func (e environment) WithRoot(target string) string {
  return fmt.Sprintf("%s%s", e.Root, target)
}

func LoadEnvironment() error {
	tmp := make(map[string]interface{})
	for k, v := range enviromentVariables {
		value, isSet := os.LookupEnv(v)
		if !isSet {
			return fmt.Errorf("Environment (%s) is required", v)
		}

		tmp[k] = value
	}

	body, err := json.Marshal(tmp)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, Environment); err != nil {
		return err
	}

  if err := getApplicationRootPath(Environment); err != nil {
    return err
  }

	return nil
}

func getApplicationRootPath(e *environment) (error) {
  var err error 
  if e.IsProduction() {
    ex, err := os.Executable()
    if err != nil {
      return err
    }

    e.Root = filepath.Dir(ex)
  }  

  e.Root, err = os.Getwd()
  if err != nil {
    return err
  }

  return nil
}


