package main

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	_ "github.com/StephanHCB/go-autumn-logging-log"
	generatorlib "github.com/StephanHCB/go-generator-lib"
	generatorapi "github.com/StephanHCB/go-generator-lib/api"
	"os"
)

func main() {
	ctx := context.TODO()
	request := &generatorapi.Request{
		SourceBaseDir: "/path/to/generator",
		TargetBaseDir: "/path/to/target",
	}
	response := generatorlib.Render(ctx, request)
	if response.Success {
		os.Exit(0)
	} else {
		aulogging.Logger.Ctx(ctx).Error().Print("there were errors, return code is 1")
		os.Exit(1)
	}
}
