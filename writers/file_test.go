package writers

import (
	"context"
	"github.com/boz/kail"
	"testing"
	"time"
)

type TestEventSource struct {
	node      string
	container string
	namespace string
	name      string
}

func (t TestEventSource) Namespace() string {
	return t.namespace
}

func (t TestEventSource) Name() string {
	return t.name
}

func (t TestEventSource) Container() string {
	return t.container
}

func (t TestEventSource) Node() string {
	return t.node
}

type TestEvent struct {
	source TestEventSource
	log    []byte
}

func (t TestEvent) Source() kail.EventSource {
	return t.source
}

func (t TestEvent) Log() []byte {
	return t.log
}

func Test(t *testing.T) {
	writer := NewFileWriter(context.Background())
	writer.Print(TestEvent{
		source: TestEventSource{
			"test-node",
			"test-container",
			"test-namespace",
			"test-pod",
		},
		log: []byte("Hello, world!"),
	})

	time.Sleep(2 * time.Second)
}
