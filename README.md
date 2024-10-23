<p align="center">
  <img alt="golangci-lint logo" src="assets/go.png" height="150" />
  <h3 align="center">godyl</h3>
  <p align="center">Asset downloader</p>
</p>

---

[![Go Reference](https://pkg.go.dev/badge/github.com/idelchi/godyl.svg)](https://pkg.go.dev/github.com/idelchi/godyl)
[![Go Report Card](https://goreportcard.com/badge/github.com/idelchi/godyl)](https://goreportcard.com/report/github.com/idelchi/godyl)
[![Build Status](https://github.com/idelchi/godyl/actions/workflows/github-actions.yml/badge.svg)](https://github.com/idelchi/godyl/actions/workflows/github-actions.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`godyl` helps with batch-downloading and installing statically compiled binaries from:

- GitHub releases
- URLs
- Go projects

As an alternative to above, custom commands can be used as well.

`godyl` will infer the platform and architecture from the system it is running on, and will attempt to download the appropriate binary.

This uses simple heuristics to select the correct binary to download, and will not work for all projects.

However, most properties can be overridden, with `hints` and `skip` used to help `godyl` make the correct decision.

> [!WARNING]
> This repo is a work in progress!
> Needing both cleaning up and documenting.

> [!CAUTION]
> This project serves as a learning exercise for Go and its surrounding ecosystem and tooling.
> As such, it might be of limited use for others, but feel free to use it if you find it useful.

> [!NOTE]
> Tested on:
>
> **Linux**: `amd64`, `arm64`, `armv7`
>
> **Windows**: `amd64`
>
> **MacOS**: `amd64`, `arm64`
>
> for tools listed in [tools.yml](./tools.yml)

Tool is inspired by [task](https://github.com/go-task/task), [dra](https://github.com/devmatteini/dra) and [ansible](https://github.com/ansible/ansible)

---

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Tools](#tools)
  - [Simple form](#simple-form)
  - [Full form](#full-form)
- [Defaults](#defaults)
- [Template overview](#template-overview)
  - [Variables](#variables)
  - [Allowed in](#allowed-in)
- [Inference](#Inference)
  - [Operating Systems](#operating-systems)
  - [Architectures](#architectures)
  - [Libraries](#libraries)
- [Notes](#notes)

---

## Installation

### From source

```sh
go install github.com/idelchi/godyl/cmd/godyl@latest
```

## From installation script

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/scripts/install.sh | sh -s -- -v v0.1-beta -o ~/.local/bin
```

## Update

```sh
godyl --update
```

## Usage

Use together with `yaml` file:

```sh
godyl tools.yml  --output ./bin
```

Or use to download a single tool:

```sh
godyl idelchi/godyl --output ./bin
```

When used to download a single tool, the `mode` will be set to `extract` by default.

Override `os` and `arch` to download a specific binary:

```sh
godyl idelchi/godyl --os linux --arch amd64 --output ./bin
```

> [!NOTE]
> Set up a GitHub API token to avoid rate limiting when using `github` as a source type.
> See [configuration](#configuration) for more information, or simply `export GODYL_GITHUB_TOKEN=<token>`

## Configuration

The tools can be configured (in order of priority) by

- flags
- environment variables
- `.env` file

The following flags and their corresponding environment variables are available:

| Flag               | Environment Variable  | Default        | Description                                    |
| ------------------ | --------------------- | -------------- | ---------------------------------------------- |
| `--help`, `-h`     | `GODYL_HELP`          | `false`        | Show help message and exit                     |
| `--version`        | `GODYL_VERSION`       | `false`        | Show version information and exit              |
| `--dot-env`        | `GODYL_DOT_ENV`       | `.env`         | Path to .env file                              |
| `--defaults`, `-d` | `GODYL_DEFAULTS`      | `defaults.yml` | Path to defaults file                          |
| `--show-config`    | `GODYL_SHOW_CONFIG`   | `false`        | Show the parsed configuration and exit         |
| `--show-defaults`  | `GODYL_SHOW_DEFAULTS` | `false`        | Show the parsed default configuration and exit |
| `--show-env`       | `GODYL_SHOW_ENV`      | `false`        | Show the parsed environment variables and exit |
| `--show-platform`  | `GODYL_SHOW_PLATFORM` | `false`        | Detect the platform and exit                   |
| `--update`         | `GODYL_UPDATE`        | `false`        | Update `godyl` itself                          |
| `--dry`            | `GODYL_DRY`           | `false`        | Run without making any changes (dry run)       |
| `--log`            | `GODYL_LOG`           | `info`         | Log level (debug, info, warn, error)           |
| `--parallel`, `-j` | `GODYL_PARALLEL`      | `0`            | Number of parallel downloads (0 is unlimited)  |
| `--output`         | `GODYL_OUTPUT`        | `""`           | Output path for the downloaded tools           |
| `--tags`, `-t`     | `GODYL_TAGS`          | `["!native"]`  | Tags to filter tools by. Use `!` to exclude    |
| `--source`         | `GODYL_SOURCE`        | `github`       | Source from which to install the tools         |
| `--strategy`       | `GODYL_STRATEGY`      | `none`         | Strategy to use for updating tools             |
| `--os`             | `GODYL_OS`            | `""`           | Operating system to use for downloading        |
| `--arch`           | `GODYL_ARCH`          | `""`           | Architecture to use for downloading            |
| `--github-token`   | `GODYL_GITHUB_TOKEN`  | `""`           | GitHub token for authentication                |

The path to the file containing the tool installation instructions is provided as a positional argument, defaulting to `tools.yml`.

An example [tools.yml](./tools.yml) is provided.

In general, settings can be set in the following ways (order of priority):

- as a field in the [tools.yml](./tools.yml) definition

  ```yaml
  output: ~/.local/bin
  ```

- as a flag to the tool

  ```sh
  godyl --output ~/.local/bin
  ```

- as an environment variable

  ```sh
  GODYL_OUTPUT=~/.local/bin godyl
  ```

- in an `.env` file

  ```
  GODYL_OUTPUT=~/.local/bin
  ```

- by setting the value in a `defaults.yml` file (see [defaults](#defaults))

  ```yaml
  output: ~/.local/bin
  ```

If none of the above are fulfilled, the default configuration embedded from [defaults.yml](./cmd/godyl/defaults.yml) will be used.

## Tools

A YAML file controls the tools to download and install. Alternatively, if the positional argument to the tool is not a YAML file, it will be treated as a single tool name or URL.

Examples are provided in [tools.yml](./tools.yml).

### Simple form

```yaml
- idelchi/godyl
```

Above is the `simple` form to attempt to download the latest release of `godyl` from `idelchi/godyl`.

If it is a URL, it will be considered as a `source.url` type.
Otherwise, it will be assumed to be a `source.github` type on the form `owner/repo`.

### Full form

```yaml
name: string
description: string
version: string
path: string
output: string
exe:
  name: string
  patterns:
    - regex
platform:
  os: string
  architecture:
    type: string
    version: string
  library: string
  extension: string
  distribution: string
aliases: []
values: {}
fallbacks: []
hints:
  - pattern: regex
    weight: string
    must: boolean
source:
  type: string
  github:
    repo: string
    owner: string
    token: string
  url:
    url: string
    token: string
  go:
    command: string
  commands: []
tags:
  - string
strategy: string
extensions:
  - string
skip:
  - condition: string
    reason: string
post: []
mode: string
env:
  key: string
```

Any field that accepts a list, can also be provided as a string.

For example:

```yaml
aliases: gd
fallbacks: go
exe:
  patterns: "{{ .Exe }}{{ .EXTENSION }}$"
```

is equivalent to:

```yaml
aliases:
  - gd
fallbacks:
  - go
exe:
  patterns:
    - "{{ .Exe }}{{ .EXTENSION }}$"
```

and

```yaml
skip:
  reason: "tool is not available on windows"
  condition: '{{ eq .OS "windows" }}'
```

is equivalent to:

```yaml
skip:
  - reason: "tool is not available on windows"
    condition: '{{ eq .OS "windows" }}'
```

### Name

![Required](https://img.shields.io/badge/Required-Yes-green)

| Template      | Templated | As Template |
| ------------- | --------- | ----------- |
| `{{ .Name }}` | ![no]     | ![yes]      |

`name` is the name of the tool to download.

#### Usage

- Used as display name and for inferring other fields (see [GitHub](#github) and [exe](#exe))

### Description

![Optional](https://img.shields.io/badge/Optional-Yes-blue)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`description` is an optional description of the tool, for documentation purposes.

### Version

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template         | Templated | As Template |
| ---------------- | --------- | ----------- |
| `{{ .Version }}` | ![yes]    | ![yes]      |

#### Usage

- As template for other fields
- To compare tool versions if `strategy` is set to `upgrade`
- Will be inferred and populated by the `source` method if not given

### Path

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![yes]    | ![no]       |

`path` is the path to the tool to download. Currently only supports URLs.

#### Usage

- Will be inferred and populated by the `source` method if not given

### Output

![Required](https://img.shields.io/badge/Required-red)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![yes]    | ![yes]      |

`output` is the path to the directory where the tool will be installed.

#### Usage

- Set according to [flags and environment variables](#configuration) or [defaults](#defaults) if not given

### Exe

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template              | Templated | As Template |
| --------------------- | --------- | ----------- |
| `{{ .Exe}}`           | ![yes]    | ![yes]      |
| `{{ .Exe.Patterns }}` | ![yes]    | ![no]       |

`exe` is a dictionary containing the name of the executable and patterns to use for finding the executable.

#### Usage

- `exe.name` is the name of the executable, inferred from `name` if not given
- `exe.patterns` is a list of patterns to use for finding the executable, inferred from `name` if not given
- Set according to [defaults](#defaults) if not given

#### Alternative form

```yaml
exe: godyl
```

is equivalent to:

```yaml
exe:
  name: godyl
```

### Platform

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template             | Templated | As Template                 |
| -------------------- | --------- | --------------------------- |
| `{{ .Platform.<> }}` | ![no]     | See [variables](#variables) |

`platform` is a dictionary containing the platform and architecture information.

#### Usage

- Any field not given will be inferred from the system
- Set according to [defaults](#defaults) if not given

### Aliases

![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`aliases` is a list of aliases for the tool. Will be used to create symlinks (or copies on `Windows`) for the tool.

### Values

![Optional](https://img.shields.io/badge/Optional-green)

| Template           | Templated | As Template |
| ------------------ | --------- | ----------- |
| `{{ .Values.<> }}` | ![no]     | ![yes]      |

`values` is an arbitrary values map, which can be used for templating in other fields.

### Fallbacks

![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`fallbacks` is a list of fallback strategies to use if the tool cannot be found.

They will be tried in order until the tool is found or all have been tried.

### Hints

![Optional](https://img.shields.io/badge/Optional-green)

| Template               | Templated | As Template |
| ---------------------- | --------- | ----------- |
| `{{ .Hints.Weight }}`  | ![yes]    | ![no]       |
| `{{ .Hints.Pattern }}` | ![yes]    | ![no]       |
| `{{ .Hints.Regex }}`   | ![no]     | ![no]       |
| `{{ .Hints.Must }}`    | ![no]     | ![no]       |

`hints` is a list of hints for matching, which can be used to help `godyl` find the correct tool.

#### Usage

- Set according to [defaults](#defaults) if not given

### Source

#### Type

![Required](https://img.shields.io/badge/Required-red)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`source.type` is the source type. Accepted values are:

- `github`
- `url`
- `go`
- `commands`

#### GitHub

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`source.github` is a dictionary containing the owner, repository and token of the tool.

#### Usage

- `repo` and `owner` will be inferred from `name` if not given, or set according to [defaults](#defaults) (not recommended)
- `token` will be set according to [flags and environment variables](#configuration) or [defaults](#defaults) if not given

#### URL

![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`source.url` is a dictionary containing the URL and token of the tool.

#### Usage

- `token` will be set according to [flags and environment variables](#configuration) if not given

#### Go

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`source.go` can be set to the relative path of the go `command` to download, if non-standard (i.e not matching `<name>`, `cmd/<name>` or `cmd`).

> [!WARNING]
> Go will be downloaded into a temporary directory `/tmp/.godyl-go` if not present.

#### Usage

#### Commands

![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![yes]    | ![no]       |

`commands` is a list of commands to run to install the tool.

### Tags

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`tags` is a list of tags to filter tools by.

#### Usage

- the `name` of the tool will always be added as a tag

### Strategy

![Required](https://img.shields.io/badge/Required-red)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`strategy` is the strategy to use for updating the tool.

Accepted values are:

- `none`
- `upgrade`
- `force`

#### Usage

- Set according to [flags and environment variables](#configuration) or [defaults](#defaults) if not given
- `none` will skip the tool if it already exists
- `upgrade` will attempt to parse the version of an existing tool and upgrade if necessary
- `force` will always download and install the tool

### Extensions

![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![yes]    | ![no]       |

`extensions` is a list of extensions to consider when matching the most suitable tool to download (used only for `github` source type).

#### Usage

- Set according to [defaults](#defaults) if not given
- Can be used to for example prefer `.zip` files for Windows

### Skip

![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![yes]    | ![no]       |

`skip` is a list of conditions under which the tool should be skipped.

#### Usage

- `reason` is a description of why the tool was skipped
- `condition` is a condition to check, and can use templating

#### Alternative form

```yaml
skip: <condition>
```

is equivalent to:

```yaml
skip:
  - condition: <condition>
```

### Post

![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![yes]    | ![no]       |

`post` is a list of commands to run after the tool has been installed.

### Mode

![Required](https://img.shields.io/badge/Required-red)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`mode` is the mode to use for the tool.

Accepted values are:

- `find`
- `extract`

#### Usage

- `find` will download, extract and find the executable
- `extract` will download the tool and extract it directly to the output directory
- Set according to [flags and environment variables](#configuration) or [defaults](#defaults) if not given
- Automatically set to `extract` if the tool is used without `tools.yml` (e.g. `godyl idelchi/godyl`)

## Defaults

A default configuration may be used to specify default settings for all tools. These will override (or extend in some case) the settings for each tool.

The following is embedded and used by default if no default configuration is provided:

[defaults.yml](./cmd/godyl/defaults.yml)

The example above defines:

- The default output directory for all tools (`~/.local/bin`)
- Patterns to use for when searching for the executable

  - `^{{ .Exe }}{{ .EXTENSION }}$`
  - `.\*/{{ .Exe }}{{ .EXTENSION }}$`

- Hints to:

  - use the executable name as a pattern (useful for repositories with multiple binaries, such as `ahmetb/kubectx`)

    ```yaml
    pattern: "{{ .Exe }}"
    weight: 1
    ```

  - prefer `.zip` files for Windows

    ```yaml
    pattern: zip
    weight: '{{ if eq .OS "windows" }}1{{ else }}0{{ end }}'
    ```

  - prefer the three last versions of the arm architecture, in descending order
    ```
    - pattern: "armv{{.ARCH_VERSION}}"
      weight: 3
    - pattern: "armv{{sub .ARCH_VERSION 1}}"
      weight: 2
    - pattern: "armv{{sub .ARCH_VERSION 2}}"
      weight: 1
    ```

- Extensions to use when filtering assets

  - `.exe` for Windows, empty for Linux and MacOS
  - `.zip` for Windows and Darwin
  - `.tar.gz` for all platforms

- `find` mode for downloading, extracting and finding the executable
- The default source type as `github`
- `none` strategy to skip tools which already exist
- Settings the environment variable `GH_TOKEN` to the value of `GODYL_GITHUB_TOKEN`

The full set of default options are:

```yaml
exe:
output:
platform:
values:
fallbacks:
hints:
source:
tags:
strategy:
env:
mode:
```

For reference full reference of what values you can set, see the [tools](#tools) section.

## Template overview

All functions available in the [slim-sprig](https://github.com/go-task/slim-sprig) library are available for use in templates.

### Variables

The following table lists the available template variables, where they may be used, and their descriptions:

| Variable              | Description                                               |
| --------------------- | --------------------------------------------------------- |
| `{{ .Name }}`         | The name of the tool or project                           |
| `{{ .Output }}`       | The output path template for built artifacts              |
| `{{ .Exe }}`          | The name of the executable                                |
| `{{ .Env.<> }}`       | Any environment variable                                  |
| `{{ .Values.<> }}`    | Custom values for templating                              |
| `{{ .Version }}`      | The version of the tool or project                        |
| `{{ .OS }}`           | The operating system (e.g., `linux`, `darwin`, `windows`) |
| `{{ .ARCH }}`         | The architecture type (e.g., `amd64`, `arm64`)            |
| `{{ .ARCH_VERSION }}` | The version of the architecture, if applicable            |
| `{{ .LIBRARY }}`      | The system library (e.g., `gnu`, `musl`)                  |
| `{{ .EXTENSION }}`    | The file extension specific to the platform               |
| `{{ .DISTRIBUTION }}` | The distribution name (e.g., `debian`, `alpine`)          |

### Allowed in

Only certain fields are templated. Below is a list of fields where templating is allowed, along with examples of how they might be used:

- `output`

  ```yaml
  output: bin/{{ .OS }}-{{ .ARCH }}
  ```

- `skip`

  ```yaml
  skip:
    reason: "tool is not available for windows"
    condition: '{{ eq .OS "windows" }}'
  ```

- `version`

  ```yaml
  version: |-
    {{- if has .OS (list "linux" "darwin") -}}
        v0.1.0
    {{- else -}}
        v0.2.0
    {{- end -}}
  ```

- `source.type`

  ```yaml
  source:
    type: |-
      {{- if has .OS (list "linux" "darwin") -}}
          github
      {{- else -}}
          go
      {{- end -}}
  ```

- `exe.patterns`

  ```yaml
  exe:
    patterns:
      - "^{{ .OS }}-{{ .Exe}}"
  ```

- `extensions`

  ```yaml
  extensions:
    - '{{ if eq .OS "windows" }}.exe{{ else }}{{ end }}'
  ```

- `source.commands`

  ```yaml
  commands:
    - pip install {{ .Exe}}=={{ .Version }}
  ```

- `hints[].pattern`

  ```yaml
  hints:
    - pattern: "{{ .OS }}"
      must: true
  ```

- `hints[].weight`

  ```yaml
  hints:
    - pattern: armhf
      weight: |-
        {{- if eq .ARCH "arm" -}}
        1
        {{- else -}}
        0
        {{- end -}}
  ```

- `path`

  ```yaml
  path: https://get.helm.sh/helm-v{{ .Version }}-{{ .OS }}-{{ .ARCH }}.tar.gz
  ```

> [!NOTE]
> Not that if `Version` is not provided, it will be evaluated to an empty string. If it is being inferred, it will be available in the following fields:
>
> - `exe.name`
> - `exe.patterns`
> - `extensions`
> - `hints[].pattern`
> - `hints[].weight`
> - `path`
> - `source.commands`
> - `post`

## Inference

### Operating Systems

| OS      | Inferred from           |
| ------- | ----------------------- |
| Linux   | linux                   |
| Darwin  | darwin, macos, mac, osx |
| Windows | windows, win            |

### Architectures

| Architecture       | Inferred from                           |
| ------------------ | --------------------------------------- |
| AMD64              | amd64, x86_64, x64, win64               |
| ARM64              | arm64, aarch64                          |
| AMD32              | amd32, x86, i386, i686, win32, 386, 686 |
| ARM32 (v7)         | armv7, armv7l, armhf                    |
| ARM32 (v6) \*      | armv6, armv6l                           |
| ARM32 (v5)         | armv5, armel                            |
| ARM32 (v<unknown>) | arm                                     |

### Libraries

| Library    | Inferred from |
| ---------- | ------------- |
| GNU        | gnu           |
| Musl       | musl          |
| MSVC       | msvc          |
| LibAndroid | android       |

## Notes

All `regex` expressions are evaluated using `search`, meaning that `^` and `$` are necessary to match the start and end of the string.
When running `32-bit` userland on a `64-bit` Kernel, there's some attempts to infer the matching `32-bit` architecture.

However, to be certain that the right binary is downloaded, it's recommended to pass the `--arch` flag to the tool.

<!-- Badges -->

[yes]: https://img.shields.io/badge/Yes-green
[no]: https://img.shields.io/badge/No-red
[inferred]: https://img.shields.io/badge/Inferred-blue
[required]: https://img.shields.io/badge/Required-red
[optional]: https://img.shields.io/badge/Optional-green
[not-implemented]: https://img.shields.io/badge/Not%20Implemented-gray
[na]: https://img.shields.io/badge/N%2FA-lightgrey

<!-- Badges -->

```

```
