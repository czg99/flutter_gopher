/// 所有 Flutter Gopher 异常的基类
/// All Flutter Gopher exceptions base class
abstract class FgException implements Exception {
  /// 异常发生的库名称
  /// The library name where the exception occurred
  final String libName;

  /// 异常发生的函数名称
  /// The function name where the exception occurred
  final String funcName;

  /// 异常的描述信息
  /// The description of the exception
  final String message;

  /// 创建一个包含指定上下文信息的 FgException
  /// Creates a FgException with the specified context information
  const FgException(this.libName, this.funcName, this.message);

  @override
  String toString();
}

/// 表示库中的可恢复错误
/// Represents a recoverable error in the library
final class FgError extends FgException {
  /// 创建一个包含指定上下文信息的可恢复错误
  /// Creates a recoverable error with the specified context information
  const FgError(super.libName, super.funcName, super.message);

  @override
  String toString() =>
      'FgError {lib: $libName, func: $funcName, message: $message}';
}

/// 表示库中的严重且不可恢复的错误
final class FgPanic extends FgException {
  /// 创建一个包含指定上下文信息的不可恢复错误
  /// Creates a unrecoverable error with the specified context information
  const FgPanic(super.libName, super.funcName, super.message);

  @override
  String toString() =>
      'FgPanic {lib: $libName, func: $funcName, message: $message}';
}
