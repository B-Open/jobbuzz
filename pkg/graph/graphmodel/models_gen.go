// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphmodel

type Job struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Title     string `json:"title"`
	Company   string `json:"company"`
	Salary    string `json:"salary"`
	Location  string `json:"location"`
}

type NewUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}