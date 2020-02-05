/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/5 0005
 * Time: 14:48
 */
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
