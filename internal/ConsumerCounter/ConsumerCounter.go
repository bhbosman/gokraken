package ConsumerCounter

import "sync/atomic"

type ConsumerCounter struct {
	count int64
}

func (self *ConsumerCounter) CanDial() bool {
	return self.count > 0
}

func (self *ConsumerCounter) DialSuccess() {

}

func (self *ConsumerCounter) RemoveConsumer() {
	atomic.AddInt64(&self.count, -1)
}

func (self *ConsumerCounter) AddConsumer() {
	atomic.AddInt64(&self.count, 1)
}


