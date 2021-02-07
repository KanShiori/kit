// Package eventbus 简单实现了一个用于发布订阅的事件总线.
package eventbus

import (
	"fmt"
	"sync"
)

// Event 为发布的事件消息
type Event struct {
	Data  interface{}
	Topic string
}

// EventHandler 为 Subscribe 注册的回调
type EventHandler interface {
	// EventHandle 接收到事件时的回调
	EventHandle(data Event)

	// Name 会作为 bucket 的 key, 用于检索
	Name() string
}

// Bus 是一个中间件, 提供订阅与发布模型
//
// 发布者通过 Bus.Publish 进行 topic 的事件发布.
// 订阅者 Bus.Subscribe 与 UnSubscribe 进行订阅与反订阅.
// 订阅意味着注册一个回调函数, 会在事件发布时执行对应的回调
type Bus interface {
	Subscribe(topic string, handler EventHandler) error

	UnSubscribe(topic string, handler EventHandler) error

	Publish(topic string, data interface{})
}

// NewEventBus 创建一个 Bus
func NewEventBus() Bus {
	return &bus{
		subscribers: make(map[string]bucket),
		mutex:       sync.Mutex{},
	}
}

// bus implement Bus
type bus struct {
	subscribers map[string]bucket
	mutex       sync.Mutex
}

func (b *bus) Subscribe(topic string, handler EventHandler) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// 检查是否需要新建 topic 对应的 bucket
	_, ok := b.subscribers[topic]
	if !ok {
		b.subscribers[topic] = newBucket()
	}

	// 向 topic 对应的 bucket 加入 handler
	bucket := b.subscribers[topic]
	err := bucket.add(handler)
	if err != nil {
		return err
	}

	return nil
}

func (b *bus) UnSubscribe(topic string, handler EventHandler) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// 检查 topic 是否存在
	bucket, ok := b.subscribers[topic]
	if !ok {
		return fmt.Errorf("topic not exist")
	}

	// topic 对应的 bucket 移除其 handler
	bucket.remove(handler)

	return nil
}

func (b *bus) Publish(topic string, data interface{}) {
	b.mutex.Lock()

	if bucket, ok := b.subscribers[topic]; ok {
		// deepcopy 一份发布
		go publish(topic, data, bucket.deepcopy())
	}

	b.mutex.Unlock()
}

func publish(topic string, data interface{}, bk bucket) {
	// 遍历 bucket 的 handler, 循环执行 handler 的回调
	for _, handler := range bk {
		handler.EventHandle(Event{
			Data:  data,
			Topic: topic,
		})
	}
}
