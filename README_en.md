# Flutter Gopher

[简体中文](https://github.com/czg99/flutter_gopher/blob/main/README.md) | English

Flutter Gopher is used to quickly create Golang-based Flutter plugins and generates convenient Flutter, Golang, and Platform bridge code.

## ✨ Features

- 🔄 Create complete Flutter plugin project structure
- 🚀 Provide seamless interoperability between Flutter, Go, and Platform
- 💻 Support multiple platforms (iOS, Android, Windows, macOS, Linux)

## 🛠️ Installation

### Prerequisites

- Go 1.23.0 or higher
- Flutter 3.10.0 or higher
- Zig 0.14.0 or higher (required for compiling libraries for Windows or Linux)

### Installation Steps

```bash
go install github.com/czg99/flutter_gopher/cmd/fgo@latest
```

## 📋 Usage

### Create a New Flutter Plugin Project

```bash
fgo create <project_name> [--example]
```

**Parameters:**
- `<project_name>`: Plugin project name (required)
- `--example`: Generate an example Flutter application using the plugin

**Examples:**
```bash
fgo create my_ffi
fgo create my_ffi --example
```

## 📁 Project Structure

The plugin project structure generated using the `create` command is as follows:

```
my_ffi/
├── android/          # Android platform code
├── ios/              # iOS platform code
├── linux/            # Linux platform code
├── macos/            # macOS platform code
├── windows/          # Windows platform code
├── lib/              # Dart code
├── src/              # Go code
├── protos/           # Protobuf code
│   ├── proto/        # Protobuf definition files
│   ├── gen_protos.sh # Script to generate Protobuf code
└── example/          # Example Flutter application (if using the --example option)
```