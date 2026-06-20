# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go implementation of the CodeCrafters ["Build Your Own Docker" challenge](https://codecrafters.io/challenges/docker). The goal is to build a program that pulls images from Docker Hub and executes commands inside them, implementing chroot, Linux kernel namespaces, and the Docker Registry API from scratch.

## Commands

**Build:**
```sh
go build -o /tmp ./app/...
```

**Test:**
```sh
go test -timeout 10s ./app/...
```

**Vet / Format:**
```sh
go vet ./...
gofmt -w .
```

**Run pre-commit hooks against all files:**
```sh
pre-commit run -a
```

**Run (Stage 1 — no Linux syscalls needed):**
```sh
./your_docker.sh run <image> <command path> [args...]
```

**Run (Stage 2+ — requires Linux syscalls, must run inside Docker):**
```sh
# Set up alias first:
alias mydocker='docker build -t docker_from_scratch . && docker run --cap-add="SYS_ADMIN" docker_from_scratch'

# Then run:
mydocker run alpine:latest /usr/local/bin/docker-explorer echo hey
```

The `--cap-add="SYS_ADMIN"` flag is required for PID namespace creation.

## Architecture

**Entry point:** `app/main.go` — `main()` reads CLI arguments in the form `run <image> <command path> [args...]`, so `os.Args[2]` is the image, `os.Args[3]` is the command, and `os.Args[4:]` are its arguments.

**Bootstrap:** `your_docker.sh` compiles the Go binary and passes all arguments through to it. The Dockerfile sets this as the container ENTRYPOINT.

**Test utility:** `docker-explorer` (downloaded into the container at `/usr/local/bin/docker-explorer`) is a CodeCrafters-provided binary used to test Docker Registry API and namespace features.

**Pre-commit hooks** (`.pre-commit-config.yaml`): automatically runs `gofmt`, `go vet`, `go build`, and `go test` on staged Go files before each commit.

## Challenge Stages

Each CodeCrafters stage adds new functionality to `app/main.go`:
1. Execute a command in the current environment (basic `exec.Command`)
2. Isolate using chroot + filesystem setup
3. PID namespace isolation
4. Pull images from Docker Hub using the Registry API
