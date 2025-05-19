import 'dart:io' as io;
import 'dart:ffi' as ffi;
import 'exception.dart';

/// A class responsible for loading native libraries for Flutter Gopher
class FgLoader {
  final String libName;
  late final ffi.DynamicLibrary _library;

  /// Creates a new FgLoader instance and loads the specified library
  ///
  /// [libName] is the base name of the library without platform-specific prefixes or extensions
  FgLoader(this.libName) {
    _loadLibrary();
  }

  /// Determines the appropriate library file name based on the current platform
  String _libraryFileName() {
    // Check if running on web (WASM)
    if (identical(0, 0.0)) {
      return '$libName.wasm';
    }

    // Platform-specific library naming
    if (io.Platform.isAndroid || io.Platform.isLinux) {
      return 'lib$libName.so';
    } else if (io.Platform.isWindows) {
      return '$libName.dll';
    } else if (io.Platform.isMacOS) {
      return 'lib$libName.dylib';
    } else {
      throw FgError('FgLoader', 'libraryFileName',
          'Unsupported platform: ${io.Platform.operatingSystem}');
    }
  }

  /// Loads the native library based on the current platform
  void _loadLibrary() {
    try {
      if (io.Platform.isIOS) {
        _library = ffi.DynamicLibrary.executable();
      } else {
        final libFileName = _libraryFileName();
        _library = ffi.DynamicLibrary.open(libFileName);
      }
    } catch (e) {
      throw FgError(
          'FgLoader', 'loadLibrary', 'Failed to load native library: $e');
    }
  }

  /// Looks up a symbol in the loaded library
  ///
  /// [symbolName] is the name of the symbol to look up
  /// Returns a pointer to the symbol
  /// Throws a [FgError] if the symbol cannot be found
  ffi.Pointer<T> lookup<T extends ffi.NativeType>(String symbolName) {
    try {
      return _library.lookup<T>(symbolName);
    } catch (e) {
      throw FgError('FgLoader', 'lookup',
          'Failed to lookup symbol: $symbolName. Error: $e');
    }
  }
}
