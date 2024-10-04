# GoFiN-template

A template for a GoFiN project.

## Getting Started

### Prerequisites

- [Bun](https://bun.sh/)
- [Go](https://go.dev/)

### Usage

1. Clone the repository

```bash
git clone https://github.com/juls0730/gofin.git
```

2. Run the project

To run the project in dev mode with hot reloading:

```bash
bun --cwd=./ui install
CompileDaemon --build="go build -tags netgo,dev -ldflags=-s" --command=./gofin --exclude-dir=data/ --exclude-dir=ui/ --graceful-kill
```

To run the project with Nuxt SSG:

```bash
bun --cwd=./ui install
RENDERING_MODE=static bun --bun --cwd=./ui run generate
go build -tags netgo -ldflags=-s
./gofin
```

To run the project with Nuxt SSR:

```bash
bun --cwd=./ui install
bun --bun --cwd=./ui run build
go build -tags netgo,ssr -ldflags=-s
./gofin
```

## Features

- [Nuxt 3](https://v3.nuxtjs.org/)
- [Fiber](https://github.com/gofiber/fiber)
- [Bun](https://bun.sh/)
- [Sonic](https://github.com/bytedance/sonic)

## License

[BSL-1.0](https://github.com/juls0730/gofin/blob/main/LICENSE)
