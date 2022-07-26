package internal

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerViewer"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	"github.com/bhbosman/goFxApp"
	app2 "github.com/bhbosman/gocommon/Providers"
	"github.com/bhbosman/gokraken/internal/krakenWS"
	"github.com/bhbosman/gokraken/internal/listener"
	"go.uber.org/fx"
	"log"
	"os"
)

func CreateFxApp() *goFxApp.TerminalAppUsingFxApp {
	settings := &AppSettings{
		Logger:                log.New(os.Stderr, "", log.LstdFlags),
		textListenerUrl:       "tcp4://127.0.0.1:3010",
		compressedListenerUrl: "tcp4://127.0.0.1:3011",
	}

	var shutDowner fx.Shutdowner
	return goFxApp.NewFxMainApplicationServices(
		"KrakenStream",
		false,
		fx.Supply(settings),
		fx.Populate(&shutDowner),
		app2.RegisterRunTimeManager(),
		fullMarketDataManagerViewer.Provide(),
		fullMarketDataManagerService.Provide(false),
		fullMarketDataHelper.Provide(),
		instrumentReference.Provide(),

		krakenWS.ProvideKrakenDialer(),
		listener.TextListener(1024, settings.textListenerUrl),
		listener.CompressedListener(1024, settings.compressedListenerUrl),
	)
}
