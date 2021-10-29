package connection

import (
	"context"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/cskr/pubsub"
	"go.uber.org/zap"
)

type Factory struct {
	crfName string
	PubSub  *pubsub.PubSub
}

func (self *Factory) Name() string {
	return self.crfName
}

func (self Factory) Values(inputValues map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	return result, nil
}

func (self Factory) Create(
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc common.ConnectionCancelFunc,
	logger *zap.Logger,
	userContext interface{}) intf.IConnectionReactor {
	return NewReactor(
		logger,
		cancelCtx,
		cancelFunc,
		connectionCancelFunc,
		userContext,
		self.PubSub)
}

func NewFactory(
	crfName string,
	PubSub *pubsub.PubSub) intf.IConnectionReactorFactory {
	return &Factory{
		crfName: crfName,
		PubSub:  PubSub,
	}
}
