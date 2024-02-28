package cefileprotocol

import (
	"context"
	"io"

	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/binding/format"
	"github.com/cloudevents/sdk-go/v2/protocol"
)

type Protocol struct {
	w io.WriteCloser
}

func New(w io.WriteCloser) (*Protocol, error) {
	return &Protocol{w: w}, nil
}

func (t *Protocol) Send(ctx context.Context, in binding.Message, transformers ...binding.Transformer) error {
	// Write the message to the file.
	return WriteMessage(ctx, in, t.w, transformers...)
}

func (t *Protocol) Close(ctx context.Context) error {
	return t.w.Close()
}

func WriteMessage(ctx context.Context, m binding.Message, writer io.Writer, transformers ...binding.Transformer) error {
	_, err := binding.Write(
		ctx,
		m,
		eventWriter{writer},
		nil,
		transformers...,
	)
	return err
}

type eventWriter struct {
	io.Writer
}

func (w eventWriter) SetStructuredEvent(ctx context.Context, format format.Format, event io.Reader) error {
	if _, err := io.Copy(w, event); err != nil {
		return err
	}
	return nil
}

var (
	_ protocol.Sender = (*Protocol)(nil)
	_ protocol.Closer = (*Protocol)(nil)

	_ binding.StructuredWriter = (*eventWriter)(nil)
)
