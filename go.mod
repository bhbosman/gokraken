module github.com/bhbosman/gokraken

go 1.24.0

require (
	github.com/bhbosman/goCommonMarketData v0.0.0-20250308054451-6dade8c5c3f9
	github.com/bhbosman/goCommsDefinitions v0.0.0-20250308000247-4306925b3dfd
	github.com/bhbosman/goCommsMultiDialer v0.0.0-20250308053512-71bb6f373484
	github.com/bhbosman/goCommsNetListener v0.0.0-20250308054504-9e394f3cd581
	github.com/bhbosman/goCommsStacks v0.0.0-20250308055233-8992b9d0c219
	github.com/bhbosman/goConn v0.0.0-20250307235008-177f4ffe3521
	github.com/bhbosman/goFxApp v0.0.0-20250308052840-c690d60e479d
	github.com/bhbosman/goFxAppManager v0.0.0-20250308052839-63033ea503a7
	github.com/bhbosman/gocommon v0.0.0-20250308052839-0ebeb121f996
	github.com/bhbosman/gocomms v0.0.0-20250308053431-cfd3e830d8a5
	github.com/bhbosman/goerrors v0.0.0-20250307194237-312d070c8e38
	github.com/bhbosman/gomessageblock v0.0.0-20250307141417-ab783e8e2eba
	github.com/bhbosman/goprotoextra v0.0.2
	github.com/cskr/pubsub v1.0.2
	github.com/emirpasic/gods v1.18.1
	github.com/golang/protobuf v1.5.4
	github.com/reactivex/rxgo/v2 v2.5.0
	go.uber.org/fx v1.23.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.27.0
	golang.org/x/net v0.37.0
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/bhbosman/goCommsNetDialer v0.0.0-20250308053613-a763632f1b3d // indirect
	github.com/bhbosman/goConnectionManager v0.0.0-20250308054322-587c786bedb0 // indirect
	github.com/bhbosman/goMessages v0.0.0-20250308052839-b11fc966f96a // indirect
	github.com/bhbosman/goUi v0.0.0-20250308064519-26f29e163722 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gdamore/encoding v1.0.1 // indirect
	github.com/gdamore/tcell/v2 v2.8.1 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/icza/gox v0.2.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/tview v0.0.0-20241227133733-17b7edb88c57 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/teivah/onecontext v1.3.0 // indirect
	go.uber.org/dig v1.18.1 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/cskr/pubsub => github.com/bhbosman/pubsub v1.0.3-0.20220802200819-029949e8a8af
	github.com/gdamore/tcell/v2 => github.com/bhbosman/tcell/v2 v2.5.2-0.20220624055704-f9a9454fab5b
	github.com/golang/mock => github.com/bhbosman/gomock v1.6.1-0.20250308071159-4cf72f668c72
	github.com/rivo/tview => github.com/bhbosman/tview v0.0.0-20250308051327-a656c1bc9cfa
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
