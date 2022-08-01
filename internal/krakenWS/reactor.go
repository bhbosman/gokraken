package krakenWS

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bhbosman/goCommonMarketData/fullMarketData"
	stream2 "github.com/bhbosman/goCommonMarketData/fullMarketData/stream"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
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
)

type Subscribe struct {
}

type registeredSubscription struct {
	instrumentReference.KrakenFeed
	status string
}

type subscriptionKey struct {
	pair string
	name string
}

type Reactor struct {
	common.BaseConnectionReactor
	messageRouter             *messageRouter.MessageRouter
	connectionID              uint64
	status                    string
	version                   string
	registeredSubscriptionMap map[subscriptionKey]registeredSubscription
	FmdService                fullMarketDataManagerService.IFmdManagerService
	otherData                 instrumentReference.KrakenReferenceData
}

func (self *Reactor) handleKrakenStreamSubscribe(_ *Subscribe) error {
	var ss []string
	messagesToFmdService := make([]interface{}, len(self.otherData.Feeds))
	for _, feed := range self.otherData.Feeds {
		ss = append(ss, feed.Pair)
		name := self.otherData.Type
		if name == "book" {
			name = fmt.Sprintf("%v-%v", name, self.otherData.Depth)
		}
		key := subscriptionKey{
			pair: feed.Pair,
			name: name,
		}
		self.registeredSubscriptionMap[key] = registeredSubscription{
			KrakenFeed: *feed,
			status:     "Busy subscribing...",
		}
		messagesToFmdService = append(
			messagesToFmdService,
			&stream2.FullMarketData_Instrument_InstrumentStatus{
				Instrument: feed.SystemName,
				Status:     "Busy subscribing...",
			},
		)
	}
	self.FmdService.MultiSend(messagesToFmdService...)

	msg := &krakenWsStream.KrakenWsMessageOutgoing{
		Event: "subscribe",
		Reqid: 1,
		Pair:  ss,
		Subscription: &krakenWsStream.KrakenSubscriptionData{
			Depth:    uint32(self.otherData.Depth),
			Interval: 0,
			Name:     self.otherData.Type,
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

func (self *Reactor) handleWebsocketDataResponse(inData websocketDataResponse) {
	if len(inData) == 4 {
		if pair, ok := inData[3].(string); ok {
			if name, ok := inData[2].(string); ok {
				key := subscriptionKey{
					pair: pair,
					name: name,
				}
				if v, ok := self.registeredSubscriptionMap[key]; ok {
					switch self.otherData.Type {
					case "book":
						self.HandleBook(v, inData[1].(map[string]interface{}))
					}
				}
			}
		}
	}
	//if channelId, ok := inData[0].(float64); ok {
	//	if data, ok := self.registeredSubscriptions[uint32(channelId)]; ok {
	//		switch data.Name {
	//		case "ticker":
	//			self.handleTicker(data, inData[1].(map[string]interface{}))
	//		//case "ohlc":
	//		//	return self.handleOhlc(data, inData[1].([]interface{}))
	//		//case "trade":
	//		//	return self.handleTrade(data, inData[1].([]interface{}))
	//		//case "spread":
	//		//	return self.handleSpread(data, inData[1].([]interface{}))
	//		case "book":
	//			self.HandleBook(data, inData[1].(map[string]interface{}))
	//		}
	//
	//		//switch data.Name {
	//		//case krakenWsTicker:
	//		//	return k.wsProcessTickers(&channelData, response[1].(map[string]interface{}))
	//		//case krakenWsOHLC:
	//		//	return k.wsProcessCandles(&channelData, response[1].([]interface{}))
	//		//case krakenWsOrderbook:
	//		//	return k.wsProcessOrderBook(&channelData, response[1].(map[string]interface{}))
	//		//case krakenWsSpread:
	//		//	k.wsProcessSpread(&channelData, response[1].([]interface{}))
	//		//case krakenWsTrade:
	//		//	k.wsProcessTrades(&channelData, response[1].([]interface{}))
	//		//default:
	//		//	return fmt.Errorf("%s Unidentified websocket data received: %+v",
	//		//		k.Name,
	//		//		response)
	//		//}
	//
	//	}
	//
	//} else if _, ok := inData[1].(string); ok {
	//	//err = k.wsHandleAuthDataResponse(dataResponse)
	//	//if err != nil {
	//	//	return err
	//	//}
	//}
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
			self.messageRouter.Route(dataResponse)
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
			self.messageRouter.Route(krakenMessage)
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
	for _, value := range self.otherData.Feeds {
		_ = self.FmdService.Send(
			&stream2.FullMarketData_Clear{
				Instrument: value.SystemName,
			},
		)
		_ = self.FmdService.Send(
			&stream2.FullMarketData_RemoveInstrumentInstruction{
				Instrument: value.SystemName,
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

	self.OnSendToReactor(&Subscribe{})
	messagesToFmdService := make([]interface{}, len(self.otherData.Feeds)*2)
	for _, value := range self.otherData.Feeds {
		messagesToFmdService = append(
			messagesToFmdService,
			&stream2.FullMarketData_Clear{
				Instrument: value.SystemName,
			},
			&stream2.FullMarketData_Instrument_InstrumentStatus{
				Instrument: value.SystemName,
				Status:     "Connection opened",
			},
		)
	}
	self.FmdService.MultiSend(messagesToFmdService...)
	return nil
}

func (self *Reactor) doNext(_ bool, i interface{}) {
	self.messageRouter.Route(i)
}

func (self *Reactor) handleSystemStatus(data krakenWsStream.ISystemStatus) {
	self.status = data.GetStatus()
	self.connectionID = data.GetConnectionID()
	self.version = data.GetVersion()
}

func (self Reactor) handlePing(data krakenWsStream.IPing) {
	_ = self.SendTextOpMessage(
		&krakenWsStream.KrakenWsMessageOutgoing{
			Event: "pong",
			Reqid: data.GetReqid(),
		},
	)
}

func (self Reactor) handleHeartbeat(_ interface{}) {
	return
}

func (self *Reactor) handleSubscriptionStatus(inData krakenWsStream.ISubscriptionStatus) {
	if inData.GetStatus() == "error" {
		for key, value := range self.registeredSubscriptionMap {
			if key.pair == inData.GetPair() {
				value.status = inData.GetStatus()
				self.registeredSubscriptionMap[key] = value
				_ = self.FmdService.Send(
					&stream2.FullMarketData_Instrument_InstrumentStatus{
						Instrument: value.SystemName,
						Status:     inData.GetStatus(),
					},
				)
				break
			}
		}
		return
	}
	key := subscriptionKey{
		pair: inData.GetPair(),
		name: inData.GetChannelName(),
	}
	if v, ok := self.registeredSubscriptionMap[key]; ok {
		v.status = inData.GetStatus()
		self.registeredSubscriptionMap[key] = v
		_ = self.FmdService.Send(
			&stream2.FullMarketData_Instrument_InstrumentStatus{
				Instrument: v.SystemName,
				Status:     inData.GetStatus(),
			},
		)
	}
}

func (self *Reactor) handleTicker(channelData *registeredSubscription, data map[string]interface{}) {
	closePrice, err := strconv.ParseFloat(data["c"].([]interface{})[0].(string), 64)
	if err != nil {
		self.Logger.Error("Error in ClosePrice convert", zap.Error(err))
		return
	}
	openPrice, err := strconv.ParseFloat(data["o"].([]interface{})[0].(string), 64)
	if err != nil {
		self.Logger.Error("Error in openPrice convert", zap.Error(err))
		return
	}
	highPrice, err := strconv.ParseFloat(data["h"].([]interface{})[0].(string), 64)
	if err != nil {
		self.Logger.Error("Error in highPrice convert", zap.Error(err))
		return
	}
	lowPrice, err := strconv.ParseFloat(data["l"].([]interface{})[0].(string), 64)
	if err != nil {
		self.Logger.Error("Error in lowPrice convert", zap.Error(err))
		return
	}
	quantity, err := strconv.ParseFloat(data["v"].([]interface{})[0].(string), 64)
	if err != nil {
		self.Logger.Error("Error in quantity convert", zap.Error(err))
		return
	}
	ask, err := strconv.ParseFloat(data["a"].([]interface{})[0].(string), 64)
	if err != nil {
		self.Logger.Error("Error in ask convert", zap.Error(err))
		return
	}
	bid, err := strconv.ParseFloat(data["b"].([]interface{})[0].(string), 64)
	if err != nil {
		self.Logger.Error("Error in bid convert", zap.Error(err))
		return
	}

	priceData := &krakenStream.Price{

		Open:   openPrice,
		Close:  closePrice,
		Volume: quantity,
		High:   highPrice,
		Low:    lowPrice,
		Bid:    bid,
		Ask:    ask,
		Pair:   channelData.SystemName,
	}
	if priceData != nil {

	}
}

func (self *Reactor) HandlePublishRxHandlerCounters(_ *model.PublishRxHandlerCounters) {}

func (self *Reactor) HandleEmptyQueue(_ *messages.EmptyQueue) {
}

func (self *Reactor) Unknown(_ interface{}) {

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
	instrumentData registeredSubscription,
	askData, bidData []interface{},
	snapShot bool,
) {
	if snapShot {
		_ = self.FmdService.Send(
			&stream2.FullMarketData_Clear{
				Instrument: instrumentData.SystemName,
			},
		)
	}
	askMessages := make([]interface{}, len(askData)*2)
	for i := range askData {
		asks := askData[i].([]interface{})
		priceAsString := asks[0].(string)
		price, err := strconv.ParseFloat(priceAsString, 64)
		if err != nil {
			return
		}
		priceAsString = fmt.Sprintf("%.12f", price)
		volumeAsString := asks[1].(string)
		volume, err := strconv.ParseFloat(volumeAsString, 64)
		if err != nil {
			return
		}
		askMessages = append(
			askMessages,
			&stream2.FullMarketData_DeleteOrderInstruction{
				Instrument: instrumentData.SystemName,
				Id:         priceAsString,
			},
		)

		if volume != 0 {
			buffer := bytes.Buffer{}
			self.preCrc(&buffer, priceAsString)
			self.preCrc(&buffer, volumeAsString)
			askMessages = append(
				askMessages,
				&stream2.FullMarketData_AddOrderInstruction{
					Instrument: instrumentData.SystemName,
					Order: &stream2.FullMarketData_AddOrder{
						Side:      stream2.OrderSide_AskOrder,
						Id:        priceAsString,
						Price:     price,
						Volume:    volume,
						ExtraData: buffer.Bytes(),
					},
				},
			)
		}
	}
	self.FmdService.MultiSend(askMessages...)

	bidMessages := make([]interface{}, len(bidData)*2)
	for i := range bidData {
		bids := bidData[i].([]interface{})
		priceAsString := bids[0].(string)
		price, err := strconv.ParseFloat(priceAsString, 64)
		if err != nil {
			return
		}
		priceAsString = fmt.Sprintf("%.12f", price)
		volumeAsString := bids[1].(string)
		volume, err := strconv.ParseFloat(volumeAsString, 64)
		if err != nil {
			return
		}

		bidMessages = append(
			bidMessages,
			&stream2.FullMarketData_DeleteOrderInstruction{
				Instrument: instrumentData.SystemName,
				Id:         priceAsString,
			},
		)

		if volume != 0 {
			buffer := bytes.Buffer{}
			self.preCrc(&buffer, priceAsString)
			self.preCrc(&buffer, volumeAsString)
			bidMessages = append(
				bidMessages,
				&stream2.FullMarketData_AddOrderInstruction{
					Instrument: instrumentData.SystemName,
					Order: &stream2.FullMarketData_AddOrder{
						Side:      stream2.OrderSide_BidOrder,
						Id:        priceAsString,
						Price:     price,
						Volume:    volume,
						ExtraData: buffer.Bytes(),
					},
				},
			)
		}
	}
	self.FmdService.MultiSend(bidMessages...)

}

func (self *Reactor) preCrc(writer *bytes.Buffer, s string) {
	addValue := false
	cb := func(r rune) bool {
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
	for _, r := range s {
		if cb(r) {
			writer.WriteRune(r)
		}
	}
}

func (self *Reactor) HandleBook(
	instrumentData registeredSubscription,
	data map[string]interface{},
) {
	askSnapshot, askSnapshotExists := data["as"].([]interface{})
	bidSnapshot, bidSnapshotExists := data["bs"].([]interface{})
	askData, asksExist := data["a"].([]interface{})
	bidData, bidsExist := data["b"].([]interface{})
	checkSumData, checkSumExist := data["c"].(string)

	if askSnapshotExists || bidSnapshotExists {
		self.wsProcessOrderBookPartial(instrumentData, askSnapshot, bidSnapshot, true)
	} else if asksExist || bidsExist {
		self.wsProcessOrderBookPartial(instrumentData, askData, bidData, false)
	}

	_ = self.FmdService.Send(
		fullMarketDataManagerService.NewCallbackMessage(
			instrumentData.SystemName,
			func(data interface{}, fullMarketOrderBook fullMarketData.IFullMarketOrderBook) []interface{} {
				var result []interface{}
				for fullMarketOrderBook.BidOrderSide().Size() > self.otherData.Depth {
					key := fullMarketOrderBook.BidOrderSide().Left().Key
					fullMarketOrderBook.BidOrderSide().Remove(key)

					keyAsFloat := key.(float64)
					priceAsString := fmt.Sprintf("%.12f", keyAsFloat)

					result = append(result,
						&stream2.FullMarketData_DeleteOrderInstruction{
							Instrument: instrumentData.SystemName,
							Id:         priceAsString,
						},
					)
				}

				for fullMarketOrderBook.AskOrderSide().Size() > self.otherData.Depth {
					key := fullMarketOrderBook.AskOrderSide().Right().Key
					fullMarketOrderBook.AskOrderSide().Remove(key)

					keyAsFloat := key.(float64)
					priceAsString := fmt.Sprintf("%.12f", keyAsFloat)

					result = append(result,
						&stream2.FullMarketData_DeleteOrderInstruction{
							Instrument: instrumentData.SystemName,
							Id:         priceAsString,
						},
					)
				}
				return result
			},
			instrumentData,
		),
	)

	if checkSumExist {
		if atoi, err := strconv.Atoi(checkSumData); err == nil {
			_ = self.FmdService.Send(
				fullMarketDataManagerService.NewCallbackMessage(
					instrumentData.SystemName,
					func(data interface{}, fullMarketOrderBook fullMarketData.IFullMarketOrderBook) []interface{} {
						if v, ok := data.(int); ok {
							crc := crc32.NewIEEE()
							var count uint32 = 0
							for node := fullMarketOrderBook.AskOrderSide().Left(); node != nil && count < 10; node = node.Next() {
								pp := node.Value.(*fullMarketData.PricePoint)
								unk, _ := pp.List.Get(0)
								ss := unk.(*fullMarketData.FullMarketOrder)
								_, _ = crc.Write(ss.ExtraData)
								count++
							}
							count = 0
							for node := fullMarketOrderBook.BidOrderSide().Right(); node != nil && count < 10; node = node.Prev() {
								pp := node.Value.(*fullMarketData.PricePoint)
								unk, _ := pp.List.Get(0)
								ss := unk.(*fullMarketData.FullMarketOrder)
								_, _ = crc.Write(ss.ExtraData)
								count++
							}
							ddd := crc.Sum32()
							if int(ddd) == v {
							} else {
								self.CancelFunc()
							}
						}
						return nil
					},
					&atoi,
				),
			)
		}
	}
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
	otherData instrumentReference.KrakenReferenceData,
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
		messageRouter:             messageRouter.NewMessageRouter(),
		connectionID:              0,
		status:                    "",
		version:                   "",
		FmdService:                FmdService,
		registeredSubscriptionMap: make(map[subscriptionKey]registeredSubscription),
		otherData:                 otherData,
	}
	_ = result.messageRouter.Add(result.handleWebSocketMessage)
	_ = result.messageRouter.Add(result.handleKrakenStreamSubscribe)
	_ = result.messageRouter.Add(result.handleKrakenWsMessageIncoming)
	_ = result.messageRouter.Add(result.handleWebsocketDataResponse)
	_ = result.messageRouter.Add(result.HandleEmptyQueue)
	_ = result.messageRouter.Add(result.HandlePublishRxHandlerCounters)

	result.messageRouter.RegisterUnknown(result.Unknown)

	return result
}
