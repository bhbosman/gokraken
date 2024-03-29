package krakenConfiguration

import (
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gokraken/internal/krakenWS"
)

type data struct {
	m             map[string]*krakenWS.KrakenConnection
	isDirty       map[string]bool
	MessageRouter messageRouter.IMessageRouter
}

func (self *data) GetAll() []*krakenWS.KrakenConnection {
	var result []*krakenWS.KrakenConnection
	for _, lunoConfiguration := range self.m {
		result = append(result, lunoConfiguration)
	}
	return result
}

func (self *data) Send(message interface{}) error {
	self.MessageRouter.Route(message)
	return nil
}

func (self *data) SomeMethod() {
}

func (self *data) ShutDown() error {
	return nil
}

func (self *data) handleEmptyQueue(*messages.EmptyQueue) {
	self.isDirty = make(map[string]bool)
}

func (self *data) handleLunoConfiguration(msg *krakenWS.KrakenConnection) {
	if msg.Name != "" {
		self.isDirty[msg.Name] = true
		self.m[msg.Name] = msg
	}
}

func newData() (IKrakenConfigurationData, error) {
	result := &data{
		MessageRouter: messageRouter.NewMessageRouter(),
		m:             make(map[string]*krakenWS.KrakenConnection),
		isDirty:       make(map[string]bool),
	}
	_ = result.MessageRouter.Add(result.handleEmptyQueue)
	_ = result.MessageRouter.Add(result.handleLunoConfiguration)
	return result, nil
}
