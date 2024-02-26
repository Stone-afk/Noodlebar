// Copyright 2023 ecodeclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integration

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ecodeclub/ecache"

	"github.com/ecodeclub/ekit/iox"
	"github.com/ecodeclub/ginx/session"
	"github.com/ecodeclub/webook/internal/question/internal/integration/startup"
	"github.com/ecodeclub/webook/internal/question/internal/repository/dao"
	"github.com/ecodeclub/webook/internal/question/internal/web"
	"github.com/ecodeclub/webook/internal/test"
	testioc "github.com/ecodeclub/webook/internal/test/ioc"
	"github.com/ego-component/egorm"
	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/server/egin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const uid = 123

type HandlerTestSuite struct {
	suite.Suite
	server         *egin.Component
	db             *egorm.Component
	rdb            ecache.Cache
	dao            dao.QuestionDAO
	questionSetDAO dao.QuestionSetDAO
}

func (s *HandlerTestSuite) TearDownSuite() {
	err := s.db.Exec("DROP TABLE `answer_elements`").Error
	require.NoError(s.T(), err)
	err = s.db.Exec("DROP TABLE `questions`").Error
	require.NoError(s.T(), err)

	err = s.db.Exec("DROP TABLE `publish_answer_elements`").Error
	require.NoError(s.T(), err)
	err = s.db.Exec("DROP TABLE `publish_questions`").Error
	require.NoError(s.T(), err)

	err = s.db.Exec("DROP TABLE `question_sets`").Error
	require.NoError(s.T(), err)

	err = s.db.Exec("DROP TABLE `question_set_questions`").Error
	require.NoError(s.T(), err)
}

func (s *HandlerTestSuite) TearDownTest() {
	err := s.db.Exec("TRUNCATE TABLE `answer_elements`").Error
	require.NoError(s.T(), err)
	err = s.db.Exec("TRUNCATE TABLE `questions`").Error
	require.NoError(s.T(), err)

	err = s.db.Exec("TRUNCATE TABLE `publish_answer_elements`").Error
	require.NoError(s.T(), err)
	err = s.db.Exec("TRUNCATE TABLE `publish_questions`").Error
	require.NoError(s.T(), err)

	err = s.db.Exec("TRUNCATE TABLE `question_sets`").Error
	require.NoError(s.T(), err)

	err = s.db.Exec("TRUNCATE TABLE `question_set_questions`").Error
	require.NoError(s.T(), err)
}

func (s *HandlerTestSuite) SetupSuite() {
	handler, err := startup.InitHandler()
	require.NoError(s.T(), err)
	econf.Set("server", map[string]any{"contextTimeout": "1s"})
	server := egin.Load("server").Build()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("_session", session.NewMemorySession(session.Claims{
			Uid: uid,
		}))
	})
	handler.PrivateRoutes(server.Engine)
	s.server = server
	s.db = testioc.InitDB()
	err = dao.InitTables(s.db)
	require.NoError(s.T(), err)
	s.dao = dao.NewGORMQuestionDAO(s.db)
	s.questionSetDAO = dao.NewGORMQuestionSetDAO(s.db)
	s.rdb = testioc.InitCache()
}

