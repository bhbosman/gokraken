package connection

import (
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
)

const FactoryName = "KrakenWSS"

func ProvideKrakenDialer(canDial netDial.ICanDial) fx.Option {
	var canDials []netDial.ICanDial
	if canDial != nil{
		canDials = append(canDials, canDial)
	}

	const KrakenDialerConst = "KrakenDialer"
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: impl.ConnectionReactorFactoryConst,
				Target: func(
					params struct {
						fx.In
						PubSub      *pubsub.PubSub `name:"Application"`
					}) (intf.IConnectionReactorFactory, error) {
					return NewFactory(KrakenDialerConst,  params.PubSub), nil
				},
			}),
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: netDial.NewNetDialApp(
					"Kraken",
					"wss://ws.kraken.com:443",
					impl.WebSocketName,
					KrakenDialerConst,
					netDial.MaxConnectionsSetting(1),
					netDial.CanDial(canDials...)),
			}))
}
