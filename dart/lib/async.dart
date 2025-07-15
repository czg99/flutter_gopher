import 'dart:async';
import 'dart:ffi';
import 'dart:isolate';
import 'dart:ui';

/// 提供 Flutter 和 Go 之间的异步通信功能
///
/// 包含在 Dart 隔离区执行代码和与 Go 代码通信的方法
class FgAsync {
  /// 监听接收端口并处理结果
  static void _listen<T>(ReceivePort receivePort, Completer<T> completer,
      [int? portId]) {
    receivePort.listen((message) {
      if (portId != null) {
        IsolateNameServer.removePortNameMapping(portId.toString());
      }
      receivePort.close();
      if (message is Exception || message is Error) {
        completer.completeError(message);
      } else {
        completer.complete(message as T);
      }
    });
  }

  /// 创建并设置接收端口和完成器
  static (ReceivePort, Completer<T>) _setupCommunication<T>() {
    final receivePort = ReceivePort();
    final completer = Completer<T>();
    return (receivePort, completer);
  }

  /// 在 Dart 隔离区中执行计算
  static Future<R> dart<P, R>(
    R Function(P params) computation,
    P params,
  ) async {
    final (receivePort, completer) = _setupCommunication<R>();
    final message = {
      'sendPort': receivePort.sendPort,
      'computation': computation,
      'params': params,
    };

    _listen(receivePort, completer);
    try {
      await Isolate.spawn(_isolateEntryPoint, message);
    } catch (e) {
      receivePort.close();
      throw Exception('Error starting isolate: ${e.toString()}');
    }

    return await completer.future;
  }

  /// 隔离区入口点函数
  static void _isolateEntryPoint(Map<String, dynamic> message) {
    final sendPort = message['sendPort'] as SendPort;
    final computation = message['computation'] as Function;
    final params = message['params'];

    try {
      final result = computation(params);
      sendPort.send(result);
    } catch (e) {
      sendPort.send(e);
    }
  }

  /// 执行 Go 代码并等待结果
  static Future<R> go<P, R>(
      void Function(P params, int portId) computation, P params,
      [int? customPortId]) async {
    final (receivePort, completer) = _setupCommunication<R>();
    final sendPort = receivePort.sendPort;
    final portId = customPortId ?? sendPort.nativePort;
    final success =
        IsolateNameServer.registerPortWithName(sendPort, portId.toString());
    if (!success) {
      receivePort.close();
      throw Exception('Failed to register port id: $portId');
    }

    _listen(receivePort, completer, portId);
    try {
      computation(params, portId);
    } catch (e) {
      receivePort.close();
      IsolateNameServer.removePortNameMapping(portId.toString());
      throw Exception('Error executing computation: ${e.toString()}');
    }

    return await completer.future;
  }

  /// 发送 Go 函数结果到 Dart
  static bool sendGoResult(int portId, dynamic result) {
    if (portId == 0) return false;
    final sendPort = IsolateNameServer.lookupPortByName(portId.toString());
    if (sendPort == null) return false;
    sendPort.send(result);
    return true;
  }
}
