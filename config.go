package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

func readConfig(file string) map[string]interface{} {
	if file == "" {
		return readConfigPipe()
	} else {
		return readConfigFile(file)
	}
}

func mergeWithEnv(config map[string]interface{}) map[string]interface{} {
	env := make(map[string]interface{})
	for _, envStr := range os.Environ() {
		ret := strings.SplitN(envStr, "=", 2)
		if len(ret) == 2 {
			env[ret[0]] = ret[1]
		}
	}
	config["env"] = env
	return config
}

func readConfigPipe() map[string]interface{} {
	config := make(map[string]interface{})
	defer os.Stdin.Close()
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		err := json.NewDecoder(os.Stdin).Decode(&config)
		if err != nil {
			log.Fatal(err)
		}
	}
	return config
}

func readConfigFile(file string) map[string]interface{} {
	configFile, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	// Parse json
	return readConfigBuffer(configFile)
}

func readConfigBuffer(buffer []byte) map[string]interface{} {
	var config map[string]interface{}

	// Parse json
	err := json.Unmarshal(buffer, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func validateJSONVaule(schema map[string]any, json map[string]any) error {
	schem, err := gojsonschema.NewSchema(gojsonschema.NewRawLoader(schema))
	if err != nil {
		return err
	}
	value := gojsonschema.NewRawLoader(json)

	ret, err := schem.Validate(value)
	if err != nil {
		return err
	}

	if !ret.Valid() {
		return fmt.Errorf("%+v", ret.Errors())
	}
	return nil
}
