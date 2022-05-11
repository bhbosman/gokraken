package listener

import (
	"encoding/json"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/gomessageblock"

	"github.com/bhbosman/gocomms/netListener"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

func TextListener(
	ConsumerCounter *netDial.CanDialDefaultImpl,
	maxConnections int,
	url string) fx.Option {
	const TextListenerConnection = "TextListenerConnection"
	crfName := "TextListenerConnection.CRF"
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
										cfr := NewFactory(
											crfName,
											params.PubSub,
											func(m proto.Message) (goprotoextra.IReadWriterSize, error) {
												bytes, err := json.Marshal(m)
												if err != nil {
													return nil, err
												}
												return gomessageblock.NewReaderWriterBlock(bytes), nil
											},
											ConsumerCounter)
										return cfr, nil
									}
								},
							}),
					)
					return netListener.NewNetListenApp(
						fxOptions,
						TextListenerConnection,
						url,
						impl.TransportFactoryEmptyName,
						crfName,
						netListener.MaxConnectionsSetting(maxConnections))(params.NetAppFuncInParams)
				},
			}),
	)
}

func CompressedListener(
	ConsumerCounter *netDial.CanDialDefaultImpl,
	maxConnections int, url string) fx.Option {
	const CompressedListenerConnection = "CompressedListenerConnection"
	crfName := "CompressedListenerConnection.CRF"
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

										cfr := NewFactory(
											crfName,
											params.PubSub,
											func(data proto.Message) (goprotoextra.IReadWriterSize, error) {
												return stream.Marshall(data)
											},
											ConsumerCounter)
										return cfr, nil
									}
								},
							}),
					)
					return netListener.NewNetListenApp(
						fxOptions,
						CompressedListenerConnection,
						url,
						impl.TransportFactoryCompressedName,
						crfName,
						netListener.MaxConnectionsSetting(maxConnections))(params.NetAppFuncInParams)
				},
			}),
	)
}
