package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"log"
	"os"
	"strings"
	
	"github.com/codegangsta/cli"
	"github.com/seizadi/template-engine/cmd/template-engine/commands"
	"github.com/spf13/viper"
)

const version = "0.1"

func init() {
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath(viper.GetString("config.source"))
	if viper.GetString("config.file") != "" {
		viper.SetConfigName(viper.GetString("config.file"))
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("cannot load configuration: %v", err)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "Template Engine CLI"
	app.Usage = ""
	app.Author = "Soheil Eizadi"
	app.Email = "seizadi@gmail.com"
	app.Version = version
	app.Flags = []cli.Flag{}
	
	fmt.Printf ("Logging Level %s\n",viper.GetString("logging.level"))
	
	app.Before = func(c *cli.Context) error {
		return nil
	}
	app.Commands = commands.Commands
	app.CommandNotFound = cmdNotFound
	
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("%v", err)
	}
}


func cmdNotFound(c *cli.Context, command string) {
	fmt.Printf(
		"%s: '%s' is not a %s command. See '%s --help'.\n",
		c.App.Name,
		command,
		c.App.Name,
		os.Args[0],
	)
	os.Exit(1)
}