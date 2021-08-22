package listener

import (
	"context"
	"fmt"
	marketDataStream "github.com/bhbosman/goMessages/marketData/stream"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/gologging"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"google.golang.org/protobuf/proto"
	"net/url"
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
	impl.BaseConnectionReactor
	ConsumerCounter      *netDial.CanDialDefaultImpl
	messageRouter        *messageRouter.MessageRouter
	SerializeData        SerializeData
	republishChannelName string
	publishChannelName   string
	dirtyMap             map[string]*DirtyMapData
	PubSub               *pubsub.PubSub
	messageOut           int
}

func (self *Reactor) Init(
	url *url.URL,
	connectionId string,
	connectionManager connectionManager.IConnectionManager__,
	toConnectionFunc goprotoextra.ToConnectionFunc,
	toConnectionReactor goprotoextra.ToReactorFunc) (intf.NextExternalFunc, error) {
	_, err := self.BaseConnectionReactor.Init(url, connectionId, connectionManager, toConnectionFunc, toConnectionReactor)
	if err != nil {
		return nil, err
	}

	self.republishChannelName = "republishChannel"
	self.publishChannelName = "publishChannel"

	ch := self.PubSub.Sub(self.publishChannelName)
	go func(ch chan interface{}, topics ...string) {
		defer self.PubSub.Unsub(ch, topics...)
		<-self.CancelCtx.Done()
	}(ch, self.publishChannelName)

	go func(ch chan interface{}, topics ...string) {
		for v := range ch {
			if self.CancelCtx.Err() == nil {
				_ = self.ToReactor(true, v)
			}
		}
	}(ch, self.publishChannelName)

	self.PubSub.Pub(&struct{}{}, self.republishChannelName)

	return self.doNext, nil
}

func (self *Reactor) doNext(external bool, i interface{}) {
	_, _ = self.messageRouter.Route(i)
}

func (self *Reactor) Open() error {
	self.ConsumerCounter.AddConsumer()
	return self.BaseConnectionReactor.Open()
}

func (self *Reactor) Close() error {
	self.ConsumerCounter.RemoveConsumer()
	return self.BaseConnectionReactor.Close()
}

func (self *Reactor) HandleEmptyQueue(top5 *RxHandlers.EmptyQueue) error {
	var deleteKeys []string
	for k, v := range self.dirtyMap {
		self.messageOut++
		marshal, err := self.SerializeData(v.data)
		if err != nil {
			return err
		}
		_ = self.ToConnection(marshal)
		deleteKeys = append(deleteKeys, k)
	}
	for _, s := range deleteKeys {
		delete(self.dirtyMap, s)
	}

	s := fmt.Sprintf("\n\r-->%v<--\r\n", self.messageOut)
	rws := gomessageblock.NewReaderWriter()
	rws.Write([]byte(s))
	_ = self.ToConnection(rws)
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
	logger *gologging.SubSystemLogger,
	name string,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	userContext interface{},
	ConsumerCounter *netDial.CanDialDefaultImpl,
	SerializeData SerializeData,
	PubSub *pubsub.PubSub) *Reactor {
	result := &Reactor{
		BaseConnectionReactor: impl.NewBaseConnectionReactor(logger, name, cancelCtx, cancelFunc, userContext),
		ConsumerCounter:       ConsumerCounter,
		messageRouter:         messageRouter.NewMessageRouter(),
		SerializeData:         SerializeData,
		dirtyMap:              make(map[string]*DirtyMapData),
		PubSub:                PubSub,
	}
	result.messageRouter.Add(result.HandlePublishTop5)
	result.messageRouter.Add(result.HandleEmptyQueue)
	return result
}
