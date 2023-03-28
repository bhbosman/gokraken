module github.com/bhbosman/gokraken

go 1.18

require (
	github.com/bhbosman/goCommonMarketData v0.0.0-20230302071905-68b90c2dabe4
	github.com/bhbosman/goCommsDefinitions v0.0.0-20230320101758-971a50fdbf8c
	github.com/bhbosman/goCommsMultiDialer v0.0.0-20230302092320-748285ee2b26
	github.com/bhbosman/goCommsNetListener v0.0.0-20230302062227-0285e8a8cd6f
	github.com/bhbosman/goCommsStacks v0.0.0-20230302185506-f2c4d856b6d3
	github.com/bhbosman/goConn v0.0.0-20230327111455-7a39299fb0aa
	github.com/bhbosman/goFxApp v0.0.0-20230302094801-5074560d188e
	github.com/bhbosman/goFxAppManager v0.0.0-20230320053418-125763bdf8d5
	github.com/bhbosman/gocommon v0.0.0-20230328140834-3b6f0ea2faaa
	github.com/bhbosman/gocomms v0.0.0-20230307212550-0918a992672c
	github.com/bhbosman/goerrors v0.0.0-20220623084908-4d7bbcd178cf
	github.com/bhbosman/gomessageblock v0.0.0-20230308173223-e8144f25444c
	github.com/bhbosman/goprotoextra v0.0.2
	github.com/cskr/pubsub v1.0.2
	github.com/emirpasic/gods v1.18.1
	github.com/golang/protobuf v1.5.0
	github.com/reactivex/rxgo/v2 v2.5.0
	go.uber.org/fx v1.19.2
	go.uber.org/multierr v1.10.0
	go.uber.org/zap v1.24.0
	golang.org/x/net v0.7.0
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/bhbosman/goCommsNetDialer v0.0.0-20230302062228-c88685956b15 // indirect
	github.com/bhbosman/goConnectionManager v0.0.0-20230328152128-966e33ba52f4 // indirect
	github.com/bhbosman/goMessages v0.0.0-20230302063433-258339efe599 // indirect
	github.com/bhbosman/goUi v0.0.0-20230327111003-b80e10f7c669 // indirect
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/gdamore/tcell/v2 v2.5.1 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.1.0 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/icza/gox v0.0.0-20220321141217-e2d488ab2fbc // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/tview v0.0.0-20220709181631-73bf2902b59a // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	github.com/teivah/onecontext v0.0.0-20200513185103-40f981bfd775 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/dig v1.16.1 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/term v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/bhbosman/gocommon => ../gocommon
	github.com/bhbosman/gocomms => ../gocomms
)

replace github.com/gdamore/tcell/v2 => github.com/bhbosman/tcell/v2 v2.5.2-0.20220624055704-f9a9454fab5b

replace github.com/golang/mock => github.com/bhbosman/gomock v1.6.1-0.20230302060806-d02c40b7514e

replace github.com/bhbosman/goCommsStacks => ../goCommsStacks

replace github.com/bhbosman/goCommsNetDialer => ../goCommsNetDialer

replace github.com/bhbosman/goCommsNetListener => ../goCommsNetListener

replace github.com/bhbosman/goCommsDefinitions => ../goCommsDefinitions

replace github.com/bhbosman/goFxApp => ../goFxApp

replace github.com/bhbosman/goUi => ../goUi

replace github.com/bhbosman/goerrors => ../goerrors

replace github.com/bhbosman/goFxAppManager => ../goFxAppManager

replace github.com/bhbosman/goConnectionManager => ../goConnectionManager

replace github.com/rivo/tview => ../tview

//replace github.com/bhbosman/goprotoextra => ../goprotoextra

replace github.com/bhbosman/goMessages => ../goMessages

replace github.com/cskr/pubsub => ../pubsub

replace github.com/bhbosman/goCommonMarketData => ../goCommonMarketData

replace github.com/bhbosman/goCommsMultiDialer => ../goCommsMultiDialer

replace github.com/reactivex/rxgo/v2 => ../goRx

replace github.com/bhbosman/goConn => ../goConn
