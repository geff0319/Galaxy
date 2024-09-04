package ytdlp

import (
	"context"
	"errors"
	"fmt"
	evbus "github.com/asaskevich/EventBus"
	gefflog "github.com/ge-fei-fan/gefflog"
	"golang.org/x/sync/semaphore"
)

const queueName = "process:pending"

type MessageQueue struct {
	concurrency int
	eventBus    evbus.Bus
	//logger      *slog.Logger
}

// Creates a new message queue.
// By default it will be created with a size equals to nthe number of logical
// CPU cores -1.
// The queue size can be set via the qs flag.
func NewMessageQueue() (*MessageQueue, error) {
	qs := YdpConfig.QueueSize

	if qs <= 0 {
		return nil, errors.New("invalid queue size")
	}

	return &MessageQueue{
		concurrency: qs,
		eventBus:    evbus.New(),
		//logger:      l,
	}, nil
}

// Publish a message to the queue and set the task to a peding state.
func (m *MessageQueue) Publish(p *Process) {
	// needs to have an id set before
	p.SetPending()

	m.eventBus.Publish(queueName, p)
}
func (m *MessageQueue) PublishByTopic(topic string, p *Process) {
	m.eventBus.Publish(topic, p)
}

func (m *MessageQueue) SetupConsumers() {
	go m.downloadConsumer()
	go m.metadataSubscriber()
}

// Setup the consumer listener which subscribes to the changes to the producer
// channel and triggers the "download" action.
func (m *MessageQueue) downloadConsumer() {
	//sem := semaphore.NewWeighted(int64(m.concurrency))
	sem := semaphore.NewWeighted(int64(m.concurrency))
	m.eventBus.SubscribeAsync("process:downloading", func(p *Process) {
		//TODO: provide valid context
		sem.Acquire(context.Background(), 1)
		defer sem.Release(1)

		gefflog.Info(fmt.Sprintf("下载开始：received process from event bus bus=%s consumer=downloadConsumer id=%s", queueName, p.getShortId()))

		if p.Progress.Status != StatusCompleted {
			p.Start()
		}
		gefflog.Info(fmt.Sprintf("started process bus=%s id=%s", queueName, p.getShortId()))
	}, false)
}

// Setup the metadata consumer listener which subscribes to the changes to the
// producer channel and adds metadata to each download.
func (m *MessageQueue) metadataSubscriber() {
	// How many concurrent metadata fetcher jobs are spawned
	// Since there's ongoing downloads, 1 job at time seems a good compromise
	sem := semaphore.NewWeighted(1)
	//sem := semaphore.NewWeighted(int64(m.concurrency))
	m.eventBus.SubscribeAsync(queueName, func(p *Process) {
		//TODO: provide valid context
		sem.Acquire(context.TODO(), 1)
		defer sem.Release(1)

		gefflog.Info(fmt.Sprintf("解析开始：received process from event bus bus=%s consumer=metadataConsumer id=%s", queueName, p.getShortId()))

		if p.Progress.Status == StatusCompleted {
			gefflog.Info(fmt.Sprintf("process has an illegal state id=%s status=%d", p.getShortId(), p.Progress.Status))
			return
		}
		// 解析失败，设置为解析失败状态
		if err := p.SetMetadata(); err != nil {
			gefflog.Err(fmt.Sprintf("failed to retrieve metadata id=%s err=%s", p.getShortId(), err.Error()))
			p.Progress.Status = StatusErrored
			YdpConfig.Mq.eventBus.Publish("notify", "error", "解析视频失败："+err.Error())
			return
		}

		if YdpConfig.Mdb.IsProcessExist(p) {
			YdpConfig.Mdb.Delete(p.Id)
			m.eventBus.Publish("notify", "error", "任务已存在")
		} else {
			m.eventBus.Publish("process:downloading", p)
		}
	}, false)
}

func (m *MessageQueue) SetConsumer(topic string, fn any) error {
	return m.eventBus.SubscribeAsync(topic, fn, false)
}
