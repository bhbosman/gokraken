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
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"io"
	"net/url"
)

type decorator struct {
	stoppedCalled  bool
	NetMultiDialer goCommsMultiDialer.INetMultiDialerService
	otherData      instrumentReference.KrakenReferenceData

	pubSub               *pubsub.PubSub
	dialApp              messages.IApp
	dialAppCancelFunc    goCommsDefinitions.ICancellationContext
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

func (self *decorator) internalStart(ctx context.Context) error {
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

	_, err = self.dialAppCancelFunc.Add(
		connectionId,
		func(dialApp messages.IApp, Logger *zap.Logger) func(goCommsDefinitions.ICancellationContext) {
			b := false
			return func(cancellationContext goCommsDefinitions.ICancellationContext) {
				if !b {
					b = true
					stopErr := dialApp.Stop(context.Background())
					if stopErr != nil {
						Logger.Error(
							"Stopping error. not really a problem. informational",
							zap.Error(stopErr))
					}
					_ = cancellationContext.Remove(connectionId)
				}
			}
		}(self.dialApp, self.Logger),
	)
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
