// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"encoding/json"
	"github.com/goharbor/harbor/src/api/artifact"
	"github.com/goharbor/harbor/src/api/repository"
	"github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/pkg/tag/model/tag"
	arttesting "github.com/goharbor/harbor/src/testing/api/artifact"
	repotesting "github.com/goharbor/harbor/src/testing/api/repository"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type tagTestSuite struct {
	suite.Suite
	originalRepoCtl repository.Controller
	repoCtl         *repotesting.FakeController
	originalArtCtl  artifact.Controller
	artCtl          *arttesting.FakeController
}

func (c *tagTestSuite) SetupSuite() {
	c.originalArtCtl = artifact.Ctl
	c.originalRepoCtl = repository.Ctl
}

func (c *tagTestSuite) SetupTest() {
	c.artCtl = &arttesting.FakeController{}
	artifact.Ctl = c.artCtl
	c.repoCtl = &repotesting.FakeController{}
	repository.Ctl = c.repoCtl
}

func (c *tagTestSuite) TearDownTest() {
}

func (c *tagTestSuite) TearDownSuite() {
	repository.Ctl = c.originalRepoCtl
	artifact.Ctl = c.originalArtCtl
}

func (c *tagTestSuite) TestListTag() {
	c.SetupTest()
	req := httptest.NewRequest(http.MethodGet, "/v2/library/hello-world/tags/list", nil)
	var w *httptest.ResponseRecorder
	c.repoCtl.On("GetByName").Return(&models.RepoRecord{
		RepositoryID: 1,
		Name:         "library/hello-world",
	}, nil)
	c.artCtl.On("ListTags").Return(2, []*artifact.Tag{
		{
			Tag: tag.Tag{
				RepositoryID: 1,
				Name:         "v1",
			},
		},
		{
			Tag: tag.Tag{
				RepositoryID: 1,
				Name:         "v2",
			},
		},
	}, nil)
	w = httptest.NewRecorder()
	newTagHandler().ServeHTTP(w, req)
	c.Equal(http.StatusOK, w.Code)
	var tagsAPIResponse struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}
	decoder := json.NewDecoder(w.Body)
	err := decoder.Decode(&tagsAPIResponse)
	c.Nil(err)
	c.Equal(2, len(tagsAPIResponse.Tags))
}

func (c *tagTestSuite) TestListTagPagination1() {
	c.SetupTest()
	req := httptest.NewRequest(http.MethodGet, "/v2/hello-world/tags/list?n=1", nil)
	var w *httptest.ResponseRecorder
	c.repoCtl.On("GetByName").Return(&models.RepoRecord{
		RepositoryID: 1,
		Name:         "hello-world",
	}, nil)
	c.artCtl.On("ListTags").Return(2, []*artifact.Tag{
		{
			Tag: tag.Tag{
				RepositoryID: 1,
				Name:         "v1",
			},
		},
		{
			Tag: tag.Tag{
				RepositoryID: 1,
				Name:         "v2",
			},
		},
	}, nil)
	w = httptest.NewRecorder()
	newTagHandler().ServeHTTP(w, req)
	c.Equal(http.StatusOK, w.Code)
	var ctlg struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}
	decoder := json.NewDecoder(w.Body)
	err := decoder.Decode(&ctlg)
	c.Nil(err)
	c.Equal(1, len(ctlg.Tags))
	c.Equal("v1", ctlg.Tags[0])
}

func (c *tagTestSuite) TestListTagPagination2() {
	c.SetupTest()
	req := httptest.NewRequest(http.MethodGet, "/v2/hello-world/tags/list?n=3", nil)
	var w *httptest.ResponseRecorder
	c.repoCtl.On("GetByName").Return(&models.RepoRecord{
		RepositoryID: 1,
		Name:         "hello-world",
	}, nil)
	c.artCtl.On("ListTags").Return(2, []*artifact.Tag{
		{
			Tag: tag.Tag{
				RepositoryID: 1,
				Name:         "v1",
			},
		},
		{
			Tag: tag.Tag{
				RepositoryID: 1,
				Name:         "v2",
			},
		},
	}, nil)
	w = httptest.NewRecorder()
	newTagHandler().ServeHTTP(w, req)
	c.Equal(http.StatusOK, w.Code)
	var ctlg struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}
	decoder := json.NewDecoder(w.Body)
	err := decoder.Decode(&ctlg)
	c.Nil(err)
	c.Equal(2, len(ctlg.Tags))
	c.Equal("v2", ctlg.Tags[1])
}

func (c *tagTestSuite) TestListTagPagination3() {
	c.SetupTest()
	req := httptest.NewRequest(http.MethodGet, "/v2/hello-world/tags/list?last=v1&n=1", nil)
	var w *httptest.ResponseRecorder
	c.repoCtl.On("GetByName").Return(&models.RepoRecord{
		RepositoryID: 1,
		Name:         "hello-world",
	}, nil)
	c.artCtl.On("ListTags").Return(2, []*artifact.Tag{
		{
			Tag: tag.Tag{
				RepositoryID: 1,
				Name:         "v1",
			},
		},
		{
			Tag: tag.Tag{
				RepositoryID: 1,
				Name:         "v2",
			},
		},
	}, nil)
	w = httptest.NewRecorder()
	newTagHandler().ServeHTTP(w, req)
	c.Equal(http.StatusOK, w.Code)
	var ctlg struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}
	decoder := json.NewDecoder(w.Body)
	err := decoder.Decode(&ctlg)
	c.Nil(err)
	c.Equal(1, len(ctlg.Tags))
	c.Equal("v2", ctlg.Tags[0])
}

func TestTagTestSuite(t *testing.T) {
	suite.Run(t, &tagTestSuite{})
}
