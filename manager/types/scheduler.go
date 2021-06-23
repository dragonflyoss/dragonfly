package types

type SchedulerParams struct {
	ID uint `uri:"id" binding:"required"`
}

type CreateSchedulerRequest struct {
	Host                string                 `json:"host" binding:"required"`
	SecurityGroupDomain string                 `json:"security_group_domain" binding:"omitempty"`
	VIPs                string                 `json:"vips" binding:"omitempty"`
	IDC                 string                 `json:"idc" binding:"required"`
	Location            string                 `json:"location" binding:"omitempty"`
	NetConfig           map[string]interface{} `json:"net_config" binding:"omitempty"`
	IP                  string                 `json:"ip" binding:"required"`
	Port                int32                  `json:"port" binding:"required"`
}

type UpdateSchedulerRequest struct {
	Host                string                 `json:"host" binding:"omitempty"`
	SecurityGroupDomain string                 `json:"security_group_domain" binding:"omitempty"`
	VIPs                string                 `json:"vips" binding:"omitempty"`
	IDC                 string                 `json:"idc" binding:"omitempty"`
	Location            string                 `json:"location" binding:"omitempty"`
	NetConfig           map[string]interface{} `json:"net_config" binding:"omitempty"`
	IP                  string                 `json:"ip" binding:"omitempty"`
	Port                int32                  `json:"port" binding:"omitempty"`
	SchedulerID         *uint                  `json:"scheduler_id" binding:"omitempty"`
}

type GetSchedulersQuery struct {
	Page     int    `form:"page" binding:"omitempty,gte=1"`
	PerPage  int    `form:"per_page" binding:"omitempty,gte=1,lte=50"`
	Host     string `form:"host" binding:"omitempty"`
	IDC      string `form:"idc" binding:"omitempty"`
	Location string `form:"location" binding:"omitempty"`
	IP       string `form:"ip" binding:"omitempty"`
	Status   string `form:"status" binding:"omitempty,oneof=active inactive"`
}
