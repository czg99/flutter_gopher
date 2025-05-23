/// 提供统一的异常处理机制，包括可恢复的错误和不可恢复的严重错误
/// 所有 Flutter Gopher 异常的基类
abstract class FgException implements Exception {
  /// 异常发生的库名称
  final String libName;

  /// 异常发生的函数名称
  final String funcName;

  /// 异常的描述信息
  final String message;

  /// 创建一个包含指定上下文信息的 FgException
  const FgException(this.libName, this.funcName, this.message);

  @override
  String toString();
}

/// 表示库中的可恢复错误
final class FgError extends FgException {
  /// 创建一个包含指定上下文信息的可恢复错误
  const FgError(super.libName, super.funcName, super.message);

  @override
  String toString() =>
      'FgError {lib: $libName, func: $funcName, message: $message}';
}

/// 表示库中的严重且不可恢复的错误
final class FgPanic extends FgException {
  /// 创建一个包含指定上下文信息的不可恢复错误
  const FgPanic(super.libName, super.funcName, super.message);

  @override
  String toString() =>
      'FgPanic {lib: $libName, func: $funcName, message: $message}';
}
