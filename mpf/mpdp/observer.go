/**
 * 观察者模式
 * User: 姜伟
 * Date: 20-2-16
 * Time: 下午4:30
 */
package mpdp

// 观察者接口
type IObserver interface {
    Notify(data interface{})
}

// 被观察者接口
type ISubject interface {
    AddObservers(observers ...IObserver)    // 添加观察者
    ClearObservers()                        // 清空观察者
    NotifyObservers(data interface{}) error // 通知观察者
}

type SubjectBasic struct {
    Observers []IObserver
}

func (subject *SubjectBasic) AddObservers(observers ...IObserver) {
    if len(observers) > 0 {
        subject.Observers = append(subject.Observers, observers...)
    }
}

func (subject *SubjectBasic) ClearObservers() {
    subject.Observers = make([]IObserver, 0)
}

func NewSubject() SubjectBasic {
    subject := SubjectBasic{}
    subject.Observers = make([]IObserver, 0)
    return subject
}
