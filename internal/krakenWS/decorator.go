package krakenWS

import (
	"fmt"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsMultiDialer"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/topStack"
	"github.com/bhbosman/goCommsStacks/websocket"
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
	stoppedCalled        bool
	NetMultiDialer       goCommsMultiDialer.INetMultiDialerService
	krakenConnection     *KrakenConnection
	pubSub               *pubsub.PubSub
	dialApp              messages.IApp
	dialAppCancelFunc    goCommsDefinitions.ICancellationContext
	Logger               *zap.Logger
	FullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper
	FmdService           fullMarketDataManagerService.IFmdManagerService
}

func NewDecorator(
	Logger *zap.Logger,
	NetMultiDialer goCommsMultiDialer.INetMultiDialerService,
	krakenConnection *KrakenConnection,
	pubSub *pubsub.PubSub,
	FullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper,
	FmdService fullMarketDataManagerService.IFmdManagerService,
) *decorator {
	return &decorator{
		NetMultiDialer:       NetMultiDialer,
		krakenConnection:     krakenConnection,
		pubSub:               pubSub,
		Logger:               Logger,
		FullMarketDataHelper: FullMarketDataHelper,
		FmdService:           FmdService,
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
	self.dialApp, self.dialAppCancelFunc, err = self.NetMultiDialer.Dial(
		false,
		nil,
		krakenUrl,
		self.reconnect,
		fmt.Sprintf("Luno.%v", self.krakenConnection.Name),
		fmt.Sprintf("Luno.%v", self.krakenConnection.Name),
		ProvideConnectionReactor(),
		goCommsDefinitions.ProvideTransportFactoryForWebSocketName(
			topStack.ProvideTopStack(),
			websocket.ProvideWebsocketStacks(),
			bottom.Provide(),
		),
		PubSub.ProvidePubSubInstance("Application", self.pubSub),
		fx.Provide(
			fx.Annotated{
				Target: func() (
					fullMarketDataHelper.IFullMarketDataHelper,
					fullMarketDataManagerService.IFmdManagerService,
					*KrakenConnection,
				) {
					return self.FullMarketDataHelper,
						self.FmdService,
						self.krakenConnection
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

	err = self.dialAppCancelFunc.Add(
		func() func() {
			b := false
			return func() {
				if !b {
					b = true
					stopErr := self.dialApp.Stop(context.Background())
					if stopErr != nil {
						self.Logger.Error(
							"Stopping error. not really a problem. informational",
							zap.Error(stopErr))
					}
				}
			}
		}(),
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