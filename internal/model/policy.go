// internal/model/policy.go
package model

import "time"

type Policy struct {
    ID                string    `json:"id" db:"id"`
    PolicyType        string    `json:"policy_type" db:"policy_type"`
    PolicyNumber      string    `json:"policy_number" db:"policy_number"`
    InsuranceProvider string    `json:"insurance_provider" db:"insurance_provider"`
    PolicyComment     string    `json:"policy_comment" db:"policy_comment"`
    StartDate         time.Time `json:"start_date" db:"start_date"`
    EndDate           time.Time `json:"end_date" db:"end_date"`
    AutomaticRenewal  bool      `json:"automatic_renewal" db:"automatic_renewal"`
    CreatedBy         string    `json:"created_by" db:"created_by"`
    Premium           float64   `json:"premium" db:"premium"`
    PaymentFrequency  int       `json:"payment_frequency" db:"payment_frequency"`
    Created           time.Time `json:"created" db:"created"`
    Updated           time.Time `json:"updated" db:"updated"`
}