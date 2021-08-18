# go-generator-cli

A golang command line utility for generating files from templates. Can be used to scaffold application code.

Uses [go-generator-lib](https://github.com/StephanHCB/go-generator-lib/).

## Command line parameters

  * `--generator=<path>` to set the path of the generator base directory.
  * `--target=<path>` to set the path of the target directory.
  * `--create[=<generator name>]` to write a specfile with defaults. If the generator name is omitted, 
     it defaults to `main`, and the file written will be called `generated-main.yaml`. 
  * `--render[=<specfile>]` to use a specfile to render. If the filename is omitted, it defaults to `generated-main.yaml`.

You will need to specify exactly one of `--create` and `--render`, while the other two arguments are mandatory.

## Generators

A generator is a directory that contains one or more `generator-*.yaml` files, called 
*generator specification files*, plus a number of 
golang [text/templates](https://golang.org/pkg/text/template/). The main generator spec is
usually called `generator-main.yaml`.

Example:
```
templates:
  - source: 'src/sub/sub.go.tmpl'
    target: 'sub/sub.go'
  - source: 'src/main.go.tmpl'
    target: 'main.go'
  - source: 'src/web/controller.go.tmpl'
    target: 'web/controller/{{ .item }}.go'
    with_items:
     - health
     - reservations
variables:
  serviceUrl:
    description: 'The URL of the service repository, to be used in imports etc.'
    default: 'github.com/StephanHCB/temp'
  serviceName:
    description: 'The name of the service to be rendered.'
    pattern: '^[a-z-]+$'
  helloMessage:
    description: 'A message to be inserted in the code.'
    default: 'hello world'
```

This defines templates like `src/sub/sub.go.tmpl` and `src/main.go.tmpl` and what target path
they'll be rendered to in the target directory. 
It also specifies which parameter variables will be available during rendering.

  * If a variable does not have a default value, it is a required parameter.
  * if a variable has a pattern set, the parameter value must regex-match that pattern. Please be advised that
    you must enclose the pattern with ^...$ if you want to force the whole value to match, otherwise
    it's enough for part of the value to match the pattern.

The idea is that you keep your generators under version control.

Note how you can create ansible-style loops using the same template to generate multiple output files using `with_items`.
In fact, the output file name is always parsed using the same template engine as the actual templates,
so you could also use other variables in it. 

If you set `with_items`, the template is used multiple times
with the `item` variable set to the value you provided under `with_items`. These values can also be 
a whole yaml data structure, you simply access it as `{{ .item.some.field }}`. 

Also note how output directories are created for you on the fly if they don't exist.
  
The [golang template language](https://golang.org/pkg/text/template/#example_Template) is pretty 
versatile, vaguely similar to the .j2 templates used by ansible. Here's a very simple example
of how to include one of the parameters in your template output:

```
fmt.Println("{{ .helloMessage }}")
```
## Render Targets

A render target is a directory that contains a yaml file which records the name of the generator used
and all parameter values. We call this a *render specification file*. If you do not specify anything,
the generator expects the file to be called `generated-main.yaml`.

Example:
```
generator: main
parameters:
  helloMessage: hello world
  serviceName: 'my-service'
  serviceUrl: github.com/StephanHCB/temp
```
# Build and test

This program uses go modules. If cloned outside your GOPATH, you can build and test it using
`go build main.go` and `go test ./...`. This will also download all required dependencies.
