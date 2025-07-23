#import <Foundation/Foundation.h>

@protocol FgBridgeDelegate <NSObject>
- (NSData*)methodHandle:(NSString*)method data:(NSData*)data;
@end


@interface FgBridge : NSObject

@property (nonatomic, weak) id<FgBridgeDelegate> delegate;

+ (instancetype)sharedInstance;

- (NSData*)callGoMethod:(NSString*)method data:(NSData*)data;

@end
