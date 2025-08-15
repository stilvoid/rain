package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws-cloudformation/rain/cft"
	"github.com/aws-cloudformation/rain/cft/diff"
	"github.com/aws-cloudformation/rain/cft/format"
	"github.com/aws-cloudformation/rain/cft/parse"
	"github.com/aws-cloudformation/rain/cft/pkg"
	"github.com/aws-cloudformation/rain/internal/cmd/merge"
	"github.com/mark3labs/mcp-go/mcp"
)

// buildTemplateHandler handles the build_template MCP tool
func buildTemplateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := request.GetArguments()

	resourcesStr, ok := params["resources"].(string)
	if !ok || strings.TrimSpace(resourcesStr) == "" {
		return mcp.NewToolResultError("missing or empty resources parameter"), nil
	}

	resourceTypes := strings.Split(resourcesStr, ",")
	for i, resourceType := range resourceTypes {
		resourceTypes[i] = strings.TrimSpace(resourceType)
	}

	outputFormat := getString(params, "format", "yaml")
	bare := getBool(params, "bare", false)

	template, err := buildTemplate(resourceTypes, outputFormat, bare)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(template), nil
}

// formatTemplateHandler handles the format_template MCP tool
func formatTemplateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := request.GetArguments()

	templateStr, ok := params["template"].(string)
	if !ok || strings.TrimSpace(templateStr) == "" {
		return mcp.NewToolResultError("missing or empty template parameter"), nil
	}

	outputFormat := getString(params, "output_format", "yaml")
	unsorted := getBool(params, "unsorted", false)

	template, err := parse.String(templateStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse template: %v", err)), nil
	}

	formatOptions := format.Options{
		JSON:     outputFormat == "json",
		Unsorted: unsorted,
	}

	formattedTemplate := format.String(template, formatOptions)
	return mcp.NewToolResultText(formattedTemplate), nil
}

// validateTemplateHandler handles the validate_template MCP tool
func validateTemplateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := request.GetArguments()

	templateStr, ok := params["template"].(string)
	if !ok || strings.TrimSpace(templateStr) == "" {
		return mcp.NewToolResultError("missing or empty template parameter"), nil
	}

	_, err := parse.String(templateStr)

	result := map[string]interface{}{
		"valid":    err == nil,
		"errors":   []string{},
		"warnings": []string{},
	}

	if err != nil {
		result["errors"] = []string{err.Error()}
	}

	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(jsonResult)), nil
}

// diffTemplatesHandler handles the diff_templates MCP tool
func diffTemplatesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := request.GetArguments()

	template1Str, ok1 := params["template1"].(string)
	template2Str, ok2 := params["template2"].(string)

	if !ok1 || !ok2 || strings.TrimSpace(template1Str) == "" || strings.TrimSpace(template2Str) == "" {
		return mcp.NewToolResultError("missing or empty template1 or template2 parameter"), nil
	}

	longFormat := getBool(params, "long_format", false)

	template1, err := parse.String(template1Str)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse template1: %v", err)), nil
	}

	template2, err := parse.String(template2Str)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse template2: %v", err)), nil
	}

	d := diff.New(template1, template2)

	if d.Mode() == diff.Unchanged {
		return mcp.NewToolResultText("No differences found"), nil
	}

	diffOutput := d.Format(longFormat)
	return mcp.NewToolResultText(diffOutput), nil
}

