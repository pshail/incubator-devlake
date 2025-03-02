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
	"strings"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractApiPrCommitsMeta = core.SubTaskMeta{
	Name:             "extractApiPullRequestCommits",
	EntryPoint:       ExtractApiPullRequestCommits,
	EnabledByDefault: true,
	Description:      "Extract raw PullRequestCommits data into tool layer table bitbucket_commits",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE_REVIEW},
}

type ApiPrCommitsResponse struct {
	Type   string    `json:"type"`
	Hash   string    `json:"hash"`
	Date   time.Time `json:"date"`
	Author struct {
		Type string                   `json:"type"`
		Raw  string                   `json:"raw"`
		User BitbucketAccountResponse `json:"user"`
	} `json:"author"`
	Message string `json:"message"`
	Summary struct {
		Type   string `json:"type"`
		Raw    string `json:"raw"`
		Markup string `json:"markup"`
		HTML   string `json:"html"`
	} `json:"summary"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
	Parents []struct {
		Type  string `json:"type"`
		Hash  string `json:"hash"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
		} `json:"links"`
	} `json:"parents"`
	Repository BitbucketApiRepo `json:"repository"`
}

func ExtractApiPullRequestCommits(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*BitbucketTaskData)
	repoId := data.Repo.BitbucketId
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: BitbucketApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			/*
				Table store raw data
			*/
			Table: RAW_PULL_REQUEST_COMMITS_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			apiPullRequestCommit := &ApiPrCommitsResponse{}
			if strings.HasPrefix(string(row.Data), "Not Found") {
				return nil, nil
			}
			err := errors.Convert(json.Unmarshal(row.Data, apiPullRequestCommit))
			if err != nil {
				return nil, err
			}
			pull := &BitbucketInput{}
			err = errors.Convert(json.Unmarshal(row.Input, pull))
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			results := make([]interface{}, 0, 3)
			bitbucketRepoCommit := &models.BitbucketRepoCommit{
				ConnectionId: data.Options.ConnectionId,
				RepoId:       repoId,
				CommitSha:    apiPullRequestCommit.Hash,
			}
			results = append(results, bitbucketRepoCommit)

			bitbucketCommit, err := convertPullRequestCommit(apiPullRequestCommit, data.Options.ConnectionId)
			if err != nil {
				return nil, err
			}
			results = append(results, bitbucketCommit)

			bitbucketPullRequestCommit := &models.BitbucketPrCommit{
				ConnectionId:  data.Options.ConnectionId,
				CommitSha:     apiPullRequestCommit.Hash,
				PullRequestId: pull.BitbucketId,
			}
			if err != nil {
				return nil, err
			}
			results = append(results, bitbucketPullRequestCommit)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertPullRequestCommit(prCommit *ApiPrCommitsResponse, connId uint64) (*models.BitbucketCommit, errors.Error) {
	bitbucketCommit := &models.BitbucketCommit{
		Sha:           prCommit.Hash,
		Message:       prCommit.Message,
		AuthorId:      prCommit.Author.User.AccountId,
		AuthorName:    prCommit.Author.User.UserName,
		AuthorEmail:   prCommit.Author.Raw,
		AuthoredDate:  prCommit.Date,
		CommittedDate: prCommit.Date,
		Url:           prCommit.Links.Self.Href,
	}
	return bitbucketCommit, nil
}
