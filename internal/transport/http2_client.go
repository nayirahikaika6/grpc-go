package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/grpcsync"
	"google.golang.org/grpc/internal/status"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
)

// http2Client implements the ClientTransport interface for HTTP2.
type http2Client struct {
	ctx        context.Context
	cancel     context.CancelFunc
	ctxDone    <-chan struct{} // same as ctx.Done(), cached for performance
	userAgent  string
	md         interface{}
	conn       net.Conn
	remoteAddr net.Addr
	localAddr  net.Addr
	authInfo   credentials.AuthInfo // it can be nil

	// Fields of the client transport that are protected by mu.
	mu            sync.Mutex
	state         transportState
	activeStreams map[uint32]*Stream
	// goAwayID is the LastStreamID received in a GOAWAY frame.
	goAwayID      uint32
	goAwayReason  GoAwayReason
	// nextChanID is the next stream ID to allocate.
	nextChanID    uint32
	// ...
}

// NewStream creates a new stream on the transport.
func (t *http2Client) NewStream(ctx context.Context, callHdr *CallHeader) (*Stream, error) {
	t.mu.Lock()
	if t.state == draining || t.state == closing {
		t.mu.Unlock()
		return nil, errConnDrain
	}
	if t.nextChanID > t.goAwayID {
		t.mu.Unlock()
		return nil, errConnDrain
	}
	// ...
	return nil, nil
}