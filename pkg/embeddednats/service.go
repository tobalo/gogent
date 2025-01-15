package embeddednats

import (
	"fmt"
	"log"
	"time"

	server "github.com/nats-io/nats-server/v2/server"
	nats "github.com/nats-io/nats.go"
	"github.com/tobalo/gogent/pkg/shared"
)

type NatsService struct {
	server *server.Server
	js     nats.JetStreamContext
	port   int
}

func NewNatsService(port int) (*NatsService, error) {
	// Configure NATS server with JetStream enabled
	opts := &server.Options{
		Port:           port,
		JetStream:      true,
		StoreDir:       "data/jetstream", // Store JetStream data
		MaxPayload:     server.MAX_PAYLOAD_SIZE,
		WriteDeadline:  10 * time.Second,
		MaxPending:     server.MAX_PENDING_SIZE,
		MaxControlLine: 4096,
		MaxConn:        64 * 1024,
		MaxSubs:        0, // unlimited
		NoLog:          false,
		NoSigs:         false,
		Debug:          true, // Enable debug logging
		Trace:          true, // Enable trace logging
	}

	s, err := server.NewServer(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create NATS server: %w", err)
	}

	// Configure server logging
	s.ConfigureLogger()

	return &NatsService{
		server: s,
		port:   port,
	}, nil
}

func (n *NatsService) Start() error {
	go n.server.Start()

	// Wait for server to be ready
	if !n.server.ReadyForConnections(10 * time.Second) {
		return fmt.Errorf("NATS server failed to start")
	}

	// Connect to the server
	nc, err := nats.Connect(fmt.Sprintf("nats://localhost:%d", n.port))
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Create JetStream context
	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Create the stream
	_, err = js.StreamInfo(shared.StreamName)
	if err != nil {
		// Create the stream if it doesn't exist
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     shared.StreamName,
			Subjects: []string{shared.SubjectName},
			Storage:  nats.FileStorage,
			MaxAge:   24 * time.Hour, // Keep messages for 24 hours
		})
		if err != nil {
			nc.Close()
			return fmt.Errorf("failed to create stream: %w", err)
		}
		log.Printf("Created NATS stream: %s", shared.StreamName)
	}

	n.js = js
	log.Printf("NATS server started on port %d with JetStream enabled", n.port)
	return nil
}

func (n *NatsService) Stop() error {
	if n.server != nil {
		n.server.Shutdown()
		log.Println("NATS server stopped")
	}
	return nil
}

func (n *NatsService) GetJetStream() (nats.JetStreamContext, error) {
	if n.js == nil {
		nc, err := nats.Connect(fmt.Sprintf("nats://localhost:%d", n.port))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to NATS: %w", err)
		}

		js, err := nc.JetStream()
		if err != nil {
			nc.Close()
			return nil, fmt.Errorf("failed to create JetStream context: %w", err)
		}
		n.js = js
	}
	return n.js, nil
}

func (n *NatsService) CreateConsumer(streamName, durableName string) error {
	js, err := n.GetJetStream()
	if err != nil {
		return err
	}

	_, err = js.AddConsumer(streamName, &nats.ConsumerConfig{
		Durable:       durableName,
		AckPolicy:     nats.AckExplicitPolicy,
		MaxDeliver:    -1,
		FilterSubject: shared.SubjectName,
		MaxAckPending: -1,
	})
	if err != nil {
		log.Printf("Error creating consumer: %v", err)
	} else {
		log.Printf("Created consumer: %s", durableName)
	}
	return err
}

func (n *NatsService) Publish(subject string, message []byte) error {
	js, err := n.GetJetStream()
	if err != nil {
		return err
	}

	_, err = js.Publish(subject, message)
	if err != nil {
		log.Printf("Error publishing message: %v", err)
	} else {
		log.Printf("Published message to %s", subject)
	}
	return err
}

func (n *NatsService) Subscribe(subject string, handler func(message []byte) error) error {
	js, err := n.GetJetStream()
	if err != nil {
		return err
	}

	// Create a durable consumer
	err = n.CreateConsumer(shared.StreamName, shared.ConsumerName)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	// Subscribe with the durable consumer
	sub, err := js.PullSubscribe(subject, shared.ConsumerName)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	log.Printf("Listening on subject: %s", subject)

	// Start processing messages
	go func() {
		for {
			msgs, err := sub.Fetch(1, nats.MaxWait(time.Second))
			if err != nil {
				if err != nats.ErrTimeout {
					log.Printf("Error fetching message: %v", err)
				}
				continue
			}

			for _, msg := range msgs {
				if err := handler(msg.Data); err != nil {
					log.Printf("Error handling message: %v", err)
				} else {
					log.Printf("Successfully processed message")
				}
				msg.Ack()
			}
		}
	}()

	return nil
}
