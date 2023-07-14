// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import "time"

// NatsItem -.
type NatsItem struct {
	Id          int       `json:"id"       example:2674`
	CampaignId  int       `json:"campaign_id"       example:2674`
	Name        string    `json:"name"       example:"Abu"`
	Description string    `json:"description"       example:"some text"`
	Priority    int       `json:"priority"       example:2674`
	Removed     bool      `json:"removed"       example:false`
	EventTime   time.Time `json:"EventTime "`
}
