package connection

import (
	"context"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/intf"
	"github.com/cskr/pubsub"
	"go.uber.org/zap"
)

type Factory struct {
	crfName           string
	PubSub            *pubsub.PubSub
	goFunctionCounter GoFunctionCounter.IService
}

func (self *Factory) Name() string {
	return self.crfName
}

func (self Factory) Values(_ map[string]interface{}) (map[string]interface{}, error) {
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
			self.PubSub,
			self.goFunctionCounter,
		),
		nil
}

func NewFactory(
	crfName string,
	PubSub *pubsub.PubSub,
	goFunctionCounter GoFunctionCounter.IService,
) (intf.IConnectionReactorFactory, error) {
	fac := &Factory{
		crfName:           crfName,
		PubSub:            PubSub,
		goFunctionCounter: goFunctionCounter,
	}
	return fac, nil
}
