# Flutter Gopher

简体中文 | [English](https://github.com/czg99/flutter_gopher/blob/main/README_en.md)

Flutter Gopher 用于快速创建基于 Golang 的 Flutter 插件，并生成了便利的 Flutter、Golang、Native 桥接代码。

## ✨ 功能特点

- 🔄 创建完整的 Flutter 插件项目结构
- 🚀 提供无缝的 Flutter、Go、Native 互操作性
- 💻 支持多平台（iOS、Android、Windows、macOS、Linux）

## 🛠️ 安装

### 前置条件

- Go 1.23.0 或更高版本
- Flutter 3.10.0 或更高版本 
- Zig 0.14.0 或更高版本 (编译为 Windows 或 Linux 的库需要)

### 安装步骤

```bash
go install github.com/czg99/flutter_gopher/cmd/fgo@latest
```

## 📋 使用方法

### 创建新的 Flutter 插件项目

```bash
fgo create <project_name> [--example]
```

**参数说明：**
- `<project_name>`：插件项目名称（必需）
- `--example`：生成使用该插件的示例 Flutter 应用

**示例：**
```bash
fgo create my_ffi
fgo create my_ffi --example
```

## 📁 项目结构

使用 `create` 命令生成的插件项目结构如下：

```
my_ffi/
├── android/          # Android 平台代码
├── ios/              # iOS 平台代码
├── linux/            # Linux 平台代码
├── macos/            # macOS 平台代码
├── windows/          # Windows 平台代码
├── lib/              # Dart 代码
├── src/              # Go 代码
├── protos/           # Protobuf 代码
│   ├── proto/        # Protobuf 定义文件
│   ├── gen_protos.sh # 生成 Protobuf 代码的脚本
└── example/          # 示例 Flutter 应用（如果使用 --example 选项）
```