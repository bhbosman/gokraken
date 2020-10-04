package stream

type ISystemStatus interface {
	GetConnectionID() uint64
	GetEvent() string
	GetStatus() string
	GetVersion() string
}

type IPing interface {
	GetReqid() uint32
}

type ISubscriptionStatus interface {
	GetChannelID() uint32
	GetChannelName() string
	GetEvent() string
	GetStatus() string
	GetReqid() uint32
	GetPair() string
	GetSubscription() *KrakenSubscriptionData
	GetErrorMessage() string
}
