package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	etcdregistry "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Data holds the database connection.
type Data struct {
	DB *gorm.DB
}

// DataOption config the Data creation.
type DataOption func(*dataOptions)

type dataOptions struct {
	plugins []gorm.Plugin
}

// WithPlugins adds GORM plugins to the database connection.
func WithPlugins(plugins ...gorm.Plugin) DataOption {
	return func(o *dataOptions) {
		o.plugins = append(o.plugins, plugins...)
	}
}

// NewData creates a new Data with database connection.
func NewData(dsn string, logger log.Logger, opts ...DataOption) (*Data, func(), error) {
	options := &dataOptions{}
	for _, opt := range opts {
		opt(options)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	for _, plugin := range options.plugins {
		if err := db.Use(plugin); err != nil {
			return nil, nil, err
		}
	}

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{DB: db}, cleanup, nil
}

// NewDiscovery creates an etcd service discovery client.
func NewDiscovery(endpoints []string) (registry.Discovery, error) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		return nil, err
	}
	return etcdregistry.New(etcdClient), nil
}
