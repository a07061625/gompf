// Package mpdp observer
// User: 姜伟
// Time: 2020-02-19 06:27:16
package mpdp

// IObserver 观察者接口
type IObserver interface {
    Notify(data interface{})
}

// ISubject 被观察者接口
type ISubject interface {
    AddObservers(observers ...IObserver)    // 添加观察者
    ClearObservers()                        // 清空观察者
    NotifyObservers(data interface{}) error // 通知观察者
}

// SubjectBasic 被观察者基础结构体
type SubjectBasic struct {
    Observers []IObserver
}

// AddObservers 添加观察者
func (subject *SubjectBasic) AddObservers(observers ...IObserver) {
    if len(observers) > 0 {
        subject.Observers = append(subject.Observers, observers...)
    }
}

// ClearObservers 清空观察者
func (subject *SubjectBasic) ClearObservers() {
    subject.Observers = make([]IObserver, 0)
}

// NewSubject 实例化被观察者
func NewSubject() SubjectBasic {
    subject := SubjectBasic{}
    subject.Observers = make([]IObserver, 0)
    return subject
}
