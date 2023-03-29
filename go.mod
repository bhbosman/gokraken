module github.com/bhbosman/gokraken

go 1.18

require (
	github.com/bhbosman/goCommonMarketData v0.0.0-20230329102141-a91b266b20e3
	github.com/bhbosman/goCommsDefinitions v0.0.0-20230329100608-a6a24c060ad8
	github.com/bhbosman/goCommsMultiDialer v0.0.0-20230329122920-e6c932f8e98f
	github.com/bhbosman/goCommsNetListener v0.0.0-20230329104212-9f1b0eafaa6b
	github.com/bhbosman/goCommsStacks v0.0.0-20230328221032-cd6c6063e9ef
	github.com/bhbosman/goConn v0.0.0-20230328145648-e6bab0984275
	github.com/bhbosman/goFxApp v0.0.0-20230329100615-cac401f3c57c
	github.com/bhbosman/goFxAppManager v0.0.0-20230329105958-8874cb25c628
	github.com/bhbosman/gocommon v0.0.0-20230329101749-40db0f52d859
	github.com/bhbosman/gocomms v0.0.0-20230329110556-946ebc6ff5f4
	github.com/bhbosman/goerrors v0.0.0-20220623084908-4d7bbcd178cf
	github.com/bhbosman/gomessageblock v0.0.0-20230308173223-e8144f25444c
	github.com/bhbosman/goprotoextra v0.0.2
	github.com/cskr/pubsub v1.0.2
	github.com/emirpasic/gods v1.18.1
	github.com/golang/protobuf v1.5.2
	github.com/reactivex/rxgo/v2 v2.5.0
	go.uber.org/fx v1.19.2
	go.uber.org/multierr v1.10.0
	go.uber.org/zap v1.24.0
	golang.org/x/net v0.7.0
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/bhbosman/goCommsNetDialer v0.0.0-20230329104213-5493957faab0 // indirect
	github.com/bhbosman/goConnectionManager v0.0.0-20230328221032-2e48eb2721d0 // indirect
	github.com/bhbosman/goMessages v0.0.0-20230328221033-d00082868706 // indirect
	github.com/bhbosman/goUi v0.0.0-20230328220052-834528d68f6d // indirect
	github.com/cenkalti/backoff/v4 v4.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/gdamore/tcell/v2 v2.5.1 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.1.0 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/icza/gox v0.0.0-20220321141217-e2d488ab2fbc // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/tview v0.0.0-20220709181631-73bf2902b59a // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	github.com/teivah/onecontext v0.0.0-20200513185103-40f981bfd775 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/dig v1.16.1 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/term v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/cskr/pubsub => github.com/bhbosman/pubsub v1.0.3-0.20220802200819-029949e8a8af
	github.com/gdamore/tcell/v2 => github.com/bhbosman/tcell/v2 v2.5.2-0.20220624055704-f9a9454fab5b
	github.com/golang/mock => github.com/bhbosman/gomock v1.6.1-0.20230302060806-d02c40b7514e
	github.com/rivo/tview => github.com/bhbosman/tview v0.0.0-20230310100135-f8b257a85d36
)

//replace (
//	github.com/bhbosman/goCommonMarketData => ../goCommonMarketData
//	github.com/bhbosman/goCommsDefinitions => ../goCommsDefinitions
//	github.com/bhbosman/goCommsMultiDialer => ../goCommsMultiDialer
//	github.com/bhbosman/goCommsNetDialer => ../goCommsNetDialer
//	github.com/bhbosman/goCommsNetListener => ../goCommsNetListener
//	github.com/bhbosman/goCommsStacks => ../goCommsStacks
//	github.com/bhbosman/goFxApp => ../goFxApp
//	github.com/bhbosman/goFxAppManager => ../goFxAppManager
//	github.com/bhbosman/goMessages => ../goMessages
//	github.com/bhbosman/gocommon => ../gocommon
//	github.com/bhbosman/gocomms => ../gocomms
//	github.com/bhbosman/gomessageblock => ../gomessageblock
//)
