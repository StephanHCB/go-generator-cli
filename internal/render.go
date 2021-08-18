package internal

import (
	"context"
	generatorlib "github.com/StephanHCB/go-generator-lib"
	generatorapi "github.com/StephanHCB/go-generator-lib/api"
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

func Perform(ctx context.Context) bool {
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
