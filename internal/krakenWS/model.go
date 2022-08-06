package krakenWS

import (
	"fmt"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
)

type FeedRegistration struct {
	Pair string
	Name string
}

func NewInstance(pair string, name string) FeedRegistration {
	return FeedRegistration{
		Pair: pair,
		Name: name,
	}
}

type KrakenConnection struct {
	Name     string
	Instance []FeedRegistration
}

func NewKrakenConnection(name string, instance ...FeedRegistration) *KrakenConnection {
	return &KrakenConnection{
		Name:     name,
		Instance: instance,
	}
}

type Subscribe struct {
}

type registeredSubscription struct {
	instrumentReference.KrakenFeed
	status string
}

type subscriptionKey struct {
	pair string
	name string
}

func (self *subscriptionKey) String() string {
	return fmt.Sprintf("%v-%v", self.pair, self.name)
}
