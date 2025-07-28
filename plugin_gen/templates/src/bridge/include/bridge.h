#include <stdlib.h>
#include <stdint.h>

#ifndef FG_BRIDGE_H
#define FG_BRIDGE_H

typedef struct {
	void* data;
	int size;
} FgData;

typedef struct {
	FgData method;
	FgData data;
} FgRequest;

typedef struct {
	FgData data;
} FgResponse;


typedef void (*FgMethodHandle)(FgRequest, FgResponse*);
static inline void call_fg_method_handle(FgMethodHandle handle, FgRequest request, FgResponse* response) {
	handle(request, response);
}

#ifdef _WIN32
    #define DLLEXPORT __declspec(dllexport)
#else
    #define DLLEXPORT __attribute__((visibility("default")))
#endif

extern DLLEXPORT FgData fg_empty_data(void);
extern DLLEXPORT FgRequest fg_empty_request(void);
extern DLLEXPORT FgResponse fg_empty_response(void);

extern DLLEXPORT void fg_init_dart_api(void* api, int64_t port);
extern DLLEXPORT void fg_init_method_handle(FgMethodHandle handle);

extern DLLEXPORT void fg_call_dart_method(FgRequest request);
extern DLLEXPORT FgResponse fg_call_go_method(FgRequest request);
extern DLLEXPORT void fg_call_go_method_async(int64_t port, FgRequest request);
extern DLLEXPORT FgResponse fg_call_native_method(FgRequest request);
extern DLLEXPORT void fg_call_native_method_async(int64_t port, FgRequest request);

extern DLLEXPORT void enforce_binding(void);

#endif