func (s *HandlerTestSuite) TestSave() {
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)
		req    web.SaveReq

		wantCode int
		wantResp test.Result[int64]
	}{
		{
			//
			name: "全部新建",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				q, eles, err := s.dao.GetByID(ctx, 1)
				require.NoError(t, err)
				s.assertQuestion(t, dao.Question{
					Uid:     uid,
					Title:   "面试题1",
					Content: "面试题内容",
				}, q)
				assert.Equal(t, 4, len(eles))
			},
			req: web.SaveReq{
				Question: web.Question{
					Title:   "面试题1",
					Content: "面试题内容",
					Answer: web.Answer{
						Analysis:     s.buildAnswerEle(0),
						Basic:        s.buildAnswerEle(1),
						Intermediate: s.buildAnswerEle(2),
						Advanced:     s.buildAnswerEle(3),
					},
				},
			},
			wantCode: 200,
			wantResp: test.Result[int64]{
				Data: 1,
			},
		},
		{
			//
			name: "部分更新",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := s.db.WithContext(ctx).Create(&dao.Question{
					Id:      2,
					Uid:     uid,
					Title:   "老的标题",
					Content: "老的内容",
					Ctime:   123,
					Utime:   234,
				}).Error
				require.NoError(t, err)
				err = s.db.Create(&dao.AnswerElement{
					Id:        1,
					Qid:       2,
					Type:      dao.AnswerElementTypeAnalysis,
					Content:   "老的分析",
					Keywords:  "老的 keyword",
					Shorthand: "老的速记",
					Highlight: "老的亮点",
					Guidance:  "老的引导点",
					Ctime:     123,
					Utime:     123,
				}).Error
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				q, eles, err := s.dao.GetByID(ctx, 2)
				require.NoError(t, err)
				s.assertQuestion(t, dao.Question{
					Uid:     uid,
					Title:   "面试题1",
					Content: "新的内容",
				}, q)
				assert.Equal(t, 4, len(eles))
				analysis := eles[0]
				s.assertAnswerElement(t, dao.AnswerElement{
					Content:   "新的分析",
					Type:      dao.AnswerElementTypeAnalysis,
					Qid:       2,
					Keywords:  "新的 keyword",
					Shorthand: "新的速记",
					Highlight: "新的亮点",
					Guidance:  "新的引导点",
				}, analysis)
			},
			req: func() web.SaveReq {
				analysis := web.AnswerElement{
					Id:        1,
					Content:   "新的分析",
					Keywords:  "新的 keyword",
					Shorthand: "新的速记",
					Highlight: "新的亮点",
					Guidance:  "新的引导点",
				}
				return web.SaveReq{
					Question: web.Question{
						Id:      2,
						Title:   "面试题1",
						Content: "新的内容",
						Answer: web.Answer{
							Analysis:     analysis,
							Basic:        s.buildAnswerEle(1),
							Intermediate: s.buildAnswerEle(2),
							Advanced:     s.buildAnswerEle(3),
						},
					},
				}
			}(),
			wantCode: 200,
			wantResp: test.Result[int64]{
				Data: 2,
			},
		},
		{
			name: "非法访问",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := s.db.WithContext(ctx).Create(&dao.Question{
					Id:      3,
					Uid:     234,
					Title:   "老的标题",
					Content: "老的内容",
					Ctime:   123,
					Utime:   234,
				}).Error
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				q, _, err := s.dao.GetByID(ctx, 3)
				require.NoError(t, err)
				s.assertQuestion(t, dao.Question{
					Uid:     234,
					Title:   "老的标题",
					Content: "老的内容",
				}, q)
			},
			req: func() web.SaveReq {
				analysis := web.AnswerElement{
					Id:        1,
					Content:   "新的分析",
					Keywords:  "新的 keyword",
					Shorthand: "新的速记",
					Highlight: "新的亮点",
					Guidance:  "新的引导点",
				}
				return web.SaveReq{
					Question: web.Question{
						Id:      3,
						Title:   "面试题1",
						Content: "新的内容",
						Answer: web.Answer{
							Analysis:     analysis,
							Basic:        s.buildAnswerEle(1),
							Intermediate: s.buildAnswerEle(2),
							Advanced:     s.buildAnswerEle(3),
						},
					},
				}
			}(),
			wantCode: 500,
			wantResp: test.Result[int64]{
				Code: 502001,
				Msg:  "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost,
				"/question/save", iox.NewJSONReader(tc.req))
			req.Header.Set("content-type", "application/json")
			require.NoError(t, err)
			recorder := test.NewJSONResponseRecorder[int64]()
			s.server.ServeHTTP(recorder, req)
			require.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantResp, recorder.MustScan())
			tc.after(t)
			// 清理掉 123 的数据
			err = s.db.Exec("TRUNCATE table `questions`").Error
			require.NoError(t, err)
			err = s.db.Exec("TRUNCATE table `answer_elements`").Error
			require.NoError(t, err)
		})
	}
}

