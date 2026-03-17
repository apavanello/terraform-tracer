package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apavanello/terraform-tracer/internal/gitclone"
	"github.com/apavanello/terraform-tracer/internal/models"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

// Parse reads all .tf and .tfvars files from the given directory and returns a Graph.
func Parse(rootDir string) (*models.Graph, error) {
	// Create a temp cache dir for git clones
	cacheDir, err := os.MkdirTemp("", "terraform-tracer-modules-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create module cache dir: %w", err)
	}
	defer os.RemoveAll(cacheDir)

	graph := &models.Graph{
		Nodes:        []models.Node{},
		Edges:        []models.Edge{},
		Environments: []models.Environment{},
		Variables:    []models.Variable{},
		Files:        []string{},
	}

	if err := parseDir(rootDir, rootDir, cacheDir, graph); err != nil {
		return nil, err
	}

	return graph, nil
}

// parseDir collects and parses all .tf and .tfvars files in a directory.
func parseDir(dir string, rootDir string, cacheDir string, graph *models.Graph) error {
	// Collect all .tf files
	tfFiles, err := collectFiles(dir, ".tf")
	if err != nil {
		return fmt.Errorf("failed to collect .tf files: %w", err)
	}
	graph.Files = append(graph.Files, tfFiles...)

	// Parse each .tf file
	for _, filePath := range tfFiles {
		if err := parseTFFile(filePath, rootDir, cacheDir, graph); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to parse %s: %v\n", filePath, err)
		}
	}

	// Collect and parse .tfvars files for environments
	tfvarsFiles, err := collectFiles(dir, ".tfvars")
	if err != nil {
		return fmt.Errorf("failed to collect .tfvars files: %w", err)
	}
	for _, filePath := range tfvarsFiles {
		if err := parseTFVarsFile(filePath, rootDir, graph); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to parse tfvars %s: %v\n", filePath, err)
		}
	}

	return nil
}

