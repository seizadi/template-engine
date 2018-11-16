package main

import (
	"github.com/spf13/pflag"
)

const (
	// configuration defaults support local development (i.e. "go run ...")
	// Server
	defaultApiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2NvdW50SUQiOjF9.GsXyFDDARjXe1t9DPo2LIBKHEal3O7t3vLI3edA7dGU"
	defaultCmdbHost = "localhost:9090"

	// Config
	defaultConfigDirectory = "config/"
	defaultConfigFile      = ""
	
	// Logging
	defaultLoggingLevel = "debug"
)

var (
	// define flag overrides
	flagServerApiKey = pflag.String("server.apiKey", defaultApiKey, "API Key to Access e.g. CMDB")
	flagServerCmdbHost = pflag.String("server.cmdbHost", defaultCmdbHost, "CMDB Host e.g. localhost:9090")
	flagLoggingLevel = pflag.String("logging.level", defaultLoggingLevel, "log level of application")
	flagConfigDirectory = pflag.String("config.source", defaultConfigDirectory, "directory of the configuration file")
	flagConfigFile      = pflag.String("config.file", defaultConfigFile, "directory of the configuration file")
)
