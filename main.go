package main

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	_ "github.com/StephanHCB/go-autumn-logging-log"
	"github.com/StephanHCB/go-generator-cli/internal"
	"os"
)

var version = ""

func main() {
	ctx := context.Background()

	if version != "" {
		aulogging.Logger.Ctx(ctx).Info().Print("welcome to go-generator-cli version: " + version)
	} else {
		aulogging.Logger.Ctx(ctx).Info().Print("welcome to go-generator-cli")
	}

	if internal.Perform(ctx) {
		aulogging.Logger.Ctx(ctx).Info().Print("success") 
		os.Exit(0)
	} else {
		aulogging.Logger.Ctx(ctx).Error().Print("there were errors, return code is 1")
		os.Exit(1)
	}
}
