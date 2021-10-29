package listener

import (
	"encoding/json"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netDial"

	"github.com/bhbosman/gocomms/netListener"
	"github.com/bhbosman/gomessageblock"
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
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: netListener.NewNetListenApp(
					TextListenerConnection,
					url,
					impl.TransportFactoryEmptyName,
					netListener.MaxConnectionsSetting(maxConnections),
					netListener.FxOption(
						fx.Provide(
							fx.Annotated{
								Target: func(pubSub *pubsub.PubSub) (intf.IConnectionReactorFactory, error) {
									cfr := NewFactory(
										pubSub,
										func(m proto.Message) (goprotoextra.IReadWriterSize, error) {
											bytes, err := json.Marshal(m)
											if err != nil {
												return nil, err
											}
											return gomessageblock.NewReaderWriterBlock(bytes), nil
										},
										ConsumerCounter)
									return cfr, nil
								},
							}))),
			}),
	)
}

func CompressedListener(
	ConsumerCounter *netDial.CanDialDefaultImpl,
	maxConnections int, url string) fx.Option {
	const CompressedListenerConnection = "CompressedListenerConnection"
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: netListener.NewNetListenApp(
					CompressedListenerConnection,
					url,
					impl.TransportFactoryCompressedName,
					netListener.MaxConnectionsSetting(maxConnections),
					netListener.FxOption(
						fx.Provide(
							fx.Annotated{
								Target: func(pubSub *pubsub.PubSub) (intf.IConnectionReactorFactory, error) {
									cfr := NewFactory(
										pubSub,
										func(data proto.Message) (goprotoextra.IReadWriterSize, error) {
											return stream.Marshall(data)
										},
										ConsumerCounter)
									return cfr, nil
								},
							}))),
			}),
	)
}
