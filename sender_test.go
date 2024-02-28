package cefileprotocol

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudevents/sdk-go/v2/protocol"
)

type noopWriteCloser struct {
	io.Writer
}

func (no noopWriteCloser) Close() error {
	return nil
}

func TestClient(t *testing.T) {
	var buf bytes.Buffer
	tp, _ := New(noopWriteCloser{&buf})

	// Create an ce client.
	c, err := client.New(tp)
	if err != nil {
		t.Errorf("error creating client: %s", err)
	}

	evt := event.New()
	evt.Context.SetID("001")
	evt.Context.SetSource("source")
	evt.Context.SetTime(time.Now())
	evt.Context.SetSubject("subject")
	evt.Context.SetType("type")
	evt.SetExtension("timestamp", time.Now())
	evt.SetExtension("app", "verzuimen")

	res := c.Send(context.Background(), evt)
	if ! protocol.IsACK(res) {
		t.Errorf("res error: %s", res.Error())
	}
	file := buf.String()
	t.Errorf("res: %s", file)
}
