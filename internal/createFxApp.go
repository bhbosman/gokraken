package internal

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerViewer"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	"github.com/bhbosman/goCommsMultiDialer"
	"github.com/bhbosman/goFxApp"
	app2 "github.com/bhbosman/gocommon/Providers"
	"github.com/bhbosman/gokraken/internal/krakenConfiguration"
	"github.com/bhbosman/gokraken/internal/listener"
	"go.uber.org/fx"
)

func CreateFxApp() *goFxApp.TerminalAppUsingFxApp {
	settings := &AppSettings{
		compressedListenerUrl: "tcp4://127.0.0.1:3011",
	}

	var shutDowner fx.Shutdowner
	return goFxApp.NewFxMainApplicationServices(
		"kraken-stream",
		false,
		fx.Supply(settings),
		fx.Populate(&shutDowner),
		app2.RegisterRunTimeManager(),
		fullMarketDataManagerViewer.Provide(),
		fullMarketDataManagerService.Provide(false),
		fullMarketDataHelper.Provide(),
		instrumentReference.Provide(),
		goCommsMultiDialer.Provide(),
		krakenConfiguration.Provide(),
		listener.CompressedListener(1024, settings.compressedListenerUrl),
	)
}
