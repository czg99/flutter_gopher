import 'dart:async';
import 'dart:ffi';
import 'dart:isolate';
import 'dart:ui';

/// 用于在Go到Dart之间传递数据的通道
///
/// FgChan是一个泛型类，用于从Go代码接收数据并转换为Dart对象。
/// 它使用Dart的[ReceivePort]和[IsolateNameServer]来实现跨语言通信。
class FgChan<T> {
  /// 接收端口，用于接收来自Go的消息
  ReceivePort? _receivePort;

  /// 指针到Dart对象的转换函数
  T Function(Pointer<Void>)? _pointerToDart;

  /// 端口ID，用于在IsolateNameServer中注册
  int? _portId;

  /// 获取当前通道的端口ID
  int get portId => _portId ?? 0;

  /// 获取Dart对象数据流
  Stream<T> get stream {
    if (_receivePort == null) {
      throw StateError('Channel is closed');
    }

    return _receivePort!.map((message) {
      if (_pointerToDart == null) {
        throw StateError(
            'PointerToDart converter is not set, please call setPointerToDart');
      }
      return _pointerToDart!(message as Pointer<Void>);
    });
  }

  /// 检查通道是否已关闭
  bool get isClosed => _receivePort == null;

  /// 创建一个通道
  FgChan() : _receivePort = ReceivePort();

  /// 监听Dart对象数据流
  StreamSubscription listen(
    void Function(T) onData, {
    Function? onError,
    void Function()? onDone,
    bool? cancelOnError,
  }) {
    return stream.listen(
      onData,
      onError: onError,
      onDone: onDone,
      cancelOnError: cancelOnError,
    );
  }

  /// 设置端口ID并注册到IsolateNameServer，不要手动调用
  void setPortId(int portId) {
    if (_receivePort == null) {
      throw StateError('Channel is closed');
    }

    // 移除旧的端口映射（如果存在）
    if (_portId != null) {
      IsolateNameServer.removePortNameMapping(_portId.toString());
    }

    // 注册新的端口映射
    final registered = IsolateNameServer.registerPortWithName(
        _receivePort!.sendPort, portId.toString());

    if (!registered) {
      throw StateError(
          'Failed to register port ID: $portId, it may already be in use');
    }

    _portId = portId;
  }

  /// 设置指针到Dart对象的转换函数，不要手动调用
  void setPointerToDart(T Function(Pointer<Void>) converter) {
    _pointerToDart = converter;
  }

  /// 关闭Dart通道并清理资源
  void close() {
    if (_portId != null) {
      IsolateNameServer.removePortNameMapping(_portId.toString());
      _portId = null;
    }
    _receivePort?.close();
    _receivePort = null;
  }
}
