package model

const (
	CDNStatusActive   = "active"
	CDNStatusInactive = "inactive"
)

type CDNInstance struct {
	Model
	Host            string        `gorm:"column:host;size:256;uniqueIndex;not null" json:"host"`
	IDC             string        `gorm:"column:idc;size:1024;not null" json:"idc"`
	Location        string        `gorm:"column:location;size:1024" json:"location"`
	IP              string        `gorm:"column:ip;size:256;not null" json:"ip"`
	Port            int32         `gorm:"column:port;not null" json:"port"`
	DownloadPort    int32         `gorm:"column:download_port;not null" json:"download_port"`
	Status          string        `gorm:"type:enum('active', 'inactive');default:'inactive'" json:"status"`
	CDNID           *uint         `json:"cdn_id"`
	CDN             CDN           `json:"cdn"`
	SecurityGroupID *uint         `json:"security_group_id"`
	SecurityGroup   SecurityGroup `json:"security_group"`
}
