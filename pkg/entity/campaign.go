// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

// Campaign -.
type Campaign struct {
	Id   int    `json:"id"       example:2674`
	Name string `json:"name"       example:"Abu"`
}
