package connection_test

import (
	"encoding/json"
	"github.com/bhbosman/gokraken/internal/krakenWS/connection"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDff(t *testing.T) {
	ss := `{
"as": [
    [ "0.05005", "0.00000500", "1582905487.684110" ],
    [ "0.05010", "0.00000500", "1582905486.187983" ],
    [ "0.05015", "0.00000500", "1582905484.480241" ],
    [ "0.05020", "0.00000500", "1582905486.645658" ],
    [ "0.05025", "0.00000500", "1582905486.859009" ],
    [ "0.05030", "0.00000500", "1582905488.601486" ],
    [ "0.05035", "0.00000500", "1582905488.357312" ],
    [ "0.05040", "0.00000500", "1582905488.785484" ],
    [ "0.05045", "0.00000500", "1582905485.302661" ],
    [ "0.05050", "0.00000500", "1582905486.157467" ] ],
"bs": [
    [ "0.05000", "0.00000500", "1582905487.439814" ],
    [ "0.04995", "0.00000500", "1582905485.119396" ],
    [ "0.04990", "0.00000500", "1582905486.432052" ],
    [ "0.04980", "0.00000500", "1582905480.609351" ],
    [ "0.04975", "0.00000500", "1582905476.793880" ],
    [ "0.04970", "0.00000500", "1582905486.767461" ],
    [ "0.04965", "0.00000500", "1582905481.767528" ],
    [ "0.04960", "0.00000500", "1582905487.378907" ],
    [ "0.04955", "0.00000500", "1582905483.626664" ],
    [ "0.04950", "0.00000500", "1582905488.509872" ] ]
}`

	var data interface{}
	err := json.Unmarshal([]byte(ss), &data)
	assert.NoError(t, err)
	book := connection.NewFullMarketOrderBook("", "book-10")
	_ = book.HandleBook(data.(map[string]interface{}))
	c := book.CalculateCheckSum()
	assert.Equal(t, uint32(974947235), c)
}
