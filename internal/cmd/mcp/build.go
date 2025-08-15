package mcp

import (
	"github.com/aws-cloudformation/rain/cft"
	"github.com/aws-cloudformation/rain/cft/format"
	"github.com/aws-cloudformation/rain/internal/cmd/build"
	"github.com/aws-cloudformation/rain/internal/cmd/forecast"
)

// buildTemplate creates a CloudFormation template using Rain's build functionality
func buildTemplate(resourceTypes []string, outputFormat string, bare bool) (string, error) {
	// Use Rain's actual build functionality
	template, err := build.BuildTemplate(resourceTypes, bare)
	if err != nil {
		return "", err
	}

	// Format the template according to the requested output format
	formatOptions := format.Options{
		JSON: outputFormat == "json",
	}

	return format.String(template, formatOptions), nil
}

// forecastTemplate uses Rain's forecast functionality to estimate deployment time
func forecastTemplate(template *cft.Template, stackExists bool) int {
	return forecast.PredictTotalEstimate(template, stackExists)
}
