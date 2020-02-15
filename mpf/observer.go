package mpf

// 观察者接口
type IObserver interface {
    Notify(data interface{})
}

// 被观察者接口
type ISubject interface {
    AddObservers(observers ...IObserver) // 添加观察者
    ClearObservers()                     // 清空观察者
    GetObservers() []IObserver           // 获取观察者列表
    NotifyObservers(data interface{})    // 通知观察者
}

type subjectBasic struct {
    observers []IObserver
}

func (subject *subjectBasic) AddObservers(observers ...IObserver) {
    if len(observers) > 0 {
        subject.observers = append(subject.observers, observers...)
    }
}

func (subject *subjectBasic) ClearObservers() {
    subject.observers = make([]IObserver, 0)
}

func (subject *subjectBasic) GetObservers() []IObserver {
    return subject.observers
}

func NewSubject() *subjectBasic {
    subject := &subjectBasic{}
    subject.observers = make([]IObserver, 0)
    return subject
}
