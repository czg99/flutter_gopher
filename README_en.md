# Flutter Gopher

[简体中文](https://github.com/czg99/flutter_gopher/blob/main/README.md) | English

Flutter Gopher is used to bridge code between Flutter, Golang, and Native platforms. Quickly create Flutter plugins based on Golang.

## ✨ Features

- 🔄 Create complete Flutter plugin project structure
- 🚀 Provide seamless interoperability between Flutter, Go, and Native
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
fgo create -n <project_name> -o <output_directory> [--example]
```

**Parameters:**
- `-n, --name`: Plugin project name (required)
- `-o, --output`: Output directory for the generated plugin project (default: <project_name>)
- `--example`: Generate an example Flutter application using the plugin

**Examples:**
```bash
fgo create -n my_ffi -o ./my_ffi
fgo create -n my_ffi --example
```

## 📁 Project Structure

The plugin project structure generated using the `create` command is as follows:

```
my_ffi/
├── android/        # Android platform code
├── ios/            # iOS platform code
├── linux/          # Linux platform code
├── macos/          # macOS platform code
├── windows/        # Windows platform code
├── lib/            # Dart code
├── src/            # Go code
└── example/        # Example Flutter application (if using the --example option)
```