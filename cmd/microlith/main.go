package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/tobalo/gogent/pkg/agent"
	embeddednats "github.com/tobalo/gogent/pkg/embeddednats"
	"github.com/tobalo/gogent/pkg/shared"
)

func main() {
	// Configure logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC)
	log.SetPrefix("[AGENT SIG] ")

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize embedded NATS server
	log.Println("Starting embedded NATS server...")
	natsService, err := embeddednats.NewNatsService(shared.DefaultNATSPort)
	if err != nil {
		log.Fatalf("Failed to create NATS service: %v", err)
	}

	if err := natsService.Start(); err != nil {
		log.Fatalf("Failed to start NATS service: %v", err)
	}
	defer natsService.Stop()
	log.Println("NATS server started successfully")

	// Give NATS server time to initialize
	time.Sleep(1 * time.Second)

	// Initialize agent service
	log.Println("Initializing agent service...")
	agentService, err := agent.NewService(agent.Config{
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		NATSUrl:      shared.DefaultNATSURL,
		AgentName:    shared.DefaultAgentName,
		Instructions: shared.DefaultAgentInstructions,
		Model:        shared.DefaultAgentModel,
	})
	if err != nil {
		log.Fatalf("Failed to create agent service: %v", err)
	}
	defer agentService.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := agentService.Start(ctx); err != nil {
		log.Fatalf("Failed to start agent service: %v", err)
	}
	log.Println("Agent service started successfully")
	log.Println("Ready to process messages on AGENT.TECHNICAL.SUPPORT")

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	log.Printf("Received signal %v, shutting down gracefully...", sig)
}
