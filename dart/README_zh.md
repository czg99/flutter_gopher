# Flutter Gopher Dart 库

[English](https://github.com/czg99/flutter_gopher/blob/main/dart/README.md) | [中文](https://github.com/czg99/flutter_gopher/blob/main/dart/README_zh.md)

Flutter Gopher Dart 库是 Flutter Gopher 工具生成的 Flutter 插件中的依赖组件，提供了 Flutter 与 Golang 原生代码的 FFI 绑定支持。

## 功能特点

- 提供动态库加载机制，支持多平台（iOS、Android、Windows、macOS、Linux）
- 实现异步调用支持，使 Go 函数可以在 Dart 中异步执行
- 提供统一的错误处理机制
- 简化 FFI 绑定代码的使用

## 安装

此库通常由 Flutter Gopher 工具自动集成到生成的 Flutter 插件中，无需手动安装。如果需要单独使用，可以在 `pubspec.yaml` 中添加依赖：

```yaml
dependencies:
  flutter_gopher: ^0.0.1
```

## Flutter Gopher 项目

Flutter Gopher 工具项目地址：[https://github.com/czg99/flutter_gopher](https://github.com/czg99/flutter_gopher)