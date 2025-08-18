package main

import (
	"context"
	"os"

	"otusgruz/cmd"
	"otusgruz/config"

	"github.com/rs/zerolog"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

	exitCode := 0

	logger.Info().Msg("application is launching")

	err = cmd.Run(ctx, conf)
	if err != nil {
		exitCode = 1
	}

	os.Exit(exitCode)
}
