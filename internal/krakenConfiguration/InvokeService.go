package krakenConfiguration

import (
	"context"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	"github.com/bhbosman/goCommsMultiDialer"
	"github.com/bhbosman/goConn"
	"github.com/bhbosman/goFxApp/Services/fileDumpService"
	fxAppManager "github.com/bhbosman/goFxAppManager/service"
	"github.com/bhbosman/gocommon"
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
				ApplicationContext         context.Context `name:"Application"`
			},
		) error {
			params.Lifecycle.Append(
				fx.Hook{
					OnStart: func(ctx context.Context) error {
						providers, err := params.InstrumentReferenceService.GetKrakenProviders()
						if err != nil {
							return err
						}
						f := func(
							name string,
							otherData instrumentReference.KrakenReferenceData,
						) func() (gocommon.IApp, gocommon.ICancellationContext, error) {
							return func() (gocommon.IApp, gocommon.ICancellationContext, error) {
								namedLogger := params.Logger.Named(name)
								ctx, cancelFunc := context.WithCancel(params.ApplicationContext)
								cancellationContext, err := goConn.NewCancellationContextNoCloser(name, cancelFunc, ctx, namedLogger)
								if err != nil {
									return nil, nil, err
								}
								dec, err := krakenWS.NewDecorator(
									namedLogger,
									params.NetMultiDialer,
									otherData,
									params.PubSub,
									params.FullMarketDataHelper,
									params.FmdService,
									params.FileDumpService,
									cancellationContext,
								)
								if err != nil {
									return nil, nil, err
								}
								return dec, cancellationContext, nil
							}
						}

						for _, configuration := range providers {
							localConfiguration := configuration
							err = multierr.Append(
								err,
								params.FxManagerService.Add(
									localConfiguration.ConnectionName,
									f(configuration.ConnectionName, localConfiguration),
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