func (s *HandlerTestSuite) TestList() {
	// 插入一百条
	data := make([]dao.PublishQuestion, 0, 100)
	for idx := 0; idx < 100; idx++ {
		data = append(data, dao.PublishQuestion{
			Uid:     uid,
			Title:   fmt.Sprintf("这是标题 %d", idx),
			Content: fmt.Sprintf("这是解析 %d", idx),
		})
	}
	err := s.db.Create(&data).Error
	require.NoError(s.T(), err)
	testCases := []struct {
		name string
		req  web.Page

		wantCode int
		wantResp test.Result[web.QuestionList]
	}{
		{
			name: "获取成功",
			req: web.Page{
				Limit:  2,
				Offset: 0,
			},
			wantCode: 200,
			wantResp: test.Result[web.QuestionList]{
				Data: web.QuestionList{
					Total: 100,
					Questions: []web.Question{
						{
							Id:      100,
							Title:   "这是标题 99",
							Content: "这是解析 99",
							Utime:   time.UnixMilli(0).Format(time.DateTime),
						},
						{
							Id:      99,
							Title:   "这是标题 98",
							Content: "这是解析 98",
							Utime:   time.UnixMilli(0).Format(time.DateTime),
						},
					},
				},
			},
		},
		{
			name: "获取部分",
			req: web.Page{
				Limit:  2,
				Offset: 99,
			},
			wantCode: 200,
			wantResp: test.Result[web.QuestionList]{
				Data: web.QuestionList{
					Total: 100,
					Questions: []web.Question{
						{
							Id:      1,
							Title:   "这是标题 0",
							Content: "这是解析 0",
							Utime:   time.UnixMilli(0).Format(time.DateTime),
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost,
				"/question/pub/list", iox.NewJSONReader(tc.req))
			req.Header.Set("content-type", "application/json")
			require.NoError(t, err)
			recorder := test.NewJSONResponseRecorder[web.QuestionList]()
			s.server.ServeHTTP(recorder, req)
			require.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantResp, recorder.MustScan())
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = s.rdb.Delete(ctx, "webook:question:total")
	require.NoError(s.T(), err)
}

func (s *HandlerTestSuite) TestSync() {
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)
		req    web.SaveReq

		wantCode int
		wantResp test.Result[int64]
	}{
		{
			//
			name: "全部新建",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				q, eles, err := s.dao.GetPubByID(ctx, 1)
				require.NoError(t, err)
				s.assertQuestion(t, dao.Question{
					Uid:     uid,
					Title:   "面试题1",
					Content: "面试题内容",
				}, dao.Question(q))
				assert.Equal(t, 4, len(eles))
			},
			req: web.SaveReq{
				Question: web.Question{
					Title:   "面试题1",
					Content: "面试题内容",
					Answer: web.Answer{
						Analysis:     s.buildAnswerEle(0),
						Basic:        s.buildAnswerEle(1),
						Intermediate: s.buildAnswerEle(2),
						Advanced:     s.buildAnswerEle(3),
					},
				},
			},
			wantCode: 200,
			wantResp: test.Result[int64]{
				Data: 1,
			},
		},
		{
			//
			name: "部分更新",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := s.db.WithContext(ctx).Create(&dao.Question{
					Id:      2,
					Uid:     uid,
					Title:   "老的标题",
					Content: "老的内容",
					Ctime:   123,
					Utime:   234,
				}).Error
				require.NoError(t, err)
				err = s.db.Create(&dao.AnswerElement{
					Id:        1,
					Qid:       2,
					Type:      dao.AnswerElementTypeAnalysis,
					Content:   "老的分析",
					Keywords:  "老的 keyword",
					Shorthand: "老的速记",
					Highlight: "老的亮点",
					Guidance:  "老的引导点",
					Ctime:     123,
					Utime:     123,
				}).Error
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				q, eles, err := s.dao.GetByID(ctx, 2)
				require.NoError(t, err)
				s.assertQuestion(t, dao.Question{
					Uid:     uid,
					Title:   "面试题1",
					Content: "新的内容",
				}, q)
				assert.Equal(t, 4, len(eles))
				analysis := eles[0]
				s.assertAnswerElement(t, dao.AnswerElement{
					Content:   "新的分析",
					Type:      dao.AnswerElementTypeAnalysis,
					Qid:       2,
					Keywords:  "新的 keyword",
					Shorthand: "新的速记",
					Highlight: "新的亮点",
					Guidance:  "新的引导点",
				}, analysis)
			},
			req: func() web.SaveReq {
				analysis := web.AnswerElement{
					Id:        1,
					Content:   "新的分析",
					Keywords:  "新的 keyword",
					Shorthand: "新的速记",
					Highlight: "新的亮点",
					Guidance:  "新的引导点",
				}
				return web.SaveReq{
					Question: web.Question{
						Id:      2,
						Title:   "面试题1",
						Content: "新的内容",
						Answer: web.Answer{
							Analysis:     analysis,
							Basic:        s.buildAnswerEle(1),
							Intermediate: s.buildAnswerEle(2),
							Advanced:     s.buildAnswerEle(3),
						},
					},
				}
			}(),
			wantCode: 200,
			wantResp: test.Result[int64]{
				Data: 2,
			},
		},
		{
			name: "非法访问",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := s.db.WithContext(ctx).Create(&dao.Question{
					Id:      3,
					Uid:     234,
					Title:   "老的标题",
					Content: "老的内容",
					Ctime:   123,
					Utime:   234,
				}).Error
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				q, _, err := s.dao.GetByID(ctx, 3)
				require.NoError(t, err)
				s.assertQuestion(t, dao.Question{
					Uid:     234,
					Title:   "老的标题",
					Content: "老的内容",
				}, q)
			},
			req: func() web.SaveReq {
				analysis := web.AnswerElement{
					Id:        1,
					Content:   "新的分析",
					Keywords:  "新的 keyword",
					Shorthand: "新的速记",
					Highlight: "新的亮点",
					Guidance:  "新的引导点",
				}
				return web.SaveReq{
					Question: web.Question{
						Id:      3,
						Title:   "面试题1",
						Content: "新的内容",
						Answer: web.Answer{
							Analysis:     analysis,
							Basic:        s.buildAnswerEle(1),
							Intermediate: s.buildAnswerEle(2),
							Advanced:     s.buildAnswerEle(3),
						},
					},
				}
			}(),
			wantCode: 500,
			wantResp: test.Result[int64]{
				Code: 502001,
				Msg:  "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost,
				"/question/publish", iox.NewJSONReader(tc.req))
			req.Header.Set("content-type", "application/json")
			require.NoError(t, err)
			recorder := test.NewJSONResponseRecorder[int64]()
			s.server.ServeHTTP(recorder, req)
			require.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantResp, recorder.MustScan())
			tc.after(t)
			// 清理掉 123 的数据
			err = s.db.Exec("TRUNCATE table `questions`").Error
			require.NoError(t, err)
			err = s.db.Exec("TRUNCATE table `answer_elements`").Error
			require.NoError(t, err)
		})
	}
}

