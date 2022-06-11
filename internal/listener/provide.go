package listener

import (
	"github.com/bhbosman/goCommsStacks"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/gocomms/netListener"
	"github.com/bhbosman/gocomms/stacks/bvisMessageBreaker"
	"github.com/bhbosman/gocomms/stacks/messageCompressor"
	"github.com/bhbosman/gocomms/stacks/messageNumber"
	"github.com/bhbosman/gocomms/stacks/pingPong"
	"github.com/bhbosman/gocomms/stacks/top"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

func CompressedListener(
	serviceIdentifier model.ServiceIdentifier,
	serviceDependentOn model.ServiceIdentifier,
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
					NetAppFuncInParams common.NetAppFuncInParams
				}) messages.CreateAppCallback {
					f := netListener.NewNetListenApp(
						CompressedListenerConnection,
						serviceIdentifier,
						serviceDependentOn,
						CompressedListenerConnection,
						url,
						common.TransportFactoryCompressedName,
						func() (intf.IConnectionReactorFactory, error) {
							cfr := NewFactory(
								crfName,
								params.PubSub,
								func(data proto.Message) (goprotoextra.IReadWriterSize, error) {
									return stream.Marshall(data)
								},
								ConsumerCounter)
							return cfr, nil
						},
						common.MaxConnectionsSetting(maxConnections),
						common.NewConnectionInstanceOptions(
							goCommsStacks.ProvideDefinedStackNames(),
							top.ProvideTopStack(),
							pingPong.ProvidePingPongStacks(),
							messageNumber.ProvideMessageNumberStack(),
							messageCompressor.ProvideMessageCompressorStack(),
							bvisMessageBreaker.ProvideBvisMessageBreakerStack(),
							bottom.ProvideBottomStack()))
					return f(
						params.NetAppFuncInParams)
				},
			}),
	)
}
