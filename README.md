# Swagger YAML to CSV Converter

This repository provides a command-line tool that converts Swagger YAML files into CSV format. The primary command is `convert`, which reads Swagger YAML files specified in the configuration and outputs the API data in a CSV format.

## Features

- Converts Swagger YAML files into CSV format.
- Extracts relevant information such as `tag`, `operationID`, `summary`, `method`, and `path`.
- Handles multiple frameworks as defined in the configuration file.

## Prerequisites

- Go 1.16 or later
- `cobra` package for command-line interface support
- `viper` package for configuration management

## Installation

To install the required dependencies, run:

```bash
go get github.com/spf13/cobra
go get github.com/spf13/viper
```

## Usage

```
    go run main.go update
    go run main.go convert
```