package respx

import "sync"

type EventType string

const (
	BeforeResponse EventType = "before_response"
	AfterResponse  EventType = "after_response"
)

type EventHandler func(w ResponseWriter, data any)

var (
	eventHandlers = make(map[EventType][]EventHandler)
	eventMu       sync.RWMutex
)

// RegisterEvent 注册一个事件处理器
func RegisterEvent(eventType EventType, handler EventHandler) {
	eventMu.Lock()
	defer eventMu.Unlock()
	eventHandlers[eventType] = append(eventHandlers[eventType], handler)
}

// triggerEvent 触发指定类型的事件
func triggerEvent(eventType EventType, w ResponseWriter, data any) {
	eventMu.RLock()
	defer eventMu.RUnlock()
	for _, handler := range eventHandlers[eventType] {
		handler(w, data)
	}
}

// ClearEventHandlers 清除所有注册的事件处理器
func ClearEventHandlers() {
	eventMu.Lock()
	defer eventMu.Unlock()
	eventHandlers = make(map[EventType][]EventHandler)
}
