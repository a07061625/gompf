package main

import (
    "fmt"
    "reflect"
)

func xxxx(data string) {
    fmt.Println("print:", data)
}

func rrrr(a int, b int) {
    fmt.Println("add: ", a+b)
}

func oooo() {
    fmt.Println("xxxx")
}

var command = map[string]interface{}{
    "f1": xxxx,
    "f2": rrrr,
    "f3": oooo,
}

func call(name string, params ...interface{}) {
    // 获取反射对象
    f := reflect.ValueOf(command[name])
    // 判断函数参数和传入的参数是否相等
    if len(params) != f.Type().NumIn() {
        return
    }
    // 然后将传入参数转为反射类型切片
    in := make([]reflect.Value, len(params))
    for k, param := range params {
        in[k] = reflect.ValueOf(param)
    }
    // 利用函数反射对象的call方法调用函数.
    f.Call(in)
}

func main() {
    call("f1", "abcd")
    call("f2", "100", 66)
    call("f3")
}
