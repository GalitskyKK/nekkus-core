# Сборка nekkus-core

## Генерация proto

Нужны плагины (один раз):

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Генерация (из корня nekkus-core):

```bash
protoc -I. --go_out=. --go_opt=module=github.com/GalitskyKK/nekkus-core \
  --go-grpc_out=. --go-grpc_opt=module=github.com/GalitskyKK/nekkus-core \
  proto/module.proto proto/hub.proto
```

Или через Task: `task proto` (если установлен [Task](https://taskfile.dev/)).

## Сборка пакетов

- **pkg/discovery, pkg/config, pkg/server** — собираются без доп. требований.
- **pkg/desktop** (окно + tray) требует:
  - **CGO_ENABLED=1**
  - **C/C++ компилятор в PATH** (Windows: MinGW-w64/gcc, Linux: gcc/build-essential)
  - **Windows:** [WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) (часто уже есть)
  - **Linux:** GTK3 и WebKit2GTK (`libgtk-3-dev libwebkit2gtk-4.0-dev` или аналог)

Полная сборка с desktop:

```bash
CGO_ENABLED=1 go build ./...
```

Без desktop (только protocol, discovery, config, server):

```bash
go build ./pkg/discovery/... ./pkg/config/... ./pkg/server/... ./pkg/protocol/...
```
