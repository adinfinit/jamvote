package datastoredb

import (
	"bytes"
	"encoding/gob"
	"sync"

	"github.com/adinfinit/jamvote/event"
	"github.com/adinfinit/jamvote/user"
)

// cache is a simple in-process typed cache backed by sync.Map.
// It stores gob-encoded copies to prevent mutation of cached values.
type cache[T any] struct {
	m sync.Map
}

var userCache = &cache[user.User]{}
var eventCache = &cache[event.Event]{}

// Get retrieves a cached value by key, gob-decoding and returning it.
// Returns the value and true on hit, or nil and false on miss.
func (c *cache[T]) Get(key string) (*T, bool) {
	v, ok := c.m.Load(key)
	if !ok {
		return nil, false
	}
	dst := new(T)
	buf := bytes.NewBuffer(v.([]byte))
	if err := gob.NewDecoder(buf).Decode(dst); err != nil {
		return nil, false
	}
	return dst, true
}

// Set stores a gob-encoded copy of val under key.
func (c *cache[T]) Set(key string, val *T) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(val); err != nil {
		return
	}
	c.m.Store(key, buf.Bytes())
}
