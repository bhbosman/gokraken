package listener

import (
	"context"
	stream2 "github.com/bhbosman/goCommonMarketData/fullMarketData/stream"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/messageRouter"
	common3 "github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/services/interfaces"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type serializeData func(m proto.Message) (goprotoextra.IReadWriterSize, error)
type reactor struct {
	common2.BaseConnectionReactor
	messageRouter          messageRouter.IMessageRouter
	SerializeData          serializeData
	UniqueReferenceService interfaces.IUniqueReferenceService
	FullMarketDataHelper   fullMarketDataHelper.IFullMarketDataHelper
	FmdService             fullMarketDataManagerService.IFmdManagerService
}

func (self *reactor) Init(params intf.IInitParams) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
	_, _, _, err := self.BaseConnectionReactor.Init(params)
	if err != nil {
		return nil, nil, nil, err
	}

	return func(i interface{}) {
			self.doNext(false, i)
		},
		func(err error) {
			self.doNext(false, err)
		},
		func() {

		},
		nil
}

func (self *reactor) doNext(_ bool, i interface{}) {
	self.messageRouter.Route(i)
}

func (self *reactor) Open() error {
	err := self.BaseConnectionReactor.Open()
	if err != nil {
		return err
	}
	return nil
}

func (self *reactor) Close() error {
	return self.BaseConnectionReactor.Close()
}

//goland:noinspection GoSnakeCaseUsage
func (self *reactor) handleFullMarketData_Instrument_RegisterWrapper(message *stream2.FullMarketData_Instrument_RegisterWrapper) {
	key := self.FullMarketDataHelper.InstrumentChannelName(message.Data.Instrument)
	self.PubSub.AddSub(self.OnSendToConnectionPubSubBag, key)
	message.SetNext(self.OnSendToConnectionPubSubBag)
	_ = self.FmdService.Send(message)
}

//goland:noinspection GoSnakeCaseUsage
func (self *reactor) handleFullMarketData_Instrument_UnregisterWrapper(message *stream2.FullMarketData_Instrument_UnregisterWrapper) {
	key := self.FullMarketDataHelper.InstrumentChannelName(message.Data.Instrument)
	self.PubSub.Unsub(self.OnSendToConnectionPubSubBag, key)
	message.SetNext(self.OnSendToConnectionPubSubBag)
	_ = self.FmdService.Send(message)
}

//goland:noinspection GoSnakeCaseUsage
func (self *reactor) handleFullMarketData_InstrumentList_SubscribeWrapper(*stream2.FullMarketData_InstrumentList_SubscribeWrapper) {
	self.PubSub.AddSub(
		self.OnSendToConnectionPubSubBag,
		self.FullMarketDataHelper.InstrumentListChannelName(),
	)
}

//goland:noinspection GoSnakeCaseUsage
func (self *reactor) handleFullMarketData_InstrumentList_RequestWrapper(message *stream2.FullMarketData_InstrumentList_RequestWrapper) {
	message.SetNext(self.OnSendToConnectionPubSubBag)
	_ = self.FmdService.Send(message)
}

func NewConnectionReactor(
	logger *zap.Logger,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc common3.ConnectionCancelFunc,
	PubSub *pubsub.PubSub,
	SerializeData serializeData,
	GoFunctionCounter GoFunctionCounter.IService,
	UniqueReferenceService interfaces.IUniqueReferenceService,
	FullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper,
	FmdService fullMarketDataManagerService.IFmdManagerService,
) (intf.IConnectionReactor, error) {
	result := &reactor{
		BaseConnectionReactor: common2.NewBaseConnectionReactor(
			logger,
			cancelCtx,
			cancelFunc,
			connectionCancelFunc,
			UniqueReferenceService.Next("ConnectionReactor"),
			PubSub,
			GoFunctionCounter,
		),
		messageRouter:          messageRouter.NewMessageRouter(),
		SerializeData:          SerializeData,
		UniqueReferenceService: UniqueReferenceService,
		FmdService:             FmdService,
		FullMarketDataHelper:   FullMarketDataHelper,
	}
	_ = result.messageRouter.Add(result.handleFullMarketData_InstrumentList_SubscribeWrapper)
	_ = result.messageRouter.Add(result.handleFullMarketData_InstrumentList_RequestWrapper)
	//
	_ = result.messageRouter.Add(result.handleFullMarketData_Instrument_RegisterWrapper)
	_ = result.messageRouter.Add(result.handleFullMarketData_Instrument_UnregisterWrapper)

	return result, nil
}
