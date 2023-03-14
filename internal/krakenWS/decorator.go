package krakenWS

import (
	"fmt"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsMultiDialer"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/topStack"
	"github.com/bhbosman/goCommsStacks/websocket"
	"github.com/bhbosman/goFxApp/Services/fileDumpService"
	"github.com/bhbosman/gocommon/fx/PubSub"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocomms/common"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"io"
	"net/url"
	"sync"
)

type decorator struct {
	stoppedCalled  bool
	NetMultiDialer goCommsMultiDialer.INetMultiDialerService
	otherData      instrumentReference.KrakenReferenceData

	pubSub               *pubsub.PubSub
	dialApp              messages.IApp
	dialAppCancelFunc    common.ICancellationContext
	Logger               *zap.Logger
	FullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper
	FmdService           fullMarketDataManagerService.IFmdManagerService
	FileDumpService      fileDumpService.IFileDumpService
}

func NewDecorator(
	Logger *zap.Logger,
	NetMultiDialer goCommsMultiDialer.INetMultiDialerService,
	otherData instrumentReference.KrakenReferenceData,
	pubSub *pubsub.PubSub,
	FullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper,
	FmdService fullMarketDataManagerService.IFmdManagerService,
	FileDumpService fileDumpService.IFileDumpService,
) *decorator {
	return &decorator{
		NetMultiDialer:       NetMultiDialer,
		pubSub:               pubSub,
		Logger:               Logger,
		FullMarketDataHelper: FullMarketDataHelper,
		FmdService:           FmdService,
		otherData:            otherData,
		FileDumpService:      FileDumpService,
	}
}

func (self *decorator) Cancel() {
	if self.dialAppCancelFunc != nil {
		self.dialAppCancelFunc.Cancel()
	}
}

func (self *decorator) Start(ctx context.Context) error {
	if !self.stoppedCalled {
		go func() {
			_ = self.internalStart(ctx)
		}()
		return nil
	}
	return io.EOF
}

func (self *decorator) Stop(ctx context.Context) error {
	if !self.stoppedCalled {
		self.stoppedCalled = true
		return self.internalStop(ctx)
	}
	return io.EOF
}

func (self *decorator) Err() error {
	if self.dialApp != nil {
		return self.dialApp.Err()
	}
	return nil
}

func (self *decorator) internalStart(context.Context) error {
	krakenUrl, _ := url.Parse("wss://ws.kraken.com:443")
	var err error
	var connectionId string
	self.dialApp, self.dialAppCancelFunc, connectionId, err = self.NetMultiDialer.Dial(
		false,
		nil,
		krakenUrl,
		self.reconnect,
		self.dialAppCancelFunc,
		fmt.Sprintf("Kraken.%v", self.otherData.ConnectionName),
		fmt.Sprintf("Kraken.%v", self.otherData.ConnectionName),
		ProvideConnectionReactor(),
		goCommsDefinitions.ProvideTransportFactoryForWebSocketName(
			topStack.Provide(),
			websocket.Provide(),
			bottom.Provide(),
		),
		PubSub.ProvidePubSubInstance("Application", self.pubSub),
		fx.Supply(self.otherData),
		fx.Provide(
			fx.Annotated{
				Target: func() (
					fullMarketDataHelper.IFullMarketDataHelper,
					fullMarketDataManagerService.IFmdManagerService,
					fileDumpService.IFileDumpService,
				) {
					return self.FullMarketDataHelper,
						self.FmdService,
						self.FileDumpService
				},
			},
		),
	)
	if err != nil {
		return err
	}
	err = self.dialApp.Start(context.Background())
	if err != nil {
		self.Logger.Error("Error in start", zap.Error(err))
	}
	_ = self.registerConnectionShutdown(connectionId, self.dialApp, self.Logger, self.dialAppCancelFunc)

	return nil
}

func (self *decorator) internalStop(ctx context.Context) error {
	if self.dialAppCancelFunc != nil {
		self.dialAppCancelFunc.Cancel()
	}
	return nil
}

func (self *decorator) reconnect() {
	go func() {
		if !self.stoppedCalled {
			err := self.internalStop(context.Background())
			if err != nil {
				return
			}
			err = self.internalStart(context.Background())
			if err != nil {
				return
			}
		}
	}()
}

func (self *decorator) registerConnectionShutdown(
	connectionId string,
	connectionApp messages.IApp,
	logger *zap.Logger,
	cancellationContext ...common.ICancellationContext,
) error {
	mutex := sync.Mutex{}
	cancelCalled := false
	cb := func() {
		mutex.Lock()
		b := cancelCalled
		cancelCalled = true
		mutex.Unlock()
		if !b {
			errInGoRoutine := connectionApp.Stop(context.Background())
			if errInGoRoutine != nil {
				logger.Error(
					"Stopping error. not really a problem. informational",
					zap.Error(errInGoRoutine))
			}
			for _, instance := range cancellationContext {
				_ = instance.Remove(connectionId)
			}
		}
	}
	var result error
	for _, ctx := range cancellationContext {
		b, err := ctx.Add(connectionId, cb)
		result = multierr.Append(result, err)
		if !b {
			cb()
			return result
		}
	}
	return result
}
