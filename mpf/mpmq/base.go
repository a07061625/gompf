// Package mpmq base
// User: 姜伟
// Time: 2020-02-19 06:40:35
package mpmq

import (
    "github.com/a07061625/gompf/mpf"
)

type baseRabbit struct {
    ExchangeName string
    TopicPrefix  string
}

func newBaseRabbit() baseRabbit {
    r := baseRabbit{"", ""}
    r.ExchangeName = "mpexchange" + mpf.EnvProjectKey()
    r.TopicPrefix = mpf.EnvProjectKey() + "."
    return r
}
