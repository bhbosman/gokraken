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
	"github.com/bhbosman/goFxApp/Services/fileDumpService"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	krakenWsStream "github.com/bhbosman/gokraken/internal/krakenWS/internal/stream"
	"github.com/bhbosman/gomessageblock"
	"github.com/cskr/pubsub"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
	"hash/crc32"
	"strconv"
)

type reactor struct {
	common.BaseConnectionReactor
	messageRouter             *messageRouter.MessageRouter
	connectionID              uint64
	status                    string
	version                   string
	registeredSubscriptionMap map[subscriptionKey]registeredSubscription
	FmdService                fullMarketDataManagerService.IFmdManagerService
	otherData                 instrumentReference.KrakenReferenceData
	FileDumpService           fileDumpService.IFileDumpService
}

func (self *reactor) handleKrakenStreamSubscribe(_ *Subscribe) error {
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

func (self *reactor) handleKrakenWsMessageIncoming(inData *krakenWsStream.KrakenWsMessageIncoming) {
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

func (self *reactor) handleWebsocketDataResponse(inData websocketDataResponse) {

	switch {
	case len(inData) == 4:
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
	case len(inData) == 5:
		if pair, ok := inData[4].(string); ok {
			if name, ok := inData[3].(string); ok {
				key := subscriptionKey{
					pair: pair,
					name: name,
				}
				if v, ok := self.registeredSubscriptionMap[key]; ok {
					switch self.otherData.Type {
					case "book":
						self.HandleBook(v, inData[1].(map[string]interface{}))
						self.HandleBook(v, inData[2].(map[string]interface{}))
					}
				}
			}
		}
	}
}

type websocketDataResponse []interface{}

func (self *reactor) handleWebSocketMessage(inData *wsmsg.WebSocketMessage) {
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
			Unmarshall := jsonpb.Unmarshaler{
				AllowUnknownFields: true,
				AnyResolver:        nil,
			}
			err := Unmarshall.Unmarshal(bytes.NewBuffer(inData.Message), krakenMessage)
			if err != nil {
				return
			}
			self.messageRouter.Route(krakenMessage)
			return
		}
	default:
	}
}

func (self *reactor) Init(params intf.IInitParams) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
	_, _, _, err := self.BaseConnectionReactor.Init(params)
	if err != nil {
		return nil, nil, nil, err
	}

	return func(i interface{}) {
			self.messageRouter.Route(i)
		},
		func(err error) {
			self.messageRouter.Route(err)
		},
		func() {

		},
		nil
}

