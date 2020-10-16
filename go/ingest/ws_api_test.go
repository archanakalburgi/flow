package ingest

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	pb "go.gazette.dev/core/broker/protocol"
)

func testCSVSimple(t *testing.T, addr string) {
	var valid = []string{
		"i,s",
		// May optionally include trailing newline.
		"32,hi\n",
		"42,bye",
		// Handles multiple records in one message.
		"72,first\n87,second",
		"23,one\n24,two\n25,three\n",
		// Quotes.
		`"62","whoot"`,
	}

	var dialer = websocket.Dialer{
		Subprotocols: []string{wsCSVProtocol},
	}
	var c, _, err = dialer.Dial("ws://"+addr+"/ingest/testing/int-string", nil)
	require.NoError(t, err)

	var whoops, processed = runWebsocket(t, c, valid)
	require.Equal(t, "", whoops)
	require.Equal(t, 8, processed)
}

func testTSVSimple(t *testing.T, addr string) {
	var valid = []string{
		"i\ts",
		"32\thi\n",
		"42\tone\n52\ttwo\n62\tthree\n",
		"72\tbye",
	}

	var dialer = websocket.Dialer{
		Subprotocols: []string{wsTSVProtocol},
	}
	var c, _, err = dialer.Dial("ws://"+addr+"/ingest/testing/int-string", nil)
	require.NoError(t, err)

	var whoops, processed = runWebsocket(t, c, valid)
	require.Equal(t, "", whoops)
	require.Equal(t, 5, processed)
}

func testCSVCollectionNotFound(t *testing.T, addr string) {
	var valid = []string{
		"i,s",
	}

	var dialer = websocket.Dialer{
		Subprotocols: []string{wsCSVProtocol},
	}
	var c, _, err = dialer.Dial("ws://"+addr+"/ingest/not/found", nil)
	require.NoError(t, err)

	var whoops, _ = runWebsocket(t, c, valid)
	require.Equal(t, "processing frame: \"not/found\" is not an ingestable collection", whoops)
}

func testCSVMalformed(t *testing.T, addr string) {
	var valid = []string{
		"i,s",
		"32,hi,extra",
	}

	var dialer = websocket.Dialer{
		Subprotocols: []string{wsCSVProtocol},
	}
	var c, _, err = dialer.Dial("ws://"+addr+"/ingest/testing/int-string", nil)
	require.NoError(t, err)

	var whoops, _ = runWebsocket(t, c, valid)
	require.Equal(t, "processing frame: record on line 2: wrong number of fields", whoops)
}

func testCSVProjectionNotFound(t *testing.T, addr string) {
	var valid = []string{
		"i,sss",
	}

	var dialer = websocket.Dialer{
		Subprotocols: []string{wsCSVProtocol},
	}
	var c, _, err = dialer.Dial("ws://"+addr+"/ingest/testing/int-string", nil)
	require.NoError(t, err)

	var whoops, _ = runWebsocket(t, c, valid)
	require.Equal(t, `processing frame: collection "testing/int-string" has no projection "sss"`, whoops)
}

func testCSVConversionError(t *testing.T, addr string) {
	var valid = []string{
		"i,s",
		"32.32,hi",
	}

	var dialer = websocket.Dialer{
		Subprotocols: []string{wsCSVProtocol},
	}
	var c, _, err = dialer.Dial("ws://"+addr+"/ingest/testing/int-string", nil)
	require.NoError(t, err)

	var whoops, _ = runWebsocket(t, c, valid)
	require.Equal(t, `processing frame: failed to parse "32.32" into [integer]: strconv.ParseInt: parsing "32.32": invalid syntax`, whoops)
}

func runWebsocket(t *testing.T, c *websocket.Conn, msgs []string) (string, int) {
	// Write all messages and initiate close.
	for _, msg := range msgs {
		require.NoError(t, c.WriteMessage(websocket.TextMessage, []byte(msg)))
	}
	var frame = websocket.FormatCloseMessage(websocket.CloseNoStatusReceived, "")
	require.NoError(t, c.WriteControl(websocket.CloseMessage, frame, time.Time{}))

	var whoops string
	var processed int

	var message struct {
		// Success fields.
		Offsets   pb.Offsets
		Processed int
		Etcd      pb.Header_Etcd

		// Error fields.
		ApproximateFrame int
		Error            string
	}
	for {
		var err = c.ReadJSON(&message)
		if err != nil {
			if whoops == "" {
				require.True(t, websocket.IsCloseError(err, websocket.CloseNormalClosure), err)
			} else {
				require.True(t, websocket.IsCloseError(err, websocket.CloseProtocolError), err)
			}
			return whoops, processed
		}

		if message.Error != "" {
			whoops = message.Error
			require.NotZero(t, message.ApproximateFrame)
		} else {
			require.NotEmpty(t, message.Offsets)
			require.NotZero(t, message.Etcd.Revision)
			require.NotZero(t, message.Processed)
			processed = message.Processed
		}
	}
}