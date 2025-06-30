# Flutter Gopher

English | [中文](https://github.com/czg99/flutter_gopher/blob/main/README_zh.md)

Flutter Gopher is used to bridge Flutter with native Golang code. It allows for rapid creation of Flutter plugins based on native Golang and automatically generates the FFI binding code.

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
go install github.com/czg99/flutter_gopher/cmd/fgo@v0.1.2
```

## 📋 Usage

Flutter Gopher provides two main commands:

### 1. Create a new Flutter plugin project

```bash
fgo create -n <project_name> -o <output_directory> [--example]
```

**Parameters:**
- `-n, --name`: Plugin project name (required)
- `-o, --output`: Output directory for the generated plugin project (default: <project_name>)
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

| Go Type                                       | Dart Type     | Description                        |
| --------------------------------------------- | ------------- | ---------------------------------- |
| `bool`                                        | `bool`        | Boolean value                      |
| `int8`, `int16`, `int32`, `int64`, `int`      | `int`         | Signed integer                     |
| `uint8`, `uint16`, `uint32`, `uint64`, `uint` | `int`         | Unsigned integer                   |
| `float32`                                     | `double`      | 32-bit floating-point number       |
| `float64`                                     | `double`      | 64-bit floating-point number       |
| `string`                                      | `String`      | String                             |
| `struct`                                      | `Class`       | Struct/Class                       |
| `[]T`                                         | `List<T>`     | Slice/Array                        |
| `chan T`                                      | `FgChan<T>`   | One-way channel from Go to Dart    |
| `*T`                                          | `T?`          | Pointer converted to nullable type |
| `error`                                       | `String?`     | Error converted to nullable string |
| `func(...)`                                   | `Future<...>` | Asynchronous function support      |


### Type Conversion Rules

1. **Structs**: Go structs are converted into corresponding Dart classes, with field names transformed to camelCase.
2. **Slices**: Go slices are converted into Dart `List`s, preserving the element type.
3. **Pointers**: Go pointer types are converted into nullable types in Dart.
4. **Channels**: Go channel types are converted into Dart's `FgChan<T>`. You can use the `listen` method to receive data from the channel.
5. **Asynchronous Support**: Each Go function generates both synchronous and asynchronous Dart methods.
6. **Error Handling**: The `error` returned by Go functions is thrown as an exception on the Dart side and should be caught using `try-catch`.

## 🔄 Development Workflow

1. Use the `fgo create` command to create a new plugin project.
2. Implement the Go API in the `src/api` directory at the root of the plugin project.
3. Run the `fgo generate` command in the root directory of the plugin project to regenerate the FFI binding code.
4. Add the plugin as a dependency in your Flutter project.

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

func Add(a, b int) int {
    return a + b
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
  print('5 + 3 = $sum');
  
  // Use asynchronous API
  final sumAsync = await api.addGoAsync(5, 3);
  print('5 + 3 = $sumAsync');
}
```

## 📝 Contribution Guidelines

Pull Requests and Issues are welcome to help improve Flutter Gopher!

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.