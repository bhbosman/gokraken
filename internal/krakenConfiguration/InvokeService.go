package krakenConfiguration

import (
	"context"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommsMultiDialer"
	fxAppManager "github.com/bhbosman/goFxAppManager/service"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gokraken/internal/krakenWS"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

func InvokeService() fx.Option {
	return fx.Invoke(
		func(
			params struct {
				fx.In
				Logger                     *zap.Logger
				KrakenConfigurationService IKrakenConfigurationService
				Lifecycle                  fx.Lifecycle
				FxManagerService           fxAppManager.IFxManagerService
				NetMultiDialer             goCommsMultiDialer.INetMultiDialerService
				PubSub                     *pubsub.PubSub `name:"Application"`
				FullMarketDataHelper       fullMarketDataHelper.IFullMarketDataHelper
				FmdService                 fullMarketDataManagerService.IFmdManagerService
			},
		) error {
			params.Lifecycle.Append(
				fx.Hook{
					OnStart: func(ctx context.Context) error {
						allLunoConfiguration := params.KrakenConfigurationService.GetAll()
						var err error
						f := func(krakenConnection *krakenWS.KrakenConnection) func() (messages.IApp, context.CancelFunc, error) {
							return func() (messages.IApp, context.CancelFunc, error) {
								dec := krakenWS.NewDecorator(
									params.Logger,
									params.NetMultiDialer,
									krakenConnection,
									params.PubSub,
									params.FullMarketDataHelper,
									params.FmdService,
								)
								return dec, dec.Cancel, nil
							}
						}

						for _, configuration := range allLunoConfiguration {
							localConfiguration := configuration
							err = multierr.Append(
								err,
								params.FxManagerService.Add(
									localConfiguration.Name,
									f(localConfiguration),
								),
							)
						}
						return nil
					},
					OnStop: nil,
				})
			return nil
		},
	)
}
