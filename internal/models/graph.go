package models

// Node represents a Terraform resource or module in the dependency graph.
type Node struct {
	ID         string            `json:"id"`         // e.g. "aws_vpc.main" or "module.vpc"
	Type       string            `json:"type"`       // e.g. "resource", "module", "data"
	Provider   string            `json:"provider"`   // e.g. "aws", "google"
	Name       string            `json:"name"`       // e.g. "main"
	Label      string            `json:"label"`      // Human-readable label for the graph
	File       string            `json:"file"`       // Source file path
	LineStart  int               `json:"lineStart"`  // Line number in source file
	Properties map[string]string `json:"properties"` // Extracted key-value pairs
	Variables  []string          `json:"variables"`  // Used variable names (e.g. "vpc_cidr")
}

// Edge represents a dependency between two nodes.
type Edge struct {
	From     string `json:"from"`     // Source node ID
	To       string `json:"to"`       // Target node ID
	EdgeType string `json:"edgeType"` // "explicit" (depends_on) or "implicit" (reference)
	Label    string `json:"label"`    // e.g. "vpc_id", "depends_on"
}

// Environment represents a set of variable values from a .tfvars file.
type Environment struct {
	Name     string            `json:"name"`     // e.g. "prod", "stg", "dev"
	FilePath string            `json:"filePath"` // e.g. "envs/prod.tfvars"
	Values   map[string]string `json:"values"`   // Flattened key-value pairs
}

// Variable represents a Terraform variable definition.
type Variable struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Default string `json:"default"`
}

// Graph is the full parsed result returned to the frontend.
type Graph struct {
	Nodes        []Node        `json:"nodes"`
	Edges        []Edge        `json:"edges"`
	Environments []Environment `json:"environments"`
	Variables    []Variable    `json:"variables"`
	Files        []string      `json:"files"` // List of all .tf files found
}
