# sink (*work in progress*)

[![ci-build](https://github.com/lostmygithubaccount/sink/workflows/ci-build/badge.svg)](https://github.com/lostmygithubaccount/sink/actions/workflows/ci-build.yaml)
[![ci-lint](https://github.com/lostmygithubaccount/sink/workflows/ci-lint/badge.svg)](https://github.com/lostmygithubaccount/sink/actions/workflows/ci-lint.yaml)
[![release](https://github.com/lostmygithubaccount/sink/workflows/release/badge.svg)](https://github.com/lostmygithubaccount/sink/actions/workflows/release.yaml)

Windows is not supported. Use the Windows Subsystem for Linux (WSL) instead.

Welcome to the sink repository! The sink CLI and library is a tool for managing GitHub repositories *en masse*. It is designed for maintainers of numerous (micro)repositories who need to sync items across them. The CLI can be used manually or built into a CI/CD pipeline.

## usage

### verify installation

```bash
$ sink --version
```

If you get an error, first [install sink](#installation).

### sync labels

```bash
sink labels --help
```

### sync issues

```
sink issue --help
```

## installation

1. Install Go
2. [Optional] Install https://taskfile.dev

You can install from source with something like:

```bash
gh repo clone lostmygithubaccount/sink
cd sink
go build -o $HOME/go/bin .
```

Or if you have task installed:

```bash
gh repo clone lostmygithubaccount/sink
cd sink
task build
```

See the [`Taskfile.dist.yaml`](Taskfile.dist.yaml) for more examples, including using the GitHub releases.

## contributing

Contributions welcome, but not yet! See the [contributing guidelines](CONTRIBUTING.md) for more information.

## license

[The MIT License](LICENSE)
