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

package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"net/http"
	"net/url"
)

const RAW_RUN_TABLE = "github_api_runs"

var CollectRunsMeta = core.SubTaskMeta{
	Name:             "collectRuns",
	EntryPoint:       CollectRuns,
	EnabledByDefault: true,
	Description:      "Collect Runs data from Github action api",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func CollectRuns(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*GithubTaskData)
	collectorWithState, err := helper.NewApiCollectorWithState(helper.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: GithubApiParams{
			ConnectionId: data.Options.ConnectionId,
			Name:         data.Options.Name,
		},
		Table: RAW_RUN_TABLE,
	}, data.CreatedDateAfter)
	if err != nil {
		return err
	}

	//incremental := collectorWithState.IsIncremental()
	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		ApiClient: data.ApiClient,
		PageSize:  30,
		//Incremental: incremental,
		UrlTemplate: "repos/{{ .Params.Name }}/actions/runs",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			// if data.CreatedDateAfter != nil, we set since once
			// There is a bug for github rest api, so temporarily commented the following code
			//if data.CreatedDateAfter != nil {
			//	startDate := data.CreatedDateAfter.Format("2006-01-02")
			//	query.Set("created", fmt.Sprintf("%s..*", startDate))
			//}
			//// if incremental == true, we overwrite it
			//if incremental {
			//	startDate := collectorWithState.LatestState.LatestSuccessStart.Format("2006-01-02")
			//	query.Set("created", fmt.Sprintf("%s..*", startDate))
			//}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			body := &GithubRawRunsResult{}
			err := helper.UnmarshalResponse(res, body)
			if err != nil {
				return nil, err
			}
			return body.GithubWorkflowRuns, nil
		},
	})

	if err != nil {
		return err
	}

	return collectorWithState.Execute()
}

type GithubRawRunsResult struct {
	TotalCount         int64             `json:"total_count"`
	GithubWorkflowRuns []json.RawMessage `json:"workflow_runs"`
}
