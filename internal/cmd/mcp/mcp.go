package mcp

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var experimental bool

// Cmd is the mcp command's entrypoint
var Cmd = &cobra.Command{
	Use:                   "mcp",
	Short:                 "Run Rain as an MCP (Model Context Protocol) server (Experimental!)",
	Long:                  "Run Rain as an MCP server that exposes CloudFormation template tools for AI assistants. This is an experimental feature that requires the -x flag to run.",
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !experimental {
			panic("Please add the --experimental arg to use this feature")
		}
		
		s := server.NewMCPServer("Rain CloudFormation MCP Server", "1.0.0")
		registerTools(s)
		return server.ServeStdio(s)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&experimental, "experimental", "x", false, "Acknowledge that this is an experimental feature")
}