func (s *HandlerTestSuite) TestPubDetail() {
	// 插入一百条
	data := make([]dao.PublishQuestion, 0, 2)
	for idx := 0; idx < 2; idx++ {
		data = append(data, dao.PublishQuestion{
			Id:      int64(idx + 1),
			Uid:     uid,
			Title:   fmt.Sprintf("这是标题 %d", idx),
			Content: fmt.Sprintf("这是解析 %d", idx),
		})
	}
	err := s.db.Create(&data).Error
	require.NoError(s.T(), err)
	testCases := []struct {
		name string

		req      web.Qid
		wantCode int
		wantResp test.Result[web.Question]
	}{
		{
			name: "查询到了数据",
			req: web.Qid{
				Qid: 2,
			},
			wantCode: 200,
			wantResp: test.Result[web.Question]{
				Data: web.Question{
					Id:      2,
					Title:   "这是标题 1",
					Content: "这是解析 1",
					Utime:   time.UnixMilli(0).Format(time.DateTime),
				},
			},
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost,
				"/question/pub/detail", iox.NewJSONReader(tc.req))
			req.Header.Set("content-type", "application/json")
			require.NoError(t, err)
			recorder := test.NewJSONResponseRecorder[web.Question]()
			s.server.ServeHTTP(recorder, req)
			require.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantResp, recorder.MustScan())
		})
	}
}

func (s *HandlerTestSuite) buildAnswerEle(idx int64) web.AnswerElement {
	return web.AnswerElement{
		Content:   fmt.Sprintf("这是解析 %d", idx),
		Keywords:  fmt.Sprintf("关键字 %d", idx),
		Shorthand: fmt.Sprintf("快速记忆法 %d", idx),
		Highlight: fmt.Sprintf("亮点 %d", idx),
		Guidance:  fmt.Sprintf("引导点 %d", idx),
	}
}

// assertQuestion 不比较 id
func (s *HandlerTestSuite) assertQuestion(t *testing.T, expect dao.Question, q dao.Question) {
	assert.True(t, q.Id > 0)
	assert.True(t, q.Ctime > 0)
	assert.True(t, q.Utime > 0)
	q.Id = 0
	q.Ctime = 0
	q.Utime = 0
	assert.Equal(t, expect, q)
}

// assertAnswerElement 不包括 Id
func (s *HandlerTestSuite) assertAnswerElement(
	t *testing.T,
	expect dao.AnswerElement,
	ele dao.AnswerElement) {
	assert.True(t, ele.Id > 0)
	ele.Id = 0
	assert.True(t, ele.Ctime > 0)
	ele.Ctime = 0
	assert.True(t, ele.Utime > 0)
	ele.Utime = 0
	assert.Equal(t, expect, ele)
}

