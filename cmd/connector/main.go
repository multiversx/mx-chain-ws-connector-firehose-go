package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/closing"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-logger-go/file"
	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	"github.com/multiversx/mx-chain-ws-connector-template-go/factory"
	"github.com/urfave/cli"
)

var log = logger.GetOrCreate("main")

const (
	configPath = "config/config.toml"

	logsPath       = "logs"
	logFilePrefix  = "ws-connector-firehose"
	logLifeSpanSec = 432000 // 5 days
	logLifeSpanMb  = 1024   // 1 GB
)

func main() {
	app := cli.NewApp()
	app.Name = "MultiversX ws connector template"
	app.Usage = "This tool will communicate with an observer/light client connected to mx-chain via " +
		"websocket outport driver and listen to incoming exported data."
	app.Flags = []cli.Flag{
		configFile,
		logLevel,
		logSaveFile,
		disableAnsiColor,
		dbMode,
	}
	app.Authors = []cli.Author{
		{
			Name:  "The MultiversX Team",
			Email: "contact@multiversx.com",
		},
	}

	app.Action = startConnector

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func startConnector(ctx *cli.Context) error {
	configFilePath := ctx.GlobalString(configFile.Name)

	absoluteConfigPath, err := filepath.Abs(configFilePath)
	if err != nil {
		return err
	}

	cfg, err := loadConfig(absoluteConfigPath)
	if err != nil {
		return err
	}

	err = initializeLogger(ctx)
	if err != nil {
		return err
	}

	var logFile closing.Closer
	withLogFile := ctx.GlobalBool(logSaveFile.Name)
	if withLogFile {
		logFile, err = createLogger()
		if err != nil {
			return err
		}
	}

	dbMode := ctx.GlobalString(dbMode.Name)
	log.Info("storer sync mode", "dbMode", dbMode)

	connectorRunner, err := connector.NewConnectorRunner(cfg, dbMode)
	if err != nil {
		return fmt.Errorf("cannot create connector runner, error: %w", err)
	}

	err = connectorRunner.Run()
	if err != nil {
		return fmt.Errorf("connector runner, error: %w", err)
	}

	if withLogFile {
		err = logFile.Close()
		log.LogIfError(err)
	}

	return nil
}

func loadConfig(filepath string) (*config.Config, error) {
	cfg := &config.Config{}
	err := core.LoadTomlFile(cfg, filepath)

	log.Info("loaded config", "path", filepath)

	return cfg, err
}

func initializeLogger(ctx *cli.Context) error {
	logLevelFlagValue := ctx.GlobalString(logLevel.Name)
	err := logger.SetLogLevel(logLevelFlagValue)
	if err != nil {
		return err
	}

	disableAnsi := ctx.GlobalBool(disableAnsiColor.Name)
	return removeANSIColorsForLoggerIfNeeded(disableAnsi)
}

func removeANSIColorsForLoggerIfNeeded(disableAnsi bool) error {
	if !disableAnsi {
		return nil
	}

	err := logger.RemoveLogObserver(os.Stdout)
	if err != nil {
		return err
	}

	return logger.AddLogObserver(os.Stdout, &logger.PlainFormatter{})
}

func createLogger() (closing.Closer, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Error("error getting working directory when trying to create logger file", "error", err)
		workingDir = ""
	}

	argsLogger := file.ArgsFileLogging{
		WorkingDir:      workingDir,
		DefaultLogsPath: logsPath,
		LogFilePrefix:   logFilePrefix,
	}
	fileLogging, err := file.NewFileLogging(argsLogger)
	if err != nil {
		return nil, fmt.Errorf("%w creating log file", err)
	}

	err = fileLogging.ChangeFileLifeSpan(time.Second*time.Duration(logLifeSpanSec), logLifeSpanMb)
	if err != nil {
		return nil, err
	}

	return fileLogging, nil
}
