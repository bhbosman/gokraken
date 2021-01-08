package listener

import (
	"context"
	"fmt"
	marketDataStream "github.com/bhbosman/goMessages/marketData/stream"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/gologging"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"github.com/reactivex/rxgo/v2"
	"google.golang.org/protobuf/proto"
	"net"
	"net/url"
	"strings"
)

type SerializeData func(m proto.Message) (goprotoextra.IReadWriterSize, error)
type Reactor struct {
	impl.BaseConnectionReactor
	ConsumerCounter      *netDial.CanDialDefaultImpl
	messageRouter        *messageRouter.MessageRouter
	SerializeData        SerializeData
	republishChannelName string
	publishChannelName   string
	dirtyMap             map[string]*marketDataStream.PublishTop5
	PubSub               *pubsub.PubSub
	messageOut           int
}

func (self *Reactor) Init(
	conn net.Conn,
	url *url.URL,
	connectionId string,
	connectionManager connectionManager.IConnectionManager,
	toConnectionFunc goprotoextra.ToConnectionFunc,
	toConnectionReactor goprotoextra.ToReactorFunc) (rxgo.NextExternalFunc, error) {
	_, err := self.BaseConnectionReactor.Init(conn, url, connectionId, connectionManager, toConnectionFunc, toConnectionReactor)
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

func (self *Reactor) HandleEmptyQueue(top5 *rxgo.EmptyQueue) error {
	for _, v := range self.dirtyMap {
		self.messageOut++
		marshal, err := self.SerializeData(v)
		if err != nil {
			return err
		}
		_ = self.ToConnection(marshal)
	}
	s := fmt.Sprintf("\n\r-->%v<--\r\n", self.messageOut)
	rws := gomessageblock.NewReaderWriter()
	rws.Write([]byte(s))
	_ = self.ToConnection(rws)
	self.dirtyMap = make(map[string]*marketDataStream.PublishTop5)
	return nil
}

func (self *Reactor) HandleTop5(top5 *marketDataStream.PublishTop5) error {
	if self.CancelCtx.Err() != nil {
		return self.CancelCtx.Err()
	}
	top5.Source = "KrakenWS"
	s := strings.Replace(fmt.Sprintf("%v.%v", top5.Source, top5.Instrument), "/", ".", -1)
	top5.UniqueName = s
	self.dirtyMap[s] = top5

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
		dirtyMap:              make(map[string]*marketDataStream.PublishTop5),
		PubSub:                PubSub,
	}
	result.messageRouter.Add(result.HandleTop5)
	result.messageRouter.Add(result.HandleEmptyQueue)
	return result
}
