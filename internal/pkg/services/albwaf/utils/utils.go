package utils

import (
	albwaf "github.com/stackitcloud/stackit-sdk-go/services/albwaf/v1api"
)

// ToCreateCustomRules converts read-model rules (from a GetCustomRuleGroup response) into the
// write-model rules used by create/update payloads.
// Server-assigned rule IDs and the behavior severity are dropped.
func ToCreateCustomRules(rules []albwaf.GetCustomRule) []albwaf.CreateCustomRule {
	if rules == nil {
		return nil
	}
	createRules := make([]albwaf.CreateCustomRule, len(rules))
	for i := range rules {
		createRules[i] = *toCreateCustomRule(&rules[i])
	}
	return createRules
}

func toCreateCustomRule(rule *albwaf.GetCustomRule) *albwaf.CreateCustomRule {
	return &albwaf.CreateCustomRule{
		Behavior:    toBehavior(rule.Behavior),
		Conditions:  rule.Conditions,
		Description: rule.Description,
	}
}

func toBehavior(behavior albwaf.GetBehavior) albwaf.Behavior {
	return albwaf.Behavior{
		Action: behavior.Action,
		Log:    new(behavior.Log),
		LogMsg: behavior.LogMsg,
	}
}
