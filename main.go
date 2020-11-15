package main

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	_ "github.com/StephanHCB/go-autumn-logging-log"
	generatorlib "github.com/StephanHCB/go-generator-lib"
	generatorapi "github.com/StephanHCB/go-generator-lib/api"
	"os"
)

var generatorPath string
var targetPath string
var createGenerator string
var renderSpecfile string

func create(ctx context.Context) bool {
	request := &generatorapi.Request{
		SourceBaseDir:  generatorPath,
		TargetBaseDir:  targetPath,
	}
	response := generatorlib.WriteRenderSpecWithDefaults(ctx, request, createGenerator)
	return response.Success
}

func render(ctx context.Context) bool {
	request := &generatorapi.Request{
		SourceBaseDir:  generatorPath,
		TargetBaseDir:  targetPath,
		RenderSpecFile: renderSpecfile,
	}
	response := generatorlib.Render(ctx, request)
	return response.Success
}

func perform(ctx context.Context) bool {
	if !parseCommandLine(ctx) {
		return false
	}

	if createGenerator != "" {
		return create(ctx)
	} else if renderSpecfile != "" {
		return render(ctx)
	} else {
		return false
	}
}

func main() {
	ctx := context.Background()
	aulogging.Logger.Ctx(ctx).Info().Print("welcome to go-generator-cli")
	if perform(ctx) {
		aulogging.Logger.Ctx(ctx).Info().Print("success")
		os.Exit(0)
	} else {
		aulogging.Logger.Ctx(ctx).Error().Print("there were errors, return code is 1")
		os.Exit(1)
	}
}
