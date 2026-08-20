package utils

import (
	"testing"

	albwaf "github.com/stackitcloud/stackit-sdk-go/services/albwaf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
)

func TestToCreateCustomRules(t *testing.T) {
	tests := []struct {
		description string
		rules       []albwaf.GetCustomRule
		expected    []albwaf.CreateCustomRule
	}{
		{
			description: "nil rules",
			rules:       nil,
			expected:    nil,
		},
		{
			description: "empty rules",
			rules:       []albwaf.GetCustomRule{},
			expected:    []albwaf.CreateCustomRule{},
		},
		{
			description: "drops id and severity, wraps log, passes conditions through",
			rules: []albwaf.GetCustomRule{
				{
					Id:          42,
					Description: utils.Ptr("block /admin"),
					Behavior: albwaf.GetBehavior{
						Action:   albwaf.ACTION_ACTION_DENY,
						Log:      true,
						LogMsg:   utils.Ptr("blocked"),
						Severity: albwaf.SEVERITY_SEVERITY_WARNING,
					},
					Conditions: []albwaf.Condition{
						{
							Operator: albwaf.ConditionOperator{
								Type:  albwaf.OPERATOR_OPERATOR_BEGINS_WITH,
								Value: utils.Ptr("/admin"),
							},
							Variable: albwaf.ConditionVariable{
								Type: albwaf.VARIABLE_VARIABLE_REQUEST_URI_RAW,
							},
							Transformations: []albwaf.Transformation{
								albwaf.TRANSFORMATION_TRANSFORMATION_LOWERCASE,
							},
						},
					},
				},
			},
			expected: []albwaf.CreateCustomRule{
				{
					Description: utils.Ptr("block /admin"),
					Behavior: albwaf.Behavior{
						Action: albwaf.ACTION_ACTION_DENY,
						Log:    utils.Ptr(true),
						LogMsg: utils.Ptr("blocked"),
					},
					Conditions: []albwaf.Condition{
						{
							Operator: albwaf.ConditionOperator{
								Type:  albwaf.OPERATOR_OPERATOR_BEGINS_WITH,
								Value: utils.Ptr("/admin"),
							},
							Variable: albwaf.ConditionVariable{
								Type: albwaf.VARIABLE_VARIABLE_REQUEST_URI_RAW,
							},
							Transformations: []albwaf.Transformation{
								albwaf.TRANSFORMATION_TRANSFORMATION_LOWERCASE,
							},
						},
					},
				},
			},
		},
		{
			description: "log false is wrapped, not omitted",
			rules: []albwaf.GetCustomRule{
				{
					Id: 1,
					Behavior: albwaf.GetBehavior{
						Action:   albwaf.ACTION_ACTION_ALLOW,
						Log:      false,
						Severity: albwaf.SEVERITY_SEVERITY_INFO,
					},
					Conditions: []albwaf.Condition{},
				},
			},
			expected: []albwaf.CreateCustomRule{
				{
					Behavior: albwaf.Behavior{
						Action: albwaf.ACTION_ACTION_ALLOW,
						Log:    utils.Ptr(false),
					},
					Conditions: []albwaf.Condition{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output := ToCreateCustomRules(tt.rules)

			diff := cmp.Diff(output, tt.expected)
			if diff != "" {
				t.Fatalf("data does not match: %s", diff)
			}
		})
	}
}
