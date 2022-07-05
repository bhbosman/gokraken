package connection

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsNetDialer"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/top"
	"github.com/bhbosman/goCommsStacks/websocket"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
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

func ProvideKrakenDialer(
	serviceIdentifier model.ServiceIdentifier,
	serviceDependentOn model.ServiceIdentifier,
	canDial goCommsNetDialer.ICanDial,
) fx.Option {
	var canDials []goCommsNetDialer.ICanDial
	if canDial != nil {
		canDials = append(canDials, canDial)
	}

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
				}) (messages.CreateAppCallback, error) {
					f := goCommsNetDialer.NewNetDialApp(
						"Kraken",
						serviceIdentifier,
						serviceDependentOn,
						"Kraken",
						false,
						nil,
						krakenUrl,
						goCommsDefinitions.WebSocketName,
						common.MaxConnectionsSetting(1),
						goCommsNetDialer.CanDial(canDials...),
						common.NewConnectionInstanceOptions(
							goCommsDefinitions.ProvideTransportFactoryForWebSocketName(
								top.ProvideTopStack(),
								websocket.ProvideWebsocketStacks(),
								bottom.Provide(),
							),
							PubSub.ProvidePubSubInstance("Application", params.PubSub),
							ProvideConnectionReactorFactory(),
						),
					)
					return f(params.NetAppFuncInParams), nil
				},
			}))
}

func ProvideConnectionReactorFactory() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						CancelCtx            context.Context
						CancelFunc           context.CancelFunc
						ConnectionCancelFunc model.ConnectionCancelFunc
						Logger               *zap.Logger
						ClientContext        interface{}    `name:"UserContext"`
						PubSub               *pubsub.PubSub `name:"Application"`
						GoFunctionCounter    GoFunctionCounter.IService
					},
				) (intf.IConnectionReactor, error) {
					return NewReactor(
							params.Logger,
							params.CancelCtx,
							params.CancelFunc,
							params.ConnectionCancelFunc,
							params.ClientContext,
							params.PubSub,
							params.GoFunctionCounter,
						),
						nil

				},
			},
		),
	)
}
