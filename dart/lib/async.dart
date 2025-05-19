import 'dart:async';
import 'dart:isolate';

/// Executes a computation function in a separate isolate.
///
/// [computation] is the function to be executed in the isolate.
/// [params] contains the parameters to be passed to the computation function.
/// Returns a [Future] that completes with the result of the computation.
Future<T> fgAsync<T>(
  T Function(Map<String, dynamic> params) computation,
  Map<String, dynamic> params,
) async {
  final completer = Completer<T>();
  final receivePort = ReceivePort();

  // Create message for the isolate
  final message = {
    'sendPort': receivePort.sendPort,
    'computation': computation,
    'params': params,
  };

  // Setup listener for isolate response
  receivePort.listen((message) {
    receivePort.close();
    if (message is Exception) {
      completer.completeError(message);
    } else {
      completer.complete(message as T);
    }
  });

  // Execute computation in isolate
  await Isolate.spawn(_isolateEntryPoint, message);

  return await completer.future;
}

/// Entry point function for the isolate.
///
/// Executes the computation function with the provided parameters
/// and sends the result back through the send port.
void _isolateEntryPoint(Map<String, dynamic> message) {
  final sendPort = message['sendPort'] as SendPort;
  final computation = message['computation'] as Function;
  final params = message['params'] as Map<String, dynamic>;

  try {
    final result = computation(params);
    sendPort.send(result);
  } catch (e) {
    sendPort.send(e);
  }
}
