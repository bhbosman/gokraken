package listener

import (
	"encoding/json"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsNetDialer"
	"github.com/bhbosman/goCommsNetListener"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/top"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
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
	crfName := "TextListenerConnection.CRF"
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
						goCommsDefinitions.TransportFactoryEmptyName,
						common.MaxConnectionsSetting(maxConnections),
						common.NewConnectionInstanceOptions(
							goCommsDefinitions.ProvideTransportFactoryForEmptyName(
								top.ProvideTopStack(),
								bottom.Provide(),
							),
							fx.Provide(
								fx.Annotated{
									Target: func() (intf.IConnectionReactorFactory, error) {
										return NewFactory(
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
									},
								},
							),
						),
					)
					return f(params.NetAppFuncInParams), nil
				},
			}),
	)
}
