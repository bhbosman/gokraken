package internal

import (
	app2 "github.com/bhbosman/gocommon/app"
	"github.com/bhbosman/gocommon/fxHelper"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/connectionManager/endpoints"
	"github.com/bhbosman/gocomms/connectionManager/view"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/gocomms/provide"
	"github.com/bhbosman/gokraken/internal/krakenWS/connection"
	"github.com/bhbosman/gokraken/internal/listener"
	"github.com/bhbosman/gologging"
	"github.com/cskr/pubsub"
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
	pubSub := pubsub.New(32)

	ConsumerCounter := netDial.NewCanDialDefaultImpl()
	var shutDowner fx.Shutdowner
	fxApp := fx.New(
		fx.Supply(settings, ConsumerCounter),
		fx.Logger(settings.Logger),
		gologging.ProvideLogFactory(settings.Logger, nil),
		fx.Populate(&shutDowner),
		app2.RegisterRootContext(pubSub),
		connectionManager.RegisterDefaultConnectionManager(),
		provide.RegisterHttpHandler(settings.HttpListenerUrl),
		endpoints.RegisterConnectionManagerEndpoint(),
		view.RegisterConnectionsHtmlTemplate(),
		connection.ProvideKrakenDialer(pubSub, ConsumerCounter),
		listener.TextListener(pubSub, ConsumerCounter, 1024, settings.textListenerUrl),
		listener.CompressedListener(pubSub, ConsumerCounter, 1024, settings.compressedListenerUrl),
		fxHelper.InvokeApps(),
	)
	return fxApp, shutDowner
}
