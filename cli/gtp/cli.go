package gtp

import (
	"fmt"
	"path/filepath"

	"github.com/iancoleman/strcase"

	"github.com/k0kubun/pp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/taskie/gtp"
	"github.com/taskie/jc"
	"github.com/taskie/osplus"
)

type Config struct {
	Output, TemplateType, Data, DataType, LogLevel string
}

var configFile string
var config Config
var (
	verbose, debug, version bool
)

const CommandName = "gtp"

func init() {
	Command.PersistentFlags().StringVarP(&configFile, "config", "c", "", `config file (default "`+CommandName+`.yml")`)
	Command.Flags().StringP("output", "o", "", "output file")
	Command.Flags().StringP("template-type", "T", "", "template type [text|html]")
	Command.Flags().StringP("data", "d", "", "data file")
	Command.Flags().StringP("data-type", "D", "", "data type")
	Command.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	Command.Flags().BoolVar(&debug, "debug", false, "debug output")
	Command.Flags().BoolVarP(&version, "version", "V", false, "show Version")

	for _, s := range []string{"output", "template-type", "data", "data-type"} {
		envKey := strcase.ToSnake(s)
		structKey := strcase.ToCamel(s)
		viper.BindPFlag(envKey, Command.Flags().Lookup(s))
		viper.RegisterAlias(structKey, envKey)
	}

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if debug {
		log.SetLevel(log.DebugLevel)
	} else if verbose {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName(CommandName)
		conf, err := osplus.GetXdgConfigHome()
		if err != nil {
			log.Info(err)
		} else {
			viper.AddConfigPath(filepath.Join(conf, CommandName))
		}
		viper.AddConfigPath(".")
	}
	viper.SetEnvPrefix(CommandName)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Debug(err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Warn(err)
	}
}

func Main() {
	Command.Execute()
}

var Command = &cobra.Command{
	Use:  CommandName + ` [INPUT...]`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := run(cmd, args)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func run(cmd *cobra.Command, args []string) error {
	if version {
		fmt.Println(gtp.Version)
		return nil
	}
	if config.LogLevel != "" {
		lv, err := log.ParseLevel(config.LogLevel)
		if err != nil {
			log.Warn(err)
		} else {
			log.SetLevel(lv)
		}
	}
	if debug {
		if viper.ConfigFileUsed() != "" {
			log.Debugf("Using config file: %s", viper.ConfigFileUsed())
		}
		log.Debug(pp.Sprint(config))
	}

	opener := osplus.NewOpener()
	r, err := opener.Open(config.Data)
	if err != nil {
		return err
	}
	defer r.Close()
	w, commit, err := opener.CreateTempFileWithDestination(config.Output, "", CommandName+"-")
	if err != nil {
		return err
	}
	defer w.Close()

	dataType := config.DataType
	if dataType == "" {
		dataType = jc.ExtToType(filepath.Ext(config.Data))
		if dataType == "" {
			dataType = "json"
		}
	}
	templateType := config.TemplateType
	if templateType == "" {
		templateType = "text"
	}

	gtp := gtp.Gtp{
		TemplateFilePaths: args,
		TemplateType:      templateType,
		DataType:          dataType,
	}

	err = gtp.Run(w, r)
	if err != nil {
		return err
	}

	commit(true)
	return nil
}
