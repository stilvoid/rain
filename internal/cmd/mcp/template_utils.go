package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/aws-cloudformation/rain/cft"
)

// extractTemplateResources extracts resources from a template with optional properties
func extractTemplateResources(template *cft.Template, includeProperties bool, outputFormat string) (string, error) {
	templateMap := template.Map()
	resourcesSection, ok := templateMap["Resources"]
	if !ok {
		return formatEmptySection("Resources", outputFormat), nil
	}

	resourcesMap, ok := resourcesSection.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid Resources section format")
	}

	result := make(map[string]interface{})
	for logicalId, resource := range resourcesMap {
		resourceMap, ok := resource.(map[string]interface{})
		if !ok {
			continue
		}

		resourceInfo := map[string]interface{}{
			"LogicalId": logicalId,
		}

		if resourceType, ok := resourceMap["Type"].(string); ok {
			resourceInfo["Type"] = resourceType
		}

		if includeProperties {
			if properties, ok := resourceMap["Properties"]; ok {
				resourceInfo["Properties"] = properties
			}
		}

		result[logicalId] = resourceInfo
	}

	return formatSection("Resources", result, outputFormat)
}

// extractTemplateOutputs extracts outputs from a template with optional descriptions
func extractTemplateOutputs(template *cft.Template, includeDescriptions bool, outputFormat string) (string, error) {
	templateMap := template.Map()
	outputsSection, ok := templateMap["Outputs"]
	if !ok {
		return formatEmptySection("Outputs", outputFormat), nil
	}

	outputsMap, ok := outputsSection.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid Outputs section format")
	}

	result := make(map[string]interface{})
	for logicalId, output := range outputsMap {
		outputMap, ok := output.(map[string]interface{})
		if !ok {
			continue
		}

		outputInfo := map[string]interface{}{
			"LogicalId": logicalId,
		}

		if value, ok := outputMap["Value"]; ok {
			outputInfo["Value"] = value
		}

		if includeDescriptions {
			if description, ok := outputMap["Description"]; ok {
				outputInfo["Description"] = description
			}
		}

		result[logicalId] = outputInfo
	}

	return formatSection("Outputs", result, outputFormat)
}

// formatSection formats a section as JSON or YAML
func formatSection(sectionName string, data map[string]interface{}, outputFormat string) (string, error) {
	if outputFormat == "json" {
		jsonResult, err := json.MarshalIndent(map[string]interface{}{sectionName: data}, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal %s: %v", sectionName, err)
		}
		return string(jsonResult), nil
	}

	// Simple YAML output
	yamlResult := fmt.Sprintf("%s:\n", sectionName)
	for logicalId, item := range data {
		yamlResult += fmt.Sprintf("  %s:\n", logicalId)
		if itemMap, ok := item.(map[string]interface{}); ok {
			for key, value := range itemMap {
				yamlResult += fmt.Sprintf("    %s: %v\n", key, value)
			}
		}
	}

	return yamlResult, nil
}

// formatEmptySection formats an empty section
func formatEmptySection(sectionName string, outputFormat string) string {
	if outputFormat == "json" {
		return fmt.Sprintf(`{"%s": {}}`, sectionName)
	}
	return fmt.Sprintf("%s: {}", sectionName)
}
