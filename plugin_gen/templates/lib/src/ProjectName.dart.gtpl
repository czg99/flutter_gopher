import 'dart:math';
import 'dart:typed_data';
import 'bridge/bridge.dart';

class {{.LibClassName}} {
  (int, Uint8List) randomPacket() {
    final r = Random(DateTime.now().millisecondsSinceEpoch);
    final method = r.nextInt(99999);
    final count = r.nextInt(10);
    if (count < 3) {
      return (method, Uint8List(0));
    }
    final data = Uint8List(count);
    for (var i = 0; i < count; i++) {
      data[i] = r.nextInt(255);
    }
    return (method, data);
  }

  String bytesToHex(Uint8List data) =>
      '[${data.map((b) => b.toRadixString(16).padLeft(2, '0')).join(',')}]';

  {{.LibClassName}}() {
    FgBridge.setMethodHandle((method, data) {
      print('[{{.LibClassName}}] Dart Received: $method ${bytesToHex(data)}');
    });
  }

  static void init() => FgBridge.init();

  (int method, String send, String recv) callGoMethod() {
    final (method, send) = randomPacket();
    final recv = FgBridge.callGoMethod(method, data: send);
    return (method, bytesToHex(send), bytesToHex(recv));
  }

  (int method, String send, String recv) callPlatformMethod() {
    final (method, send) = randomPacket();
    final recv = FgBridge.callPlatformMethod(method, data: send);
    return (method, bytesToHex(send), bytesToHex(recv));
  }

  Future<(int method, String send, String recv)> callGoMethodAsync() async {
    final (method, send) = randomPacket();
    final recv = await FgBridge.callGoMethodAsync(method, data: send);
    return (method, bytesToHex(send), bytesToHex(recv));
  }

  Future<(int method, String send, String recv)>
      callPlatformMethodAsync() async {
    final (method, send) = randomPacket();
    final recv = await FgBridge.callPlatformMethodAsync(method, data: send);
    return (method, bytesToHex(send), bytesToHex(recv));
  }
}
