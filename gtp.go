package main

import (
	"flag"
	"fmt"
	"github.com/ugorji/go/codec"
	"os"
	"text/template"
)

func main() {
	confJSON := flag.String("j", "", "specify JSON path")
	confMP := flag.String("m", "", "specify MessagePack path")
	flag.Parse()

	var handle codec.Handle
	var conf string
	var jh codec.JsonHandle
	var mh codec.MsgpackHandle

	if *confJSON == "" && *confMP == "" {
		fmt.Fprintln(os.Stderr, "no config file specified.\nusage: gtp [OPTION] [FILE]...")
		flag.PrintDefaults()
		os.Exit(1)
	} else if *confJSON != "" && *confMP != "" {
		fmt.Fprintln(os.Stderr, "cannot specify both of JSON and MessagePack")
		os.Exit(1)
	} else if *confJSON != "" {
		conf = *confJSON
		handle = &jh
	} else {
		conf = *confMP
		mh.RawToString = true
		handle = &mh
	}

	confFile, err := os.Open(conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var confData interface{}
	confDecoder := codec.NewDecoder(confFile, handle)
	if err := confDecoder.Decode(&confData); err != nil {
		fmt.Fprintln(os.Stderr, "parse error: ", err)
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
