package listener

import (
	"context"
	"fmt"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	marketDataStream "github.com/bhbosman/goMessages/marketData/stream"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"strings"
)

type DirtyMapData struct {
	data *marketDataStream.PublishTop5
}

func NewDirtyMapData(data *marketDataStream.PublishTop5) *DirtyMapData {
	return &DirtyMapData{
		data: data,
	}
}

type SerializeData func(m proto.Message) (goprotoextra.IReadWriterSize, error)
type Reactor struct {
	common.BaseConnectionReactor
	messageRouter        *messageRouter.MessageRouter
	SerializeData        SerializeData
	republishChannelName string
	publishChannelName   string
	dirtyMap             map[string]*DirtyMapData
	messageOut           int
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

	self.PubSub.Pub(&struct{}{}, self.republishChannelName)

	return func(i interface{}) {
			self.doNext(false, i)
		},
		func(err error) {
			self.doNext(false, err)
		},
		func() {

		}, nil
}

func (self *Reactor) doNext(_ bool, i interface{}) {
	_, _ = self.messageRouter.Route(i)
}

func (self *Reactor) Open() error {
	return self.BaseConnectionReactor.Open()
}

func (self *Reactor) Close() error {
	return self.BaseConnectionReactor.Close()
}

func (self *Reactor) HandleEmptyQueue(_ *messages.EmptyQueue) error {
	var deleteKeys []string
	for k, v := range self.dirtyMap {
		self.messageOut++
		marshal, err := self.SerializeData(v.data)
		if err != nil {
			return err
		}
		self.OnSendToConnection(marshal)
		deleteKeys = append(deleteKeys, k)
	}
	for _, s := range deleteKeys {
		delete(self.dirtyMap, s)
	}

	s := fmt.Sprintf("\n\r-->%v<--\r\n", self.messageOut)
	rws := gomessageblock.NewReaderWriter()
	_, _ = rws.Write([]byte(s))
	self.OnSendToConnection(rws)
	return nil
}

func (self *Reactor) HandlePublishTop5(top5 *marketDataStream.PublishTop5) error {
	if self.CancelCtx.Err() != nil {
		return self.CancelCtx.Err()
	}
	top5.Source = "KrakenWS"
	s := strings.Replace(fmt.Sprintf("%v.%v", top5.Source, top5.Instrument), "/", ".", -1)
	top5.UniqueName = s
	if dirtyData, ok := self.dirtyMap[s]; ok {
		dirtyData.data = top5
	} else {
		self.dirtyMap[s] = NewDirtyMapData(top5)
	}
	return nil
}

func NewReactor(
	logger *zap.Logger,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc model.ConnectionCancelFunc,
	SerializeData SerializeData,
	PubSub *pubsub.PubSub,
	UniqueReferenceService interfaces.IUniqueReferenceService,
	GoFunctionCounter GoFunctionCounter.IService,
	FullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper,
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
			GoFunctionCounter,
		),
		messageRouter: messageRouter.NewMessageRouter(),
		SerializeData: SerializeData,
		dirtyMap:      make(map[string]*DirtyMapData),
	}
	result.messageRouter.Add(result.HandlePublishTop5)
	result.messageRouter.Add(result.HandleEmptyQueue)
	return result
}
