package main

import (
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/urfave/cli"
)

const (
	defaultConfigPath = "./config/config.toml"
)

var (
	configFile = cli.StringFlag{
		Name:  "config",
		Usage: "This flag specifies the config.toml to be used for connecting to the node.",
		Value: defaultConfigPath,
	}

	logLevel = cli.StringFlag{
		Name: "log-level",
		Usage: "This flag specifies the logger `level(s)`. It can contain multiple comma-separated value. For example" +
			", if set to *:INFO the logs for all packages will have the INFO level. However, if set to *:INFO,api:DEBUG" +
			" the logs for all packages will have the INFO level, excepting the api package which will receive a DEBUG" +
			" log level.",
		Value: "*:" + logger.LogDebug.String(),
	}

	logSaveFile = cli.BoolFlag{
		Name:  "log-save",
		Usage: "Boolean option for enabling log saving. If set, it will automatically save all the logs into a file.",
	}

	disableAnsiColor = cli.BoolFlag{
		Name:  "disable-ansi-color",
		Usage: "Boolean option for disabling ANSI colors in the logging system.",
	}

	dbMode = cli.StringFlag{
		Name:  "db-mode",
		Usage: "Option for specifying db mode. Available options: `full-persister`, `import-db`, `optimized-persister`",
		Value: "full-persister",
	}

	enableGrpcServer = cli.BoolFlag{
		Name:  "enable-grpc-server",
		Usage: "Option for enabling grpc server",
	}
)
