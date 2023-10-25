package main

import (
	"bytes"
	"log"
	"os"
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
)

func renderTemplateFile(config map[string]interface{}, file string, isStr bool) []byte {
	var (
		templateFile []byte
		err          error
	)
	if isStr {
		templateFile = []byte(file)
	} else {
		// read Template
		templateFile, err = os.ReadFile(file)
		if err != nil {
			log.Fatalf("Cant read template file %s - %s", file, err)
		}
	}

	return renderTemplateBuffer(config, templateFile)
}

func renderTemplateBuffer(config map[string]interface{}, templateData []byte) []byte {
	buf := &bytes.Buffer{}

	// generate template
	coreTemplate := template.New("tempIO").Funcs(sprig.TxtFuncMap())
	template.Must(coreTemplate.Parse(string(templateData)))

	// render
	err := coreTemplate.Execute(buf, config)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}
