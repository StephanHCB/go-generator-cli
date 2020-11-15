package main

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/spf13/pflag"
)

// handling and validation of command line arguments

func parseCommandLine(ctx context.Context) bool {
	aulogging.Logger.Ctx(ctx).Info().Print("parsing command line arguments")

	pflag.StringVar(&generatorPath, "generator", "", "path to the generator base directory, with no trailing slash")
	pflag.StringVar(&targetPath, "target", "", "path to the target directory, with no trailing slash")
	pflag.StringVar(&createGenerator, "create", "", "write a specfile filled with generator defaults, called generated-<value>.yaml. If argument is omitted, it defaults to main.")
	pflag.Lookup("create").NoOptDefVal = "main"
	pflag.StringVar(&renderSpecfile, "render", "", "render according to a specfile. If specfile filename is omitted, it defaults to generated-main.yaml")
	pflag.Lookup("render").NoOptDefVal = "generated-main.yaml"
	pflag.Parse()

	success := true
	opscount := 0
	aulogging.Logger.Ctx(ctx).Info().Print("using arguments:")
	aulogging.Logger.Ctx(ctx).Info().Printf(" generator=%v", generatorPath)
	aulogging.Logger.Ctx(ctx).Info().Printf(" target=%v", targetPath)
	if createGenerator != "" {
		aulogging.Logger.Ctx(ctx).Info().Printf(" create=%v", createGenerator)
		aulogging.Logger.Ctx(ctx).Info().Printf("will write a default render spec to file 'generated-%v.yaml' in target for generator %v", createGenerator, createGenerator)
		opscount++
	}
	if renderSpecfile != "" {
		aulogging.Logger.Ctx(ctx).Info().Printf(" render=%v", renderSpecfile)
		aulogging.Logger.Ctx(ctx).Info().Printf("will render according to render spec file '%v'", renderSpecfile)
		opscount++
	}
	if opscount != 1 {
		aulogging.Logger.Ctx(ctx).Error().Print("invalid arguments, exactly one of `create` or 'render' is required")
		success = false
	}

	if generatorPath == "" {
		aulogging.Logger.Ctx(ctx).Error().Print("missing argument, `generator` is required (path to generator base directory with no trailing slash)")
		success = false
	}
	if targetPath == "" {
		aulogging.Logger.Ctx(ctx).Error().Print("missing argument, `target` is required (path to output base directory with no trailing slash)")
		success = false
	}
	return success
}
