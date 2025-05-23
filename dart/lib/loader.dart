import 'dart:io' as io;
import 'dart:ffi' as ffi;
import 'exception.dart';

/// Flutter Gopher 原生库加载器
class FgLoader {
  final String libName;
  late final ffi.DynamicLibrary _library;

  /// 创建新的 FgLoader 实例并加载指定的库
  ///
  /// [libName] 是库的基本名称，不包含平台特定的前缀或扩展名
  FgLoader(this.libName) {
    _loadLibrary();
  }

  /// 根据当前平台确定适当的库文件名
  String _libraryFileName() {
    // 检查是否在 Web 环境运行（WASM）
    if (identical(0, 0.0)) {
      return '$libName.wasm';
    }

    // 根据平台返回对应的库文件名
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

  /// 根据当前平台加载原生库
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

  /// 在加载的库中查找符号
  ///
  /// [symbolName] 是要查找的符号名称
  /// 返回指向该符号的指针
  /// 如果找不到符号，则抛出 [FgError]
  ffi.Pointer<T> lookup<T extends ffi.NativeType>(String symbolName) {
    try {
      return _library.lookup<T>(symbolName);
    } catch (e) {
      throw FgError('FgLoader', 'lookup',
          'Failed to lookup symbol: $symbolName. Error: $e');
    }
  }
}
