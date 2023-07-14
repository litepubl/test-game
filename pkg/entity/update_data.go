package entity

// UpdateData  промежуточный объект для передачи данных -.
type UpdateData struct {
	Id          int     `json:"id"       example:2674`
	CampaignId  int     `json:"campaign_id"       example:2674`
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Priority    int
}
