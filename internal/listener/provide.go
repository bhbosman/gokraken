package listener

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsNetDialer"
	"github.com/bhbosman/goCommsNetListener"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/bvisMessageBreaker"
	"github.com/bhbosman/goCommsStacks/messageCompressor"
	"github.com/bhbosman/goCommsStacks/messageNumber"
	"github.com/bhbosman/goCommsStacks/pingPong"
	"github.com/bhbosman/goCommsStacks/top"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

func CompressedListener(
	serviceIdentifier model.ServiceIdentifier,
	serviceDependentOn model.ServiceIdentifier,
	ConsumerCounter *goCommsNetDialer.CanDialDefaultImpl,
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
					f := goCommsNetListener.NewNetListenApp(
						CompressedListenerConnection,
						serviceIdentifier,
						serviceDependentOn,
						CompressedListenerConnection,
						url,
						gocommon.TransportFactoryCompressedName,
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
							goCommsDefinitions.ProvideTransportFactoryForCompressedName(
								top.ProvideTopStack(),
								pingPong.ProvidePingPongStacks(),
								messageCompressor.ProvideMessageCompressorStack(),
								messageNumber.ProvideMessageNumberStack(),
								bvisMessageBreaker.ProvideBvisMessageBreakerStack(),
								bottom.ProvideBottomStack(),
							),
						),
					)
					return f(
						params.NetAppFuncInParams)
				},
			}),
	)
}
