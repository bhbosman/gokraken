// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bhbosman/gokraken/internal/krakenConfiguration (interfaces: IKrakenConfiguration)

// Package krakenConfiguration is a generated GoMock package.
package krakenConfiguration

import (
	fmt "fmt"
	krakenConnectionService "github.com/bhbosman/gokraken/internal/krakenWS"

	errors "github.com/bhbosman/gocommon/errors"
	"golang.org/x/net/context"
)

// Interface A Comment
// Interface github.com/bhbosman/gokraken/internal/krakenConfiguration
// Interface IKrakenConfiguration
// Interface IKrakenConfiguration, Method: GetAll
type IKrakenConfigurationGetAllIn struct {
}

type IKrakenConfigurationGetAllOut struct {
	Args0 []*krakenConnectionService.KrakenConnection
}
type IKrakenConfigurationGetAllError struct {
	InterfaceName string
	MethodName    string
	Reason        string
}

func (self *IKrakenConfigurationGetAllError) Error() string {
	return fmt.Sprintf("error in data coming back from %v::%v. Reason: %v", self.InterfaceName, self.MethodName, self.Reason)
}

type IKrakenConfigurationGetAll struct {
	inData         IKrakenConfigurationGetAllIn
	outDataChannel chan IKrakenConfigurationGetAllOut
}

func NewIKrakenConfigurationGetAll(waitToComplete bool) *IKrakenConfigurationGetAll {
	var outDataChannel chan IKrakenConfigurationGetAllOut
	if waitToComplete {
		outDataChannel = make(chan IKrakenConfigurationGetAllOut)
	} else {
		outDataChannel = nil
	}
	return &IKrakenConfigurationGetAll{
		inData:         IKrakenConfigurationGetAllIn{},
		outDataChannel: outDataChannel,
	}
}

func (self *IKrakenConfigurationGetAll) Wait(onError func(interfaceName string, methodName string, err error) error) (IKrakenConfigurationGetAllOut, error) {
	data, ok := <-self.outDataChannel
	if !ok {
		generatedError := &IKrakenConfigurationGetAllError{
			InterfaceName: "IKrakenConfiguration",
			MethodName:    "GetAll",
			Reason:        "Channel for IKrakenConfiguration::GetAll returned false",
		}
		if onError != nil {
			err := onError("IKrakenConfiguration", "GetAll", generatedError)
			return IKrakenConfigurationGetAllOut{}, err
		} else {
			return IKrakenConfigurationGetAllOut{}, generatedError
		}
	}
	return data, nil
}

func (self *IKrakenConfigurationGetAll) Close() error {
	close(self.outDataChannel)
	return nil
}
func CallIKrakenConfigurationGetAll(context context.Context, channel chan<- interface{}, waitToComplete bool) (IKrakenConfigurationGetAllOut, error) {
	if context != nil && context.Err() != nil {
		return IKrakenConfigurationGetAllOut{}, context.Err()
	}
	data := NewIKrakenConfigurationGetAll(waitToComplete)
	if waitToComplete {
		defer func(data *IKrakenConfigurationGetAll) {
			err := data.Close()
			if err != nil {
			}
		}(data)
	}
	if context != nil && context.Err() != nil {
		return IKrakenConfigurationGetAllOut{}, context.Err()
	}
	channel <- data
	var err error
	var v IKrakenConfigurationGetAllOut
	if waitToComplete {
		v, err = data.Wait(func(interfaceName string, methodName string, err error) error {
			return err
		})
	} else {
		err = errors.NoWaitOperationError
	}
	if err != nil {
		return IKrakenConfigurationGetAllOut{}, err
	}
	return v, nil
}

// Interface IKrakenConfiguration, Method: Send
type IKrakenConfigurationSendIn struct {
	arg0 interface{}
}

type IKrakenConfigurationSendOut struct {
	Args0 error
}
type IKrakenConfigurationSendError struct {
	InterfaceName string
	MethodName    string
	Reason        string
}

func (self *IKrakenConfigurationSendError) Error() string {
	return fmt.Sprintf("error in data coming back from %v::%v. Reason: %v", self.InterfaceName, self.MethodName, self.Reason)
}

type IKrakenConfigurationSend struct {
	inData         IKrakenConfigurationSendIn
	outDataChannel chan IKrakenConfigurationSendOut
}

func NewIKrakenConfigurationSend(waitToComplete bool, arg0 interface{}) *IKrakenConfigurationSend {
	var outDataChannel chan IKrakenConfigurationSendOut
	if waitToComplete {
		outDataChannel = make(chan IKrakenConfigurationSendOut)
	} else {
		outDataChannel = nil
	}
	return &IKrakenConfigurationSend{
		inData: IKrakenConfigurationSendIn{
			arg0: arg0,
		},
		outDataChannel: outDataChannel,
	}
}

func (self *IKrakenConfigurationSend) Wait(onError func(interfaceName string, methodName string, err error) error) (IKrakenConfigurationSendOut, error) {
	data, ok := <-self.outDataChannel
	if !ok {
		generatedError := &IKrakenConfigurationSendError{
			InterfaceName: "IKrakenConfiguration",
			MethodName:    "Send",
			Reason:        "Channel for IKrakenConfiguration::Send returned false",
		}
		if onError != nil {
			err := onError("IKrakenConfiguration", "Send", generatedError)
			return IKrakenConfigurationSendOut{}, err
		} else {
			return IKrakenConfigurationSendOut{}, generatedError
		}
	}
	return data, nil
}

func (self *IKrakenConfigurationSend) Close() error {
	close(self.outDataChannel)
	return nil
}
func CallIKrakenConfigurationSend(context context.Context, channel chan<- interface{}, waitToComplete bool, arg0 interface{}) (IKrakenConfigurationSendOut, error) {
	if context != nil && context.Err() != nil {
		return IKrakenConfigurationSendOut{}, context.Err()
	}
	data := NewIKrakenConfigurationSend(waitToComplete, arg0)
	if waitToComplete {
		defer func(data *IKrakenConfigurationSend) {
			err := data.Close()
			if err != nil {
			}
		}(data)
	}
	if context != nil && context.Err() != nil {
		return IKrakenConfigurationSendOut{}, context.Err()
	}
	channel <- data
	var err error
	var v IKrakenConfigurationSendOut
	if waitToComplete {
		v, err = data.Wait(func(interfaceName string, methodName string, err error) error {
			return err
		})
	} else {
		err = errors.NoWaitOperationError
	}
	if err != nil {
		return IKrakenConfigurationSendOut{}, err
	}
	return v, nil
}

func ChannelEventsForIKrakenConfiguration(next IKrakenConfiguration, event interface{}) (bool, error) {
	switch v := event.(type) {
	case *IKrakenConfigurationGetAll:
		data := IKrakenConfigurationGetAllOut{}
		data.Args0 = next.GetAll()
		if v.outDataChannel != nil {
			v.outDataChannel <- data
		}
		return true, nil
	case *IKrakenConfigurationSend:
		data := IKrakenConfigurationSendOut{}
		data.Args0 = next.Send(v.inData.arg0)
		if v.outDataChannel != nil {
			v.outDataChannel <- data
		}
		return true, nil
	default:
		return false, nil
	}
}