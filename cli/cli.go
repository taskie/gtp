package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/taskie/jc"

	"github.com/jessevdk/go-flags"
	"github.com/taskie/gtp"
	"github.com/taskie/osplus"
)

type Options struct {
	TemplateType string `short:"T" long:"templateType" default:"text" description:"template type [text|html]"`
	DataFilePath string `short:"d" long:"data" description:"data file path"`
	DataType     string `short:"D" long:"dataType" description:"data type [json|toml|msgpack]"`
	NoColor      bool   `long:"noColor" env:"NO_COLOR" description:"NOT colorize output"`
	Verbose      bool   `short:"v" long:"verbose" description:"show verbose output"`
	Version      bool   `short:"V" long:"version" description:"show version"`
}

func Main() {
	err := mainImpl()
	if err != nil {
		log.Fatal(err)
	}
}

func mainImpl() error {
	var opts Options
	args, err := flags.ParseArgs(&opts, os.Args)
	if opts.Version {
		if opts.Verbose {
			fmt.Println("Version: ", gtp.Version)
		} else {
			fmt.Println(gtp.Version)
		}
		return nil
	}
	if len(args) <= 1 {
		return fmt.Errorf("you must specify template file paths")
	}

	opener := osplus.NewOpener()
	rc, err := opener.Open(opts.DataFilePath)
	defer rc.Close()

	dataType := opts.DataType
	if dataType == "" {
		dataType = jc.ExtToType(filepath.Ext(opts.DataFilePath))
	}
	if dataType == "" {
		dataType = "json"
	}

	gtp := gtp.Gtp{
		TemplateFilePaths: args[1:],
		TemplateType:      opts.TemplateType,
		DataType:          dataType,
	}

	err = gtp.Run(os.Stdout, rc)
	if err != nil {
		return err
	}
	return nil
}
