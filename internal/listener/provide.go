package listener

import (
	"encoding/json"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netDial"

	//"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/gocomms/netListener"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

func TextListener(maxConnections int, url string) fx.Option {
	const TextListenerConnection = "TextListenerConnection"
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: impl.ConnectionReactorFactoryConst,
				Target: func(params struct {
					fx.In
					ConsumerCounter *netDial.CanDialDefaultImpl
					PubSub          *pubsub.PubSub `name:"Application"`
				}) (intf.IConnectionReactorFactory, error) {
					return NewFactory(
						TextListenerConnection,
						params.PubSub,
						func(m proto.Message) (goprotoextra.IReadWriterSize, error) {
							bytes, err := json.Marshal(m)
							if err != nil {
								return nil, err
							}
							return gomessageblock.NewReaderWriterBlock(bytes), nil
						},
						params.ConsumerCounter), nil
				},
			}),
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: netListener.NewNetListenApp(
					TextListenerConnection,
					url,
					impl.CreateEmptyStack,
					TextListenerConnection,
					netListener.MaxConnectionsSetting(maxConnections)),
			}),
	)
}

func CompressedListener(maxConnections int, url string) fx.Option {
	const CompressedListenerConnection = "CompressedListenerConnection"
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: impl.ConnectionReactorFactoryConst,
				Target: func(params struct {
					fx.In
					PubSub          *pubsub.PubSub `name:"Application"`
					ConsumerCounter *netDial.CanDialDefaultImpl
				}) (intf.IConnectionReactorFactory, error) {
					return NewFactory(
						CompressedListenerConnection,
						params.PubSub,
						func(data proto.Message) (goprotoextra.IReadWriterSize, error) {
							return stream.Marshall(data)
						},
						params.ConsumerCounter), nil
				},
			}),
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: netListener.NewNetListenApp(
					CompressedListenerConnection,
					url,
					impl.CreateCompressedStack,
					CompressedListenerConnection,
					netListener.MaxConnectionsSetting(maxConnections)),
			}),
	)
}
