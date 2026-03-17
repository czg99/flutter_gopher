package {{.PackageName}}

import io.flutter.embedding.engine.plugins.FlutterPlugin

class {{.PluginClassName}} : FlutterPlugin, FgBridgeDelegate {
    override fun onAttachedToEngine(flutterPluginBinding: FlutterPlugin.FlutterPluginBinding) {
        FgBridge.delegate = this
    }

    override fun onDetachedFromEngine(binding: FlutterPlugin.FlutterPluginBinding) {
    }

    // Implementation of FgBridgeDelegate
    override fun methodHandle(method: Int, data: ByteArray?): Result<ByteArray?> {
        println("[{{.LibClassName}}] Platform Received: $method ${data?.joinToString(separator = " ") { "%02x".format(it) } ?: "null"}")
        return Result.success(data)
    }
}
