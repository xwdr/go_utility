package utils

import (
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 数据库字典
var dbs = make(map[string]*gorm.DB)

// 读写锁
var lock = sync.RWMutex{}

// Database 数据库配置
type DatabaseConfig struct {
	DbMode      bool
	DbName      string
	Type        string
	Dsn         string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
}

// NewInstance 根据配置名字创建并返回 DB 连接池对象
func NewInstance(dbConf *DatabaseConfig, conf ...*gorm.Config) *gorm.DB {
	lock.RLock()
	db := dbs[dbConf.DbName]
	lock.RUnlock()
	// 存在链接则直接返回
	if db != nil {
		return db
	}

	// 数据库类型
	if len(dbConf.Type) == 0 { // 默认为mysql
		dbConf.Type = "mysql"
	}

	// 设置gorm config
	c := &gorm.Config{}
	if len(conf) > 0 {
		c = conf[0]
	}

	// 初始化
	db, err := gorm.Open(mysql.Open(dbConf.Dsn), c)
	if err != nil {
		panic(err.Error())
	}

	// 判断是否需要打印数据库debug日志
	if dbConf.DbMode {
		db = db.Debug()
	}

	// 初始化连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(err.Error())
	}
	sqlDB.SetMaxIdleConns(dbConf.MaxIdle)
	sqlDB.SetMaxOpenConns(dbConf.MaxActive)
	sqlDB.SetConnMaxLifetime(time.Duration(dbConf.IdleTimeout) * time.Second)

	// 设置连接对象
	lock.Lock()
	dbs[dbConf.DbName] = db
	lock.Unlock()
	return db
}

// ReleaseInstance 关闭所有 DB 连接
// 新调用 Get 方法时会使用最新 DB 配置创建连接
func ReleaseInstance() {
	for k, _ := range dbs {
		delete(dbs, k)
	}
}
