package krakenWS

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bhbosman/goCommonMarketData/fullMarketData"
	stream2 "github.com/bhbosman/goCommonMarketData/fullMarketData/stream"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
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
	"github.com/cskr/pubsub"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
	"hash/crc32"
	"strconv"
	"strings"
)

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
	ReqId uint32
	Pair  string
	Name  string
	Depth uint32
}

type registeredSubscription struct {
	channelName  string
	channelId    uint32
	Reqid        uint32
	Pair         string
	Name         string
	LastCheckSum uint32
	depth        uint32
}

type Reactor struct {
	common.BaseConnectionReactor
	messageRouter            *messageRouter.MessageRouter
	connectionID             uint64
	status                   string
	version                  string
	outstandingSubscriptions map[uint32]outstandingSubscription
	registeredSubscriptions  map[uint32]*registeredSubscription
	pairs                    map[registrationKey]*registrationValue
	reqid                    uint32
	republishChannelName     string
	publishChannelName       string
	FmdService               fullMarketDataManagerService.IFmdManagerService
}

func (self *Reactor) handleKrakenStreamSubscribe(inData *krakenStream.Subscribe) error {
	outstandingSubscriptionInstance := outstandingSubscription{
		ReqId: inData.Reqid,
		Pair:  inData.Pair,
		Name:  inData.Name,
		Depth: 100,
	}

	self.outstandingSubscriptions[inData.Reqid] = outstandingSubscriptionInstance
	msg := &krakenWsStream.KrakenWsMessageOutgoing{
		Event: "subscribe",
		Reqid: inData.Reqid,
		Pair:  []string{inData.Pair},
		Subscription: &krakenWsStream.KrakenSubscriptionData{
			Depth:    outstandingSubscriptionInstance.Depth,
			Interval: 0,
			Name:     inData.Name,
			Snapshot: false,
			Token:    "",
		},
	}

	return self.SendTextOpMessage(msg)
}

