package listener

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/services/interfaces"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
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
						UniqueReferenceService interfaces.IUniqueReferenceService
						GoFunctionCounter      GoFunctionCounter.IService
						FullMarketDataHelper   fullMarketDataHelper.IFullMarketDataHelper
						FmdService             fullMarketDataManagerService.IFmdManagerService
					},
				) (intf.IConnectionReactor, error) {
					return NewConnectionReactor(
						params.Logger,
						params.CancelCtx,
						params.CancelFunc,
						params.ConnectionCancelFunc,
						params.PubSub,
						func(data proto.Message) (goprotoextra.IReadWriterSize, error) {
							return stream.Marshall(data)
						},
						params.GoFunctionCounter,
						params.UniqueReferenceService,
						params.FullMarketDataHelper,
						params.FmdService,
					)
				},
			},
		),
	)
}
