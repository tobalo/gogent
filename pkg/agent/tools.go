package agent

import (
	"fmt"
	"time"

	swarmgo "github.com/prathyushnallamothu/swarmgo"
)

// ExternalSystemsConfig holds configuration for external system connections
type ExternalSystemsConfig struct {
	PalantirConfig struct {
		BaseURL    string `json:"baseUrl"`
		Token      string `json:"token"`
		Dataset    string `json:"dataset"`
		Project    string `json:"project"`
		MaxRetries int    `json:"maxRetries"`
	} `json:"palantir"`

	ServiceNowConfig struct {
		Instance     string `json:"instance"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		ClientID     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
	} `json:"serviceNow"`

	SplunkConfig struct {
		Host  string `json:"host"`
		Port  int    `json:"port"`
		Token string `json:"token"`
		Index string `json:"index"`
		SSL   bool   `json:"ssl"`
	} `json:"splunk"`

	JiraConfig struct {
		URL      string `json:"url"`
		Username string `json:"username"`
		APIToken string `json:"apiToken"`
		Project  string `json:"project"`
	} `json:"jira"`
}

// PalantirFoundry Tools

func uploadToPalantir(args map[string]interface{}, contextVariables map[string]interface{}) swarmgo.Result {
	dataset := args["dataset"].(string)
	data := args["data"].(map[string]interface{})

	// Example implementation
	return swarmgo.Result{
		Success: true,
		Data: map[string]interface{}{
			"message":     fmt.Sprintf("Successfully uploaded data to dataset %s", dataset),
			"timestamp":   time.Now().String(),
			"recordCount": len(data),
		},
	}
}

func createPalantirAnalysis(args map[string]interface{}, contextVariables map[string]interface{}) swarmgo.Result {
	analysisName := args["name"].(string)
	analysisType := args["type"].(string)
	parameters := args["parameters"].(map[string]interface{})

	return swarmgo.Result{
		Success: true,
		Data: map[string]interface{}{
			"analysisId": "pa_123456",
			"status":     "created",
			"name":       analysisName,
			"type":       analysisType,
			"parameters": parameters,
		},
	}
}

// ServiceNow Tools

func createServiceNowIncident(args map[string]interface{}, contextVariables map[string]interface{}) swarmgo.Result {
	shortDescription := args["shortDescription"].(string)
	priority := args["priority"].(string)
	assignmentGroup := args["assignmentGroup"].(string)

	return swarmgo.Result{
		Success: true,
		Data: map[string]interface{}{
			"incidentNumber":   "INC0010234",
			"status":           "New",
			"shortDescription": shortDescription,
			"priority":         priority,
			"assignmentGroup":  assignmentGroup,
			"createdAt":        time.Now().String(),
		},
	}
}

func updateServiceNowTicket(args map[string]interface{}, contextVariables map[string]interface{}) swarmgo.Result {
	ticketNumber := args["ticketNumber"].(string)
	status := args["status"].(string)
	notes := args["workNotes"].(string)

	return swarmgo.Result{
		Success: true,
		Data: map[string]interface{}{
			"ticketNumber": ticketNumber,
			"newStatus":    status,
			"updateTime":   time.Now().String(),
			"workNotes":    notes,
		},
	}
}

// Splunk Tools

func querySplunk(args map[string]interface{}, contextVariables map[string]interface{}) swarmgo.Result {
	searchQuery := args["query"].(string)
	timeRange := args["timeRange"].(string)

	return swarmgo.Result{
		Success: true,
		Data: map[string]interface{}{
			"searchId":      "search_789012",
			"query":         searchQuery,
			"timeRange":     timeRange,
			"resultCount":   150,
			"executionTime": "2.5s",
		},
	}
}

func createSplunkAlert(args map[string]interface{}, contextVariables map[string]interface{}) swarmgo.Result {
	alertName := args["name"].(string)
	query := args["query"].(string)
	threshold := args["threshold"].(float64)

	return swarmgo.Result{
		Success: true,
		Data: map[string]interface{}{
			"alertId":   "alert_345678",
			"name":      alertName,
			"query":     query,
			"threshold": threshold,
			"status":    "enabled",
		},
	}
}

// Jira Tools

func createJiraIssue(args map[string]interface{}, contextVariables map[string]interface{}) swarmgo.Result {
	summary := args["summary"].(string)
	description := args["description"].(string)
	issueType := args["issueType"].(string)
	priority := args["priority"].(string)

	return swarmgo.Result{
		Success: true,
		Data: map[string]interface{}{
			"issueKey":    "PROJ-1234",
			"summary":     summary,
			"description": description,
			"type":        issueType,
			"priority":    priority,
			"status":      "Open",
			"created":     time.Now().String(),
		},
	}
}

func updateJiraIssue(args map[string]interface{}, contextVariables map[string]interface{}) swarmgo.Result {
	issueKey := args["issueKey"].(string)
	status := args["status"].(string)
	comment := args["comment"].(string)

	return swarmgo.Result{
		Success: true,
		Data: map[string]interface{}{
			"issueKey":   issueKey,
			"newStatus":  status,
			"comment":    comment,
			"updateTime": time.Now().String(),
		},
	}
}

// ExampleEnterpriseAgent provides an example configuration for enterprise system integration
var ExampleEnterpriseAgent = &swarmgo.Agent{
	Name: "Enterprise Integration Agent",
	Instructions: `You are an AI agent capable of interacting with enterprise systems including 
		PalantirFoundry, ServiceNow, Splunk, and Jira. You can create and update tickets, 
		perform data analysis, manage incidents, and handle various integration tasks.`,
	Functions: []swarmgo.AgentFunction{
		{
			Name:        "uploadToPalantir",
			Description: "Upload data to a Palantir Foundry dataset",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"dataset": map[string]interface{}{
						"type":        "string",
						"description": "The target dataset reference",
					},
					"data": map[string]interface{}{
						"type":        "object",
						"description": "The data to upload",
					},
				},
				"required": []string{"dataset", "data"},
			},
			Function: uploadToPalantir,
		},
		{
			Name:        "createServiceNowIncident",
			Description: "Create a new incident in ServiceNow",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"shortDescription": map[string]interface{}{
						"type":        "string",
						"description": "Brief description of the incident",
					},
					"priority": map[string]interface{}{
						"type":        "string",
						"description": "Incident priority (1-5)",
					},
					"assignmentGroup": map[string]interface{}{
						"type":        "string",
						"description": "Group to assign the incident to",
					},
				},
				"required": []string{"shortDescription", "priority", "assignmentGroup"},
			},
			Function: createServiceNowIncident,
		},
		{
			Name:        "querySplunk",
			Description: "Execute a search query in Splunk",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Splunk search query",
					},
					"timeRange": map[string]interface{}{
						"type":        "string",
						"description": "Time range for the search (e.g., '-24h')",
					},
				},
				"required": []string{"query", "timeRange"},
			},
			Function: querySplunk,
		},
		{
			Name:        "createJiraIssue",
			Description: "Create a new issue in Jira",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"summary": map[string]interface{}{
						"type":        "string",
						"description": "Issue summary",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Detailed description",
					},
					"issueType": map[string]interface{}{
						"type":        "string",
						"description": "Type of issue (Bug, Story, Task)",
					},
					"priority": map[string]interface{}{
						"type":        "string",
						"description": "Issue priority",
					},
				},
				"required": []string{"summary", "description", "issueType", "priority"},
			},
			Function: createJiraIssue,
		},
	},
	Model: "gpt-4",
}
