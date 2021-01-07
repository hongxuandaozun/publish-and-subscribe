package pubsub

import (
	"sync"
	"time"
)

type (
	subscriber chan interface{}
	topicFunc func(v interface{}) bool
)
// 发布者对象
type Publisher struct {
	 sync.RWMutex // 读写锁
	buffer int     // 订阅队列的缓存大小
	timeout time.Duration   // 发布超时时间
	subscribers map[subscriber]topicFunc   //订阅者消息
}
// 构建一个发布者对象,可以设置发布超时时间和缓存队列的长度
func NewPublisher(publishTimeout time.Duration,buffer int)*Publisher  {
	return &Publisher{
		buffer: buffer,
		timeout: publishTimeout,
		subscribers: make(map[subscriber]topicFunc),
	}
}
// 提阿难捱一个新的订阅者,订阅过滤器筛选后的主题
func (p *Publisher) SubscribeTopic(topic topicFunc)chan interface{}  {
	ch := make(chan interface{},p.buffer)
	p.Lock()
	p.subscribers[ch] = topic
	p.Unlock()
	return ch
}
// 添加一个新的订阅者,订阅全部主题
func (p *Publisher) Subscribe()chan interface{} {
	return p.SubscribeTopic(nil)
}

// 发送主题,可以容忍一定的超时
func (p *Publisher) SendTopic(sub subscriber,topic topicFunc,v interface{} ,wg *sync.WaitGroup)  {
	defer wg.Done()
	if topic != nil && !topic(v){
		return
	}
	select {
		case sub <- v:
		case <- time.After(p.timeout):
	}
}

// 发布一个主题

func (p *Publisher) Publish(v interface{})  {
	p.RLock()
	defer p.RUnlock()
	var wg sync.WaitGroup
	for sub, topic := range p.subscribers{
		wg.Add(1)
		go p.SendTopic(sub,topic,v,&wg)
	}
	wg.Wait()
}
// 关闭发布者对象,同时关闭所有的订阅者管道

func (p *Publisher) Close(){
	p.Lock()
	defer p.Unlock()
	for sub := range p.subscribers {
		delete(p.subscribers,sub)
		close(sub)
	}
}