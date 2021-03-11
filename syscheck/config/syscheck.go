// Create package in v.1.0.0.
// config package contains App global variable with config value about syscheck from environment variable or config file
// App return field value from method having same name with that field name

// config.go is file that define syscheckConfig type which is type of App
// Also, App implement various config interface each of package in syscheck domain by declaring method

package config

import (
	"github.com/spf13/viper"
)

// App is the application config about syscheck domain
var App *syscheckConfig

// syscheckConfig having config value and implement various interface about Config by declaring method
type syscheckConfig struct {
	// indexName represent name of elasticsearch index including syscheck history document
	indexName *string

	// indexShardNum represent shard number of elasticsearch index storing syscheck history document
	indexShardNum *int

	// indexReplicaNum represent replica number of elasticsearch index to replace index when node become unable
	indexReplicaNum *int
}

func (sc *syscheckConfig) IndexName() string {
	var key = "domain.syscheck.elasticsearch.index.name"
	if sc.indexName == nil {
		sc.indexName = _string(viper.GetString(key))
	}
	return *sc.indexName
}

func (sc *syscheckConfig) IndexShardNum() int {
	var key = "domain.syscheck.elasticsearch.index.shardNum"
	if sc.indexShardNum == nil {
		sc.indexShardNum = _int(viper.GetInt(key))
	}
	return *sc.indexShardNum
}

func (sc *syscheckConfig) IndexReplicaNum() int {
	var key = "domain.syscheck.elasticsearch.index.replicaNum"
	if sc.indexReplicaNum == nil {
		sc.indexReplicaNum = _int(viper.GetInt(key))
	}
	return *sc.indexReplicaNum
}

func init() {
	App = &syscheckConfig{}
}

func _string(s string) *string { return &s }
func _int(i int) *int {return &i}
