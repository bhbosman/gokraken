package connection

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsNetDialer"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/top"
	"github.com/bhbosman/goCommsStacks/websocket"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"net/url"
)

const FactoryName = "KrakenWSS"

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
	crfName := "KrakenDialer.CRF"

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
							fx.Provide(
								fx.Annotated{
									Target: func() (intf.IConnectionReactorFactory, error) {
										return NewFactory(
											crfName,
											params.PubSub,
											params.GoFunctionCounter,
										)
									},
								},
							),
						),
					)
					return f(params.NetAppFuncInParams), nil
				},
			}))
}
