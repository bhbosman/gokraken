package listener

import (
	"encoding/json"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/netDial"

	"github.com/bhbosman/gocomms/netListener"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

func TextListener(
	pubSub *pubsub.PubSub,
	ConsumerCounter *netDial.CanDialDefaultImpl,
	maxConnections int,
	url string) fx.Option {
	const TextListenerConnection = "TextListenerConnection"
	cfr := NewFactory(
		TextListenerConnection,
		pubSub,
		func(m proto.Message) (goprotoextra.IReadWriterSize, error) {
			bytes, err := json.Marshal(m)
			if err != nil {
				return nil, err
			}
			return gomessageblock.NewReaderWriterBlock(bytes), nil
		},
		ConsumerCounter)
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: netListener.NewNetListenApp(
					TextListenerConnection,
					url,
					impl.CreateEmptyStack,
					cfr,
					netListener.MaxConnectionsSetting(maxConnections)),
			}),
	)
}

func CompressedListener(
	pubSub *pubsub.PubSub,
	ConsumerCounter *netDial.CanDialDefaultImpl,
	maxConnections int, url string) fx.Option {
	const CompressedListenerConnection = "CompressedListenerConnection"
	cfr := NewFactory(
		CompressedListenerConnection,
		pubSub,
		func(data proto.Message) (goprotoextra.IReadWriterSize, error) {
			return stream.Marshall(data)
		},
		ConsumerCounter)
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: netListener.NewNetListenApp(
					CompressedListenerConnection,
					url,
					impl.CreateCompressedStack,
					cfr,
					netListener.MaxConnectionsSetting(maxConnections)),
			}),
	)
}
