/*
Copyright 2023 cuisongliu@qq.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gh

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/json"
	"os"
)

var GlobalsGithubVar *GithubVar

type GithubVar struct {
	RunnerID            string
	SafeRepo            string
	IssueOrPRNumber     string
	CommentBody         string
	SenderOrCommentUser string
}

func (g *GithubVar) String() string {
	return "RunnerID: " + g.RunnerID + " SafeRepo: " + g.SafeRepo + " IssueOrPRNumber: " + g.IssueOrPRNumber + " CommentBody: " + g.CommentBody + " SenderOrCommentUser: " + g.SenderOrCommentUser
}

func GetGHEnvToVar() (*GithubVar, error) {
	gVar := new(GithubVar)
	gVar.RunnerID = os.Getenv("GITHUB_RUN_ID")
	//gVar.SafeRepo = os.Getenv("GITHUB_REPOSITORY")
	path := os.Getenv("GITHUB_EVENT_PATH")
	if path == "" {
		return nil, errors.New("GITHUB_EVENT_PATH is empty")
	}
	eventData, _ := os.ReadFile(path)
	var mData map[string]interface{}
	if err := json.Unmarshal(eventData, &mData); err != nil {
		return nil, errors.Wrap(err, "unmarshal github event data")
	}
	id, ok, _ := unstructured.NestedString(mData, "issue", "number")
	if !ok {
		id, _, _ = unstructured.NestedString(mData, "pull_request", "number")
	}
	gVar.IssueOrPRNumber = id
	gVar.SafeRepo, _, _ = unstructured.NestedString(mData, "repository", "full_name")
	gVar.CommentBody, _, _ = unstructured.NestedString(mData, "comment", "body")

	user, ok, _ := unstructured.NestedString(mData, "comment", "user", "login")
	if !ok {
		user, _, _ = unstructured.NestedString(mData, "comment", "user", "login")
	}
	gVar.SenderOrCommentUser = user
	return gVar, nil
}