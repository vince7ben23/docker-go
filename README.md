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

## Troubleshooting: `bad CPU type in executable`                                                                                                                                                     
如果你在 Apple Silicon 的 macOS host 上直接執行 Stage 2+（而不是在 Docker 容器內執行），可能會看到：

```
Err: fork/exec /usr/local/bin/docker-explorer: bad CPU type in executable
```

這是因為 `docker-explorer` 是 x86_64 的執行檔。在 Apple Silicon 上，macOS 原本會透過 Rosetta 2 自動轉譯執行它，但程式呼叫 `syscall.Chroot` 之後，Rosetta 無法再存取轉譯所需的系統路徑，導致執行失敗。請務必照上方說明在 Linux 容器內執行 Stage 2+，不要直接在 macOS host 上執行。

## Development

```sh
go build -o /tmp ./app/...
go test -timeout 10s ./app/...
go vet ./...
```

Pre-commit hooks run `gofmt`, `go vet`, `go build`, and `go test` automatically.
