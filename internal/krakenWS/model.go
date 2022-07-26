package krakenWS

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
