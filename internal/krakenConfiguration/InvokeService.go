package krakenConfiguration

import (
	"context"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	"github.com/bhbosman/goCommsMultiDialer"
	"github.com/bhbosman/goFxApp/Services/fileDumpService"
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
				InstrumentReferenceService instrumentReference.IInstrumentReferenceService
				FileDumpService            fileDumpService.IFileDumpService
			},
		) error {
			params.Lifecycle.Append(
				fx.Hook{
					OnStart: func(ctx context.Context) error {
						providers, err := params.InstrumentReferenceService.GetKrakenProviders()
						if err != nil {
							return err
						}
						//allLunoConfiguration := params.KrakenConfigurationService.GetAll()
						f := func(
							otherData instrumentReference.KrakenReferenceData,
						) func() (messages.IApp, context.CancelFunc, error) {
							return func() (messages.IApp, context.CancelFunc, error) {
								dec := krakenWS.NewDecorator(
									params.Logger,
									params.NetMultiDialer,
									otherData,
									params.PubSub,
									params.FullMarketDataHelper,
									params.FmdService,
									params.FileDumpService,
								)
								return dec, dec.Cancel, nil
							}
						}

						for _, configuration := range providers {
							localConfiguration := configuration
							err = multierr.Append(
								err,
								params.FxManagerService.Add(
									localConfiguration.ConnectionName,
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
