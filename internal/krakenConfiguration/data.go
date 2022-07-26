package krakenConfiguration

import (
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gokraken/internal/krakenWS"
)

type data struct {
	m             map[string]*krakenWS.KrakenConnection
	isDirty       map[string]bool
	MessageRouter *messageRouter.MessageRouter
}

func (self *data) GetAll() []*krakenWS.KrakenConnection {
	var result []*krakenWS.KrakenConnection
	for _, lunoConfiguration := range self.m {
		result = append(result, lunoConfiguration)
	}
	return result
}

func (self *data) Send(message interface{}) error {
	_, err := self.MessageRouter.Route(message)
	return err
}

func (self *data) SomeMethod() {
}

func (self *data) ShutDown() error {
	return nil
}

func (self *data) handleEmptyQueue(msg *messages.EmptyQueue) {
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
	result.MessageRouter.Add(result.handleEmptyQueue)
	result.MessageRouter.Add(result.handleLunoConfiguration)
	return result, nil
}
