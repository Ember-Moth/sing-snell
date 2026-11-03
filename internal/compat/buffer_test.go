package compat

import (
	"bytes"
	"testing"

	"github.com/metacubex/sing/common/buf"
	N "github.com/metacubex/sing/common/network"
)

func TestLargeRecordHeadroom(t *testing.T) {
	options := N.ReadWaitOptions{FrontHeadroom: 2048, RearHeadroom: 256, MTU: 1200}
	payload := bytes.Repeat([]byte{0x6a}, 65535)
	buffer := CopyReadBuffer(options, buf.As(payload))
	defer buffer.Release()
	if !bytes.Equal(buffer.Bytes(), payload) {
		t.Fatal("adding headroom truncated the record")
	}
	if buffer.Start() < options.FrontHeadroom || buffer.FreeLen() < options.RearHeadroom {
		t.Fatal("missing requested headroom")
	}
}

func TestRecordAllocationIgnoresSmallerMTU(t *testing.T) {
	options := N.ReadWaitOptions{FrontHeadroom: 2048, RearHeadroom: 256, MTU: 1200}
	buffer := NewBufferSize(options, 65535)
	defer buffer.Release()
	if buffer.FreeLen() != 65535 {
		t.Fatalf("record capacity = %d", buffer.FreeLen())
	}
	buffer.Extend(65535)
	options.PostReturn(buffer)
	if buffer.FreeLen() != options.RearHeadroom {
		t.Fatalf("rear headroom = %d", buffer.FreeLen())
	}
}