func (s *HandlerTestSuite) TestQuestionSet_Create() {
	var testCases = []struct {
		name  string
		after func(t *testing.T)
		req   web.CreateQuestionSetReq

		wantCode int
		wantResp test.Result[int64]
	}{
		{
			name: "创建成功1",
			after: func(t *testing.T) {
				t.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				qs, err := s.questionSetDAO.GetByID(ctx, 1)
				assert.NoError(t, err)

				s.assertQuestionSetEqual(t, dao.QuestionSet{
					Uid:         uid,
					Title:       "mysql",
					Description: "mysql相关面试题",
				}, qs)
			},
			req: web.CreateQuestionSetReq{
				Title:       "mysql",
				Description: "mysql相关面试题",
			},
			wantCode: 200,
			wantResp: test.Result[int64]{
				Data: 1,
			},
		},
		{
			name: "创建成功2",
			after: func(t *testing.T) {
				t.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				qs, err := s.questionSetDAO.GetByID(ctx, 2)
				assert.NoError(t, err)

				s.assertQuestionSetEqual(t, dao.QuestionSet{
					Uid:         uid,
					Title:       "mq",
					Description: "mq相关面试题",
				}, qs)
			},
			req: web.CreateQuestionSetReq{
				Title:       "mq",
				Description: "mq相关面试题",
			},
			wantCode: 200,
			wantResp: test.Result[int64]{
				Data: 2,
			},
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			targeURL := "/question-sets/create"
			req, err := http.NewRequest(http.MethodPost, targeURL, iox.NewJSONReader(tc.req))
			req.Header.Set("content-type", "application/json")
			require.NoError(t, err)

			recorder := test.NewJSONResponseRecorder[int64]()

			s.server.ServeHTTP(recorder, req)
			require.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantResp, recorder.MustScan())

			tc.after(t)
		})
	}
}

// assertQuestionSetEqual 不比较 id
func (s *HandlerTestSuite) assertQuestionSetEqual(t *testing.T, expect dao.QuestionSet, actual dao.QuestionSet) {
	assert.True(t, actual.Id > 0)
	assert.True(t, actual.Ctime > 0)
	assert.True(t, actual.Utime > 0)
	actual.Id = 0
	actual.Ctime = 0
	actual.Utime = 0
	assert.Equal(t, expect, actual)
}

