[![CI](https://github.com/qba73/gov/actions/workflows/go.yml/badge.svg)](https://github.com/qba73/gov/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/qba73/gov)](https://goreportcard.com/report/github.com/qba73/gov)
[![Scorecard](https://github.com/qba73/gov/actions/workflows/scorecard.yml/badge.svg)](https://github.com/qba73/gov/actions/workflows/scorecard.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/qba73/gov)
![GitHub License](https://img.shields.io/github/license/qba73/gov)
![Maintenance](https://img.shields.io/badge/maintenance-actively--developed-brightgreen.svg)



# gov
Go version parser and formatter

`gov` is a Go package and a command-line utility for parsing output from the `go version -v -m <binary>` command.

It parses packages included in the Go binary and outputs the report in JSON format.
The JSON format makes it easier to process dependencies further and embed them in, for example, security reports such as [SLSA](https://slsa.dev/spec/v1.2-rc1/build-provenance) or [in-toto](https://in-toto.io).

Every time you need to retrieve information about the packages in your binary and you require this information in JSON format, you can use the `gov` utility.

