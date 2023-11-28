package data

import (
	"github.com/chainxx/bitx/internal/conf"
	"github.com/google/wire"
	"gorm.io/gorm"
	
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
)

var ProviderSet = wire.NewSet(NewData, NewWalletDataSource, NewMarketDataResource)

// Data .
type Data struct {
	// TODO warpped database client
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	// database
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	
	if err != nil {
		_ = logger.Log(log.LevelError, "failed opening connection to mysql: %v", err)
		return nil, nil, err
	}
	
	d := &Data{
		db: db.Debug(),
	}
	
	// 清理资源
	cleanup := func() {
		_ = logger.Log(log.LevelInfo, "closing the data resources")
	}
	
	return d, cleanup, nil
}
