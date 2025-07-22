#include <stdlib.h>
#include <stdint.h>

typedef struct {
	int64_t id;
	char* method;
	int method_len;
	void* data;
	int data_len;
} FgPacket;

typedef FgPacket (*FgMethodHandle)(FgPacket);
static inline FgPacket call_fg_method_handle(FgMethodHandle handle, FgPacket packet) {
	return handle(packet);
}

#ifdef _WIN32
    #define DLLEXPORT __declspec(dllexport)
#else
    #define DLLEXPORT
#endif

extern DLLEXPORT FgPacket fg_empty_packet(void);
extern DLLEXPORT FgPacket fg_packet_loop(void);
extern DLLEXPORT int64_t fg_next_port_id(void);
extern DLLEXPORT FgPacket fg_call_method(FgPacket packet);
extern DLLEXPORT void fg_call_method_async(FgPacket packet);
extern DLLEXPORT void fg_init_method_handle(FgMethodHandle handle);