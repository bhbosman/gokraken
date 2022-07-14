package connection

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bhbosman/goCommsStacks/webSocketMessages/wsmsg"
	krakenStream "github.com/bhbosman/goMessages/kraken/stream"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/common"
	krakenWsStream "github.com/bhbosman/gokraken/internal/krakenWS/internal/stream"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"github.com/golang/protobuf/jsonpb"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
	"strconv"
)

type RePublishMessage struct {
}

type registrationKey struct {
	pair string
	name string
}

func newRegistrationKey(pair string, name string) registrationKey {
	return registrationKey{pair: pair, name: name}
}

type registrationValue struct {
	reqid uint32
	pair  string
	name  string
}

func newRegistrationValue(reqid uint32, pair, name string) *registrationValue {
	return &registrationValue{
		reqid: reqid,
		pair:  pair,
		name:  name,
	}
}

type outstandingSubscription struct {
	Reqid uint32
	Pair  string
	Name  string
}

type registeredSubscription struct {
	channelName string
	channelId   uint32
	Reqid       uint32
	Pair        string
	Name        string
}

type Reactor struct {
	common.BaseConnectionReactor
	messageRouter            *messageRouter.MessageRouter
	connectionID             uint64
	status                   string
	version                  string
	PubSub                   *pubsub.PubSub
	outstandingSubscriptions map[uint32]outstandingSubscription
	registeredSubscriptions  map[uint32]registeredSubscription
	FullMarketOrderBook      map[string]*FullMarketOrderBook
	pairs                    map[registrationKey]*registrationValue
	reqid                    uint32
	republishChannelName     string
	publishChannelName       string
	goFunctionCounter        GoFunctionCounter.IService
}

func (self *Reactor) handleKrakenStreamSubscribe(inData *krakenStream.Subscribe) error {
	self.outstandingSubscriptions[inData.Reqid] = outstandingSubscription{
		Reqid: inData.Reqid,
		Pair:  inData.Pair,
		Name:  inData.Name,
	}
	msg := &krakenWsStream.KrakenWsMessageOutgoing{
		Event: "subscribe",
		Reqid: inData.Reqid,
		Pair:  []string{inData.Pair},
		Subscription: &krakenWsStream.KrakenSubscriptionData{
			Depth:    0,
			Interval: 0,
			Name:     inData.Name,
			Snapshot: false,
			Token:    "",
		},
	}
	if inData.Name == "book" {
		msg.Subscription.Depth = 25
	}

	return SendTextOpMessage(msg, self.ToConnection, self.ToConnectionFuncReplacement)
}

func (self *Reactor) handleKrakenWsMessageIncoming(inData *krakenWsStream.KrakenWsMessageIncoming) error {
	switch inData.Event {
	case "heartbeat":
		return self.handleHeartbeat(inData)
	case "ping":
		return self.handlePing(inData)
	case "systemStatus":
		return self.handleSystemStatus(inData)
	case "subscriptionStatus":
		return self.handleSubscriptionStatus(inData)
	default:
		self.Logger.Info(fmt.Sprintf("Unhandled message: %v\nData: %v",
			zap.String("event", inData.Event),
			zap.String("data", inData.String())))
		return nil
	}
}

func (self *Reactor) handleWebsocketDataResponse(inData websocketDataResponse) error {
	if channelId, ok := inData[0].(float64); ok {
		if data, ok := self.registeredSubscriptions[uint32(channelId)]; ok {
			switch data.Name {
			case "ticker":
				return self.handleTicker(data, inData[1].(map[string]interface{}))
			//case "ohlc":
			//	return self.handleOhlc(data, inData[1].([]interface{}))
			//case "trade":
			//	return self.handleTrade(data, inData[1].([]interface{}))
			//case "spread":
			//	return self.handleSpread(data, inData[1].([]interface{}))
			case "book":
				var FullMarketOrderBook *FullMarketOrderBook
				if storedValue, ok := self.FullMarketOrderBook[data.Pair]; ok {
					FullMarketOrderBook = storedValue
				} else {
					FullMarketOrderBook = NewFullMarketOrderBook(data.Pair, inData[2].(string))
					self.FullMarketOrderBook[data.Pair] = FullMarketOrderBook

				}
				return FullMarketOrderBook.HandleBook(inData[1].(map[string]interface{}))
			}

			//switch data.Name {
			//case krakenWsTicker:
			//	return k.wsProcessTickers(&channelData, response[1].(map[string]interface{}))
			//case krakenWsOHLC:
			//	return k.wsProcessCandles(&channelData, response[1].([]interface{}))
			//case krakenWsOrderbook:
			//	return k.wsProcessOrderBook(&channelData, response[1].(map[string]interface{}))
			//case krakenWsSpread:
			//	k.wsProcessSpread(&channelData, response[1].([]interface{}))
			//case krakenWsTrade:
			//	k.wsProcessTrades(&channelData, response[1].([]interface{}))
			//default:
			//	return fmt.Errorf("%s Unidentified websocket data received: %+v",
			//		k.Name,
			//		response)
			//}

		}

	} else if _, ok := inData[1].(string); ok {
		//err = k.wsHandleAuthDataResponse(dataResponse)
		//if err != nil {
		//	return err
		//}
	}
	return nil
}

