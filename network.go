package snell

import (
	"github.com/metacubex/sing/common"
	"github.com/metacubex/sing/common/buf"
	"github.com/metacubex/sing/common/bufio"
	M "github.com/metacubex/sing/common/metadata"
	N "github.com/metacubex/sing/common/network"
)

// CloseHandlerFunc completes a logical request, which may share a reusable
// transport. Closing the transport is not a substitute for this callback.
type CloseHandlerFunc = func(error)

type VectorisedWriteCreator interface {
	CreateVectorisedWriter() (N.VectorisedWriter, bool)
}

type PacketBatchWriter interface {
	WritePacketBatch(buffers []*buf.Buffer, destinations []M.Socksaddr) error
}

type PacketBatchWriteCreator interface {
	CreatePacketBatchWriter() (PacketBatchWriter, bool)
}

// CreateVectorisedWriter preserves Snell's lazy record writer while using
// MetaCubeX/sing for buffers and the underlying scatter/gather implementation.
func CreateVectorisedWriter(writer any) (N.VectorisedWriter, bool) {
	if creator, ok := writer.(VectorisedWriteCreator); ok {
		return creator.CreateVectorisedWriter()
	}
	if vectorised, ok := bufio.CreateVectorisedWriter(writer); ok {
		return vectorised, true
	}
	if replaceable, ok := writer.(N.WriterWithUpstream); ok && replaceable.WriterReplaceable() {
		if upstream, ok := writer.(N.WithUpstreamWriter); ok {
			return CreateVectorisedWriter(upstream.UpstreamWriter())
		}
		if upstream, ok := writer.(common.WithUpstream); ok {
			return CreateVectorisedWriter(upstream.Upstream())
		}
	}
	return nil, false
}

func CreatePacketBatchWriter(writer N.PacketWriter) (PacketBatchWriter, bool) {
	if creator, ok := writer.(PacketBatchWriteCreator); ok {
		return creator.CreatePacketBatchWriter()
	}
	if batch, ok := writer.(PacketBatchWriter); ok {
		return batch, true
	}
	if replaceable, ok := writer.(N.WriterWithUpstream); ok && replaceable.WriterReplaceable() {
		if upstream, ok := writer.(N.WithUpstreamWriter); ok {
			return CreatePacketBatchWriter(upstream.UpstreamWriter().(N.PacketWriter))
		}
		if upstream, ok := writer.(common.WithUpstream); ok {
			return CreatePacketBatchWriter(upstream.Upstream().(N.PacketWriter))
		}
	}
	return nil, false
}
