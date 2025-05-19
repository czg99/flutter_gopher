/// Base class for all Flutter Gopher exceptions.
///
/// This abstract class serves as the foundation for the exception hierarchy
/// in the Flutter Gopher library. All specific exceptions should extend this class.
abstract class FgException implements Exception {
  /// The name of the library where the exception occurred
  final String libName;

  /// The name of the function where the exception occurred
  final String funcName;

  /// Descriptive message about the exception
  final String message;

  /// Creates a new FgException with the specified context information
  const FgException(this.libName, this.funcName, this.message);

  @override
  String toString();
}

/// Represents a recoverable error in library
final class FgError extends FgException {
  /// Creates a new recoverable error with the specified context information
  const FgError(super.libName, super.funcName, super.message);

  @override
  String toString() =>
      'FgError {lib: $libName, func: $funcName, message: $message}';
}

/// Represents a critical, non-recoverable error in library
final class FgPanic extends FgException {
  /// Creates a new non-recoverable error with the specified context information
  const FgPanic(super.libName, super.funcName, super.message);

  @override
  String toString() =>
      'FgPanic {lib: $libName, func: $funcName, message: $message}';
}
