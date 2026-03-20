# Flutter Gopher

[简体中文](https://github.com/czg99/flutter_gopher/blob/main/README.md) | English

Flutter Gopher is used to quickly create Golang-based Flutter plugins and generates convenient Flutter, Golang, and Platform bridge code.

## ✨ Features

- 🔄 Create complete Flutter plugin project structure
- 🚀 Provide seamless interoperability between Flutter, Go, and Platform
- 💻 Support multiple platforms (iOS, Android, HarmonyOS, Windows, macOS, Linux)

## 🛠️ Installation

### Prerequisites

- Go 1.23.0 or higher
- Flutter 3.22.0 or higher
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
├── darwin/           # iOS and macOS platform code
├── ohos/             # HarmonyOS platform code
├── linux/            # Linux platform code
├── windows/          # Windows platform code
├── lib/              # Dart code
├── gosrc/            # Go code
├── protos/           # Protobuf definition files
├── scripts/          # Script files (compile dynamic libraries, generate protobuf code)
└── example/          # Example Flutter application (if using the --example option)
```

## 🔧 Configuration

### Configure Android ProGuard Rules in the Main Project

1. Add the following rules to the `android/app/proguard-rules.pro` file:

```
-keep class * extends com.google.protobuf.** {*;}
```

2. Modify the `android.buildTypes` section in the `android/app/build.gradle` file as follows:

```
    buildTypes {
        release {
            signingConfig = signingConfigs.release
            proguardFiles getDefaultProguardFile('proguard-android.txt'), 'proguard-rules.pro'
        }
    }
```

