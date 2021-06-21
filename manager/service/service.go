package service

import (
	"d7y.io/dragonfly/v2/manager/database"
	"d7y.io/dragonfly/v2/manager/model"
	"d7y.io/dragonfly/v2/manager/types"
	rdbcache "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Service interface {
	CreateCDN(types.CreateCDNRequest) (*model.CDN, error)
	DestroyCDN(string) error
	UpdateCDN(string, types.UpdateCDNRequest) (*model.CDN, error)
	GetCDN(string) (*model.CDN, error)
	GetCDNs(types.GetCDNsQuery) (*[]model.CDN, error)
	CDNTotalCount() (int64, error)

	CreateCDNInstance(types.CreateCDNInstanceRequest) (*model.CDNInstance, error)
	DestroyCDNInstance(string) error
	UpdateCDNInstance(string, types.UpdateCDNInstanceRequest) (*model.CDNInstance, error)
	GetCDNInstance(string) (*model.CDNInstance, error)
	GetCDNInstances(types.GetCDNInstancesQuery) (*[]model.CDNInstance, error)
	CDNInstanceTotalCount() (int64, error)

	CreateScheduler(types.CreateSchedulerRequest) (*model.Scheduler, error)
	DestroyScheduler(string) error
	UpdateScheduler(string, types.UpdateSchedulerRequest) (*model.Scheduler, error)
	GetScheduler(string) (*model.Scheduler, error)
	GetSchedulers(types.GetSchedulersQuery) (*[]model.Scheduler, error)
	SchedulerTotalCount() (int64, error)

	CreateSecurityGroup(types.CreateSecurityGroupRequest) (*model.SecurityGroup, error)
	DestroySecurityGroup(string) error
	UpdateSecurityGroup(string, types.UpdateSecurityGroupRequest) (*model.SecurityGroup, error)
	GetSecurityGroup(string) (*model.SecurityGroup, error)
	GetSecurityGroups(types.GetSecurityGroupsQuery) (*[]model.SecurityGroup, error)
	SecurityGroupTotalCount() (int64, error)
}

type service struct {
	db    *gorm.DB
	rdb   *redis.Client
	cache *rdbcache.Cache
}

// Option is a functional option for service
type Option func(s *service)

// WithDatabase set the database client
func WithDatabase(database *database.Database) Option {
	return func(s *service) {
		s.db = database.DB
		s.rdb = database.RDB
		s.cache = database.Cache
	}
}

// New returns a new Service instence
func New(options ...Option) Service {
	s := &service{}

	for _, opt := range options {
		opt(s)
	}

	return s
}
