package krakenConfiguration

import (
	"github.com/bhbosman/gocommon/services/IDataShutDown"
	"github.com/bhbosman/gocommon/services/IFxService"
	"github.com/bhbosman/gocommon/services/ISendMessage"
)

type IKrakenConfiguration interface {
	ISendMessage.ISendMessage
	//GetAll() []*krakenWS.KrakenConnection
}

type IKrakenConfigurationService interface {
	IKrakenConfiguration
	IFxService.IFxServices
}

type IKrakenConfigurationData interface {
	IKrakenConfiguration
	IDataShutDown.IDataShutDown
	ISendMessage.ISendMessage
}
