package shared

// NATS Stream Constants
const (
	// StreamName is the name of the NATS stream
	StreamName = "AGENT_STREAM"
	// ConsumerName is the name of the NATS consumer
	ConsumerName = "AGENT_CONSUMER"
	// SubjectName is the NATS subject for agent technical support messages
	SubjectName = "agent.technical.support"
)

// NATS Configuration Constants
const (
	// DefaultNATSPort is the default port for NATS server
	DefaultNATSPort = 4222
	// DefaultNATSURL is the default URL for NATS connection
	DefaultNATSURL = "nats://localhost:4222"
)

// Agent Configuration Constants
const (
	// DefaultAgentName is the default name for the agent
	DefaultAgentName = "Agent Sig"
	// DefaultAgentModel is the default model used by the agent
	DefaultAgentModel = "gemini-1.5-flash-8b"
	// DefaultAgentInstructions are the default instructions for the agent
	DefaultAgentInstructions = "You are a technical analyst that executes natural language reporting from technical information and raw SIGINT data. Analyze system logs and provide concise, actionable insights."
)
