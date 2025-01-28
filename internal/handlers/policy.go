// internal/handlers/policy.go
package handlers

import (
	"database/sql"
	"goapi/internal/dto"
	"goapi/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PolicyHandler struct {
	db *sql.DB
}

func NewPolicyHandler(db *sql.DB) *PolicyHandler {
	return &PolicyHandler{db: db}
}

// getCustomRules returns the custom validation rules for policies
func getCustomRules() []dto.ValidationRule {
	return []dto.ValidationRule{
		{
			Field: "premium",
			Rule: func(i interface{}) bool {
				premium := i.(float64)
				return premium >= 1 && premium <= 100
			},
			Message: "Premium must be between 1 and 100",
		},
		{
			Field: "payment_frequency",
			Rule: func(i interface{}) bool {
				freq := i.(int)
				// Only allow yearly (1), semi-annual (2), quarterly (4), or monthly (12)
				validFrequencies := map[int]bool{1: true, 2: true, 4: true, 12: true}
				return validFrequencies[freq]
			},
			Message: "Payment frequency must be one of: 12 (monthly), 4 (quarterly), 2 (semi-annual), or 1 (yearly)",
		},
	}
}

func (h *PolicyHandler) CreatePolicy(c *fiber.Ctx) error {
	var rawData map[string]interface{}
	var policyDTO dto.PolicyDTO

	// Parse raw JSON payload into a map for type validation
	if err := c.BodyParser(&rawData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request format",
			"details": "Please check your input format and try again",
		})
	}

	// Parse raw JSON into the DTO
	if err := c.BodyParser(&policyDTO); err != nil {
		// Handle unmarshaling errors with clear messages
		switch err.Error() {
		case "json: cannot unmarshal number into Go struct field PolicyDTO.policy_number of type string":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid field format",
				"details": "Policy number must be a text value",
				"field":   "policy_number",
				"example": "\"12345\"",
			})
		case "json: cannot unmarshal string into Go struct field PolicyDTO.premium of type float64":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid field format",
				"details": "Premium must be a number",
				"field":   "premium",
				"example": 99.99,
			})
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid request format",
				"details": "Please check the format of all fields",
			})
		}
	}

	// First validate types
	if typeErrors := policyDTO.ValidateTypes(rawData); len(typeErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Invalid field types",
			"fields": typeErrors,
		})
	}

	// Then validate business rules
	if ruleErrors := policyDTO.ValidateBusinessRules(getCustomRules()); len(ruleErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Business rule violations",
			"fields": ruleErrors,
		})
	}

	// Prepare the policy object
	var startDate time.Time
	if policyDTO.StartDate != nil {
		startDate = *policyDTO.StartDate
	}

	policy := model.Policy{
		ID:                uuid.New().String(),
		PolicyType:        policyDTO.PolicyType,
		PolicyNumber:      policyDTO.PolicyNumber,
		InsuranceProvider: policyDTO.InsuranceProvider,
		PolicyComment:     policyDTO.PolicyComment,
		StartDate:         startDate,
		EndDate:           policyDTO.EndDate,
		AutomaticRenewal:  policyDTO.AutomaticRenewal,
		CreatedBy:         policyDTO.CreatedBy,
		Premium:           policyDTO.Premium,
		PaymentFrequency:  policyDTO.PaymentFrequency,
		Created:           time.Now(),
		Updated:           time.Now(),
	}

	// Save policy to the database
	query := `INSERT INTO policies (
		id, policy_type, policy_number, insurance_provider, 
		policy_comment, start_date, end_date, automatic_renewal, created_by,
		premium, payment_frequency, created, updated
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := h.db.Exec(query,
		policy.ID, policy.PolicyType, policy.PolicyNumber, policy.InsuranceProvider,
		policy.PolicyComment, policy.StartDate, policy.EndDate,
		policy.AutomaticRenewal, policy.CreatedBy, policy.Premium, policy.PaymentFrequency,
		policy.Created, policy.Updated,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create policy",
			"details": "An error occurred while saving the policy",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(policy)
}

func (h *PolicyHandler) UpdatePolicy(c *fiber.Ctx) error {
	id := c.Params("id")
	var rawData map[string]interface{}
	var policyDTO dto.PolicyDTO

	// Parse raw JSON payload into a map for type validation
	if err := c.BodyParser(&rawData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request format",
			"details": "Please check your input format and try again",
		})
	}

	// Parse raw JSON into the DTO
	if err := c.BodyParser(&policyDTO); err != nil {
		// Handle unmarshaling errors with clear messages
		switch err.Error() {
		case "json: cannot unmarshal number into Go struct field PolicyDTO.policy_number of type string":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid field format",
				"details": "Policy number must be a text value",
				"field":   "policy_number",
				"example": "\"12345\"",
			})
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid request format",
				"details": "Please check the format of all fields",
			})
		}
	}

	// Validate both types and business rules
	if typeErrors := policyDTO.ValidateTypes(rawData); len(typeErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Invalid field types",
			"fields": typeErrors,
		})
	}

	if ruleErrors := policyDTO.ValidateBusinessRules(getCustomRules()); len(ruleErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Business rule violations",
			"fields": ruleErrors,
		})
	}

	var startDate time.Time
	if policyDTO.StartDate != nil {
		startDate = *policyDTO.StartDate
	}

	query := `UPDATE policies SET
		policy_type = ?, policy_number = ?, insurance_provider = ?, 
		policy_comment = ?, start_date = ?, end_date = ?, automatic_renewal = ?,
		created_by = ?, premium = ?, payment_frequency = ?, updated = ?
		WHERE id = ?`

	result, err := h.db.Exec(query,
		policyDTO.PolicyType, policyDTO.PolicyNumber, policyDTO.InsuranceProvider,
		policyDTO.PolicyComment, startDate,
		policyDTO.EndDate, policyDTO.AutomaticRenewal, policyDTO.CreatedBy,
		policyDTO.Premium, policyDTO.PaymentFrequency, time.Now(), id,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update policy",
			"details": "An error occurred while updating the policy",
		})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Policy not found",
			"details": "The requested policy does not exist",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Policy updated successfully",
	})
}

func (h *PolicyHandler) GetPolicy(c *fiber.Ctx) error {
	id := c.Params("id")
	var policy model.Policy

	query := `SELECT * FROM policies WHERE id = ?`
	err := h.db.QueryRow(query, id).Scan(
		&policy.ID, &policy.PolicyType, &policy.PolicyNumber, &policy.InsuranceProvider,
		&policy.PolicyComment, &policy.StartDate, &policy.EndDate,
		&policy.AutomaticRenewal, &policy.CreatedBy, &policy.Premium, &policy.PaymentFrequency,
		&policy.Created, &policy.Updated,
	)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Policy not found",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch policy",
		})
	}

	return c.JSON(policy)
}

func (h *PolicyHandler) DeletePolicy(c *fiber.Ctx) error {
	id := c.Params("id")

	query := `DELETE FROM policies WHERE id = ?`
	result, err := h.db.Exec(query, id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete policy",
		})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Policy not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Policy deleted successfully",
	})
}
