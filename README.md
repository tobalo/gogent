# Gogent - Agentic Worker AI in Golang

Gogent is a distributed log analysis system that uses embedded NATS messaging and LLM-powered agents to process and analyze system logs in real-time.

## Architecture

```mermaid
flowchart LR
    T[test.sh] -->|agent.technical.support| N[Embedded NATS Server]
    subgraph Gogent Service
        S[AGENT_STREAM] -->|AGENT_CONSUMER| A[Agent Sig]
        A -->|Format Message| L[gemini-1.5-flash-8b]
        L -->|Analysis| A
    end
    N -->|JetStream| S
    A -->|Response| N
```

## Core Components

- **Embedded NATS Server**: Handles message queuing and distribution
- **Agent Service**: Processes messages using LLM
- **Gemini Integration**: Provides AI-powered log analysis
- **JetStream**: Persistent message storage

## Message Flow

1. Log messages are published to `agent.technical.support` subject
2. Agent subscribes to messages and formats them for LLM processing
3. Gemini API analyzes the log content using gemini-1.5-flash-8b model
4. Analysis results are sent back through NATS if reply subject exists
5. Messages are persisted in `AGENT_STREAM` with `AGENT_CONSUMER` subscription

## Technical Details

### Agent Configuration

```go
type Config struct {
    GeminiAPIKey string
    NATSUrl      string
    AgentName    string
    Instructions string
    Model        string
}
```

### Default Configuration
```go
// NATS Configuration
StreamName    = "AGENT_STREAM"
ConsumerName  = "AGENT_CONSUMER"
SubjectName   = "agent.technical.support"
DefaultNATSPort = 4222
DefaultNATSURL  = "nats://localhost:4222"

// Agent Configuration
DefaultAgentName = "Agent Sig"
DefaultAgentModel = "gemini-1.5-flash-8b"
DefaultAgentInstructions = "You are a technical analyst that executes natural language reporting from technical information and raw SIGINT data. Analyze system logs and provide concise, actionable insights."
```

### Message Structure
```go
type LogMessage struct {
    Timestamp string
    Hostname  string
    Severity  string
    Service   string
    Message   string
    Context   map[string]interface{}
}
```

## Setup

### Prerequisites

- Go 1.23.4 or later
- Gemini API key from Google AI Studio
- Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/tobalo/gogent.git
cd gogent
```

2. Install dependencies:
```bash
go mod download
```

3. Create environment file:
```bash
cp .env.example .env
```

4. Configure your .env file:
```sh
GEMINI_API_KEY=your_api_key_here    # Required: API key from Google AI Studio
```

### Running

1. Start the agent:
```bash
go run cmd/microlith/main.go
```

2. The agent will initialize with:
   - Embedded NATS server on port 4222
   - JetStream enabled for message persistence
   - Agent subscribed to agent.technical.support
   - 30-second timeout for LLM processing

3. Monitor the startup logs:
```sh
[AGENT SIG] 2025/01/16 10:08:14.123 UTC Starting embedded NATS server...
[AGENT SIG] 2025/01/16 10:08:14.234 UTC NATS server started successfully
[AGENT SIG] 2025/01/16 10:08:15.345 UTC Agent service started successfully
[AGENT SIG] 2025/01/16 10:08:15.456 UTC Ready to process messages
```

## Sample Usage

### Publishing Messages

Messages can be published using the NATS CLI:

```bash
nats pub agent.technical.support '{
    "timestamp": "2025-01-15T02:14:23.123Z",
    "hostname": "web-server-01",
    "severity": "ERROR",
    "service": "nginx",
    "message": "Failed to bind to port 80: Address already in use",
    "context": {
        "pid": 1234,
        "user": "www-data"
    }
}'
```

## Features

- Real-time log processing
- AI-powered log analysis
- Distributed message handling
- Persistent message storage via JetStream
- Configurable agent behavior
- Automatic message formatting for LLM processing
- Response handling with original context