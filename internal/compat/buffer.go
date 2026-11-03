package compat

import (
	"github.com/metacubex/sing/common/buf"
	N "github.com/metacubex/sing/common/network"
)

// NewBufferSize reserves space for the entire record and the caller's framing.
// ReadWaitOptions.NewBuffer may instead allocate the caller's smaller MTU.
func NewBufferSize(options N.ReadWaitOptions, size int) *buf.Buffer {
	buffer := buf.NewSize(size + options.FrontHeadroom + options.RearHeadroom)
	if options.FrontHeadroom > 0 {
		buffer.Resize(options.FrontHeadroom, 0)
	}
	if options.RearHeadroom > 0 {
		buffer.Reserve(options.RearHeadroom)
	}
	return buffer
}

// CopyReadBuffer retains large records when adding headroom. MetaCubeX/sing's
// generic Copy uses a fixed packet allocation, which can truncate these records.
func CopyReadBuffer(options N.ReadWaitOptions, buffer *buf.Buffer) *buf.Buffer {
	if buffer.Start() >= options.FrontHeadroom && buffer.FreeLen() >= options.RearHeadroom {
		return buffer
	}
	copy := NewBufferSize(options, buffer.Len())
	_, _ = copy.Write(buffer.Bytes())
	options.PostReturn(copy)
	buffer.Release()
	return copy
}
