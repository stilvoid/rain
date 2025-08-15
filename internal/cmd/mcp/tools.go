package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerTools registers all Rain MCP tools
func registerTools(s *server.MCPServer) {
	s.AddTool(mcp.NewTool("build_template",
		mcp.WithDescription("Generate CloudFormation template from resource types using Rain's build functionality"),
		mcp.WithString("resources", mcp.Required(), mcp.Description("Comma-separated AWS resource types")),
		mcp.WithString("format", mcp.Description("Output format: yaml or json (default: yaml)")),
		mcp.WithBoolean("bare", mcp.Description("Generate minimal template without metadata (default: false)")),
	), buildTemplateHandler)

	s.AddTool(mcp.NewTool("format_template",
		mcp.WithDescription("Format and standardize CloudFormation templates"),
		mcp.WithString("template", mcp.Required(), mcp.Description("CloudFormation template content")),
		mcp.WithString("output_format", mcp.Description("Output format: yaml or json (default: yaml)")),
		mcp.WithBoolean("unsorted", mcp.Description("Skip sorting of template elements (default: false)")),
	), formatTemplateHandler)

	s.AddTool(mcp.NewTool("validate_template",
		mcp.WithDescription("Validate CloudFormation template syntax and structure"),
		mcp.WithString("template", mcp.Required(), mcp.Description("CloudFormation template content")),
	), validateTemplateHandler)

	s.AddTool(mcp.NewTool("diff_templates",
		mcp.WithDescription("Compare two CloudFormation templates and show differences"),
		mcp.WithString("template1", mcp.Required(), mcp.Description("First CloudFormation template content")),
		mcp.WithString("template2", mcp.Required(), mcp.Description("Second CloudFormation template content")),
		mcp.WithBoolean("long_format", mcp.Description("Show detailed diff output (default: false)")),
	), diffTemplatesHandler)

	s.AddTool(mcp.NewTool("forecast_template",
		mcp.WithDescription("Estimate CloudFormation template deployment time using Rain's forecast functionality"),
		mcp.WithString("template", mcp.Required(), mcp.Description("CloudFormation template content")),
		mcp.WithString("stack_name", mcp.Description("Stack name (if provided, assumes update operation)")),
	), forecastTemplateHandler)

	s.AddTool(mcp.NewTool("package_template",
		mcp.WithDescription("Package CloudFormation templates with Rain directives and S3 uploads"),
		mcp.WithString("template_path", mcp.Required(), mcp.Description("Path to the CloudFormation template file")),
		mcp.WithString("s3_bucket", mcp.Description("S3 bucket for uploading assets")),
		mcp.WithString("s3_prefix", mcp.Description("S3 prefix for uploaded assets")),
		mcp.WithString("output_format", mcp.Description("Output format: yaml or json (default: yaml)")),
	), packageTemplateHandler)

	s.AddTool(mcp.NewTool("merge_templates",
		mcp.WithDescription("Combine multiple CloudFormation templates into a single template"),
		mcp.WithArray("templates", mcp.Required(), mcp.Description("Array of CloudFormation template contents to merge")),
		mcp.WithString("output_format", mcp.Description("Output format: yaml or json (default: yaml)")),
	), mergeTemplatesHandler)

	s.AddTool(mcp.NewTool("get_template_resources",
		mcp.WithDescription("Extract and list all resources from a CloudFormation template"),
		mcp.WithString("template", mcp.Required(), mcp.Description("CloudFormation template content")),
		mcp.WithString("output_format", mcp.Description("Output format: yaml or json (default: yaml)")),
		mcp.WithBoolean("include_properties", mcp.Description("Include resource properties in output (default: false)")),
	), getTemplateResourcesHandler)

	s.AddTool(mcp.NewTool("get_template_outputs",
		mcp.WithDescription("Extract and list all outputs from a CloudFormation template"),
		mcp.WithString("template", mcp.Required(), mcp.Description("CloudFormation template content")),
		mcp.WithString("output_format", mcp.Description("Output format: yaml or json (default: yaml)")),
		mcp.WithBoolean("include_descriptions", mcp.Description("Include output descriptions in output (default: false)")),
	), getTemplateOutputsHandler)
}

// Helper functions for parameter extraction
func getString(params map[string]interface{}, key string, defaultValue string) string {
	if val, ok := params[key].(string); ok {
		return val
	}
	return defaultValue
}

func getBool(params map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := params[key].(bool); ok {
		return val
	}
	return defaultValue
}
