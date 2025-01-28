package dto

import (
	"fmt"
	"time"
)

// ValidationRule represents a custom validation rule
type ValidationRule struct {
	Field   string
	Rule    func(interface{}) bool
	Message string
}

type PolicyDTO struct {
	PolicyType        string     `json:"policy_type"`
	PolicyNumber      string     `json:"policy_number"`
	InsuranceProvider string     `json:"insurance_provider"`
	PolicyComment     string     `json:"policy_comment"`
	StartDate         *time.Time `json:"start_date"`
	EndDate           time.Time  `json:"end_date"`
	AutomaticRenewal  bool       `json:"automatic_renewal"`
	CreatedBy         string     `json:"created_by"`
	Premium           float64    `json:"premium"`
	PaymentFrequency  int        `json:"payment_frequency"`
}

// ValidateTypes performs type validation only
func (p *PolicyDTO) ValidateTypes(rawData map[string]interface{}) map[string]string {
	errors := make(map[string]string)

	typeChecks := map[string]struct {
		required     bool
		expectedType string
		message      string
	}{
		"policy_type": {
			required:     true,
			expectedType: "string",
			message:      "Policy type must be a text value",
		},
		"policy_number": {
			required:     true,
			expectedType: "string",
			message:      "Policy number must be a text value",
		},
		"insurance_provider": {
			required:     true,
			expectedType: "string",
			message:      "Insurance provider must be a text value",
		},
		"policy_comment": {
			required:     true,
			expectedType: "string",
			message:      "Policy comment must be a text value",
		},
		"start_date": {
			required:     false,
			expectedType: "string",
			message:      "Start date must be in ISO format (e.g., 2024-03-01T00:00:00Z)",
		},
		"end_date": {
			required:     true,
			expectedType: "string",
			message:      "End date must be in ISO format (e.g., 2024-03-01T00:00:00Z)",
		},
		"automatic_renewal": {
			required:     true,
			expectedType: "bool",
			message:      "Automatic renewal must be true or false",
		},
		"created_by": {
			required:     true,
			expectedType: "string",
			message:      "Created by must be a text value",
		},
		"premium": {
			required:     true,
			expectedType: "float64",
			message:      "Premium must be a number (e.g., 199.99)",
		},
		"payment_frequency": {
			required:     true,
			expectedType: "float64",
			message:      "Payment frequency must be a whole number (e.g., 12)",
		},
	}

	// Check required fields and types
	for field, check := range typeChecks {
		value, exists := rawData[field]
		if !exists {
			if check.required {
				errors[field] = fmt.Sprintf("%s is required", field)
			}
			continue
		}

		if value == nil {
			if check.required {
				errors[field] = fmt.Sprintf("%s cannot be null", field)
			}
			continue
		}

		if !isType(value, check.expectedType) {
			errors[field] = check.message
		}
	}

	return errors
}

// ValidateBusinessRules performs business rule validation
func (p *PolicyDTO) ValidateBusinessRules(rules []ValidationRule) map[string]string {
	errors := make(map[string]string)

	// Default business rules
	defaultRules := []ValidationRule{
		{
			Field: "end_date",
			Rule: func(i interface{}) bool {
				return !p.EndDate.IsZero() && (p.StartDate == nil || !p.EndDate.Before(*p.StartDate))
			},
			Message: "End date cannot be before start date",
		},
		{
			Field: "premium",
			Rule: func(i interface{}) bool {
				return p.Premium > 0
			},
			Message: "Premium must be greater than 0",
		},
		{
			Field: "payment_frequency",
			Rule: func(i interface{}) bool {
				return p.PaymentFrequency > 0
			},
			Message: "Payment frequency must be greater than 0",
		},
	}

	// Combine default rules with custom rules
	allRules := append(defaultRules, rules...)

	// Apply all rules
	for _, rule := range allRules {
		var value interface{}
		switch rule.Field {
		case "premium":
			value = p.Premium
		case "payment_frequency":
			value = p.PaymentFrequency
		case "end_date":
			value = p.EndDate
		// Add other fields as needed
		}

		if !rule.Rule(value) {
			errors[rule.Field] = rule.Message
		}
	}

	return errors
}

// Validate performs both type and business rule validation
func (p *PolicyDTO) Validate(rawData map[string]interface{}, customRules []ValidationRule) map[string]string {
	// First perform type validation
	errors := p.ValidateTypes(rawData)
	if len(errors) > 0 {
		return errors
	}

	// If types are valid, perform business rule validation
	return p.ValidateBusinessRules(customRules)
}

func isType(value interface{}, expectedType string) bool {
	switch expectedType {
	case "string":
		_, ok := value.(string)
		return ok
	case "float64":
		switch value.(type) {
		case float64, int, int64:
			return true
		default:
			return false
		}
	case "bool":
		_, ok := value.(bool)
		return ok
	default:
		return false
	}
}