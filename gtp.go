package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/template"
)

func main() {
	conf := flag.String("j", "", "JSON file path")
	flag.Parse()

	if *conf == "" {
		fmt.Fprintln(os.Stderr, "no JSON file specified.\nusage: ")
		flag.PrintDefaults()
		os.Exit(1)
	}

	confFile, err := os.Open(*conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	confDecoder := json.NewDecoder(confFile)
	var confData interface{}
	if err := confDecoder.Decode(&confData); err != nil {
		fmt.Fprintln(os.Stderr, "JSON parse error:", err)
		os.Exit(1)
	}

	errNo := 0
	for _, arg := range flag.Args() {
		tpl := template.Must(template.ParseFiles(arg))
		if err := tpl.Execute(os.Stdout, confData); err != nil {
			fmt.Fprintln(os.Stderr, err)
			errNo = 1
		}
	}
	os.Exit(errNo)
}
