import 'dart:isolate';
import 'dart:async';

/// 循环回调函数类型定义，返回布尔值决定是否继续循环
typedef FgLoopCallback = bool Function(Map<String, dynamic>? params);

/// 存储单例循环实例的映射
final _onceLoopInstances = <String, FgLoop>{};

/// 提供在隔离区中循环执行代码的功能
class FgLoop {
  /// 隔离区实例
  Isolate? _isolate;

  /// 接收端口
  ReceivePort? _receivePort;

  /// 循环运行状态
  bool _isRunning = false;

  /// 单例模式的键
  String? _key;

  /// 重试相关配置
  static const int _maxRetries = 3;
  static const int _retryDelayBaseMs = 100;

  /// 获取循环运行状态
  bool get isRunning => _isRunning;

  /// 默认构造函数
  FgLoop();

  /// 创建并立即启动循环的工厂构造函数
  ///
  /// [callback] 循环执行的回调函数
  /// [params] 传递给回调函数的参数
  factory FgLoop.run(FgLoopCallback callback, [Map<String, dynamic>? params]) {
    final loop = FgLoop();
    loop.start(callback, params);
    return loop;
  }

  /// 创建单例循环的工厂构造函数
  ///
  /// [key] 单例的唯一标识符
  /// [callback] 循环执行的回调函数
  /// [params] 传递给回调函数的参数
  factory FgLoop.once(String key, FgLoopCallback callback,
      [Map<String, dynamic>? params]) {
    if (_onceLoopInstances.containsKey(key)) {
      return _onceLoopInstances[key]!;
    }
    final loop = FgLoop.run(callback, params);
    loop._key = key;
    _onceLoopInstances[key] = loop;
    return loop;
  }

  /// 停止循环隔离区
  void stop() {
    _cleanupResources();
    _removeFromOnceInstances();
    _isRunning = false;
  }

  /// 清理隔离区和端口资源
  void _cleanupResources() {
    if (_isolate != null) {
      _isolate!.kill(priority: Isolate.immediate);
      _isolate = null;
    }
    if (_receivePort != null) {
      _receivePort!.close();
      _receivePort = null;
    }
  }

  /// 从单例映射中移除实例
  void _removeFromOnceInstances() {
    if (_key != null) {
      _onceLoopInstances.remove(_key);
      _key = null;
    }
  }

  /// 启动循环隔离区
  ///
  /// [callback] 循环执行的回调函数
  /// [params] 传递给回调函数的参数
  Future<void> start(FgLoopCallback callback,
      [Map<String, dynamic>? params]) async {
    if (_isRunning) return;
    _isRunning = true;

    final receivePort = ReceivePort();
    _receivePort = receivePort;
    final message = <String, dynamic>{
      'callback': callback,
      'params': params,
      'sendPort': receivePort.sendPort,
    };

    _setupMessageListener(receivePort);

    try {
      _isolate = await Isolate.spawn(_loopIsolateEntryPoint, message);
    } catch (e) {
      stop();
      throw Exception('Error starting loop isolate: $e');
    }
  }

  /// 设置消息监听器
  void _setupMessageListener(ReceivePort receivePort) {
    receivePort.listen((message) {
      if (message is Exception || message is Error) {
        print('Error in loop isolate: $message');
      } else if (message == 'stop') {
        _isolate = null;
        stop();
      }
    });
  }

  /// 隔离区入口点函数
  static void _loopIsolateEntryPoint(Map<String, dynamic> message) {
    final callback =
        message['callback'] as bool Function(Map<String, dynamic>? params);
    final params = message['params'] as Map<String, dynamic>?;
    final sendPort = message['sendPort'] as SendPort;

    int retryCount = 0;

    while (true) {
      try {
        final isContinue = callback(params);
        if (!isContinue) {
          sendPort.send('stop');
          break;
        }
        // 重置重试计数
        retryCount = 0;
      } catch (e) {
        sendPort.send(e);
        retryCount++;
        if (retryCount >= _maxRetries) {
          // 达到最大重试次数后停止
          sendPort.send('stop');
          break;
        }
        // 延迟重试
        Future.delayed(Duration(milliseconds: _retryDelayBaseMs * retryCount));
      }
    }
  }
}