func (self *reactor) Close() error {
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
func (self *reactor) Open() error {
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

func (self *reactor) handleSystemStatus(data krakenWsStream.ISystemStatus) {
	self.status = data.GetStatus()
	self.connectionID = data.GetConnectionID()
	self.version = data.GetVersion()
}

func (self *reactor) handlePing(data krakenWsStream.IPing) {
	_ = self.SendTextOpMessage(
		&krakenWsStream.KrakenWsMessageOutgoing{
			Event: "pong",
			Reqid: data.GetReqid(),
		},
	)
}

func (self *reactor) handleHeartbeat(_ interface{}) {
	return
}

func (self *reactor) handleSubscriptionStatus(inData krakenWsStream.ISubscriptionStatus) {
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

func (self *reactor) HandlePublishRxHandlerCounters(_ *model.PublishRxHandlerCounters) {}

func (self *reactor) HandleEmptyQueue(_ *messages.EmptyQueue) {
}

func (self *reactor) Unknown(_ interface{}) {

}

func (self *reactor) SendTextOpMessage(
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

func (self *reactor) wsProcessOrderBookPartial(
	instrumentData registeredSubscription,
	askData, bidData []interface{},
	snapShot bool,
) []interface{} {
	result := make([]interface{}, 0, len(askData)*2+len(bidData)*2+1)
	if snapShot {
		result = append(
			result,
			&stream2.FullMarketData_Clear{
				Instrument: instrumentData.SystemName,
			},
		)
	}
	result = append(result, self.iterateOrderData(instrumentData, "ASK", askData, stream2.OrderSide_AskOrder)...)
	result = append(result, self.iterateOrderData(instrumentData, "BID", bidData, stream2.OrderSide_BidOrder)...)

	return result
}

func (self *reactor) iterateOrderData(
	instrumentData registeredSubscription,
	orderIdPrefix string, data []interface{}, orderSide stream2.OrderSide) []interface{} {
	result := make([]interface{}, 0, len(data)*2)
	for i := range data {
		lineData := data[i].([]interface{})
		priceStr, _ := lineData[0].(string)
		volumeStr, _ := lineData[1].(string)
		priceAsFloat, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			self.CancelFunc()
			self.Logger.Error("could not read price as float", zap.Error(err))
			return nil
		}
		volumeAsFloat, err := strconv.ParseFloat(volumeStr, 64)
		if err != nil {
			self.CancelFunc()
			self.Logger.Error("could not read volume as float", zap.Error(err))
			return nil
		}

		orderId := fmt.Sprintf("%v_%v", orderIdPrefix, priceStr)
		result = append(
			result,
			&stream2.FullMarketData_DeleteOrderInstruction{
				Instrument: instrumentData.SystemName,
				Id:         orderId,
			},
		)
		if volumeStr != "0.00000000" {
			buffer := bytes.Buffer{}
			preCrc(&buffer, priceStr)
			preCrc(&buffer, volumeStr)
			result = append(
				result,
				&stream2.FullMarketData_AddOrderInstruction{
					Instrument: instrumentData.SystemName,
					Order: &stream2.FullMarketData_AddOrder{
						Side:      orderSide,
						Id:        orderId,
						Price:     priceAsFloat,
						Volume:    volumeAsFloat,
						ExtraData: buffer.Bytes(),
					},
				},
			)
		}
	}
	return result
}

func (self *reactor) HandleBook(
	instrumentData registeredSubscription,
	data map[string]interface{},
) {
	askSnapshot, askSnapshotExists := data["as"].([]interface{})
	bidSnapshot, bidSnapshotExists := data["bs"].([]interface{})
	askData, asksExist := data["a"].([]interface{})
	bidData, bidsExist := data["b"].([]interface{})
	checkSumData, checkSumExist := data["c"].(string)

	var messagesToFmd []interface{}
	if askSnapshotExists || bidSnapshotExists {
		messagesToFmd = append(
			messagesToFmd,
			self.wsProcessOrderBookPartial(instrumentData, askSnapshot, bidSnapshot, true)...,
		)
	} else if asksExist || bidsExist {
		messagesToFmd = append(
			messagesToFmd,
			self.wsProcessOrderBookPartial(instrumentData, askData, bidData, false)...,
		)
	}
	messagesToFmd = append(
		messagesToFmd,
		fullMarketDataManagerService.NewCallbackMessage(
			instrumentData.SystemName,
			self.doManageOrderBookSize,
			instrumentData,
		),
	)
	if checkSumExist {
		if checkSumAsInt, err := strconv.Atoi(checkSumData); err == nil {
			messagesToFmd = append(
				messagesToFmd,
				fullMarketDataManagerService.NewCallbackMessage(
					instrumentData.SystemName,
					self.doCrcCheck,
					&callbackData{
						inputCrc:   uint32(checkSumAsInt),
						data:       data,
						SystemName: instrumentData.SystemName,
					},
				),
			)
		}
	}
	self.FmdService.MultiSend(messagesToFmd...)
}

func (self *reactor) doManageOrderBookSize(cbData interface{}, fullMarketOrderBook fullMarketData.IFullMarketOrderBook) []interface{} {
	if instrumentData, ok := cbData.(registeredSubscription); ok {
		var result []interface{}
		walkTree := func(tree *avltree.Tree, firstNode func(tree *avltree.Tree) *avltree.Node, nextNode func(node *avltree.Node) *avltree.Node) {
			for tree.Size() > self.otherData.Depth {
				key := firstNode(tree).Key
				keyAsFloat := key.(float64)
				if unk01, ok := tree.Get(keyAsFloat); ok {
					pp := unk01.(*fullMarketData.PricePoint)
					if unk02, ok := pp.List.Get(0); ok {
						order := unk02.(*fullMarketData.FullMarketOrder)
						result = append(result,
							&stream2.FullMarketData_DeleteOrderInstruction{
								Instrument: instrumentData.SystemName,
								Id:         order.Id,
							},
						)
					} else {
						self.Logger.Error("could not find order at index 0", zap.Float64("key", keyAsFloat))
						self.CancelFunc()
					}
				} else {
					self.Logger.Error("could not find order to delete", zap.Float64("key", keyAsFloat))
					self.CancelFunc()
				}
				tree.Remove(key)
			}
		}
		walkTree(
			fullMarketOrderBook.BidOrderSide(),
			func(tree *avltree.Tree) *avltree.Node {
				return tree.Left()
			},
			func(node *avltree.Node) *avltree.Node {
				return node.Next()
			},
		)
		walkTree(
			fullMarketOrderBook.AskOrderSide(),
			func(tree *avltree.Tree) *avltree.Node {
				return tree.Right()
			},
			func(node *avltree.Node) *avltree.Node {
				return node.Prev()
			},
		)

		return result
	}
	return nil
}

func (self *reactor) doCrcCheck(cbData interface{}, fullMarketOrderBook fullMarketData.IFullMarketOrderBook) []interface{} {
	if v, ok := cbData.(*callbackData); ok {
		crc := crc32.NewIEEE()
		walkTree := func(
			tree *avltree.Tree,
			firstNode func(tree *avltree.Tree) *avltree.Node,
			nextNode func(node *avltree.Node) *avltree.Node,
		) {
			var count uint32 = 0
			for node := firstNode(tree); node != nil && count < 10; node = nextNode(node) {
				pp := node.Value.(*fullMarketData.PricePoint)
				unk, _ := pp.List.Get(0)
				ss := unk.(*fullMarketData.FullMarketOrder)
				_, _ = crc.Write(ss.ExtraData)
				count++
			}
		}
		walkTree(
			fullMarketOrderBook.AskOrderSide(),
			func(tree *avltree.Tree) *avltree.Node {
				return tree.Left()
			},
			func(node *avltree.Node) *avltree.Node {
				return node.Next()
			},
		)
		walkTree(
			fullMarketOrderBook.BidOrderSide(),
			func(tree *avltree.Tree) *avltree.Node {
				return tree.Right()
			},
			func(node *avltree.Node) *avltree.Node {
				return node.Prev()
			},
		)
		calculatedCrc := crc.Sum32()
		if calculatedCrc != v.inputCrc {
			self.Logger.Error("Closing connection due to CRC error")
			self.CancelFunc()
		}
	}
	return nil
}

type callbackData struct {
	inputCrc   uint32
	data       map[string]interface{}
	SystemName string
}

func newReactor(
	logger *zap.Logger,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc model.ConnectionCancelFunc,
	PubSub *pubsub.PubSub,
	goFunctionCounter GoFunctionCounter.IService,
	UniqueReferenceService interfaces.IUniqueReferenceService,
	FmdService fullMarketDataManagerService.IFmdManagerService,
	otherData instrumentReference.KrakenReferenceData,
	FileDumpService fileDumpService.IFileDumpService,
) intf.IConnectionReactor {
	result := &reactor{
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
		FileDumpService:           FileDumpService,
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

func preCrc(writer *bytes.Buffer, s string) {
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
