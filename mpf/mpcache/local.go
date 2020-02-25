// Package mpcache memcache
// User: 姜伟
// Time: 2020-02-19 06:15:40
package mpcache

import (
    "time"

    "github.com/patrickmn/go-cache"
)

var (
    insLocal *cache.Cache
)

func init() {
    insLocal = cache.New(5*time.Minute, 10*time.Minute)
}

// NewLocal NewLocal
func NewLocal() *cache.Cache {
    return insLocal
}
