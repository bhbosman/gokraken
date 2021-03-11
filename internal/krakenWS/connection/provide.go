package connection

import (
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
)

const FactoryName = "KrakenWSS"

func ProvideKrakenDialer(
	pubSub *pubsub.PubSub,
	canDial netDial.ICanDial) fx.Option {
	var canDials []netDial.ICanDial
	if canDial != nil {
		canDials = append(canDials, canDial)
	}

	const KrakenDialerConst = "KrakenDialer"
	cfr := NewFactory(KrakenDialerConst, pubSub)
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: netDial.NewNetDialApp(
					"Kraken",
					"wss://ws.kraken.com:443",
					impl.CreateWebSocketStack,
					KrakenDialerConst,
					cfr,
					netDial.MaxConnectionsSetting(1),
					netDial.CanDial(canDials...)),
			}))
}
