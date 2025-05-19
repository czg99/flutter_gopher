# Flutter Gopher

[English](https://github.com/czg99/flutter_gopher/blob/main/README.md) | [中文](https://github.com/czg99/flutter_gopher/blob/main/README_zh.md)

Flutter Gopher is a command-line tool for creating seamless integration between Flutter plugins and Golang backends. It automatically generates FFI binding code between them, making it very simple for Dart to call Go native code.

## ✨ Features

- 🔄 Create complete Flutter plugin project structure
- 🔌 Generate FFI binding code between Go and Dart
- 🚀 Provide seamless Flutter-Go interoperability
- 💻 Support multiple platforms (iOS, Android, Windows, macOS, Linux)

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

Flutter Gopher provides two main commands:

### 1. Create a new Flutter plugin project

```bash
fgo create -n <project_name> -o <output_directory> [--example]
```

**Parameters:**
- `-n, --name`: Plugin project name (required)
- `-o, --output`: Output directory for the generated plugin project (defaults to current directory)
- `--example`: Generate an example Flutter application using the plugin

**Examples:**
```bash
fgo create -n my_api -o ./my_api
fgo create -n payment_service --example
```

### 2. Generate Go and Dart FFI binding code

```bash
cd <fgo_created_project>
fgo generate
```

## 📁 Project Structure

The plugin project structure generated using the `create` command is as follows:

```
my_api/
├── android/        # Android platform-specific code
├── ios/            # iOS platform-specific code
├── linux/          # Linux platform-specific code
├── macos/          # macOS platform-specific code
├── windows/        # Windows platform-specific code
├── lib/            # Dart API code
│   └── my_api.dart
├── src/            # Go source code
│   ├── api/        # User-implemented API
│   └── api.go      # Generated Go FFI code
└── example/        # Example Flutter application (if --example option is used)
```

## 📊 Supported Data Types

Flutter Gopher supports converting the following data types between Go and Dart:

| Go Type | Dart Type | Description |
|---------|-----------|-------------|
| `bool` | `bool` | Boolean value |
| `int8`, `int16`, `int32` | `int` | Signed integers |
| `uint8`, `uint16`, `uint32` | `int` | Unsigned integers |
| `int64`, `uint64` | `int` | 64-bit integers |
| `int`, `uint` | `int` | Platform-dependent integers |
| `float32` | `double` | 32-bit floating point |
| `float64` | `double` | 64-bit floating point |
| `string` | `String` | String |
| `struct` | Go struct | Dart class |
| `[]T` | `List<T>` | Slice/Array |
| `[]*T` | `List<T?>` | Pointer slice/Nullable element list |
| `*T` | `T?` | Pointer converted to nullable type |
| `error` | `String?` | Error converted to nullable string |
| `func(...)` | `Future<...>` | Async function support |

### Type Conversion Rules

1. **Basic Types**: Go's basic numeric types are automatically mapped to Dart's `int` or `double`
2. **Structs**: Go structs generate corresponding Dart classes, with field names converted to camelCase
3. **Slices**: Go slices are converted to Dart `List`, preserving element types
4. **Error Handling**: Errors returned by Go functions are converted to nullable `String` in Dart
5. **Async Support**: All Go functions generate both synchronous and asynchronous (returning `Future`) versions of Dart methods

## 🔄 Development Workflow

1. Use the `create` command to create a new plugin project
2. Implement Go API in the `src/api` directory
3. Use the `generate` command to regenerate FFI binding code
4. Use the plugin in your Flutter application

## 🌟 Example

### Creating a Simple Calculator Plugin

#### 1. Create a plugin project:

```bash
fgo create -n calculator -o ./calculator --example
```

#### 2. Implement calculator API in the `src/api` directory:

```go
// src/api/calculator.go
package api

import "errors"

// Add returns the sum of two numbers
func Add(a, b int) int {
    return a + b
}

// Multiply returns the product of two numbers
func Multiply(a, b float64) float64 {
    return a * b
}

// CalculateWithPrecision calculates with specified precision
func CalculateWithPrecision(values []float64) (result float64, err error) {
    if len(values) == 0 {
        return 0, errors.New("empty array")
    }
    
    // Implement calculation logic
    return values[0], nil
}
```

#### 3. Generate FFI binding code:

```bash
cd calculator
fgo generate
```

#### 4. Use the plugin in a Flutter application:

```dart
import 'package:calculator/calculator.dart';

void main() async {
  // Use synchronous API
  final api = Calculator();
  final sum = api.add(5, 3);
  print('5 + 3 = $sum'); // Output: 5 + 3 = 8
  
  final product = api.multiply(2.5, 3.0);
  print('2.5 * 3.0 = $product'); // Output: 2.5 * 3.0 = 7.5
  
  // Use asynchronous API
  try {
    final result = await api.calculateWithPrecisionAsync([1.1, 2.2, 3.3]);
    print('Calculation result: $result');
  } catch (e) {
    print('Calculation error: $e');
  }
}
```

## 🔍 Advanced Usage

### Error Handling

All errors returned by Go functions are thrown as exceptions in Dart and can be caught using try-catch.

### Concurrency Handling

Go's concurrency features can be used through Dart's `Future` and `async/await` pattern.

## 📝 Contribution Guidelines

Pull Requests and Issues are welcome to help improve Flutter Gopher!

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.