type websocketDataResponse []interface{}

func (self *Reactor) handleWebSocketMessageWrapper(inData *wsmsg.WebSocketMessageWrapper) error {
	return self.handleWebSocketMessage(inData.Data)
}
func (self *Reactor) handleWebSocketMessage(inData *wsmsg.WebSocketMessage) error {
	switch inData.OpCode {
	case wsmsg.WebSocketMessage_OpText:
		if len(inData.Message) > 0 && inData.Message[0] == '[' { //type WebsocketDataResponse []interface{}
			var dataResponse websocketDataResponse
			err := json.Unmarshal(inData.Message, &dataResponse)
			if err != nil {
				return err
			}
			_, _ = self.messageRouter.Route(dataResponse)
			return nil

		} else {
			krakenMessage := &krakenWsStream.KrakenWsMessageIncoming{}
			unMarshaler := jsonpb.Unmarshaler{
				AllowUnknownFields: true,
				AnyResolver:        nil,
			}
			err := unMarshaler.Unmarshal(bytes.NewBuffer(inData.Message), krakenMessage)
			if err != nil {
				return err
			}
			_, _ = self.messageRouter.Route(krakenMessage)
			return nil
		}
	case wsmsg.WebSocketMessage_OpEndLoop:
		return nil

	case wsmsg.WebSocketMessage_OpStartLoop:
		return nil
	default:
		return nil
	}
}

func (self Reactor) handleMessageBlockReaderWriter(inData *gomessageblock.ReaderWriter) error {
	marshal, err := stream.UnMarshal(
		inData,
		func(i interface{}) {
			self.ToReactor(false, i)
		},
		func(i interface{}) {
			if unk, ok := i.(goprotoextra.ReadWriterSize); ok {
				self.ToConnection(unk)
			}
		},
	)
	if err != nil {
		println(err.Error())
		return err
	}

	_, err = self.messageRouter.Route(marshal)
	if err != nil {
		return err
	}

	return nil
}

