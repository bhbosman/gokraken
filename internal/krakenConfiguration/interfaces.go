package krakenConfiguration

import (
	"github.com/bhbosman/gocommon/Services/IDataShutDown"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/services/ISendMessage"
	"github.com/bhbosman/gokraken/internal/krakenWS"
)

type IKrakenConfiguration interface {
	ISendMessage.ISendMessage
	GetAll() []*krakenWS.KrakenConnection
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
