package types

type SecurityGroupParams struct {
	ID uint `uri:"id" binding:"required"`
}

type AddSchedulerInstanceToSecurityGroupParams struct {
	ID         uint `uri:"id" binding:"required"`
	InstanceID uint `uri:"instance_id" binding:"required"`
}

type AddCDNInstanceToSecurityGroupParams struct {
	ID         uint `uri:"id" binding:"required"`
	InstanceID uint `uri:"instance_id" binding:"required"`
}

type CreateSecurityGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	BIO         string `json:"bio" binding:"omitempty"`
	Domain      string `json:"domain" binding:"required"`
	ProxyDomain string `json:"proxy_domain" binding:"omitempty"`
}

type UpdateSecurityGroupRequest struct {
	Name        string `json:"name" binding:"omitempty"`
	BIO         string `json:"bio" binding:"omitempty"`
	Domain      string `json:"domain" binding:"omitempty"`
	ProxyDomain string `json:"proxy_domain" binding:"omitempty"`
}

type GetSecurityGroupsQuery struct {
	Page    int    `form:"page" binding:"omitempty,gte=1"`
	PerPage int    `form:"per_page" binding:"omitempty,gte=1,lte=50"`
	Name    string `form:"name" binding:"omitempty"`
	Domain  string `form:"domain" binding:"omitempty"`
}
