# Flutter Gopher Dart 库

[English](https://github.com/czg99/flutter_gopher/blob/main/dart/README.md) | [中文](https://github.com/czg99/flutter_gopher/blob/main/dart/README_zh.md)

Flutter Gopher Dart 库是 Flutter Gopher 工具生成的 Flutter 插件中的依赖组件，提供了 Flutter 与 Golang 原生代码的 FFI 绑定支持。

## 功能特点

- 提供动态库加载机制，支持多平台（iOS、Android、Windows、macOS、Linux）
- 实现异步调用支持，使 Go 函数可以在 Dart 中异步执行
- 实现通道监听机制，使 Dart 可以实时监听 Go 的通道数据
- 提供统一的错误处理机制
- 简化 FFI 绑定代码的使用

## 安装

此库由 Flutter Gopher 工具自动集成到生成的 Flutter 插件中，无需手动安装。

## Flutter Gopher 工具项目

Flutter Gopher 工具项目地址：[https://github.com/czg99/flutter_gopher](https://github.com/czg99/flutter_gopher)