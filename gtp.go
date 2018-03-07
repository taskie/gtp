package gtp

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/taskie/jc"
	htmltemplate "html/template"
	"io"
	"os"
	"text/template"
)

var (
	version  string
	revision string
)

type Gtp struct {
	TemplateFilePaths []string
	TemplateType      string
	DataType          string
}

func (gtp *Gtp) Execute(w io.Writer, data interface{}) error {
	switch gtp.TemplateType {
	case "text":
		tpl, err := template.ParseFiles(gtp.TemplateFilePaths...)
		if err != nil {
			return err
		}
		err = tpl.Execute(w, data)
		return err
	case "html":
		tpl, err := htmltemplate.ParseFiles(gtp.TemplateFilePaths...)
		if err != nil {
			return err
		}
		err = tpl.Execute(w, data)
		return err
	default:
		return fmt.Errorf("Invalid template type: %s", gtp.TemplateType)
	}
}

func (gtp *Gtp) Run(r io.Reader, w io.Writer) error {
	if len(gtp.TemplateFilePaths) == 0 {
		return fmt.Errorf("Please specify template file path")
	}
	jco := jc.Jc{
		FromType: gtp.DataType,
	}
	var data interface{}
	err := jco.Decode(r, &data)
	if err != nil {
		return err
	}
	err = gtp.Execute(w, data)
	return err
}

type Options struct {
	TemplateFilePaths []string `short:"t" long:"template" description:"template file path"`
	TemplateType      string   `short:"T" long:"template-type" default:"text" description:"template type [text|html]"`
	DataType          string   `short:"D" long:"data-type" default:"json" description:"data type [json|toml|msgpack]"`
	NoColor           bool     `long:"no-color" env:"NO_COLOR" description:"NOT colorize output"`
	Verbose           bool     `short:"v" long:"verbose" description:"show verbose output"`
	Version           bool     `short:"V" long:"version" description:"show version"`
}

func Main(args []string) {
	var opts Options
	args, err := flags.ParseArgs(&opts, args)
	if opts.Version {
		if opts.Verbose {
			fmt.Println("Version: ", version)
			fmt.Println("Revision: ", revision)
		} else {
			fmt.Println(version)
		}
		os.Exit(0)
	}

	gtp := Gtp{
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
