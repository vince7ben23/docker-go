# docker-from-scratch

A Docker implementation built from scratch in Go. Pulls images from Docker Hub and executes commands inside isolated environments using chroot, Linux kernel namespaces, and the Docker Registry API.

## Usage

**Stage 1 — run a command directly:**
```sh
./bootstrap.sh run <image> <command> [args...]
```

**Stage 2+ — requires Linux syscalls (must run inside Docker):**
```sh
alias mydocker='docker build -t docker_from_scratch . && docker run --cap-add="SYS_ADMIN" docker_from_scratch'

mydocker run alpine:latest /usr/local/bin/docker-explorer echo hey

# means you execute
docker build -t docker_from_scratch . && docker run --cap-add="SYS_ADMIN" docker_from_scratch run alpine:latest /usr/local/bin/docker-explorer echo hey
```

The `--cap-add="SYS_ADMIN"` flag is required for [PID namespace](https://man7.org/linux/man-pages/man7/pid_namespaces.7.html) creation.

## Development

```sh
go build -o /tmp ./app/...
go test -timeout 10s ./app/...
go vet ./...
```

Pre-commit hooks run `gofmt`, `go vet`, `go build`, and `go test` automatically.
