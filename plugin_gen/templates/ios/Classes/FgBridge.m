#import "FgBridge.h"

#include "../../src/bridge/include/bridge.h"

@implementation FgBridge

void methodHandle(FgRequest request, FgResponse* response) {
    [[FgBridge sharedInstance] methodHandle:request response:response];
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

- (FgData)mapFgDataFromNSData:(NSData*)from {
    FgData result = {};
    if (from != nil) {
        NSUInteger dataLen = [from length];
        void* data = malloc(dataLen);
        [from getBytes:data length:dataLen];
        result.data = data;
        result.size = (int)dataLen;
    }
    return result;
}

- (NSData*)mapFgDataToNSData:(FgData)from {
    if (from.data == nil) return nil;
    NSData* result = [[NSData alloc] initWithBytes:from.data length:from.size];
    [self freeFgData:&from];
    return result;
}

- (FgData)mapFgDataFromNSString:(NSString*)from {
    NSData* data = from != nil ? [from dataUsingEncoding:NSUTF8StringEncoding] : nil;
    return [self mapFgDataFromNSData:data];
}

- (NSString*)mapFgDataToNSString:(FgData)from {
    if (from.data == nil) return @"";
    NSString* result = [[NSString alloc] initWithBytes:from.data length:from.size encoding:NSUTF8StringEncoding];
    [self freeFgData:&from];
    return result;
}

- (void)freeFgData:(FgData*)value {
    if (value->data != nil) {
        free(value->data);
        value->data = nil;
        value->size = 0;
    }
}

- (void)methodHandle:(FgRequest)request response:(FgResponse*)response {
    NSString* method = [self mapFgDataToNSString:request.method];
    NSData* data = [self mapFgDataToNSData:request.data];
    
    NSData* handleData = nil;
    if (self.delegate != nil) handleData = [self.delegate methodHandle:method data:data];
    
    response->data = [self mapFgDataFromNSData:handleData];
}

- (NSData*)callGoMethod:(NSString*)method data:(NSData*)data {
    if (method == nil) {
        return nil;
    }
    
    FgRequest request = {
        .method = [self mapFgDataFromNSString:method],
        .data = [self mapFgDataFromNSData:data],
    };
    
    FgResponse response = fg_call_go_method(request);
    
    NSData* result = [self mapFgDataToNSData:response.data];
    return result;
}

- (void)callDartMethod:(NSString*)method data:(NSData*)data {
    if (method == nil) {
        return;
    }
    
    FgRequest request = {
        .method = [self mapFgDataFromNSString:method],
        .data = [self mapFgDataFromNSData:data],
    };
    
    fg_call_go_method(request);
}

@end