func (s *HandlerTestSuite) TestQuestionSet_AddQuestions() {
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)
		req    web.AddQuestionsToQuestionSetReq

		wantCode int
		wantResp test.Result[int64]
	}{
		{
			name: "空题集_添加多个问题",
			before: func(t *testing.T) {
				t.Helper()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				// 创建一个空题集
				id, err := s.questionSetDAO.Create(ctx, dao.QuestionSet{
					Id:          5,
					Uid:         uid,
					Title:       "oss",
					Description: "oss题集",
				})
				require.Equal(t, int64(5), id)
				require.NoError(t, err)

				// 创建问题
				questions := []dao.Question{
					{
						Id:      4,
						Uid:     uid + 1,
						Title:   "oss问题1",
						Content: "oss问题1",
						Ctime:   123,
						Utime:   234,
					},
					{
						Id:      5,
						Uid:     uid + 2,
						Title:   "oss问题2",
						Content: "oss问题2",
						Ctime:   1234,
						Utime:   2345,
					},
				}
				for _, q := range questions {
					require.NoError(t, s.db.WithContext(ctx).Create(q).Error)
				}

				// 题集中题目为0
				qs, err := s.questionSetDAO.GetQuestionsByID(ctx, id)
				require.NoError(t, err)
				require.Equal(t, 0, len(qs))
			},
			after: func(t *testing.T) {
				t.Helper()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				expected := []dao.Question{
					{
						Uid:     uid + 1,
						Title:   "oss问题1",
						Content: "oss问题1",
					},
					{
						Uid:     uid + 2,
						Title:   "oss问题2",
						Content: "oss问题2",
					},
				}

				actual, err := s.questionSetDAO.GetQuestionsByID(ctx, 5)
				require.NoError(t, err)
				require.Equal(t, len(expected), len(actual))

				for i := 0; i < len(expected); i++ {
					s.assertQuestion(t, expected[i], actual[i])
				}

			},
			req: web.AddQuestionsToQuestionSetReq{
				QSID: 5,
				Questions: []web.Question{
					{
						Id:      4,
						Title:   "oss问题1",
						Content: "oss问题1",
					},
					{
						Id:      5,
						Title:   "oss问题2",
						Content: "oss问题2",
					},
				},
			},
			wantCode: 200,
			wantResp: test.Result[int64]{},
		},
		{
			name: "非空题集_添加多个问题",
			before: func(t *testing.T) {
				t.Helper()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				// 创建一个空题集
				id, err := s.questionSetDAO.Create(ctx, dao.QuestionSet{
					Id:          7,
					Uid:         uid,
					Title:       "Go",
					Description: "Go题集",
				})
				require.Equal(t, int64(7), id)
				require.NoError(t, err)

				// 创建问题
				questions := []dao.Question{
					{
						Id:      14,
						Uid:     uid + 1,
						Title:   "Go问题1",
						Content: "Go问题1",
						Ctime:   123,
						Utime:   234,
					},
					{
						Id:      15,
						Uid:     uid + 2,
						Title:   "Go问题2",
						Content: "Go问题2",
						Ctime:   1234,
						Utime:   2345,
					},
					{
						Id:      16,
						Uid:     uid + 3,
						Title:   "Go问题3",
						Content: "Go问题3",
						Ctime:   1234,
						Utime:   2345,
					},
				}
				for _, q := range questions {
					require.NoError(t, s.db.WithContext(ctx).Create(q).Error)
				}

				require.NoError(t, s.questionSetDAO.AddQuestionsByID(ctx, id, questions[:1]))

				// 题集中题目为1
				qs, err := s.questionSetDAO.GetQuestionsByID(ctx, id)
				require.NoError(t, err)
				require.Equal(t, 1, len(qs))
			},
			after: func(t *testing.T) {
				t.Helper()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				expected := []dao.Question{
					{
						Uid:     uid + 1,
						Title:   "Go问题1",
						Content: "Go问题1",
					},
					{
						Uid:     uid + 2,
						Title:   "Go问题2",
						Content: "Go问题2",
					},
					{
						Uid:     uid + 3,
						Title:   "Go问题3",
						Content: "Go问题3",
					},
				}

				actual, err := s.questionSetDAO.GetQuestionsByID(ctx, 7)
				require.NoError(t, err)
				require.Equal(t, len(expected), len(actual))

				for i := 0; i < len(expected); i++ {
					s.assertQuestion(t, expected[i], actual[i])
				}

			},
			req: web.AddQuestionsToQuestionSetReq{
				QSID: 7,
				Questions: []web.Question{
					{
						Id:      15,
						Title:   "Go问题2",
						Content: "Go问题2",
					},
					{
						Id:      16,
						Title:   "Go问题3",
						Content: "Go问题3",
					},
				},
			},
			wantCode: 200,
			wantResp: test.Result[int64]{},
		},
		{
			name: "非空题集_添加多个问题_含已添加问题",
			before: func(t *testing.T) {
				t.Helper()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				// 创建一个空题集
				id, err := s.questionSetDAO.Create(ctx, dao.QuestionSet{
					Id:          9,
					Uid:         uid,
					Title:       "js",
					Description: "js题集",
				})
				require.Equal(t, int64(9), id)
				require.NoError(t, err)

				// 创建问题
				questions := []dao.Question{
					{
						Id:      114,
						Uid:     uid + 1,
						Title:   "js问题1",
						Content: "js问题1",
						Ctime:   123,
						Utime:   234,
					},
					{
						Id:      115,
						Uid:     uid + 2,
						Title:   "js问题2",
						Content: "js问题2",
						Ctime:   1234,
						Utime:   2345,
					},
					{
						Id:      116,
						Uid:     uid + 3,
						Title:   "js问题3",
						Content: "js问题3",
						Ctime:   1234,
						Utime:   2345,
					},
				}
				for _, q := range questions {
					require.NoError(t, s.db.WithContext(ctx).Create(q).Error)
				}

				require.NoError(t, s.questionSetDAO.AddQuestionsByID(ctx, id, questions[:1]))

				// 题集中题目为1
				qs, err := s.questionSetDAO.GetQuestionsByID(ctx, id)
				require.NoError(t, err)
				require.Equal(t, 1, len(qs))
			},
			after: func(t *testing.T) {
				t.Helper()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				expected := []dao.Question{
					{
						Uid:     uid + 1,
						Title:   "js问题1",
						Content: "js问题1",
					},
				}

				actual, err := s.questionSetDAO.GetQuestionsByID(ctx, 9)
				require.NoError(t, err)
				require.Equal(t, len(expected), len(actual))

				for i := 0; i < len(expected); i++ {
					s.assertQuestion(t, expected[i], actual[i])
				}

			},
			req: web.AddQuestionsToQuestionSetReq{
				QSID: 9,
				Questions: []web.Question{
					{
						Id:      114,
						Title:   "js问题1",
						Content: "js问题1",
					},
					{
						Id:      116,
						Title:   "js问题3",
						Content: "js问题3",
					},
				},
			},
			// todo: 为什么是500?
			wantCode: 500,
			wantResp: test.Result[int64]{
				Code: 402001,
				Msg:  "部分题目已添加",
			},
		},
		// {
		// 	name: "题集不存在",
		// 	before: func(t *testing.T) {
		// 		t.Helper()
		//
		// 	},
		// 	after: func(t *testing.T) {
		// 		t.Helper()
		// 	},
		// 	req:      web.AddQuestionsToQuestionSetReq{},
		// 	wantCode: 500,
		// 	wantResp: test.Result[int64]{},
		// },
		// {
		// 	name: "当前用户并非题集的创建者",
		// 	before: func(t *testing.T) {
		// 		t.Helper()
		//
		// 	},
		// 	after: func(t *testing.T) {
		// 		t.Helper()
		// 	},
		// 	req:      web.AddQuestionsToQuestionSetReq{},
		// 	wantCode: 200,
		// 	wantResp: test.Result[int64]{},
		// },
		// {
		// 	name: "待添加/删除的问题不存在",
		// 	before: func(t *testing.T) {
		// 		t.Helper()
		//
		// 	},
		// 	after: func(t *testing.T) {
		// 		t.Helper()
		// 	},
		// 	req:      web.AddQuestionsToQuestionSetReq{},
		// 	wantCode: 500,
		// 	wantResp: test.Result[int64]{
		// 		Code: 502001,
		// 		Msg:  "系统错误",
		// 	},
		// },
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost,
				"/question-sets/add", iox.NewJSONReader(tc.req))
			req.Header.Set("content-type", "application/json")
			require.NoError(t, err)
			recorder := test.NewJSONResponseRecorder[int64]()
			s.server.ServeHTTP(recorder, req)
			require.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantResp, recorder.MustScan())
			tc.after(t)
			// // 清理掉 123 的数据
			// err = s.db.Exec("TRUNCATE table `question_sets`").Error
			// require.NoError(t, err)
			// err = s.db.Exec("TRUNCATE table `answer_elements`").Error
			// require.NoError(t, err)
		})
	}
}

