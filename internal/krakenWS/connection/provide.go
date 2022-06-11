package connection

import (
	"github.com/bhbosman/goCommsStacks"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/gocomms/stacks/top"
	"github.com/bhbosman/gocomms/stacks/websocket"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
)

const FactoryName = "KrakenWSS"

func ProvideKrakenDialer(
	serviceIdentifier model.ServiceIdentifier,
	serviceDependentOn model.ServiceIdentifier,
	canDial netDial.ICanDial) fx.Option {
	var canDials []netDial.ICanDial
	if canDial != nil {
		canDials = append(canDials, canDial)
	}

	const KrakenDialerConst = "KrakenDialer"
	crfName := "KrakenDialer.CRF"
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: func(params struct {
					fx.In
					PubSub             *pubsub.PubSub `name:"Application"`
					NetAppFuncInParams common.NetAppFuncInParams
				}) messages.CreateAppCallback {
					f := netDial.NewNetDialApp(
						"Kraken",
						serviceIdentifier,
						serviceDependentOn,
						"Kraken",
						"wss://ws.kraken.com:443",
						common.WebSocketName,
						func() (intf.IConnectionReactorFactory, error) {
							cfr := NewFactory(crfName, params.PubSub)
							return cfr, nil
						},
						common.MaxConnectionsSetting(1),
						netDial.CanDial(canDials...),
						common.NewConnectionInstanceOptions(
							goCommsStacks.ProvideDefinedStackNames(),
							bottom.ProvideBottomStack(),
							top.ProvideTopStack(),
							websocket.ProvideWebsocketStacks()))
					return f(
						params.NetAppFuncInParams)
				},
			}))
}
