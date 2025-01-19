// internal/dto/policy.dto.go
package dto

import (
    "time"
    "errors"
)

type PolicyDTO struct {
    PolicyType        string    `json:"policy_type"`
    PolicyNumber      string    `json:"policy_number"`
    InsuranceProvider string    `json:"insurance_provider"`
    PolicyComment     string    `json:"policy_comment"`
    StartDate         time.Time `json:"start_date"`
    EndDate           time.Time `json:"end_date"`
    AutomaticRenewal  bool      `json:"automatic_renewal"`
    CreatedBy         string    `json:"created_by"`
    Premium           float64   `json:"premium"`
    PaymentFrequency  int       `json:"payment_frequency"`
}

func (p *PolicyDTO) Validate() error {
    if p.PolicyType == "" {
        return errors.New("policy type is required")
    }
    if p.PolicyNumber == "" {
        return errors.New("policy number is required")
    }
    if p.InsuranceProvider == "" {
        return errors.New("insurance provider is required")
    }
    if p.StartDate.IsZero() {
        return errors.New("start date is required")
    }
    if p.EndDate.IsZero() {
        return errors.New("end date is required")
    }
    if p.EndDate.Before(p.StartDate) {
        return errors.New("end date cannot be before start date")
    }
    if p.Premium <= 0 {
        return errors.New("premium must be greater than 0")
    }
    if p.PaymentFrequency <= 0 {
        return errors.New("payment frequency must be greater than 0")
    }
    return nil
}