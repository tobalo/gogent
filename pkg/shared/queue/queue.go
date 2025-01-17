package queue

import (
	"context"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// BatchProcessor defines the interface for processing message batches
type BatchProcessor interface {
	ProcessBatch(ctx context.Context, msgs []*nats.Msg) error
}

// Config holds queue configuration
type Config struct {
	QueueSize    int
	BatchSize    int
	ProcessDelay time.Duration
}

// MessageQueue handles batch message processing
type MessageQueue struct {
	config    Config
	processor BatchProcessor
	msgChan   chan *nats.Msg
	done      chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
	running   bool
}

// NewMessageQueue creates a new message queue
func NewMessageQueue(cfg Config, processor BatchProcessor) *MessageQueue {
	return &MessageQueue{
		config:    cfg,
		processor: processor,
		msgChan:   make(chan *nats.Msg, cfg.QueueSize),
		done:      make(chan struct{}),
	}
}

// Start begins processing messages
func (q *MessageQueue) Start(ctx context.Context) error {
	q.mu.Lock()
	if q.running {
		q.mu.Unlock()
		return nil
	}
	q.running = true
	q.mu.Unlock()

	q.wg.Add(1)
	go q.processingLoop(ctx)

	return nil
}

// Stop gracefully shuts down the queue
func (q *MessageQueue) Stop() error {
	q.mu.Lock()
	if !q.running {
		q.mu.Unlock()
		return nil
	}
	q.running = false
	close(q.done)
	q.mu.Unlock()

	q.wg.Wait()
	return nil
}

// Add adds a message to the queue
func (q *MessageQueue) Add(msg *nats.Msg) error {
	select {
	case q.msgChan <- msg:
		return nil
	default:
		// Queue is full, wait for processing delay
		time.Sleep(q.config.ProcessDelay)
		q.msgChan <- msg
		return nil
	}
}

func (q *MessageQueue) processingLoop(ctx context.Context) {
	defer q.wg.Done()

	batch := make([]*nats.Msg, 0, q.config.BatchSize)
	ticker := time.NewTicker(q.config.ProcessDelay)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Process remaining messages before exit
			if len(batch) > 0 {
				q.processor.ProcessBatch(context.Background(), batch)
			}
			return
		case <-q.done:
			// Process remaining messages before exit
			if len(batch) > 0 {
				q.processor.ProcessBatch(context.Background(), batch)
			}
			return
		case msg, ok := <-q.msgChan:
			if !ok {
				return
			}
			batch = append(batch, msg)
			if len(batch) >= q.config.BatchSize {
				q.processor.ProcessBatch(ctx, batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				q.processor.ProcessBatch(ctx, batch)
				batch = batch[:0]
			}
		}
	}
}
