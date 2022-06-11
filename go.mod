module github.com/bhbosman/gokraken

go 1.15

require (
	github.com/bhbosman/goCommsNetListener v0.0.0-20220611182354-46d10d89b8e1 // indirect
	github.com/bhbosman/goCommsDefinitions v0.0.0-20220612125316-9a91bacb2b7c // indirect

	github.com/bhbosman/goCommsNetDialer v0.0.0-20220611181910-f23644b1b31a // indirect
	github.com/bhbosman/goCommsStacks v0.0.0-20220611100605-a342f0b42054
	github.com/bhbosman/goMessages v0.0.0-20210414134625-4d7166d206a6
	github.com/bhbosman/gocommon v0.0.0-20220608193411-270a4a7c70cc
	github.com/bhbosman/gocomms v0.0.0-20220611042959-112035f663a7
	github.com/bhbosman/goerrors v0.0.0-20210201065523-bb3e832fa9ab
	github.com/bhbosman/gologging v0.0.0-20200921180328-d29fc55c00bc
	github.com/bhbosman/gomessageblock v0.0.0-20210901070622-be36a3f8d303
	github.com/bhbosman/goprotoextra v0.0.2-0.20210817141206-117becbef7c7
	github.com/cskr/pubsub v1.0.2
	github.com/emirpasic/gods v1.12.0
	github.com/golang/protobuf v1.4.2
	github.com/kardianos/service v1.1.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/fx v1.17.1
	go.uber.org/zap v1.21.0
	google.golang.org/protobuf v1.25.0
)

replace (
	github.com/bhbosman/goMessages => ../goMessages
	github.com/bhbosman/gocommon => ../gocommon
	github.com/bhbosman/gocomms => ../gocomms
)

replace github.com/bhbosman/goCommsStacks => ../goCommsStacks
replace github.com/bhbosman/goCommsNetDialer => ../goCommsNetDialer
replace github.com/bhbosman/goCommsNetListener => ../goCommsNetListener
replace github.com/bhbosman/goCommsDefinitions => ../goCommsDefinitions