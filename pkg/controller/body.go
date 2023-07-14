package controller

type (
	createBody struct {
		Name string `json:"name" binding:"required"`
	}

	removedResponse struct {
		Id         int  `json:"id"`
		CampaignId int  `json:"campaignId"`
		Removed    bool `json:"removed"`
	}
)
