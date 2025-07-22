package {{.PackageName}}

interface BridgeDelegate {
    fun methodHandle(method: String, data: ByteArray?): ByteArray?
}

class Bridge private constructor() {
    companion object {
        init {
            System.loadLibrary("{{.LibName}}")
        }

        @Volatile
        private var instance: Bridge? = null

        @JvmStatic
        fun getInstance(): Bridge {
            return instance ?: synchronized(this) {
                instance ?: Bridge().also { instance = it }
            }
        }
    }

    init {
        fgInit()
    }

    var delegate: BridgeDelegate? = null

    fun callGoMethod(method: String, data: ByteArray? = null): ByteArray? {
        val result = fgCallMethod(FgPacket(method, data))
        return result.data
    }

    private fun methodHandle(packet: FgPacket): FgPacket {
        if (delegate != null) {
            val result = delegate!!.methodHandle(packet.method, packet.data)
            return FgPacket(packet.method, result)
        }
        return FgPacket(packet.method)
    }

    private external fun fgInit()
    private external fun fgCallMethod(packet: FgPacket): FgPacket
}

private class FgPacket(
    @JvmField
    var method: String = "",
    @JvmField
    var data: ByteArray? = null
)