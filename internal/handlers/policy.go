// internal/handlers/policy.go
package handlers

import (
    "goapi/internal/dto"
    "goapi/internal/model"
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "database/sql"
    "time"
)

type PolicyHandler struct {
    db *sql.DB
}

func NewPolicyHandler(db *sql.DB) *PolicyHandler {
    return &PolicyHandler{db: db}
}

func (h *PolicyHandler) CreatePolicy(c *fiber.Ctx) error {
    var policyDTO dto.PolicyDTO
    if err := c.BodyParser(&policyDTO); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request payload",
        })
    }

    if err := policyDTO.Validate(); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    policy := model.Policy{
        ID:                uuid.New().String(),
        PolicyType:        policyDTO.PolicyType,
        PolicyNumber:      policyDTO.PolicyNumber,
        InsuranceProvider: policyDTO.InsuranceProvider,
        PolicyComment:     policyDTO.PolicyComment,
        StartDate:         policyDTO.StartDate,
        EndDate:           policyDTO.EndDate,
        AutomaticRenewal:  policyDTO.AutomaticRenewal,
        CreatedBy:         policyDTO.CreatedBy,
        Premium:           policyDTO.Premium,
        PaymentFrequency:  policyDTO.PaymentFrequency,
        Created:           time.Now(),
        Updated:           time.Now(),
    }

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
            "error": "Failed to create policy",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(policy)
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

func (h *PolicyHandler) UpdatePolicy(c *fiber.Ctx) error {
    id := c.Params("id")
    var policyDTO dto.PolicyDTO

    if err := c.BodyParser(&policyDTO); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request payload",
        })
    }

    if err := policyDTO.Validate(); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    query := `UPDATE policies SET
        policy_type = ?, policy_number = ?, insurance_provider = ?, 
        policy_comment = ?, start_date = ?, end_date = ?, automatic_renewal = ?,
        created_by = ?, premium = ?, payment_frequency = ?, updated = ?
        WHERE id = ?`

    result, err := h.db.Exec(query,
        policyDTO.PolicyType, policyDTO.PolicyNumber, policyDTO.InsuranceProvider,
        policyDTO.PolicyComment, policyDTO.StartDate,
        policyDTO.EndDate, policyDTO.AutomaticRenewal, policyDTO.CreatedBy,
        policyDTO.Premium, policyDTO.PaymentFrequency, time.Now(), id,
    )

    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to update policy",
        })
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil || rowsAffected == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Policy not found",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Policy updated successfully",
    })
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