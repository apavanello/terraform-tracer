package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/apavanello/terraform-tracer/internal/api"
	"github.com/apavanello/terraform-tracer/internal/parser"
	"github.com/spf13/cobra"
)

var port string

var rootCmd = &cobra.Command{
	Use:   "tracer",
	Short: "Terraform Tracer – Visual dependency tracer for Terraform resources",
	Long:  "A CLI tool that parses Terraform files and serves an interactive dependency graph in your browser.",
}

var startCmd = &cobra.Command{
	Use:   "start [directory]",
	Short: "Parse a Terraform directory and start the web UI",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}

		// Validate directory
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			return fmt.Errorf("'%s' is not a valid directory", dir)
		}

		fmt.Printf("🔍 Parsing Terraform files in: %s\n", dir)
		start := time.Now()

		graph, err := parser.Parse(dir)
		if err != nil {
			return fmt.Errorf("parse failed: %w", err)
		}

		elapsed := time.Since(start)
		fmt.Printf("✅ Found %d nodes, %d edges, %d environments in %v\n",
			len(graph.Nodes), len(graph.Edges), len(graph.Environments), elapsed.Round(time.Millisecond))

		addr := ":" + port
		fmt.Printf("🚀 Starting Terraform Tracer on http://localhost%s\n", addr)

		// Try to open browser
		go func() {
			time.Sleep(500 * time.Millisecond)
			openBrowser("http://localhost" + addr)
		}()

		server := api.NewServer(graph)
		if err := server.Listen(addr); err != nil {
			log.Fatalf("server error: %v", err)
		}
		return nil
	},
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	}
	if cmd != nil {
		_ = cmd.Start()
	}
}

func init() {
	startCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to serve the web UI on")
	rootCmd.AddCommand(startCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
