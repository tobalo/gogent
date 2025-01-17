package agent

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	swarmgo "github.com/prathyushnallamothu/swarmgo"
	llm "github.com/prathyushnallamothu/swarmgo/llm"
	"github.com/tobalo/gogent/pkg/db"
	"github.com/tobalo/gogent/pkg/shared"
)

// Config holds the configuration for the agent service
type Config struct {
	APIKey       string // Optional: required for non-Ollama providers
	NATSUrl      string
	AgentName    string
	Instructions string
	Model        string
	Provider     string // LLM provider (ollama, openai, azure, etc.)
	DBPath       string // Path to SQLite database
}

// Service manages the agent and its NATS connection
type Service struct {
	config Config
	agent  *swarmgo.Agent
	swarm  *swarmgo.Swarm
	nc     *nats.Conn
	js     nats.JetStreamContext
	dbConn *sql.DB
}

// LogMessage represents the structure of log messages received
type LogMessage struct {
	Timestamp string                 `json:"timestamp"`
	Hostname  string                 `json:"hostname"`
	Severity  string                 `json:"severity"`
	Service   string                 `json:"service"`
	Message   string                 `json:"message"`
	Context   map[string]interface{} `json:"context"`
}

// NewService creates a new agent service
func NewService(cfg Config) (*Service, error) {
	// Set defaults if not provided
	if cfg.Provider == "" {
		cfg.Provider = shared.Provider
	}
	if cfg.Model == "" {
		cfg.Model = shared.AgentModel
	}
	if cfg.DBPath == "" {
		cfg.DBPath = filepath.Join("data", "agent.db")
	}

	// Convert provider to uppercase for LLMProvider matching
	provider := strings.ToUpper(cfg.Provider)

	// Validate provider-specific requirements
	if provider != shared.ProviderOllama && cfg.APIKey == "" {
		return nil, fmt.Errorf("API key is required for %s provider", provider)
	}

	// Initialize database
	dbConn, err := db.InitDB(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Create swarm instance with appropriate provider
	var llmProvider llm.LLMProvider
	switch provider {
	case shared.ProviderOpenAI:
		llmProvider = llm.LLMProvider(shared.ProviderOpenAI)
	case shared.ProviderAzure:
		llmProvider = llm.LLMProvider(shared.ProviderAzure)
	case shared.ProviderAzureAD:
		llmProvider = llm.LLMProvider(shared.ProviderAzureAD)
	case shared.ProviderCloudflareAzure:
		llmProvider = llm.LLMProvider(shared.ProviderCloudflareAzure)
	case shared.ProviderGemini:
		llmProvider = llm.LLMProvider(shared.ProviderGemini)
	case shared.ProviderClaude:
		llmProvider = llm.LLMProvider(shared.ProviderClaude)
	case shared.ProviderOllama:
		llmProvider = llm.LLMProvider(shared.ProviderOllama)
	case shared.ProviderDeepSeek:
		llmProvider = llm.LLMProvider(shared.ProviderDeepSeek)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	swarm := swarmgo.NewSwarm(cfg.APIKey, llmProvider)
	agent := &swarmgo.Agent{
		Name:         cfg.AgentName,
		Instructions: cfg.Instructions,
		Model:        cfg.Model,
	}

	// Connect to NATS
	nc, err := nats.Connect(cfg.NATSUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Create JetStream context
	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	return &Service{
		config: cfg,
		agent:  agent,
		swarm:  swarm,
		nc:     nc,
		js:     js,
		dbConn: dbConn,
	}, nil
}

// Start begins listening for messages on the configured subject
func (s *Service) Start(ctx context.Context) error {
	// Subscribe directly to the subject
	sub, err := s.nc.Subscribe(shared.SubjectName, func(msg *nats.Msg) {
		s.handleMessage(ctx, msg)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	// Ensure subscription is properly cleaned up
	go func() {
		<-ctx.Done()
		sub.Unsubscribe()
	}()

	log.Printf("Agent service started with %s provider and %s model, listening on %s",
		strings.ToUpper(s.config.Provider), s.config.Model, shared.SubjectName)
	return nil
}

// handleMessage processes a single message through the agent
func (s *Service) handleMessage(ctx context.Context, msg *nats.Msg) {
	var logMsg LogMessage
	if err := json.Unmarshal(msg.Data, &logMsg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	// Format the message for the agent
	prompt := fmt.Sprintf(`Analyze this technical log entry and provide insights:
Timestamp: %s
Host: %s
Severity: %s
Service: %s
Message: %s
Additional Context: %v`,
		logMsg.Timestamp,
		logMsg.Hostname,
		logMsg.Severity,
		logMsg.Service,
		logMsg.Message,
		logMsg.Context,
	)

	messages := []llm.Message{
		{Role: llm.RoleUser, Content: prompt},
	}

	log.Printf("Processing log message from %s [%s] %s", logMsg.Hostname, logMsg.Severity, logMsg.Message)

	// Add timeout for agent processing
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	response, err := s.swarm.Run(ctx, s.agent, messages, nil, "", false, false, 5, true)
	if err != nil {
		log.Printf("Error processing message: %v", err)
		return
	}

	analysis := response.Messages[len(response.Messages)-1].Content

	// Store log in database
	contextJSON, err := json.Marshal(logMsg.Context)
	if err != nil {
		log.Printf("Error marshaling context: %v", err)
		return
	}

	logEntry := db.LogEntry{
		Timestamp: logMsg.Timestamp,
		Hostname:  logMsg.Hostname,
		Severity:  logMsg.Severity,
		Service:   logMsg.Service,
		Message:   logMsg.Message,
		Context:   string(contextJSON),
		Analysis:  analysis,
	}

	if err := db.InsertLogEntry(logEntry); err != nil {
		log.Printf("Error storing log in database: %v", err)
	}

	log.Printf("Analysis complete for %s: %s", logMsg.Service, analysis[:100]+"...")

	// Prepare response
	responseData, err := json.Marshal(map[string]interface{}{
		"original_message": logMsg,
		"analysis":         analysis,
		"timestamp":        time.Now().Format(time.RFC3339),
	})
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return
	}

	// Send response if reply subject is provided
	if msg.Reply != "" {
		if err := msg.Respond(responseData); err != nil {
			log.Printf("Error sending response: %v", err)
		}
	}
}

// Stop gracefully shuts down the service
func (s *Service) Stop() error {
	if s.nc != nil {
		s.nc.Close()
	}
	if s.dbConn != nil {
		s.dbConn.Close()
	}
	return nil
}
