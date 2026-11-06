package snowflake

import (
	"strconv"
	"sync"
	"time"
)

const (
	idBits            = 10
	seqBits           = 12
	timeBits          = 41
	idMask            = (1 << idBits) - 1
	seqMask           = (1 << seqBits) - 1
	ErrNodeIdOverflow = SnowFlakeError("node ids must be in range 0 - 1023")
)

var (
	epoch = time.UnixMilli(1288834974657) // default to twitter's Epoch
)

type SnowFlakeError string

func (e SnowFlakeError) Error() string {
	return string(e)
}

type ID int64

type Generator struct {
	//nodeId int64
	count  int64
	lastMs int64
	mu     sync.Mutex
}

func (g *Generator) MustID(timestamp time.Time, nodeId int64) ID {
	id, err := g.NewID(timestamp, nodeId)
	if err != nil {
		panic(err)
	}
	return id
}

func (g *Generator) NewID(timestamp time.Time, nodeId int64) (ID, error) {
	if nodeId < 0 || nodeId > idMask {
		return -1, ErrNodeIdOverflow
	}
	g.mu.Lock()
	ms := timestamp.UnixMilli()
	if ms < g.lastMs {
		ms = g.lastMs
	}
	if ms == g.lastMs {
		g.count = g.count + 1
	} else {
		g.count = 1
		g.lastMs = ms
	}
	// prevent sequence from overflowing 12 bit limit
	// "borrow" the next ms
	if g.count > seqMask {
		g.lastMs++
		ms = g.lastMs
		g.count = 1
	}
	sequence := g.count
	g.mu.Unlock()
	t := ms - epoch.UnixMilli()
	packed := t<<(idBits+seqBits) | nodeId<<seqBits | sequence
	return ID(packed), nil
}

func (id ID) Timestamp() time.Time {
	ms := id >> (idBits + seqBits)
	return epoch.Add(time.Duration(int64(ms)) * time.Millisecond)
}

func (id ID) NodeID() int64 {
	return int64(id>>seqBits) & idMask
}

func (id ID) SequenceId() int64 {
	return int64(id) & seqMask
}

func (id ID) String() string {
	return strconv.FormatInt(int64(id), 10)
}
