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
	// NATSPort is the port for NATS server
	NATSPort = 4222
	// NATSURL is the URL for NATS connection
	NATSURL = "nats://localhost:4222"
)

// LLM Provider Types
const (
	// ProviderOllama represents the Ollama LLM provider
	ProviderOllama = "OLLAMA"
	// ProviderOpenAI represents the OpenAI provider
	ProviderOpenAI = "OPEN_AI"
	// ProviderAzure represents the Azure OpenAI provider
	ProviderAzure = "AZURE"
	// ProviderAzureAD represents the Azure AD provider
	ProviderAzureAD = "AZURE_AD"
	// ProviderCloudflareAzure represents the Cloudflare Azure provider
	ProviderCloudflareAzure = "CLOUDFLARE_AZURE"
	// ProviderGemini represents the Google Gemini provider
	ProviderGemini = "GEMINI"
	// ProviderClaude represents the Anthropic Claude provider
	ProviderClaude = "CLAUDE"
	// ProviderDeepSeek represents the DeepSeek provider
	ProviderDeepSeek = "DEEPSEEK"
)

// Model Names
const (
	// ModelPhi35 is the Ollama Phi-3.5 model
	ModelPhi35 = "phi-3.5"
	// ModelGemini is the Gemini model name
	ModelGemini = "gemini-1.5-flash-8b"
	// ModelGPT4 is OpenAI's GPT-4 model
	ModelGPT4 = "gpt-4"
	// ModelClaude3 is Claude's model
	ModelClaude3 = "claude-3"
)

// Agent Configuration Constants
const (
	// AgentName is the name for the agent
	AgentName = "Agent Sig"
	// Provider is the LLM provider
	Provider = ProviderOllama
	// AgentModel is the model used by the agent
	AgentModel = ModelPhi35
	// AgentInstructions are the instructions for the agent
	AgentInstructions = "You are a technical analyst that executes natural language reporting from technical information and raw SIGINT data. Analyze system logs and provide concise, actionable insights."
)