func collectFiles(rootDir string, ext string) ([]string, error) {
	var files []string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip hidden directories, node_modules, .terraform
		base := filepath.Base(path)
		if info.IsDir() && (strings.HasPrefix(base, ".") || base == "node_modules") {
			return filepath.SkipDir
		}
		if !info.IsDir() && strings.HasSuffix(path, ext) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func parseTFFile(filePath string, rootDir string, cacheDir string, graph *models.Graph) error {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	file, diags := hclsyntax.ParseConfig(src, filePath, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return fmt.Errorf("HCL parse error: %s", diags.Error())
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return fmt.Errorf("unexpected body type in %s", filePath)
	}

	relPath, _ := filepath.Rel(rootDir, filePath)

	for _, block := range body.Blocks {
		switch block.Type {
		case "resource":
			if len(block.Labels) >= 2 {
				props, vars := extractAttributes(block.Body)
				node := models.Node{
					ID:         block.Labels[0] + "." + block.Labels[1],
					Type:       "resource",
					Provider:   extractProvider(block.Labels[0]),
					Name:       block.Labels[1],
					Label:      block.Labels[0] + "." + block.Labels[1],
					File:       relPath,
					LineStart:  block.DefRange().Start.Line,
					Properties: props,
					Variables:  vars,
				}
				graph.Nodes = append(graph.Nodes, node)

				// Extract edges
				extractEdges(block, node.ID, graph)
			}

		case "data":
			if len(block.Labels) >= 2 {
				node := models.Node{
					ID:        "data." + block.Labels[0] + "." + block.Labels[1],
					Type:      "data",
					Provider:  extractProvider(block.Labels[0]),
					Name:      block.Labels[1],
					Label:     "data." + block.Labels[0] + "." + block.Labels[1],
					File:      relPath,
					LineStart: block.DefRange().Start.Line,
				}
				graph.Nodes = append(graph.Nodes, node)
			}

		case "module":
			if len(block.Labels) >= 1 {
				props, vars := extractAttributes(block.Body)
				node := models.Node{
					ID:         "module." + block.Labels[0],
					Type:       "module",
					Provider:   "terraform",
					Name:       block.Labels[0],
					Label:      "module." + block.Labels[0],
					File:       relPath,
					LineStart:  block.DefRange().Start.Line,
					Properties: props,
					Variables:  vars,
				}
				graph.Nodes = append(graph.Nodes, node)
				extractEdges(block, node.ID, graph)

				// If the module source is a git repository, clone and parse it
				if source, ok := props["source"]; ok && gitclone.IsGitSource(source) {
					localPath, err := gitclone.CloneModule(source, cacheDir)
					if err != nil {
						fmt.Fprintf(os.Stderr, "warning: failed to clone module %s: %v\n", source, err)
					} else {
						// Recursively parse the cloned module
						if err := parseDir(localPath, rootDir, cacheDir, graph); err != nil {
							fmt.Fprintf(os.Stderr, "warning: failed to parse cloned module %s: %v\n", source, err)
						}
					}
				}
			}

		case "variable":
			if len(block.Labels) >= 1 {
				v := models.Variable{
					Name: block.Labels[0],
				}
				attrs, _ := extractAttributes(block.Body)
				if t, ok := attrs["type"]; ok {
					v.Type = t
				}
				if d, ok := attrs["default"]; ok {
					v.Default = d
				}
				graph.Variables = append(graph.Variables, v)
			}
		}
	}

	return nil
}

func extractProvider(resourceType string) string {
	parts := strings.SplitN(resourceType, "_", 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}

func extractAttributes(body *hclsyntax.Body) (map[string]string, []string) {
	attrs := make(map[string]string)
	varSet := make(map[string]bool)
	var vars []string

	for name, attr := range body.Attributes {
		val, diags := attr.Expr.Value(nil)
		if !diags.HasErrors() {
			attrs[name] = ctyToString(val)
		}

		for _, traversal := range attr.Expr.Variables() {
			if len(traversal) >= 2 {
				root, ok := traversal[0].(hcl.TraverseRoot)
				if ok && root.Name == "var" {
					if attrName, ok := traversal[1].(hcl.TraverseAttr); ok {
						if !varSet[attrName.Name] {
							varSet[attrName.Name] = true
							vars = append(vars, attrName.Name)
						}
					}
				}
			}
		}
	}
	return attrs, vars
}

func ctyToString(val cty.Value) string {
	if !val.IsKnown() || val.IsNull() {
		return ""
	}
	switch val.Type() {
	case cty.String:
		return val.AsString()
	case cty.Bool:
		if val.True() {
			return "true"
		}
		return "false"
	case cty.Number:
		bf := val.AsBigFloat()
		return bf.Text('f', -1)
	default:
		if val.Type().IsObjectType() || val.Type().IsMapType() {
			var parts []string
			for it := val.ElementIterator(); it.Next(); {
				k, v := it.Element()
				parts = append(parts, fmt.Sprintf("%s: %s", ctyToString(k), ctyToString(v)))
			}
			return "{" + strings.Join(parts, ", ") + "}"
		}

		if val.Type().IsTupleType() || val.Type().IsListType() || val.Type().IsSetType() {
			var parts []string
			for it := val.ElementIterator(); it.Next(); {
				_, v := it.Element()
				parts = append(parts, ctyToString(v))
			}
			return "[" + strings.Join(parts, ", ") + "]"
		}

		return fmt.Sprintf("%#v", val)
	}
}

func extractEdges(block *hclsyntax.Block, fromID string, graph *models.Graph) {
	body := block.Body

	// 1) Explicit depends_on
	if attr, exists := body.Attributes["depends_on"]; exists {
		for _, traversal := range attr.Expr.Variables() {
			targetID := traversalToID(traversal)
			if targetID != "" {
				graph.Edges = append(graph.Edges, models.Edge{
					From:     fromID,
					To:       targetID,
					EdgeType: "explicit",
					Label:    "depends_on",
				})
			}
		}
	}

	// 2) Implicit dependencies from all attribute expressions
	for name, attr := range body.Attributes {
		if name == "depends_on" {
			continue
		}
		for _, traversal := range attr.Expr.Variables() {
			targetID := traversalToID(traversal)
			if targetID != "" && targetID != fromID {
				graph.Edges = append(graph.Edges, models.Edge{
					From:     fromID,
					To:       targetID,
					EdgeType: "implicit",
					Label:    name,
				})
			}
		}
	}

	// 3) Also walk nested blocks (e.g. inline provisioners, lifecycle, etc.)
	for _, nestedBlock := range body.Blocks {
		for _, attr := range nestedBlock.Body.Attributes {
			for _, traversal := range attr.Expr.Variables() {
				targetID := traversalToID(traversal)
				if targetID != "" && targetID != fromID {
					graph.Edges = append(graph.Edges, models.Edge{
						From:     fromID,
						To:       targetID,
						EdgeType: "implicit",
						Label:    attr.Name,
					})
				}
			}
		}
	}
}

// traversalToID converts an HCL traversal to a Terraform resource ID.
// e.g. aws_vpc.main.id → "aws_vpc.main"
func traversalToID(traversal hcl.Traversal) string {
	if len(traversal) < 2 {
		return ""
	}
	root, ok := traversal[0].(hcl.TraverseRoot)
	if !ok {
		return ""
	}

	rootName := root.Name

	// Known Terraform top-level keywords that are NOT resource references
	skip := map[string]bool{
		"var": true, "local": true, "each": true,
		"self": true, "path": true, "terraform": true,
		"count": true, "null": true,
	}
	if skip[rootName] {
		return ""
	}

	// Handle "module.xxx" references
	if rootName == "module" {
		if len(traversal) >= 2 {
			if attr, ok := traversal[1].(hcl.TraverseAttr); ok {
				return "module." + attr.Name
			}
		}
		return ""
	}

	// Handle "data.type.name" references
	if rootName == "data" {
		if len(traversal) >= 3 {
			t, ok1 := traversal[1].(hcl.TraverseAttr)
			n, ok2 := traversal[2].(hcl.TraverseAttr)
			if ok1 && ok2 {
				return "data." + t.Name + "." + n.Name
			}
		}
		return ""
	}

	// Standard resource reference: aws_vpc.main.id → aws_vpc.main
	if attr, ok := traversal[1].(hcl.TraverseAttr); ok {
		return rootName + "." + attr.Name
	}

	return ""
}

func parseTFVarsFile(filePath string, rootDir string, graph *models.Graph) error {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	file, diags := hclsyntax.ParseConfig(src, filePath, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return fmt.Errorf("HCL parse error in tfvars: %s", diags.Error())
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return fmt.Errorf("unexpected body type in %s", filePath)
	}

	relPath, _ := filepath.Rel(rootDir, filePath)

	// Derive environment name from filename or parent directory
	envName := deriveEnvName(relPath)

	values := make(map[string]string)
	for name, attr := range body.Attributes {
		val, diags := attr.Expr.Value(nil)
		if !diags.HasErrors() {
			values[name] = ctyToString(val)
		}
	}

	if len(values) > 0 {
		graph.Environments = append(graph.Environments, models.Environment{
			Name:     envName,
			FilePath: relPath,
			Values:   values,
		})
	}

	return nil
}

func deriveEnvName(relPath string) string {
	// Try parent directory name first (e.g. envs/prod/terraform.tfvars → "prod")
	dir := filepath.Dir(relPath)
	base := filepath.Base(dir)
	if base != "." && base != "/" {
		return base
	}
	// Fall back to filename without extension (e.g. prod.tfvars → "prod")
	name := filepath.Base(relPath)
	name = strings.TrimSuffix(name, ".tfvars")
	name = strings.TrimSuffix(name, ".auto")
	return name
}
