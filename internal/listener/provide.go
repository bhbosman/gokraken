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
	"github.com/bhbosman/goCommsStacks/topStack"
	"github.com/bhbosman/gocommon/fx/PubSub"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
	"net/url"
)

func CompressedListener(
	serviceIdentifier model.ServiceIdentifier,
	serviceDependentOn model.ServiceIdentifier,
	ConsumerCounter *goCommsNetDialer.CanDialDefaultImpl,
	maxConnections int,
	urlAsText string,
) fx.Option {
	const CompressedListenerConnection = "CompressedListenerConnection"
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: func(
					params struct {
						fx.In
						PubSub             *pubsub.PubSub `name:"Application"`
						NetAppFuncInParams common.NetAppFuncInParams
					},
				) (messages.CreateAppCallback, error) {
					compressedUrl, err := url.Parse(urlAsText)
					if err != nil {
						return messages.CreateAppCallback{}, err
					}
					f := goCommsNetListener.NewNetListenApp(
						CompressedListenerConnection,
						serviceIdentifier,
						serviceDependentOn,
						CompressedListenerConnection,
						false,
						nil,
						compressedUrl,
						goCommsDefinitions.TransportFactoryCompressedName,
						common.MaxConnectionsSetting(maxConnections),
						common.NewConnectionInstanceOptions(
							goCommsDefinitions.ProvideTransportFactoryForCompressedName(
								topStack.ProvideTopStack(),
								pingPong.ProvidePingPongStacks(),
								messageCompressor.Provide(),
								messageNumber.ProvideMessageNumberStack(),
								bvisMessageBreaker.Provide(),
								bottom.Provide(),
							),
							ProvideConnectionReactorFactory2(),
							PubSub.ProvidePubSubInstance("Application", params.PubSub),
							fx.Provide(
								fx.Annotated{
									Target: func() *goCommsNetDialer.CanDialDefaultImpl {
										return ConsumerCounter
									},
								},
							),
						),
					)
					return f(params.NetAppFuncInParams), nil
				},
			},
		),
	)
}

func ProvideConnectionReactorFactory2() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						CancelCtx            context.Context
						CancelFunc           context.CancelFunc
						ConnectionCancelFunc model.ConnectionCancelFunc
						Logger               *zap.Logger
						ClientContext        interface{}    `name:"UserContext"`
						PubSub               *pubsub.PubSub `name:"Application"`
						ConsumerCounter      *goCommsNetDialer.CanDialDefaultImpl
					},
				) (intf.IConnectionReactor, error) {
					return NewReactor(
							params.Logger,
							params.CancelCtx,
							params.CancelFunc,
							params.ConnectionCancelFunc,
							params.ClientContext,
							params.ConsumerCounter,
							func(data proto.Message) (goprotoextra.IReadWriterSize, error) {
								return stream.Marshall(data)
							},
							params.PubSub),
						nil
				},
			},
		),
	)
}
