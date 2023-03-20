package krakenWS

import (
	"context"
	"fmt"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsMultiDialer"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/topStack"
	"github.com/bhbosman/goCommsStacks/websocket"
	"github.com/bhbosman/goConn"
	"github.com/bhbosman/goFxApp/Services/fileDumpService"
	"github.com/bhbosman/gocommon/fx/PubSub"
	"github.com/bhbosman/gocommon/messages"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/url"
)

type decorator struct {
	NetMultiDialer                goCommsMultiDialer.INetMultiDialerService
	otherData                     instrumentReference.KrakenReferenceData
	pubSub                        *pubsub.PubSub
	Logger                        *zap.Logger
	FullMarketDataHelper          fullMarketDataHelper.IFullMarketDataHelper
	FmdService                    fullMarketDataManagerService.IFmdManagerService
	FileDumpService               fileDumpService.IFileDumpService
	decoratorCancellationContext  goConn.ICancellationContext
	connectionCancellationContext goConn.ICancellationContext
}

func (self *decorator) Start(context.Context) error {
	if self.decoratorCancellationContext.Err() != nil {
		return self.decoratorCancellationContext.Err()
	}
	go func() {
		_ = self.internalStart()
	}()
	return nil
}

func (self *decorator) Stop(context.Context) error {
	self.decoratorCancellationContext.Cancel("ASDADSA")
	return self.decoratorCancellationContext.Err()
}

func (self *decorator) Err() error {
	return self.decoratorCancellationContext.Err()
}

func (self *decorator) internalStart() error {
	krakenUrl, _ := url.Parse("wss://ws.kraken.com:443")
	dialApp, dialAppCancelFunc, connectionId, err := self.NetMultiDialer.Dial(
		false,
		nil,
		krakenUrl,
		self.reconnect,
		self.decoratorCancellationContext,
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
	err = dialApp.Start(context.Background())
	if err != nil {
		self.Logger.Error("Error in start", zap.Error(err))
	}
	self.connectionCancellationContext = dialAppCancelFunc

	return goConn.RegisterConnectionShutdown(
		connectionId,
		func(
			connectionApp messages.IApp,
			logger *zap.Logger,
			connectionCancellationContext goConn.ICancellationContext,
		) func() {
			return func() {
				errInGoRoutine := connectionApp.Stop(context.Background())
				connectionCancellationContext.Cancel("asadasdas")
				if errInGoRoutine != nil {
					logger.Error(
						"Stopping error. not really a problem. informational",
						zap.Error(errInGoRoutine))
				}
			}
		}(dialApp, self.Logger, dialAppCancelFunc),
		self.decoratorCancellationContext,
		dialAppCancelFunc,
	)
}

func (self *decorator) internalStop() error {
	self.connectionCancellationContext.Cancel("adsdasdas")
	self.connectionCancellationContext = nil
	return self.decoratorCancellationContext.Err()
}

func (self *decorator) reconnect() {
	if self.decoratorCancellationContext.Err() != nil {
		return
	}
	go func() {
		err := self.internalStop()
		if err != nil {
			return
		}
		err = self.internalStart()
		if err != nil {
			return
		}
	}()
}

func NewDecorator(
	Logger *zap.Logger,
	NetMultiDialer goCommsMultiDialer.INetMultiDialerService,
	otherData instrumentReference.KrakenReferenceData,
	pubSub *pubsub.PubSub,
	FullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper,
	FmdService fullMarketDataManagerService.IFmdManagerService,
	FileDumpService fileDumpService.IFileDumpService,
	decoratorCancellationContext goConn.ICancellationContext,
) (messages.IApp, error) {
	return &decorator{
		NetMultiDialer:               NetMultiDialer,
		pubSub:                       pubSub,
		Logger:                       Logger,
		FullMarketDataHelper:         FullMarketDataHelper,
		FmdService:                   FmdService,
		otherData:                    otherData,
		FileDumpService:              FileDumpService,
		decoratorCancellationContext: decoratorCancellationContext,
	}, nil
}
