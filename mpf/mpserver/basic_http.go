/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/8 0008
 * Time: 13:21
 */
package mpserver

import (
    "github.com/spf13/viper"
)

type basicHttp struct {
    basic
}

func (s *basicHttp) bootServer() {
    s.bootBasic()
}

func (s *basicHttp) StartServer() {
    s.bootServer()
    s.startBasic()
}

func NewBasicHttp(conf *viper.Viper) *basicHttp {
    s := &basicHttp{newBasic(conf)}
    return s
}
