package listener

import (
	"context"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/gologging"
	"github.com/cskr/pubsub"
)

type Factory struct {
	name            string
	pubSub          *pubsub.PubSub
	SerializeData   SerializeData
	ConsumerCounter *netDial.CanDialDefaultImpl
}

func (self *Factory) Values(inputValues map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	return result, nil
}

func (self *Factory) Name() string {
	return self.name
}

func (self *Factory) Create(
	name string, cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	logger *gologging.SubSystemLogger,
	userContext interface{}) intf.IConnectionReactor {
	return NewReactor(logger, name, cancelCtx, cancelFunc, userContext, self.ConsumerCounter, self.SerializeData, self.pubSub)
}

func NewFactory(
	name string,
	pubSub *pubsub.PubSub,
	SerializeData SerializeData,
	ConsumerCounter *netDial.CanDialDefaultImpl) *Factory {
	return &Factory{
		name:            name,
		pubSub:          pubSub,
		SerializeData:   SerializeData,
		ConsumerCounter: ConsumerCounter,
	}
}
