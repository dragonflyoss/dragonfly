package service

import (
	"d7y.io/dragonfly/v2/manager/model"
	"d7y.io/dragonfly/v2/manager/types"
)

func (s *service) CreateSecurityGroup(json types.CreateSecurityGroupRequest) (*model.SecurityGroup, error) {
	securityGroup := model.SecurityGroup{
		Name:        json.Name,
		BIO:         json.BIO,
		Domain:      json.Domain,
		ProxyDomain: json.ProxyDomain,
	}

	if err := s.db.Create(&securityGroup).Error; err != nil {
		return nil, err
	}

	return &securityGroup, nil
}

func (s *service) DestroySecurityGroup(id string) error {
	if err := s.db.Unscoped().Delete(&model.SecurityGroup{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateSecurityGroup(id string, json types.UpdateSecurityGroupRequest) (*model.SecurityGroup, error) {
	securityGroup := model.SecurityGroup{}
	if err := s.db.First(&securityGroup, id).Updates(model.SecurityGroup{
		Name:        json.Name,
		BIO:         json.BIO,
		Domain:      json.Domain,
		ProxyDomain: json.ProxyDomain,
	}).Error; err != nil {
		return nil, err
	}

	return &securityGroup, nil
}

func (s *service) GetSecurityGroup(id string) (*model.SecurityGroup, error) {
	securityGroup := model.SecurityGroup{}
	if err := s.db.First(&securityGroup, id).Error; err != nil {
		return nil, err
	}

	return &securityGroup, nil
}

func (s *service) GetSecurityGroups(q types.GetSecurityGroupsQuery) (*[]model.SecurityGroup, error) {
	securityGroups := []model.SecurityGroup{}
	if err := s.db.Scopes(model.Paginate(q.Page, q.PerPage)).Where(&model.SecurityGroup{
		Name:   q.Name,
		Domain: q.Domain,
	}).Find(&securityGroups).Error; err != nil {
		return nil, err
	}

	return &securityGroups, nil
}

func (s *service) SecurityGroupTotalCount(q types.GetSecurityGroupsQuery) (int64, error) {
	var count int64
	if err := s.db.Model(&model.SecurityGroup{}).Where(&model.SecurityGroup{
		Name:   q.Name,
		Domain: q.Domain,
	}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *service) AddSchedulerInstanceToSecurityGroup(id, schedulerID string) error {
	securityGroup := model.SecurityGroup{}
	if err := s.db.First(&securityGroup, id).Error; err != nil {
		return err
	}

	schedulerInstance := model.SchedulerInstance{}
	if err := s.db.First(&schedulerInstance, schedulerID).Error; err != nil {
		return err
	}

	if err := s.db.Model(&securityGroup).Association("SecurityGroup").Append(&schedulerInstance); err != nil {
		return err
	}

	return nil
}

func (s *service) AddCDNInstanceToSecurityGroup(id, cdnID string) error {
	securityGroup := model.SecurityGroup{}
	if err := s.db.First(&securityGroup, id).Error; err != nil {
		return err
	}

	cdnInstance := model.CDNInstance{}
	if err := s.db.First(&cdnInstance, cdnID).Error; err != nil {
		return err
	}

	if err := s.db.Model(&securityGroup).Association("SecurityGroup").Append(&cdnInstance); err != nil {
		return err
	}

	return nil
}
