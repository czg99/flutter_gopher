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

extern DLLEXPORT FgData fg_empty_data_{{.Timestamp}}(void);
extern DLLEXPORT FgRequest fg_empty_request_{{.Timestamp}}(void);
extern DLLEXPORT FgResponse fg_empty_response_{{.Timestamp}}(void);

extern DLLEXPORT void fg_init_dart_api_{{.Timestamp}}(void* api, int64_t port);
extern DLLEXPORT void fg_init_method_handle_{{.Timestamp}}(FgMethodHandle handle);

extern DLLEXPORT void fg_call_dart_method_{{.Timestamp}}(FgRequest request);
extern DLLEXPORT FgResponse fg_call_go_method_{{.Timestamp}}(FgRequest request);
extern DLLEXPORT void fg_call_go_method_async_{{.Timestamp}}(int64_t port, FgRequest request);
extern DLLEXPORT FgResponse fg_call_native_method_{{.Timestamp}}(FgRequest request);
extern DLLEXPORT void fg_call_native_method_async_{{.Timestamp}}(int64_t port, FgRequest request);

extern DLLEXPORT void enforce_binding_{{.Timestamp}}(void);

#endif