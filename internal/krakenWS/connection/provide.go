package connection

import (
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
)

const FactoryName = "KrakenWSS"

func ProvideKrakenDialer(
	canDial netDial.ICanDial) fx.Option {
	var canDials []netDial.ICanDial
	if canDial != nil {
		canDials = append(canDials, canDial)
	}

	const KrakenDialerConst = "KrakenDialer"
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: netDial.NewNetDialApp(
					"Kraken",
					"wss://ws.kraken.com:443",
					impl.WebSocketName,
					netDial.MaxConnectionsSetting(1),
					netDial.CanDial(canDials...),
					netDial.FxOption(
						fx.Provide(
							fx.Annotated{
								Target: func(pubSub *pubsub.PubSub) (intf.IConnectionReactorFactory, error) {
									cfr := NewFactory(pubSub)
									return cfr, nil
								},
							}))),
			}))
}
