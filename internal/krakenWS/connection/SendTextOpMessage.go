package connection

import (
	"github.com/bhbosman/goCommsStacks/webSocketMessages/wsmsg"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/reactivex/rxgo/v2"
)

func SendTextOpMessage(
	message proto.Message,
	ToConnection goprotoextra.ToConnectionFunc,
	ToConnectionReplacement rxgo.NextFunc,
) error {
	rws := gomessageblock.NewReaderWriter()
	m := jsonpb.Marshaler{
		OrigName:     false,
		EnumsAsInts:  false,
		EmitDefaults: false,
		Indent:       "",
		AnyResolver:  nil,
	}
	var err error
	err = m.Marshal(rws, message)
	if err != nil {
		return err
	}

	var flatten []byte
	flatten, err = rws.Flatten()
	if err != nil {
		return err
	}

	WebSocketMessage := wsmsg.WebSocketMessage{
		OpCode:  wsmsg.WebSocketMessage_OpText,
		Message: flatten,
	}
	readWriterSize, err := stream.Marshall(&WebSocketMessage)
	if err != nil {
		return err
	}

	return ToConnection(readWriterSize)
}
