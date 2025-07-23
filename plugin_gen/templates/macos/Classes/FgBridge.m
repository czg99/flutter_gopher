#import "FgBridge.h"

#include "../../src/bridge/bridge.h"

@implementation FgBridge

FgPacket methodHandle(FgPacket packet) {
    return [[FgBridge sharedInstance] methodHandle:packet];
}

+ (instancetype)sharedInstance {
    static FgBridge *sharedInstance = nil;
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        sharedInstance = [[self alloc] init];
        enforce_binding();
        fg_init_method_handle(methodHandle);
    });
    return sharedInstance;
}

- (void*)nsdataToC:(NSData*)data {
    if (data == nil) {
        return nil;
    }
    NSUInteger dataLen = [data length];
    void* cData = malloc(dataLen);
    [data getBytes:cData length:dataLen];
    return cData;
}

- (FgPacket)methodHandle:(FgPacket)packet {
    FgPacket result = {};
    if (packet.method == nil) {
        free(packet.data);
        return result;
    }
    
    NSString* method = @"";
    NSData* data = nil;
    if (packet.method != nil) {
        method = [[NSString alloc] initWithBytes:packet.method length:packet.method_len encoding:NSUTF8StringEncoding];
    }
    
    if (packet.data != nil) {
        data = [[NSData alloc] initWithBytes:packet.data length:packet.data_len];
        free(packet.data);
    }
    
    NSData* handleData = nil;
    if (self.delegate != nil) handleData = [self.delegate methodHandle:method data:data];
    
    result.method = packet.method;
    result.method_len = packet.method_len;
    
    result.data = [self nsdataToC:handleData];
    if (result.data != nil) {
        result.data_len = (int)[handleData length];
    }
    
    return result;
}

- (NSData*)callGoMethod:(NSString*)method data:(NSData*)data {
    if (method == nil) {
        return nil;
    }
    
    NSUInteger method_len = [method lengthOfBytesUsingEncoding:NSUTF8StringEncoding];
    char* c_method = calloc(method_len + 1, method_len + 1);
    [method getCString:c_method maxLength:method_len+1 encoding:NSUTF8StringEncoding];
    
    FgPacket packet = {
        .method = c_method,
        .method_len = (int)method_len,
    };
    
    packet.data = [self nsdataToC:data];
    if (packet.data != nil) {
        packet.data_len = (int)[data length];
    }
    
    FgPacket c_result = fg_call_method(packet);
    free(c_result.method);
    
    if (c_result.data != nil) {
        NSData* result = [NSData dataWithBytes:c_result.data length:c_result.data_len];
        free(c_result.data);
        return result;
    }
    return nil;
}


@end
