package listener

import (
	"context"
	"github.com/bhbosman/goCommsNetDialer"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/intf"
	"github.com/cskr/pubsub"
	"go.uber.org/zap"
)

type Factory struct {
	pubSub          *pubsub.PubSub
	SerializeData   SerializeData
	ConsumerCounter *goCommsNetDialer.CanDialDefaultImpl
}

func (self *Factory) Values(_ map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	return result, nil
}

func (self *Factory) Create(
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc model.ConnectionCancelFunc,
	logger *zap.Logger,
	userContext interface{},
) (intf.IConnectionReactor, error) {
	return NewReactor(
			logger,
			cancelCtx,
			cancelFunc,
			connectionCancelFunc,
			userContext,
			self.ConsumerCounter,
			self.SerializeData,
			self.pubSub),
		nil
}

func NewFactory(
	pubSub *pubsub.PubSub,
	SerializeData SerializeData,
	ConsumerCounter *goCommsNetDialer.CanDialDefaultImpl,
) (*Factory, error) {
	fac := &Factory{
		pubSub:          pubSub,
		SerializeData:   SerializeData,
		ConsumerCounter: ConsumerCounter,
	}
	return fac, nil
}
