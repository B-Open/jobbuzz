// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphmodel

import (
	"time"
)

type GenericError interface {
	IsGenericError()
}

type NewUser interface {
	IsNewUser()
}

type PaginatedOutput interface {
	IsPaginatedOutput()
}

// A single company item.
type Company struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Name      string `json:"name"`
	Provider  int    `json:"provider"`
}

type CompanyOutput struct {
	// Company output in Paginated Format
	To          int            `json:"to"`
	From        int            `json:"from"`
	PerPage     int            `json:"per_page"`
	CurrentPage int            `json:"current_page"`
	TotalPage   int            `json:"total_page"`
	Total       int            `json:"total"`
	Data        []*Company     `json:"data"`
	Error       *StandardError `json:"error"`
}

func (CompanyOutput) IsPaginatedOutput() {}

// A single job listing item.
type Job struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Title       string    `json:"title"`
	Provider    int       `json:"provider"`
	Company     *Company  `json:"company"`
	Salary      string    `json:"salary"`
	Location    string    `json:"location"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
}

type JobOutput struct {
	// Job output in Paginated Format
	To          int            `json:"to"`
	From        int            `json:"from"`
	PerPage     int            `json:"per_page"`
	CurrentPage int            `json:"current_page"`
	TotalPage   int            `json:"total_page"`
	Total       int            `json:"total"`
	Data        []*Job         `json:"data"`
	Error       *StandardError `json:"error"`
}

func (JobOutput) IsPaginatedOutput() {}

// The login result
type LoginResult struct {
	AccessToken string `json:"accessToken"`
}

func (LoginResult) IsNewUser() {}

// The input for registering a new user.
type NewUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PaginationInput struct {
	// Standard Pagination Inputs
	Size int `json:"size"`
	Page int `json:"page"`
}

type StandardError struct {
	// A standard error with just a simple message
	Message string `json:"message"`
}

func (StandardError) IsGenericError() {}
func (StandardError) IsNewUser()      {}

type StringFilterInput struct {
	// Standard String Type Filters
	Contains    *string `json:"contains"`
	NotContains *string `json:"notContains"`
	BeginsWith  *string `json:"beginsWith"`
}

// The registered User.
type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
