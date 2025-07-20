[![Go Reference](https://pkg.go.dev/badge/github.com/qba73/gov.svg)](https://pkg.go.dev/github.com/qba73/gov)
[![CI](https://github.com/qba73/gov/actions/workflows/go.yml/badge.svg)](https://github.com/qba73/gov/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/qba73/gov)](https://goreportcard.com/report/github.com/qba73/gov)
[![Scorecard](https://github.com/qba73/gov/actions/workflows/scorecard.yml/badge.svg)](https://github.com/qba73/gov/actions/workflows/scorecard.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/qba73/gov)
![GitHub License](https://img.shields.io/github/license/qba73/gov)
![Maintenance](https://img.shields.io/badge/maintenance-actively--developed-brightgreen.svg)



# gov
Go dependencies parser and formatter

`gov` is a Go package and a command-line utility for parsing output from the `go version -v -m <binary>` command.

`gov` parses packages included in the Go binary and outputs the report in JSON format.
The JSON format makes it easier to process dependencies further and embed them in, for example, security reports such as [SLSA](https://slsa.dev/spec/v1.2-rc1/build-provenance) or [in-toto](https://in-toto.io).

Every time you need to retrieve information about the packages in your binary and you require this information in JSON format, you can use the `gov` utility.

## How it works

Let's say we build a Go binary and check what dependencies it includes. As an example, we will build NGINX [Kubernetes Ingress Controller](https://github.com/nginx/kubernetes-ingress).

1. Run `make build` from the root dir.

```shell
make build
Docker version 28.1.1, build 4eba377
go version go1.24.5 darwin/arm64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w -X main.version=5.2.0-SNAPSHOT -X main.telemetryEndpoint=oss.edge.df.f5.com:443" -o nginx-ingress github.com/nginx/kubernetes-ingress/cmd/nginx-ingress
```

2. Check if the `nginx-ingress` is in the root dir.

```shell
ls | grep nginx
nginx-ingress
```

3. List Go packages included in the `nginx-ingress` binary.

```shell
go version -v -m nginx-ingress
```

```shell
nginx-ingress: go1.24.5
	path	github.com/nginx/kubernetes-ingress/cmd/nginx-ingress
	mod	github.com/nginx/kubernetes-ingress	v1.12.1-0.20250718155242-1c47ba285898+dirty
	dep	github.com/beorn7/perks	v1.0.1	h1:VlbKKnNfV8bJzeqoa4cOKqO6bYr3WgKZxO8Z16+hsOM=
	dep	github.com/blang/semver/v4	v4.0.0	h1:1PFHFE6yCCTv8C1TeyNNarDzntLi7wMI5i/pzqYIsAM=
// ...
	dep	sigs.k8s.io/yaml	v1.5.0	h1:M10b2U7aEUY6hRtU870n2VTPgR5RZiL/I6Lcc2F4NUQ=
	build	-buildmode=exe
// ...
```

4. List dependencies (names, versions and digests) in JSON format.

```shell
go version -v -m nginx-ingress | gov | jq .
```

```json
[
  {
    "name": "github.com/beorn7/perks",
    "version": "v1.0.1",
    "digest": "h1:VlbKKnNfV8bJzeqoa4cOKqO6bYr3WgKZxO8Z16+hsOM="
  },
  {
    "name": "github.com/blang/semver/v4",
    "version": "v4.0.0",
    "digest": "h1:1PFHFE6yCCTv8C1TeyNNarDzntLi7wMI5i/pzqYIsAM="
  },
  {
    "name": "github.com/cenkalti/backoff/v5",
    "version": "v5.0.2",
    "digest": "h1:rIfFVxEf1QsI7E1ZHfp/B4DF/6QBAUhmgkxc0H7Zss8="
  },
  // ...
   {
    "name": "github.com/cert-manager/cert-manager",
    "version": "v1.18.2",
    "digest": "h1:H2P75ycGcTMauV3gvpkDqLdS3RSXonWF2S49QGA1PZE="
  },
]
```

## Installation

Install Go binary
```shell
go install github.com/qba73/gov/cmd/gov@latest
```
Get help
```shell
gov -h
Usage: gov [files...]
Parses Go dependencies from files or standard input and returns JSON output.
Input is an output from the Go command: `go version -v -m <go-binary>`
```

## Usage

### Piping output from go version command

Get raw JSON:

```shell
go version -v -m <binary> | gov
```

Pipe JSON to `jq` for further formatting or filtering:

```shell
go version -v -m <binary> | gov | jq .
```

### Parsing a file containing output from go version command

Save output to a file

```shell
go version -v -m nginx-ingress > nginx-ingress.out
```

Parse the file using `gov`

```shell
gov nginx-ingress.out
```

Parse the file using `gov` and pipe to `jq`

```shell
gov nginx-ingress.out | jq .
```

### Parsing multiple files at once:

Parse multiple files and generate one JSON report

```shell
gov nginx-ingress.out kubernetes-gateway.out | jq .
```
