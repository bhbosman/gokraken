package listener

import (
	"encoding/json"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsNetDialer"
	"github.com/bhbosman/goCommsNetListener"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/topStack"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/fx/PubSub"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
	"net/url"
)

func TextListener(
	serviceIdentifier model.ServiceIdentifier,
	serviceDependentOn model.ServiceIdentifier,
	ConsumerCounter *goCommsNetDialer.CanDialDefaultImpl,
	maxConnections int,
	urlAsText string) fx.Option {
	const TextListenerConnection = "TextListenerConnection"
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: func(params struct {
					fx.In
					PubSub             *pubsub.PubSub `name:"Application"`
					NetAppFuncInParams common.NetAppFuncInParams
				}) (messages.CreateAppCallback, error) {
					textListenerUrl, err := url.Parse(urlAsText)
					if err != nil {
						return messages.CreateAppCallback{}, err
					}
					f := goCommsNetListener.NewNetListenApp(
						TextListenerConnection,
						serviceIdentifier,
						serviceDependentOn,
						TextListenerConnection,
						false, nil,
						textListenerUrl,
						//goCommsDefinitions.TransportFactoryEmptyName,
						common.MaxConnectionsSetting(maxConnections),
						common.NewConnectionInstanceOptions(
							goCommsDefinitions.ProvideTransportFactoryForEmptyName(
								topStack.ProvideTopStack(),
								bottom.Provide(),
							),
							PubSub.ProvidePubSubInstance("Application", params.PubSub),
							fx.Provide(
								fx.Annotated{
									Target: func() *goCommsNetDialer.CanDialDefaultImpl {
										return ConsumerCounter
									},
								},
							),
							ProvideConnectionReactorFactory(),
						),
					)
					return f(params.NetAppFuncInParams), nil
				},
			}),
	)
}

func ProvideConnectionReactorFactory() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						CancelCtx              context.Context
						CancelFunc             context.CancelFunc
						ConnectionCancelFunc   model.ConnectionCancelFunc
						Logger                 *zap.Logger
						ClientContext          interface{}    `name:"UserContext"`
						PubSub                 *pubsub.PubSub `name:"Application"`
						ConsumerCounter        *goCommsNetDialer.CanDialDefaultImpl
						UniqueReferenceService interfaces.IUniqueReferenceService
					},
				) (intf.IConnectionReactor, error) {
					return NewReactor(
							params.Logger,
							params.CancelCtx,
							params.CancelFunc,
							params.ConnectionCancelFunc,
							params.ClientContext,
							params.ConsumerCounter,
							func(m proto.Message) (goprotoextra.IReadWriterSize, error) {
								bytes, err := json.Marshal(m)
								if err != nil {
									return nil, err
								}
								return gomessageblock.NewReaderWriterBlock(bytes), nil
							},
							params.PubSub,
							params.UniqueReferenceService,
						),
						nil
				},
			},
		),
	)
}