func (self *Reactor) Init(
	onSend goprotoextra.ToConnectionFunc,
	toConnectionReactor goprotoextra.ToReactorFunc,
	onSendReplacement rxgo.NextFunc,
	toConnectionReactorReplacement rxgo.NextFunc,
) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, chan interface{}, error) {
	_, _, _, _, err := self.BaseConnectionReactor.Init(
		onSend,
		toConnectionReactor,
		onSendReplacement,
		toConnectionReactorReplacement,
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	self.republishChannelName = "republishChannel"
	self.publishChannelName = "publishChannel"

	republishChannel := self.PubSub.Sub(self.republishChannelName)

	self.goFunctionCounter.GoRun(
		"Kraken.Init",
		func() {
			<-self.CancelCtx.Done()
			self.PubSub.Unsub(republishChannel, self.republishChannelName)
		},
	)
	// Todo: Register function
	self.goFunctionCounter.GoRun(
		"Kraken Read Republish Channel",
		func() {
			for range republishChannel {
				if self.CancelCtx.Err() == nil {
					_ = self.ToReactor(false, &RePublishMessage{})
				}
			}
		},
	)

	return func(i interface{}) {
			self.doNext(false, i)
		},
		func(err error) {
			self.doNext(false, err)
		},
		func() {

		}, nil, nil
}

func (self *Reactor) Close() error {
	for _, v := range self.FullMarketOrderBook {
		v.Clear()
		v.Publish(true)
	}
	return self.BaseConnectionReactor.Close()
}
func (self *Reactor) Open() error {
	err := self.BaseConnectionReactor.Open()
	if err != nil {
		return err
	}
	self.sendAllRegistration()
	return nil

}

func (self *Reactor) sendAllRegistration() {
	for _, value := range self.pairs {
		message := &krakenStream.Subscribe{
			Reqid: value.reqid,
			Pair:  value.pair,
			Name:  value.name,
		}
		//self.messageRouter.Route(message)
		_ = self.ToReactor(true, message)
	}
}

func (self *Reactor) doNext(_ bool, i interface{}) {
	_, _ = self.messageRouter.Route(i)
}

func (self *Reactor) handleSystemStatus(data krakenWsStream.ISystemStatus) error {
	self.status = data.GetStatus()
	self.connectionID = data.GetConnectionID()
	self.version = data.GetVersion()
	return nil
}

func (self Reactor) handlePing(data krakenWsStream.IPing) error {
	outgoing := &krakenWsStream.KrakenWsMessageOutgoing{}
	outgoing.Event = "pong"
	outgoing.Reqid = data.GetReqid()

	return SendTextOpMessage(outgoing, self.ToConnection, self.ToConnectionFuncReplacement)
}

func (self Reactor) handleHeartbeat(_ interface{}) error {
	return nil
}

func (self *Reactor) handleSubscriptionStatus(inData krakenWsStream.ISubscriptionStatus) error {
	if data, ok := self.outstandingSubscriptions[inData.GetReqid()]; ok {
		if inData.GetStatus() == "error" {
			err := fmt.Errorf(inData.GetErrorMessage())
			self.Logger.Error("subscription failed", zap.Error(err))
			return err
		}
		delete(self.outstandingSubscriptions, data.Reqid)

		self.registeredSubscriptions[inData.GetChannelID()] = registeredSubscription{
			channelName: inData.GetChannelName(),
			channelId:   inData.GetChannelID(),
			Reqid:       inData.GetReqid(),
			Pair:        inData.GetPair(),
			Name:        inData.GetSubscription().Name,
		}
	}
	return nil
}

func (self *Reactor) handleTicker(channelData registeredSubscription, data map[string]interface{}) error {
	closePrice, err := strconv.ParseFloat(data["c"].([]interface{})[0].(string), 64)
	if err != nil {
		return err
	}
	openPrice, err := strconv.ParseFloat(data["o"].([]interface{})[0].(string), 64)
	if err != nil {
		return err
	}
	highPrice, err := strconv.ParseFloat(data["h"].([]interface{})[0].(string), 64)
	if err != nil {
		return err
	}
	lowPrice, err := strconv.ParseFloat(data["l"].([]interface{})[0].(string), 64)
	if err != nil {
		return err
	}
	quantity, err := strconv.ParseFloat(data["v"].([]interface{})[0].(string), 64)
	if err != nil {
		return err
	}
	ask, err := strconv.ParseFloat(data["a"].([]interface{})[0].(string), 64)
	if err != nil {
		return err
	}
	bid, err := strconv.ParseFloat(data["b"].([]interface{})[0].(string), 64)
	if err != nil {
		return err
	}

	priceData := &krakenStream.Price{

		Open:   openPrice,
		Close:  closePrice,
		Volume: quantity,
		High:   highPrice,
		Low:    lowPrice,
		Bid:    bid,
		Ask:    ask,
		Pair:   channelData.Pair,
	}
	if priceData != nil {

	}
	return nil
}

func (self *Reactor) HandlePublishMessage(_ *RePublishMessage) error {
	return self.publishData(true)
}

func (self *Reactor) HandleEmptyQueue(_ *messages.EmptyQueue) error {
	return self.publishData(false)
}

func (self *Reactor) publishData(forcePublish bool) error {
	for _, v := range self.FullMarketOrderBook {
		top5 := v.Publish(forcePublish)
		if top5 != nil {
			self.PubSub.Pub(top5, self.publishChannelName)
		}
	}
	return nil
}

func (self *Reactor) Register(pair string, name string) error {
	key := newRegistrationKey(pair, name)
	if _, ok := self.pairs[key]; !ok {
		self.reqid++
		value := newRegistrationValue(self.reqid, pair, name)
		self.pairs[key] = value
	}
	return nil
}

func NewReactor(
	logger *zap.Logger,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc model.ConnectionCancelFunc,
	//userContext interface{},
	PubSub *pubsub.PubSub,
	goFunctionCounter GoFunctionCounter.IService,
	UniqueReferenceService interfaces.IUniqueReferenceService,
) *Reactor {
	result := &Reactor{
		BaseConnectionReactor: common.NewBaseConnectionReactor(
			logger,
			cancelCtx,
			cancelFunc,
			connectionCancelFunc,
			//userContext,
			UniqueReferenceService.Next("ConnectionReactor"),
		),
		messageRouter:            messageRouter.NewMessageRouter(),
		connectionID:             0,
		status:                   "",
		version:                  "",
		PubSub:                   PubSub,
		outstandingSubscriptions: make(map[uint32]outstandingSubscription),
		registeredSubscriptions:  make(map[uint32]registeredSubscription),
		FullMarketOrderBook:      make(map[string]*FullMarketOrderBook),
		pairs:                    make(map[registrationKey]*registrationValue),
		goFunctionCounter:        goFunctionCounter,
	}
	_ = result.messageRouter.Add(result.handleMessageBlockReaderWriter)
	_ = result.messageRouter.Add(result.handleWebSocketMessage)
	_ = result.messageRouter.Add(result.handleWebSocketMessageWrapper)
	_ = result.messageRouter.Add(result.handleKrakenStreamSubscribe)
	_ = result.messageRouter.Add(result.handleKrakenWsMessageIncoming)
	_ = result.messageRouter.Add(result.handleWebsocketDataResponse)
	_ = result.messageRouter.Add(result.HandleEmptyQueue)
	_ = result.messageRouter.Add(result.HandlePublishMessage)

	_ = result.Register("XBT/USD", "book")
	_ = result.Register("XBT/EUR", "book")
	_ = result.Register("XBT/CAD", "book")
	_ = result.Register("EUR/USD", "book")
	_ = result.Register("GBP/USD", "book")
	_ = result.Register("USD/CAD", "book")

	//_ = result.Register("XBT/USD", "ohlc")
	//_ = result.Register("XBT/USD", "spread")
	//_ = result.Register("XBT/USD", "ticker")
	//_ = result.Register("XBT/USD", "trade")

	return result
}
