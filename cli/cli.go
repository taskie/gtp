package cli

import (
	"github.com/taskie/gtp"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

type Options struct {
	TemplateFilePaths []string `short:"t" long:"template" description:"template file path"`
	TemplateType      string   `short:"T" long:"templateType" default:"text" description:"template type [text|html]"`
	DataType          string   `short:"D" long:"dataType" default:"json" description:"data type [json|toml|msgpack]"`
	NoColor           bool     `long:"noColor" env:"NO_COLOR" description:"NOT colorize output"`
	Verbose           bool     `short:"v" long:"verbose" description:"show verbose output"`
	Version           bool     `short:"V" long:"version" description:"show version"`
}

func Main() {
	var opts Options
	_, err := flags.ParseArgs(&opts, os.Args)
	if opts.Version {
		if opts.Verbose {
			fmt.Println("Version: ", gtp.Version)
		} else {
			fmt.Println(gtp.Version)
		}
		os.Exit(0)
	}

	gtp := gtp.Gtp{
		TemplateFilePaths: opts.TemplateFilePaths,
		TemplateType:      opts.TemplateType,
		DataType:          opts.DataType,
	}
	err = gtp.Run(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
