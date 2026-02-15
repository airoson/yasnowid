package yasnowid

import (
	"errors"
	"sync/atomic"
	"time"
)

// Generator generates Snowflake IDs using a specific node ID.
// The node ID typically combines a data center ID and a machine ID.
type Generator struct {
	i        atomic.Int64
	nodeMask int64
}

const (
	epoch   = int64(1288834974657)
	cmpMask = int64(^0x3FFFFF)
)

var (
	ErrNodeIDTooLarge = errors.New("node id can't exceed 1023 (10 bits in total)")
)

// NewGenerator returns a new generator for the specified nodeID.
func NewGenerator(nodeID int) (*Generator, error) {
	if nodeID > 1023 {
		return nil, ErrNodeIDTooLarge
	}

	return &Generator{
		i:        atomic.Int64{},
		nodeMask: int64(nodeID) << 12,
	}, nil
}

func (g *Generator) timestamp() int64 {
	return (time.Now().UnixMilli() - epoch) << 22
}

func (g *Generator) compareTimestamps(id, timestamp int64) int64 {
	return timestamp - (id & cmpMask)
}

// ID returns the next unique Snowflake ID.
// It handles clock drift by waiting until the system time catches up
// with the timestamp stored in the last generated ID.
func (g *Generator) ID() int64 {
	for {
		timestamp := g.timestamp()
		current := g.i.Load()

		var next int64
		var seq int64
		cmp := g.compareTimestamps(current, timestamp)
		if cmp < 0 {
			// Clock drift
			continue
		}
		if cmp == 0 {
			seq = (current & 0xFFF) + 1
			if seq > 0xFFF {
				// Very unluckily: overflow (more than 4kk RPS are needed to reach this)
				continue
			}
		}
		next = timestamp | g.nodeMask | seq

		if g.i.CompareAndSwap(current, next) {
			return next
		}
	}
}

// JoinIDs function helps to create nodeID from its subparts: dataCenterID and machineID
func JoinIDs(dataCenterID, machineID int) int {
	return dataCenterID<<5 | machineID
}
