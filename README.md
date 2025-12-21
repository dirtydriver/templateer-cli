# templateer-cli

A powerful and flexible project generator CLI tool written in Go that helps you quickly scaffold new projects from templates.

## Features

- Generate project skeletons from customizable templates
- Support for multiple project types
- Flexible parameter system via command line flags or parameter files
- Template parameter inspection
- Customizable output directory

## Installation

### From Source
```bash
# Clone the repository
git clone https://github.com/dirtydriver/templateer-cli.git
cd templateer-cli

# Build and install
go install
```

### Via Go Install
```bash
go install github.com/dirtydriver/templateer-cli@latest
```

## Usage

templateer-cli follows Unix-style command structure with subcommands. Two flags are required for all operations:
- `--template-dir`: Path to the template directory
- `--type`: Type of project (e.g., maven, gradle, angular)

### Available Commands

```bash
# Generate a new project
templateer-cli --template-dir <dir> --type <type> generate [flags]

# Inspect template parameters
templateer-cli --template-dir <dir> --type <type> inspect

# Show version
templateer-cli version
```

### Generate Command Flags

- `-n, --name`: Name of the project (can also be provided via --parameter name=value)
- `-o, --out`: Output directory (default: current directory)
- `-p, --parameter`: Additional parameters in key=value format (can be used multiple times)
- `-f, --file`: Path to a parameters file

### Examples

1. Generate a new project:
```bash
templateer-cli --template-dir ./templates --type maven generate --name my-project
```

2. Generate with custom parameters:
```bash
templateer-cli --template-dir ./templates --type angular generate \
  --parameter name=my-app \
  --parameter version=1.0.0 \
  --parameter author="John Doe"
```

3. Use a parameters file:
```bash
templateer-cli --template-dir ./templates --type gradle generate \
  --name my-lib \
  --file params.file
```

4. Check template parameters:
```bash
templateer-cli --template-dir ./templates --type maven inspect
```

## Template System

templateer-cli uses a powerful templating system that allows you to create and customize project templates. Templates are stored in the `templates` directory and use the `.tmpl` extension.

### Template Structure
```
templates/
├── maven/           # Template for Maven projects
│   ├── pom.xml.tmpl
│   └── src/
├── gradle/          # Template for Gradle projects
│   ├── build.gradle.tmpl
│   └── settings.gradle.tmpl
└── angular/         # Template for Angular projects
    └── ...
```

### Go Template Syntax
templateer-cli uses Go's built-in template engine. For detailed documentation, visit the [official text/template package documentation](https://pkg.go.dev/text/template).

Common template syntax used in project generation:
```
{{.name}}           # Access a parameter value
{{.group_id}}       # Use snake_case for parameter names
{{if .test}}        # Conditional block
  {{.test}}
{{end}}
{{range .items}}    # Loop through array/slice
  {{.}}
{{end}}
{{$var := .name}}   # Assign to variable
{{title .name}}     # Use 'title' function to capitalize
```

### Parameter Files
You can create parameter files to store commonly used values. Parameter files use YAML format:
```yaml
name: my-project
version: 1.0.0
author: John Doe
description: My awesome project
# You can also use nested structures
metadata:
  team: backend
  priority: high
dependencies:
  - mysql
  - redis
```

### Creating Custom Templates
1. Create a new directory in `templates/` for your project type
2. Add template files with the `.tmpl` extension
3. Use Go template syntax for variable substitution: `{{.variable_name}}`
4. Use `templateer-cli inspect` to check required parameters

## Project Structure

```
templateer-cli/
├── cmd/          # Command line interface implementation
├── filescheck/   # File system operations and checks
├── project/      # Project generation logic
├── templater/    # Template processing and rendering
├── utils/        # Utility functions
└── version/      # Version information
```

## Development

To contribute to templateer-cli:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the terms of the LICENSE file included in the repository.
