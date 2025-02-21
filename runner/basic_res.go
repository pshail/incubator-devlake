/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package runner

import (
	"fmt"
	"sync"

	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/impl"
	"github.com/apache/incubator-devlake/impl/dalgorm"
	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/plugins/core"
	"gorm.io/gorm"
)

var app_lock sync.Mutex
var app_inited bool

// CreateAppBasicRes returns a application level BasicRes instance based on .env/environment variables
// it is useful because multiple places need BasicRes including `main.go` `directrun` and `worker`
// keep in mind this function can be called only once
func CreateAppBasicRes() core.BasicRes {
	app_lock.Lock()
	if app_inited {
		panic(fmt.Errorf("CreateAppBasicRes can be called once"))
	}
	app_inited = true
	app_lock.Unlock()
	cfg := config.GetConfig()
	log := logger.Global
	db, err := NewGormDb(cfg, logger.Global)
	if err != nil {
		panic(err)
	}
	dalgorm.Init(cfg.GetString(core.EncodeKeyEnvStr))
	return CreateBasicRes(cfg, log, db)
}

// CreateBasicRes returns a BasicRes based on what was given
func CreateBasicRes(cfg core.ConfigReader, log core.Logger, db *gorm.DB) core.BasicRes {
	return impl.NewDefaultBasicRes(cfg, log, dalgorm.NewDalgorm(db))
}
