#import "{{.PluginClassName}}.h"

@implementation {{.PluginClassName}}

+ (void)registerWithRegistrar:(id<FlutterPluginRegistrar>)registrar {
    static {{.PluginClassName}} *instance = nil;
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        instance = [[self alloc] init];
    });
    [FgBridge setDelegate:instance];
}

// Implementation of FgBridgeDelegate
- (NSData *)methodHandle:(int)method data:(NSData *)data error:(NSError**)error {
    NSString *hexString = @"null";
    
    if (data != nil) {
        NSMutableString *mutableHexString = [NSMutableString string];
        const unsigned char *bytes = [data bytes];
        for (NSInteger i = 0; i < [data length]; i++) {
            if (i > 0) {
                [mutableHexString appendString:@" "];
            }
            [mutableHexString appendFormat:@"%02x", bytes[i]];
        }
        hexString = mutableHexString;
    }
    
    NSLog(@"[{{.LibClassName}}] Platform Received: %d %@", method, hexString);
    return data;
}

@end
