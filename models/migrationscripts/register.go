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

package migrationscripts

import (
	"github.com/apache/incubator-devlake/plugins/core"
)

// All return all the migration scripts of framework
func All() []core.MigrationScript {
	return []core.MigrationScript{
		new(addFrameworkTables),
		new(renamePipelineStepToStage),
		new(addSubtaskToTaskTable),
		new(addBlueprintMode),
		new(renameTasksToPlan),
		new(resetDomainTables),
		new(addCommitFileComponent),
		new(removeNotes),
		new(addProjectMapping),
		new(renameColumnsOfPullRequestIssue),
		new(addNoPKModelToCommitParent),
		new(addSubtasksTable),
		new(addCICDTables),
		new(renameColumnsOfPrCommentIssueComment),
		new(modifyTablesForDora),
		new(addTypeToBoard),
		new(encryptBlueprint),
		new(encryptPipeline),
		new(modifyCicdPipeline),
		new(modifyCICDTasks),
		new(addOriginChangeValueForPr),
		new(fixCommitFileIdTooLong),
		new(addRawDataOriginToBoardRepos),
		new(renamePipelineCommits),
		new(commitLineChange),
		new(changeLeadTimeMinutesToInt64),
		new(addRepoSnapshot),
		new(createCollectorState),
		new(removeCicdPipelineRelation),
		new(addCicdScopeDropBuildsJobs),
		new(addSkipOnFail),
		new(modifyCommitsDiffs),
		new(addProjectPrMetric),
		new(addProjectTables),
		new(addProjectToBluePrint),
		new(addProjectIssueMetric),
		new(addLabels),
		new(renameFiledsInProjectPrMetric),
		new(addEnableToProjectMetric),
		new(addCollectorMeta20221125),
		new(addOriginalProject),
		new(addErrorName),
		new(encryptTask221221),
		new(renameProjectMetrics),
		new(addOriginalTypeToIssue221230),
	}
}
