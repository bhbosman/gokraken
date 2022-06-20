package internal

import (
	"github.com/bhbosman/goCommsNetDialer"
	"github.com/bhbosman/gocommon/FxWrappers"
	app2 "github.com/bhbosman/gocommon/Providers"
	"github.com/bhbosman/gokraken/internal/krakenWS/connection"
	"github.com/bhbosman/gokraken/internal/listener"
	"go.uber.org/fx"
	"log"
	"os"
)

func CreateFxApp() *FxWrappers.TerminalAppUsingFxApp {
	settings := &AppSettings{
		Logger:                log.New(os.Stderr, "", log.LstdFlags),
		textListenerUrl:       "tcp4://127.0.0.1:3010",
		compressedListenerUrl: "tcp4://127.0.0.1:3011",
	}

	ConsumerCounter := goCommsNetDialer.NewCanDialDefaultImpl()
	var shutDowner fx.Shutdowner
	return FxWrappers.NewFxMainApplicationServices(
		"KrakenStream",
		false,
		fx.Supply(settings, ConsumerCounter),
		fx.Populate(&shutDowner),
		app2.RegisterRunTimeManager(),
		connection.ProvideKrakenDialer(0, 0, ConsumerCounter),
		listener.TextListener(0, 0, ConsumerCounter, 1024, settings.textListenerUrl),
		listener.CompressedListener(0, 0, ConsumerCounter, 1024, settings.compressedListenerUrl),
	)
}
