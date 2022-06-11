package listener

import (
	"encoding/json"
	"github.com/bhbosman/goCommsStacks"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/gocomms/netListener"
	"github.com/bhbosman/gocomms/stacks/bottom"
	"github.com/bhbosman/gocomms/stacks/top"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

func TextListener(
	serviceIdentifier model.ServiceIdentifier,
	serviceDependentOn model.ServiceIdentifier,
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
					NetAppFuncInParams common.NetAppFuncInParams
				}) messages.CreateAppCallback {
					f := netListener.NewNetListenApp(
						TextListenerConnection,
						serviceIdentifier,
						serviceDependentOn,
						TextListenerConnection,
						url,
						common.TransportFactoryEmptyName,
						func() (intf.IConnectionReactorFactory, error) {
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
						},
						common.MaxConnectionsSetting(maxConnections),
						common.NewConnectionInstanceOptions(
							goCommsStacks.ProvideDefinedStackNames(),
							bottom.ProvideBottomStack(),
							top.ProvideTopStack()))
					return f(
						params.NetAppFuncInParams)
				},
			}),
	)
}