// forecastTemplateHandler handles the forecast_template MCP tool
func forecastTemplateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := request.GetArguments()

	templateStr, ok := params["template"].(string)
	if !ok || strings.TrimSpace(templateStr) == "" {
		return mcp.NewToolResultError("missing or empty template parameter"), nil
	}

	template, err := parse.String(templateStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse template: %v", err)), nil
	}

	// Use Rain's forecast functionality for time estimation
	stackExists := false
	if stackName := getString(params, "stack_name", ""); stackName != "" {
		stackExists = true // Assume it's an update if stack name is provided
	}

	estimatedSeconds := forecastTemplate(template, stackExists)

	// Determine deployment action for messaging
	action := "deployment"
	if stackExists {
		action = "update"
	}

	result := map[string]interface{}{
		"status":         "success",
		"message":        fmt.Sprintf("Template is valid for %s. Estimated deployment time: %d seconds", action, estimatedSeconds),
		"estimated_time": fmt.Sprintf("%d seconds", estimatedSeconds),
		"stack_action":   action,
		"issues":         []interface{}{},
	}

	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(jsonResult)), nil
}

// packageTemplateHandler handles the package_template MCP tool
func packageTemplateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := request.GetArguments()

	templatePathStr, ok := params["template_path"].(string)
	if !ok || strings.TrimSpace(templatePathStr) == "" {
		return mcp.NewToolResultError("missing or empty template_path parameter"), nil
	}

	outputFormat := getString(params, "output_format", "yaml")

	template, err := pkg.File(templatePathStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to package template: %v", err)), nil
	}

	formatOptions := format.Options{
		JSON:     outputFormat == "json",
		Unsorted: false,
	}

	packagedTemplate := format.String(template, formatOptions)
	return mcp.NewToolResultText(packagedTemplate), nil
}

// mergeTemplatesHandler handles the merge_templates MCP tool
func mergeTemplatesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := request.GetArguments()

	templatesArray, ok := params["templates"].([]interface{})
	if !ok || len(templatesArray) < 2 {
		return mcp.NewToolResultError("templates parameter must be an array with at least 2 templates"), nil
	}

	templateStrings := make([]string, len(templatesArray))
	for i, template := range templatesArray {
		templateStr, ok := template.(string)
		if !ok || strings.TrimSpace(templateStr) == "" {
			return mcp.NewToolResultError(fmt.Sprintf("template at index %d is invalid", i)), nil
		}
		templateStrings[i] = templateStr
	}

	outputFormat := getString(params, "output_format", "yaml")

	// Parse all templates
	templates := make([]*cft.Template, len(templateStrings))
	for i, content := range templateStrings {
		template, err := parse.String(content)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to parse template %d: %v", i+1, err)), nil
		}
		templates[i] = template
	}

	// Merge templates
	merged := templates[0]
	for i := 1; i < len(templates); i++ {
		mergedTemplate, err := merge.MergeTemplates(merged, templates[i])
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to merge template %d: %v", i+1, err)), nil
		}
		merged = mergedTemplate
	}

	formatOptions := format.Options{
		JSON:     outputFormat == "json",
		Unsorted: false,
	}

	mergedTemplateStr := format.String(merged, formatOptions)
	return mcp.NewToolResultText(mergedTemplateStr), nil
}

// getTemplateResourcesHandler handles the get_template_resources MCP tool
func getTemplateResourcesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := request.GetArguments()

	templateStr, ok := params["template"].(string)
	if !ok || strings.TrimSpace(templateStr) == "" {
		return mcp.NewToolResultError("missing or empty template parameter"), nil
	}

	outputFormat := getString(params, "output_format", "yaml")
	includeProperties := getBool(params, "include_properties", false)

	template, err := parse.String(templateStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse template: %v", err)), nil
	}

	result, err := extractTemplateResources(template, includeProperties, outputFormat)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(result), nil
}

// getTemplateOutputsHandler handles the get_template_outputs MCP tool
func getTemplateOutputsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := request.GetArguments()

	templateStr, ok := params["template"].(string)
	if !ok || strings.TrimSpace(templateStr) == "" {
		return mcp.NewToolResultError("missing or empty template parameter"), nil
	}

	outputFormat := getString(params, "output_format", "yaml")
	includeDescriptions := getBool(params, "include_descriptions", false)

	template, err := parse.String(templateStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse template: %v", err)), nil
	}

	result, err := extractTemplateOutputs(template, includeDescriptions, outputFormat)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(result), nil
}
