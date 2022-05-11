package connection

import (
	"github.com/bhbosman/gocommon/messages"
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
	crfName := "KrakenDialer.CRF"
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: func(params struct {
					fx.In
					PubSub             *pubsub.PubSub `name:"Application"`
					NetAppFuncInParams impl.NetAppFuncInParams
				}) messages.CreateAppCallback {
					fxOptions := fx.Options(
						fx.Provide(fx.Annotated{Name: "Application", Target: func() *pubsub.PubSub { return params.PubSub }}),
						fx.Provide(
							fx.Annotated{
								Target: func(params struct {
									fx.In
									PubSub *pubsub.PubSub `name:"Application"`
								}) intf.ConnectionReactorFactoryCallback {
									return func() (intf.IConnectionReactorFactory, error) {
										cfr := NewFactory(crfName, params.PubSub)
										return cfr, nil
									}
								},
							}),
					)
					return netDial.NewNetDialAppNoCrfName(
						fxOptions,
						"Kraken",
						"wss://ws.kraken.com:443",
						impl.WebSocketName,
						//crfName,
						netDial.MaxConnectionsSetting(1),
						netDial.CanDial(canDials...))(params.NetAppFuncInParams)
				},
			}))
}
