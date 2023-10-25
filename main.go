package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var TempioVersion string

func main() {
	var config map[string]interface{}
	configFile := flag.String("conf", "", "Config json file, can be omitted if used in a pipe")
	templateFile := flag.String("template", "", "Template file or string")
	outFile := flag.String("out", "", "Output file, if not defined output will be to console")
	templateSchemaFile := flag.String("schema", "", "Json schema file can use validate template vaule")

	flag.Usage = func() {
		writer := flag.CommandLine.Output()

		fmt.Fprintf(writer, "Version: %s\n", TempioVersion)
		fmt.Fprintf(writer, "Documentation: https://github.com/home-assistant/tempio\n\n")

		flag.PrintDefaults()
	}
	// parse command lines
	flag.Parse()
	if *templateFile == "" {
		log.Fatal("Missing template argument")
	}

	_, err := os.Stat(*templateFile)

	// Get config
	config = mergeWithEnv(readConfig(*configFile))

	if templateSchemaFile != nil {
		if err = validateJSONVaule(readConfigFile(*templateSchemaFile), config); err != nil {
			log.Fatal(err)
		}
	}
	// create & write corefile
	data := renderTemplateFile(config, *templateFile, err != nil)
	if *outFile == "" {
		fmt.Println(string(data))
	} else {
		err := ioutil.WriteFile(*outFile, data, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
