package datastoredb

import (
	"bytes"
	"encoding/gob"
	"sync"
)

// cache is a simple in-process cache backed by sync.Map.
// It stores gob-encoded copies to prevent mutation of cached values.
type cache struct {
	m sync.Map
}

var appCache = &cache{}

// Get retrieves a cached value by key, gob-decoding it into dst.
// Returns true if the key was found and decoded successfully.
func (c *cache) Get(key string, dst interface{}) bool {
	v, ok := c.m.Load(key)
	if !ok {
		return false
	}
	buf := bytes.NewBuffer(v.([]byte))
	if err := gob.NewDecoder(buf).Decode(dst); err != nil {
		return false
	}
	return true
}

// Set stores a gob-encoded copy of val under key.
func (c *cache) Set(key string, val interface{}) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(val); err != nil {
		return
	}
	c.m.Store(key, buf.Bytes())
}
