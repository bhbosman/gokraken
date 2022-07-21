package connection

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsNetDialer"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/topStack"
	"github.com/bhbosman/goCommsStacks/websocket"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/fx/PubSub"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net/url"
)

func ProvideKrakenDialer() fx.Option {
	var canDials []goCommsNetDialer.ICanDial

	const KrakenDialerConst = "KrakenDialer"

	krakenUrl, err := url.Parse("wss://ws.kraken.com:443")
	if err != nil {
		fx.Error(err)
	}
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: func(params struct {
					fx.In
					PubSub             *pubsub.PubSub `name:"Application"`
					NetAppFuncInParams common.NetAppFuncInParams
					GoFunctionCounter  GoFunctionCounter.IService
					FmdService         fullMarketDataManagerService.IFmdManagerService
				}) (messages.CreateAppCallback, error) {
					f := goCommsNetDialer.NewSingleNetDialApp(
						"Kraken",
						"Kraken",
						common.MoreOptions(
							goCommsDefinitions.ProvideUrl("ConnectionUrl", krakenUrl),
							goCommsDefinitions.ProvideUrl("ProxyUrl", nil),
							goCommsDefinitions.ProvideBool("UseProxy", false),
						),
						common.MaxConnectionsSetting(1),
						goCommsNetDialer.CanDial(canDials...),
						common.NewConnectionInstanceOptions(
							goCommsDefinitions.ProvideTransportFactoryForWebSocketName(
								topStack.ProvideTopStack(),
								websocket.ProvideWebsocketStacks(),
								bottom.Provide(),
							),
							PubSub.ProvidePubSubInstance("Application", params.PubSub),
							fx.Provide(
								fx.Annotated{
									Target: func() fullMarketDataManagerService.IFmdManagerService {
										return params.FmdService
									},
								}),
							ProvideConnectionReactorFactory(),
						),
					)
					return f(params.NetAppFuncInParams), nil
				},
			},
		),
	)
}

func ProvideConnectionReactorFactory() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						CancelCtx              context.Context
						CancelFunc             context.CancelFunc
						ConnectionCancelFunc   model.ConnectionCancelFunc
						Logger                 *zap.Logger
						PubSub                 *pubsub.PubSub `name:"Application"`
						GoFunctionCounter      GoFunctionCounter.IService
						UniqueReferenceService interfaces.IUniqueReferenceService
						FmdService             fullMarketDataManagerService.IFmdManagerService
					},
				) (intf.IConnectionReactor, error) {
					return NewReactor(
							params.Logger,
							params.CancelCtx,
							params.CancelFunc,
							params.ConnectionCancelFunc,
							params.PubSub,
							params.GoFunctionCounter,
							params.UniqueReferenceService,
							params.FmdService,
						),
						nil

				},
			},
		),
	)
}
