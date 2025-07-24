#include <stdlib.h>
#include <stdint.h>

#ifndef FG_BRIDGE_H
#define FG_BRIDGE_H

typedef struct {
	void* data;
	int len;
} FgData;

typedef struct {
	int64_t id;
	FgData method;
	FgData data;
} FgPacket;

typedef void (*FgMethodHandle)(FgPacket, FgPacket*);
static inline void call_fg_method_handle(FgMethodHandle handle, FgPacket packet, FgPacket* result) {
	handle(packet, result);
}

#ifdef _WIN32
    #define DLLEXPORT __declspec(dllexport)
#else
    #define DLLEXPORT
#endif

extern DLLEXPORT FgData fg_empty_data(void);
extern DLLEXPORT FgPacket fg_empty_packet(void);
extern DLLEXPORT FgPacket fg_packet_loop(void);
extern DLLEXPORT int64_t fg_next_port_id(void);
extern DLLEXPORT FgPacket fg_call_go_method(FgPacket packet);
extern DLLEXPORT void fg_call_go_method_async(FgPacket packet);
extern DLLEXPORT FgPacket fg_call_native_method(FgPacket packet);
extern DLLEXPORT void fg_call_native_method_async(FgPacket packet);
extern DLLEXPORT void fg_init_method_handle(FgMethodHandle handle);

extern DLLEXPORT void enforce_binding(void);

#endif