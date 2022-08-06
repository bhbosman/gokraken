package krakenWS

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	"github.com/bhbosman/goFxApp/Services/fileDumpService"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/intf"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func ProvideConnectionReactor() fx.Option {
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
						PubSub                 *pubsub.PubSub `name:"Application"`
						GoFunctionCounter      GoFunctionCounter.IService
						UniqueReferenceService interfaces.IUniqueReferenceService
						FmdService             fullMarketDataManagerService.IFmdManagerService
						OtherData              instrumentReference.KrakenReferenceData
						FileDumpService        fileDumpService.IFileDumpService
					},
				) (
					intf.IConnectionReactor,
					error,
				) {
					return newReactor(
							params.Logger,
							params.CancelCtx,
							params.CancelFunc,
							params.ConnectionCancelFunc,
							params.PubSub,
							params.GoFunctionCounter,
							params.UniqueReferenceService,
							params.FmdService,
							params.OtherData,
							params.FileDumpService,
						),
						nil
				},
			},
		),
	)
}