func (self *Reactor) handleKrakenWsMessageIncoming(inData *krakenWsStream.KrakenWsMessageIncoming) {
	switch inData.Event {
	case "heartbeat":
		self.handleHeartbeat(inData)
	case "ping":
		self.handlePing(inData)
	case "systemStatus":
		self.handleSystemStatus(inData)
	case "subscriptionStatus":
		self.handleSubscriptionStatus(inData)
	default:
		self.Logger.Info(fmt.Sprintf("Unhandled message: %v\nData: %v",
			zap.String("event", inData.Event),
			zap.String("data", inData.String())))
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
				//var FullMarketOrderBook *FullMarketOrderBook
				//if storedValue, ok := self.FullMarketOrderBook[data.Pair]; ok {
				//	FullMarketOrderBook = storedValue
				//} else {
				//	FullMarketOrderBook = NewFullMarketOrderBook(data.Pair, inData[2].(string))
				//	self.FullMarketOrderBook[data.Pair] = FullMarketOrderBook
				//
				//}

				return self.HandleBook(data, inData[1].(map[string]interface{}))
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

func (self *Reactor) handleWebSocketMessage(inData *wsmsg.WebSocketMessage) {
	switch inData.OpCode {
	case wsmsg.WebSocketMessage_OpText:
		if len(inData.Message) > 0 && inData.Message[0] == '[' { //type WebsocketDataResponse []interface{}
			var dataResponse websocketDataResponse

			err := json.Unmarshal(inData.Message, &dataResponse)
			if err != nil {
				return
			}
			_, _ = self.messageRouter.Route(dataResponse)
			return

		} else {
			krakenMessage := &krakenWsStream.KrakenWsMessageIncoming{}
			unMarshaler := jsonpb.Unmarshaler{
				AllowUnknownFields: true,
				AnyResolver:        nil,
			}
			err := unMarshaler.Unmarshal(bytes.NewBuffer(inData.Message), krakenMessage)
			if err != nil {
				return
			}
			_, _ = self.messageRouter.Route(krakenMessage)
			return
		}
	default:
	}
}

func (self *Reactor) Init(
	onSendToReactor rxgo.NextFunc,
	onSendToConnection rxgo.NextFunc,
) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
	_, _, _, err := self.BaseConnectionReactor.Init(
		onSendToReactor,
		onSendToConnection,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	self.republishChannelName = "republishChannel"
	self.publishChannelName = "publishChannel"

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

func (self *Reactor) Close() error {
	for _, value := range self.pairs {
		_ = self.FmdService.Send(
			&stream2.FullMarketData_Clear{
				Instrument: value.pair,
			},
		)
	}
	return self.BaseConnectionReactor.Close()
}
func (self *Reactor) Open() error {
	err := self.BaseConnectionReactor.Open()
	if err != nil {
		return err
	}

	for _, value := range self.pairs {
		message := &krakenStream.Subscribe{
			Reqid: value.reqid,
			Pair:  value.pair,
			Name:  value.name,
		}
		self.OnSendToReactor(message)

		_ = self.FmdService.Send(
			&stream2.FullMarketData_Clear{
				Instrument: value.pair,
			},
		)
	}
	return nil
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
	_ = self.SendTextOpMessage(
		&krakenWsStream.KrakenWsMessageOutgoing{
			Event: "pong",
			Reqid: data.GetReqid(),
		},
	)
	return nil
}

func (self Reactor) handleHeartbeat(_ interface{}) {
	return
}

func (self *Reactor) handleSubscriptionStatus(inData krakenWsStream.ISubscriptionStatus) {
	if data, ok := self.outstandingSubscriptions[inData.GetReqid()]; ok {
		if inData.GetStatus() == "error" {
			// TODO: logging and some status
			err := fmt.Errorf(inData.GetErrorMessage())
			self.Logger.Error("subscription failed", zap.Error(err))
			return
		}
		delete(self.outstandingSubscriptions, data.ReqId)

		self.registeredSubscriptions[inData.GetChannelID()] = &registeredSubscription{
			channelName: inData.GetChannelName(),
			channelId:   inData.GetChannelID(),
			Reqid:       inData.GetReqid(),
			Pair:        inData.GetPair(),
			Name:        inData.GetSubscription().Name,
			depth:       data.Depth,
		}
	}
}

func (self *Reactor) handleTicker(channelData *registeredSubscription, data map[string]interface{}) error {
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

func (self *Reactor) HandlePublishRxHandlerCounters(_ *model.PublishRxHandlerCounters) {}

func (self *Reactor) HandleEmptyQueue(_ *messages.EmptyQueue) {
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

func (self *Reactor) Unknown(i interface{}) {

}

func (self *Reactor) SendTextOpMessage(
	message proto.Message,
) error {
	rws := gomessageblock.NewReaderWriter()
	m := jsonpb.Marshaler{
		OrigName:     false,
		EnumsAsInts:  false,
		EmitDefaults: false,
		Indent:       "",
		AnyResolver:  nil,
	}
	var err error
	err = m.Marshal(rws, message)
	if err != nil {
		return err
	}

	var flatten []byte
	flatten, err = rws.Flatten()
	if err != nil {
		return err
	}

	WebSocketMessage := wsmsg.WebSocketMessage{
		OpCode:  wsmsg.WebSocketMessage_OpText,
		Message: flatten,
	}
	readWriterSize, err := stream.Marshall(&WebSocketMessage)
	if err != nil {
		return err
	}

	self.OnSendToConnection(readWriterSize)
	return nil
}

func (self *Reactor) wsProcessOrderBookPartial(
	instrumentData *registeredSubscription,
	askData, bidData []interface{},
	snapShot bool,
) {
	if snapShot {
		_ = self.FmdService.Send(
			&stream2.FullMarketData_Clear{
				Instrument: instrumentData.Pair,
			},
		)
	}
	for i := range askData {
		asks := askData[i].([]interface{})
		priceAsString := asks[0].(string)
		price, err := strconv.ParseFloat(priceAsString, 64)
		if err != nil {
			return
		}
		volumeAsString := asks[1].(string)
		volume, err := strconv.ParseFloat(volumeAsString, 64)
		if err != nil {
			return
		}

		_ = self.FmdService.Send(
			&stream2.FullMarketData_DeleteOrderInstruction{
				Instrument: instrumentData.Pair,
				Id:         priceAsString,
			},
		)

		if volume != 0 {
			aa := StateTrimmer(priceAsString) + StateTrimmer(volumeAsString)
			_ = self.FmdService.Send(
				&stream2.FullMarketData_AddOrderInstruction{
					Instrument: instrumentData.Pair,
					Order: &stream2.FullMarketData_AddOrder{
						Side:      stream2.OrderSide_AskOrder,
						Id:        priceAsString,
						Price:     price,
						Volume:    volume,
						ExtraData: aa,
					},
				},
			)
		}
	}

	for i := range bidData {
		bids := bidData[i].([]interface{})
		priceAsString := bids[0].(string)
		price, err := strconv.ParseFloat(priceAsString, 64)
		if err != nil {
			return
		}
		volumeAsString := bids[1].(string)
		volume, err := strconv.ParseFloat(volumeAsString, 64)
		if err != nil {
			return
		}

		_ = self.FmdService.Send(
			&stream2.FullMarketData_DeleteOrderInstruction{
				Instrument: instrumentData.Pair,
				Id:         priceAsString,
			},
		)

		if volume != 0 {
			aa := StateTrimmer(priceAsString) + StateTrimmer(volumeAsString)
			_ = self.FmdService.Send(
				&stream2.FullMarketData_AddOrderInstruction{
					Instrument: instrumentData.Pair,
					Order: &stream2.FullMarketData_AddOrder{
						Side:      stream2.OrderSide_BidOrder,
						Id:        priceAsString,
						Price:     price,
						Volume:    volume,
						ExtraData: aa,
					},
				},
			)
		}
	}
}

func (self *Reactor) HandleBook(
	instrumentData *registeredSubscription,
	data map[string]interface{},
) error {
	askSnapshot, askSnapshotExists := data["as"].([]interface{})
	bidSnapshot, bidSnapshotExists := data["bs"].([]interface{})
	if askSnapshotExists || bidSnapshotExists {
		self.wsProcessOrderBookPartial(instrumentData, askSnapshot, bidSnapshot, true)
	} else {
		askData, asksExist := data["a"].([]interface{})
		bidData, bidsExist := data["b"].([]interface{})
		if asksExist || bidsExist {
			self.wsProcessOrderBookPartial(instrumentData, askData, bidData, false)
		}
	}
	_ = self.FmdService.Send(
		fullMarketDataManagerService.NewCallbackMessage(
			instrumentData.Pair,
			func(data interface{}, fullMarketOrderBook fullMarketData.IFullMarketOrderBook) {
				if v, ok := data.(*registeredSubscription); ok {
					for fullMarketOrderBook.BidOrderSide().Size() > int(v.depth) {
						fullMarketOrderBook.BidOrderSide().Remove(fullMarketOrderBook.BidOrderSide().Left().Key)
					}

					for fullMarketOrderBook.AskOrderSide().Size() > int(v.depth) {
						fullMarketOrderBook.AskOrderSide().Remove(fullMarketOrderBook.AskOrderSide().Right().Key)
					}
				}
			},
			instrumentData,
		),
	)

	if checkSumData, checkSumExist := data["c"].(string); checkSumExist {
		if atoi, err := strconv.Atoi(checkSumData); err == nil {
			_ = self.FmdService.Send(
				fullMarketDataManagerService.NewCallbackMessage(
					instrumentData.Pair,
					func(data interface{}, fullMarketOrderBook fullMarketData.IFullMarketOrderBook) {
						if v, ok := data.(*dddd); ok {
							crc := crc32.NewIEEE()
							var count uint32 = 0
							for node := fullMarketOrderBook.AskOrderSide().Left(); node != nil && count < 10; node = node.Next() {
								pp := node.Value.(*fullMarketData.PricePoint)
								unk, _ := pp.List.Get(0)
								ss := unk.(*fullMarketData.FullMarketOrder)
								_, _ = crc.Write([]byte(StateTrimmer(ss.ExtraData)))
								count++
							}
							count = 0
							for node := fullMarketOrderBook.BidOrderSide().Right(); node != nil && count < 10; node = node.Prev() {
								pp := node.Value.(*fullMarketData.PricePoint)
								unk, _ := pp.List.Get(0)
								ss := unk.(*fullMarketData.FullMarketOrder)
								_, _ = crc.Write([]byte(StateTrimmer(ss.ExtraData)))
								count++
							}
							ddd := crc.Sum32()
							if ddd == uint32(v.crc) {
							} else {

							}
						}
					},
					&dddd{
						dept: instrumentData.depth,
						crc:  uint32(atoi),
					},
				),
			)
		}
	}
	return nil
}

type dddd struct {
	dept uint32
	crc  uint32
}

func NewReactor(
	logger *zap.Logger,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc model.ConnectionCancelFunc,
	PubSub *pubsub.PubSub,
	goFunctionCounter GoFunctionCounter.IService,
	UniqueReferenceService interfaces.IUniqueReferenceService,
	FmdService fullMarketDataManagerService.IFmdManagerService,
) *Reactor {
	result := &Reactor{
		BaseConnectionReactor: common.NewBaseConnectionReactor(
			logger,
			cancelCtx,
			cancelFunc,
			connectionCancelFunc,
			UniqueReferenceService.Next("ConnectionReactor"),
			PubSub,
			goFunctionCounter,
		),
		messageRouter:            messageRouter.NewMessageRouter(),
		connectionID:             0,
		status:                   "",
		version:                  "",
		outstandingSubscriptions: make(map[uint32]outstandingSubscription),
		registeredSubscriptions:  make(map[uint32]*registeredSubscription),
		pairs:                    make(map[registrationKey]*registrationValue),
		FmdService:               FmdService,
	}
	_ = result.messageRouter.Add(result.handleWebSocketMessage)
	_ = result.messageRouter.Add(result.handleKrakenStreamSubscribe)
	_ = result.messageRouter.Add(result.handleKrakenWsMessageIncoming)
	_ = result.messageRouter.Add(result.handleWebsocketDataResponse)
	_ = result.messageRouter.Add(result.HandleEmptyQueue)
	_ = result.messageRouter.Add(result.HandlePublishRxHandlerCounters)

	result.messageRouter.RegisterUnknown(result.Unknown)
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

func StateTrimmer(s string) string {
	addValue := false
	ddd := func(r rune) bool {
		switch r {
		case '0':
			return addValue
		case '.':
			return false
		default:
			addValue = true
			return true
		}
	}
	sb := strings.Builder{}
	for _, r := range s {
		if ddd(r) {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}