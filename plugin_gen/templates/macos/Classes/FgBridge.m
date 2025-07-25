#import "FgBridge.h"

#include "../../src/bridge/bridge.h"

@implementation FgBridge

void methodHandle(FgPacket packet, FgPacket* result) {
    [[FgBridge sharedInstance] methodHandle:packet result:result];
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
    return [[NSData alloc] initWithBytes:from.data length:from.size];
}

- (FgData)mapFgDataFromNSString:(NSString*)from {
    NSData* data = from != nil ? [from dataUsingEncoding:NSUTF8StringEncoding] : nil;
    return [self mapFgDataFromNSData:data];
}

- (NSString*)mapFgDataToNSString:(FgData)from {
    if (from.data == nil) return @"";
    return [[NSString alloc] initWithBytes:from.data length:from.size encoding:NSUTF8StringEncoding];
}

- (void)freeFgData:(FgData*)value {
    if (value->data != nil) {
        free(value->data);
        value->data = nil;
        value->size = 0;
    }
}

- (void)methodHandle:(FgPacket)packet result:(FgPacket*)result {
    NSString* method = [self mapFgDataToNSString:packet.method];
    NSData* data = [self mapFgDataToNSData:packet.data];
    [self freeFgData:&packet.data];
    
    NSData* handleData = nil;
    if (self.delegate != nil) handleData = [self.delegate methodHandle:method data:data];
    
    result->method = packet.method;
    result->data = [self mapFgDataFromNSData:handleData];
}

- (NSData*)callGoMethod:(NSString*)method data:(NSData*)data {
    if (method == nil) {
        return nil;
    }
    
    FgPacket packet = {
        .method = [self mapFgDataFromNSString:method],
        .data = [self mapFgDataFromNSData:data],
    };
    
    FgPacket cResult = fg_call_go_method(packet);
    
    NSData* result = [self mapFgDataToNSData:cResult.data];
    [self freeFgData:&cResult.method];
    [self freeFgData:&cResult.data];
    return result;
}

@end