func (s *HandlerTestSuite) TestQuestionSet_DeleteQuestions() {
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)
		req    web.DeleteQuestionsFromQuestionSetReq

		wantCode int
		wantResp test.Result[int64]
	}{
		{
			name: "非空题集_删除全部问题",
			before: func(t *testing.T) {
				t.Helper()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				// 创建一个空题集
				id, err := s.questionSetDAO.Create(ctx, dao.QuestionSet{
					Id:          217,
					Uid:         uid,
					Title:       "Go",
					Description: "Go题集",
				})
				require.Equal(t, int64(217), id)
				require.NoError(t, err)

				// 创建问题
				questions := []dao.Question{
					{
						Id:      214,
						Uid:     uid + 1,
						Title:   "Go问题1",
						Content: "Go问题1",
						Ctime:   123,
						Utime:   234,
					},
					{
						Id:      215,
						Uid:     uid + 2,
						Title:   "Go问题2",
						Content: "Go问题2",
						Ctime:   1234,
						Utime:   2345,
					},
					{
						Id:      216,
						Uid:     uid + 2,
						Title:   "Go问题3",
						Content: "Go问题3",
						Ctime:   1234,
						Utime:   2345,
					},
				}
				for _, q := range questions {
					require.NoError(t, s.db.WithContext(ctx).Create(q).Error)
				}

				err = s.questionSetDAO.AddQuestionsByID(ctx, id, questions)
				require.NoError(t, err)

				qs, err := s.questionSetDAO.GetQuestionsByID(ctx, id)
				require.NoError(t, err)
				require.Equal(t, len(questions), len(qs))

			},
			after: func(t *testing.T) {
				t.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				qs, err := s.questionSetDAO.GetQuestionsByID(ctx, 217)
				require.NoError(t, err)
				require.Equal(t, 0, len(qs))
			},
			req: web.DeleteQuestionsFromQuestionSetReq{
				QSID: 217,
				Questions: []web.Question{
					{
						Id:      214,
						Title:   "Go问题1",
						Content: "Go问题1",
					},
					{
						Id:      215,
						Title:   "Go问题2",
						Content: "Go问题2",
					},
					{
						Id:      216,
						Title:   "Go问题3",
						Content: "Go问题3",
					},
				},
			},
			wantCode: 200,
			wantResp: test.Result[int64]{},
		},
		{
			name: "非空题集_删除部分问题",
			before: func(t *testing.T) {
				t.Helper()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				// 创建一个空题集
				id, err := s.questionSetDAO.Create(ctx, dao.QuestionSet{
					Id:          218,
					Uid:         uid,
					Title:       "Go",
					Description: "Go题集",
				})
				require.Equal(t, int64(218), id)
				require.NoError(t, err)

				// 创建问题
				questions := []dao.Question{
					{
						Id:      314,
						Uid:     uid + 1,
						Title:   "Go问题1",
						Content: "Go问题1",
						Ctime:   123,
						Utime:   234,
					},
					{
						Id:      315,
						Uid:     uid + 2,
						Title:   "Go问题2",
						Content: "Go问题2",
						Ctime:   1234,
						Utime:   2345,
					},
					{
						Id:      316,
						Uid:     uid + 2,
						Title:   "Go问题3",
						Content: "Go问题3",
						Ctime:   1234,
						Utime:   2345,
					},
				}
				for _, q := range questions {
					require.NoError(t, s.db.WithContext(ctx).Create(q).Error)
				}

				err = s.questionSetDAO.AddQuestionsByID(ctx, id, questions)
				require.NoError(t, err)

				qs, err := s.questionSetDAO.GetQuestionsByID(ctx, id)
				require.NoError(t, err)
				require.Equal(t, len(questions), len(qs))
			},
			after: func(t *testing.T) {
				t.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				qs, err := s.questionSetDAO.GetQuestionsByID(ctx, 218)
				require.NoError(t, err)
				require.Equal(t, 1, len(qs))
				s.assertQuestion(t, dao.Question{
					Uid:     uid + 2,
					Title:   "Go问题2",
					Content: "Go问题2",
				}, qs[0])

			},
			req: web.DeleteQuestionsFromQuestionSetReq{
				QSID: 218,
				Questions: []web.Question{
					{
						Id:      314,
						Title:   "Go问题1",
						Content: "Go问题1",
					},
					{
						Id:      316,
						Title:   "Go问题3",
						Content: "Go问题3",
					},
				},
			},
			wantCode: 200,
			wantResp: test.Result[int64]{},
		},
		// {
		// 	name: "题集不存在",
		// 	before: func(t *testing.T) {
		// 		t.Helper()
		//
		// 	},
		// 	after: func(t *testing.T) {
		// 		t.Helper()
		// 	},
		// 	req:      web.AddQuestionsToQuestionSetReq{},
		// 	wantCode: 500,
		// 	wantResp: test.Result[int64]{},
		// },
		// {
		// 	name: "当前用户并非题集的创建者",
		// 	before: func(t *testing.T) {
		// 		t.Helper()
		//
		// 	},
		// 	after: func(t *testing.T) {
		// 		t.Helper()
		// 	},
		// 	req:      web.AddQuestionsToQuestionSetReq{},
		// 	wantCode: 200,
		// 	wantResp: test.Result[int64]{},
		// },
		// {
		// 	name: "待添加/删除的问题不存在",
		// 	before: func(t *testing.T) {
		// 		t.Helper()
		//
		// 	},
		// 	after: func(t *testing.T) {
		// 		t.Helper()
		// 	},
		// 	req:      web.AddQuestionsToQuestionSetReq{},
		// 	wantCode: 500,
		// 	wantResp: test.Result[int64]{
		// 		Code: 502001,
		// 		Msg:  "系统错误",
		// 	},
		// },
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost,
				"/question-sets/delete", iox.NewJSONReader(tc.req))
			req.Header.Set("content-type", "application/json")
			require.NoError(t, err)
			recorder := test.NewJSONResponseRecorder[int64]()
			s.server.ServeHTTP(recorder, req)
			require.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantResp, recorder.MustScan())
			tc.after(t)
		})
	}
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
