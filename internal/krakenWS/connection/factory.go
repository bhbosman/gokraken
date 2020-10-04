package connection

import (
	"context"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gologging"
	"github.com/cskr/pubsub"
)

type Factory struct {
	name        string
	PubSub      *pubsub.PubSub
}

func (self Factory) Name() string {
	return self.name
}

func (self Factory) Create(
	name string,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	logger *gologging.SubSystemLogger,
	userContext interface{}) intf.IConnectionReactor {
	return NewReactor(
		logger,
		name,
		cancelCtx,
		cancelFunc,
		userContext,
		self.PubSub)
}

func NewFactory(name string,
	PubSub *pubsub.PubSub) *Factory {
	return &Factory{
		name:        name,
		PubSub:      PubSub,
	}
}
