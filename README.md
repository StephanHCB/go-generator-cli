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
  - source: 'src/some.text.file.copied.verbatim'
    target: 'text.txt'
    just_copy: true
  - source: 'src/web/controller.go.tmpl'
    target: 'web/controller/{{ .item }}.go'
    condition: '{{ if eq .item "skipped" }}false{{ end }}'
    with_items:
     - health
     - reservations
     - skipped
  - source: '{{ .file }}'
    target: 'output/{{ .file | replace ".tmpl" "" }}'
    with_files:
      - 'src/sub/*.tmpl'
      - '*/*.go.tmpl'
variables:
  serviceUrl:
    description: 'The URL of the service repository, to be used in imports etc.'
    default: 'github.com/StephanHCB/temp'
  serviceName:
    description: 'The name of the service to be rendered.'
    pattern: '^[a-z-]+$'
  helloMessage:
    description: 'A message to be inserted in the code.'
    default: 'hello {{ "world" }}'
```

This defines templates like `src/sub/sub.go.tmpl` and `src/main.go.tmpl` and what target path
they'll be rendered to in the target directory. 
It also specifies which parameter variables will be available during rendering.

  * If a variable does not have a default value, it is a required parameter.
  * default values are evaluated as templates, too, but you will not be able to refer to other variables  
  * if a variable has a pattern set, the parameter value must regex-match that pattern. Please be advised that
    you must enclose the pattern with ^...$ if you want to force the whole value to match, otherwise
    it's enough for part of the value to match the pattern.
  * variables are assumed to be string-valued by default, but the template generator actually allows any
    valid yaml structure (lists and maps, even nested) both as default values and as variable values.
    There is no type checking whatsoever, parsing templates that access missing fields or list items
    will fail, so it is not recommended to overuse this feature. Also, you should definitely provide
    a default value for any list or map typed variable, for else how will your users know what structure
    you are assuming?

The idea is that you keep your generators under version control.

You can create ansible-style loops using the same template to generate multiple output files using `with_items`.
In fact, the output file name is always parsed using the same template engine as the actual templates,
so you could also use other variables in it. 

If you set `with_items`, the template is used multiple times
with the `item` variable set to the value you provided under `with_items`. These values can also be 
a whole yaml data structure, you simply access it as `{{ .item.some.field }}`. 

_At this time, it is not possible to dynamically assign the full list in with_items from a variable, 
so you can not dynamically determine the number of render runs._

If you set `with_files` to a list of file globs relative to the template base directory, the source
and target pair is used multiple times, with the `file` variable set to the current relative path of
any files in the template that match the glob patterns. For `with_files`, both the source and target
path are parsed as templates, with the `file` parameter set to the relative path of that matched the glob.
You can even use `file` inside your template, it's just set as a parameter for template parsing.

_Note that globs that navigate outside the template directory are forbidden for security reasons._

You can also add a `condition` that will be evaluated for the template. Inside it, you can use
variables, or even `item`. If the condition evaluates to any one of `0`, `false`, `skip`, `no` the template will not be 
rendered. Note that the empty string counts as true, that means that if you do not specify a condition,
the template is rendered.

Any output directories are created for you on the fly if they don't exist.

### Template Language
  
The [golang template language](https://golang.org/pkg/text/template/#example_Template) is pretty 
versatile, vaguely similar to the .j2 templates used by ansible. Here's a very simple example
of how to include one of the parameters in your template output:

```
fmt.Println("{{ .helloMessage }}")
```

Assuming `someList` is a list variable, access the second entry as follows:

```
{{ index .someList 1 }}
```

Assuming `someMap` is a map variable with a field `message`, access it as follows:

```
{{ .someList.message }}
```

You can combine the two for structures with nested lists: `{{ (index .someList 0).someField }}`.

### Additional Template Functions

We include [Masterminds/sprig](https://github.com/Masterminds/sprig) when parsing any template,
which offers a collection of useful template functions, so you can do stuff like

```
{{ .helloMessage | upper }}
```

Read the sprig documentation, it adds much of what you would otherwise miss compared to ansible
j2 templates.

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

## Build and test

This program uses go modules. If cloned outside your GOPATH, you can build and test it using
`go build main.go` and `go test ./...`. This will also download all required dependencies.

### Acceptance Tests (give you examples)

[go-generator-lib](https://github.com/StephanHCB/go-generator-lib/tree/master/test/acceptance) 
has almost complete coverage with BDD-style acceptance tests.

In the course of the test runs, several generator specs and templates are read from the
[test resources](https://github.com/StephanHCB/go-generator-lib/tree/master/test/resources),
and a number of render specs are written to the
[test output directory](https://github.com/StephanHCB/go-generator-lib/tree/master/test/output).
