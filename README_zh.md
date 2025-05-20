# Flutter Gopher

[English](https://github.com/czg99/flutter_gopher/blob/main/README.md) | [中文](https://github.com/czg99/flutter_gopher/blob/main/README_zh.md)

Flutter Gopher 用于桥接 Flutter 与 Golang 原生代码。快速创建基于 Golang 原生的 Flutter 插件，并自动生成 FFI 绑定代码。

## ✨ 功能特点

- 🔄 创建完整的 Flutter 插件项目结构
- 🔌 自动生成 Go 和 Dart 之间的 FFI 绑定代码
- 🚀 提供无缝的 Flutter-Go 互操作性
- 💻 支持多平台（iOS、Android、Windows、macOS、Linux）

## 🛠️ 安装

### 前置条件

- Go 1.23.0 或更高版本
- Flutter 3.22.0 或更高版本 
- Zig 0.14.0 或更高版本 (编译为 Windows 或 Linux 的库需要)

### 安装步骤

```bash
go install github.com/czg99/flutter_gopher/cmd/fgo@latest
```

## 📋 使用方法

Flutter Gopher 提供了两个主要命令：

### 1. 创建新的 Flutter 插件项目

```bash
fgo create -n <项目名称> -o <输出目录> [--example]
```

**参数说明：**
- `-n, --name`：插件项目名称（必需）
- `-o, --output`：生成的插件项目的输出目录（默认为当前目录）
- `--example`：生成使用该插件的示例 Flutter 应用

**示例：**
```bash
fgo create -n my_api -o ./my_api
fgo create -n payment_service --example
```

### 2. 生成 Go 和 Dart FFI 绑定代码

```bash
cd <fgo创建的项目>
fgo generate
```

## 📁 项目结构

使用 `create` 命令生成的插件项目结构如下：

```
my_api/
├── android/        # Android 平台特定代码
├── ios/            # iOS 平台特定代码
├── linux/          # Linux 平台特定代码
├── macos/          # macOS 平台特定代码
├── windows/        # Windows 平台特定代码
├── lib/            # Dart API 代码
│   └── my_api.dart
├── src/            # Go 源代码
│   ├── api/        # 用户实现的 API
│   └── api.go      # 生成的 Go FFI 代码
└── example/        # 示例 Flutter 应用（如果使用 --example 选项）
```

## 📊 支持的数据类型

Flutter Gopher 支持在 Go 和 Dart 之间转换以下数据类型：

| Go 类型 | Dart 类型 | 说明 |
|---------|-----------|------|
| `bool` | `bool` | 布尔值 |
| `int8`, `int16`, `int32` | `int` | 有符号整数 |
| `uint8`, `uint16`, `uint32` | `int` | 无符号整数 |
| `int64`, `uint64` | `int` | 64位整数 |
| `int`, `uint` | `int` | 平台相关整数 |
| `float32` | `double` | 32位浮点数 |
| `float64` | `double` | 64位浮点数 |
| `string` | `String` | 字符串 |
| `struct` | Go结构体 | Dart 类 |
| `[]T` | `List<T>` | 切片/数组 |
| `[]*T` | `List<T?>` | 指针切片/可空元素列表 |
| `*T` | `T?` | 指针转换为可空类型 |
| `error` | `String?` | 错误转换为可空字符串 |
| `func(...)` | `Future<...>` | 异步函数支持 |

### 类型转换规则

1. **基本类型**：Go 的基本数值类型会自动映射到 Dart 的 `int` 或 `double`
2. **结构体**：Go 结构体会生成对应的 Dart 类，字段名称会转换为驼峰式命名
3. **切片**：Go 切片会转换为 Dart 的 `List`，并保留元素类型
4. **错误处理**：Go 函数返回的 `error` 会转换为 Dart 的可空 `String`
5. **异步支持**：所有 Go 函数都会生成同步和异步（返回 `Future`）两个版本的 Dart 方法

## 🔄 开发流程

1. 使用 `create` 命令创建新的插件项目
2. 在 `src/api` 目录中实现 Go API
3. 使用 `generate` 命令重新生成 FFI 绑定代码
4. 在 Flutter 应用中使用该插件

## 🌟 示例

### 创建一个简单的计算器插件

#### 1. 创建插件项目：

```bash
fgo create -n calculator -o ./calculator --example
```

#### 2. 在 `src/api` 目录中实现计算器 API：

```go
// src/api/calculator.go
package api

import "errors"

// Add 返回两个数的和
func Add(a, b int) int {
    return a + b
}

// Multiply 返回两个数的乘积
func Multiply(a, b float64) float64 {
    return a * b
}

// CalculateWithPrecision 使用指定精度计算
func CalculateWithPrecision(values []float64) (result float64, err error) {
    if len(values) == 0 {
        return 0, errors.New("空数组")
    }
    
    // 实现计算逻辑
    return values[0], nil
}
```

#### 3. 生成 FFI 绑定代码：

```bash
cd calculator
fgo generate
```

#### 4. 在 Flutter 应用中使用该插件：

```dart
import 'package:calculator/calculator.dart';

void main() async {
  // 使用同步 API
  final api = Calculator();
  final sum = api.add(5, 3);
  print('5 + 3 = $sum'); // 输出: 5 + 3 = 8
  
  final product = api.multiply(2.5, 3.0);
  print('2.5 * 3.0 = $product'); // 输出: 2.5 * 3.0 = 7.5
  
  // 使用异步 API
  try {
    final result = await api.calculateWithPrecisionAsync([1.1, 2.2, 3.3]);
    print('计算结果: $result');
  } catch (e) {
    print('计算错误: $e');
  }
}
```

## 🔍 高级用法

### 错误处理

所有 Go 函数返回的错误都会在 Dart 端作为异常抛出，可以使用 try-catch 捕获。

### 并发处理

Go 的并发特性可以通过 Dart 的 `Future` 和 `async/await` 模式使用。

## 📝 贡献指南

欢迎提交 Pull Request 或创建 Issue 来帮助改进 Flutter Gopher！

## 📄 许可证

本项目采用 MIT 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。