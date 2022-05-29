package internal

import (
	"github.com/bhbosman/gocommon/FxWrappers"
	app2 "github.com/bhbosman/gocommon/Providers"
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

func CreateFxApp() *FxWrappers.TerminalAppUsingFxApp {
	settings := &AppSettings{
		Logger:                log.New(os.Stderr, "", log.LstdFlags),
		textListenerUrl:       "tcp4://127.0.0.1:3010",
		compressedListenerUrl: "tcp4://127.0.0.1:3011",
		HttpListenerUrl:       "http://127.0.0.1:8081",
	}

	ConsumerCounter := netDial.NewCanDialDefaultImpl()
	var shutDowner fx.Shutdowner
	return FxWrappers.NewFxMainApplicationServices(
		"KrakenStream",
		false,
		fx.Supply(settings, ConsumerCounter),
		gologging.ProvideLogFactory(settings.Logger, nil),
		fx.Populate(&shutDowner),
		app2.RegisterRunTimeManager(),
		provide.RegisterHttpHandler(settings.HttpListenerUrl),
		endpoints.RegisterConnectionManagerEndpoint(),
		view.RegisterConnectionsHtmlTemplate(),
		connection.ProvideKrakenDialer(0, 0, ConsumerCounter),
		listener.TextListener(0, 0, ConsumerCounter, 1024, settings.textListenerUrl),
		listener.CompressedListener(0, 0, ConsumerCounter, 1024, settings.compressedListenerUrl),
	)
}
