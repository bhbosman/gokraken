package listener

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsNetListener"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/bvisMessageBreaker"
	"github.com/bhbosman/goCommsStacks/messageCompressor"
	"github.com/bhbosman/goCommsStacks/messageNumber"
	"github.com/bhbosman/goCommsStacks/pingPong"
	"github.com/bhbosman/goCommsStacks/protoBuf"
	"github.com/bhbosman/goCommsStacks/topStack"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/fx/PubSub"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocomms/common"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"net/url"
)

func CompressedListener(
	maxConnections int,
	urlAsText string,
) fx.Option {
	const CompressedListenerConnection = "Kraken Client Connection Manager(Compressed)"
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: func(
					params struct {
						fx.In
						PubSub               *pubsub.PubSub `name:"Application"`
						GoFunctionCounter    GoFunctionCounter.IService
						NetAppFuncInParams   common.NetAppFuncInParams
						FullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper
						FmdService           fullMarketDataManagerService.IFmdManagerService
					},
				) (messages.CreateAppCallback, error) {
					compressedUrl, err := url.Parse(urlAsText)
					if err != nil {
						return messages.CreateAppCallback{}, err
					}
					f := goCommsNetListener.NewNetListenApp(
						CompressedListenerConnection,
						CompressedListenerConnection,
						false,
						nil,
						compressedUrl,
						common.MaxConnectionsSetting(maxConnections),
						common.NewConnectionInstanceOptions(
							goCommsDefinitions.ProvideTransportFactoryForCompressedName(
								topStack.Provide(),
								pingPong.Provide(),
								protoBuf.Provide(),
								messageCompressor.Provide(),
								messageNumber.Provide(),
								bvisMessageBreaker.Provide(),
								bottom.Provide(),
							),
							ProvideConnectionReactor(),
							PubSub.ProvidePubSubInstance("Application", params.PubSub),
							fx.Provide(
								fx.Annotated{
									Target: func() fullMarketDataHelper.IFullMarketDataHelper {
										return params.FullMarketDataHelper
									},
								},
							),
							fx.Provide(
								fx.Annotated{
									Target: func() fullMarketDataManagerService.IFmdManagerService {
										return params.FmdService
									},
								},
							),
						),
					)
					return f(params.NetAppFuncInParams), nil
				},
			},
		),
	)
}
