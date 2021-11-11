package internal

import (
	app2 "github.com/bhbosman/gocommon/app"
	"github.com/bhbosman/gocommon/fxHelper"
	"github.com/bhbosman/gocommon/logSettings"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/connectionManager/endpoints"
	"github.com/bhbosman/gocomms/connectionManager/view"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/gocomms/provide"
	"github.com/bhbosman/gokraken/internal/krakenWS/connection"
	"github.com/bhbosman/gokraken/internal/listener"
	"github.com/bhbosman/gologging"
	"go.uber.org/fx"
	"log"
	"os"
)

func CreateFxApp() (*fx.App, fx.Shutdowner) {
	settings := &AppSettings{
		Logger:                log.New(os.Stderr, "", log.LstdFlags),
		textListenerUrl:       "tcp4://127.0.0.1:3010",
		compressedListenerUrl: "tcp4://127.0.0.1:3011",
		HttpListenerUrl:       "http://127.0.0.1:8081",
	}

	ConsumerCounter := netDial.NewCanDialDefaultImpl()
	var shutDowner fx.Shutdowner
	fxApp := fx.New(
		fx.Supply(settings, ConsumerCounter),
		logSettings.ProvideZapConfig(),

		gologging.ProvideLogFactory(settings.Logger, nil),
		fx.Populate(&shutDowner),
		app2.RegisterRootContext(),
		connectionManager.RegisterDefaultConnectionManager(),
		provide.RegisterHttpHandler(settings.HttpListenerUrl),
		endpoints.RegisterConnectionManagerEndpoint(),
		view.RegisterConnectionsHtmlTemplate(),
		connection.ProvideKrakenDialer(ConsumerCounter),
		listener.TextListener(ConsumerCounter, 1024, settings.textListenerUrl),
		listener.CompressedListener(ConsumerCounter, 1024, settings.compressedListenerUrl),
		fxHelper.InvokeApps(),
	)
	return fxApp, shutDowner
}
