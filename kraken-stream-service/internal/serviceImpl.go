package internal

import (
	"github.com/bhbosman/gokraken/internal"
	"github.com/kardianos/service"
	"go.uber.org/fx"
)

type ServiceInterfaceImpl struct {
	app        *fx.App
	shutDowner fx.Shutdowner
}

func (self *ServiceInterfaceImpl) Start(srv service.Service) error {
	self.app, self.shutDowner = internal.CreateFxApp()
	if self.app.Err() != nil {
		return self.app.Err()
	}

	go func() {
		self.app.Run()
	}()
	return nil
}

func (self *ServiceInterfaceImpl) Stop(srv service.Service) error {
	return self.shutDowner.Shutdown()
}
