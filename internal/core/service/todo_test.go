package service_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"

	"learning/hyssa-learn/generated/todo_service"

	"github.com/go-testfixtures/testfixtures/v3"
)

func (s *TestSuite) TestCreateTodo() {
	requestBody := `{"title": "test_name"}`
	res, err := s.server.Client().Post(s.server.URL+"/v1/todo", "", bytes.NewBufferString(requestBody))
	s.NoError(err)

	defer res.Body.Close()

	s.Equal(http.StatusOK, res.StatusCode)

	response := todo_service.Todo{}
	err = json.NewDecoder(res.Body).Decode(&response)
	s.NoError(err)

	s.Equal("test_name", response.Title)
}

func (s *TestSuite) TestGetTodo() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("./../../../pkg/fixtures"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())

	res, err := s.server.Client().Get(s.server.URL + "/v1/todo/1")
	s.Require().NoError(err)

	defer res.Body.Close()

	s.Require().Equal(http.StatusOK, res.StatusCode)

	response := todo_service.Todo{}
	err = json.NewDecoder(res.Body).Decode(&response)
	s.Require().NoError(err)

	s.Assert().Equal(int32(1), response.Id)
	s.Assert().Equal("test_name", response.Title)
	s.Assert().Equal(false, response.Completed)
}

func (s *TestSuite) TestUpdateTodo() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("./../../../pkg/fixtures"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())

	resBody := map[string]interface{}{
		"id":        1,
		"title":     "patched",
		"completed": true,
	}

	bodyBytes, err := json.Marshal(resBody)
	s.Require().NoError(err)

	req, err := http.NewRequest(http.MethodPatch, s.server.URL+"/v1/todo", bytes.NewBuffer(bodyBytes))
	s.Require().NoError(err)

	res, err := s.server.Client().Do(req)
	s.Require().NoError(err)

	defer res.Body.Close()

	s.Require().Equal(http.StatusOK, res.StatusCode)

	response := todo_service.Todo{}
	err = json.NewDecoder(res.Body).Decode(&response)
	s.Require().NoError(err)

	s.Assert().Equal(int32(1), response.Id)
	s.Assert().Equal("patched", response.Title)
	s.Assert().Equal(true, response.Completed)
}

func (s *TestSuite) TestGetAllTodos() {
	res, err := s.server.Client().Get(s.server.URL + "/v1/todos")
	s.Require().NoError(err)

	defer res.Body.Close()

	s.Require().Equal(http.StatusOK, res.StatusCode)

	response := todo_service.GetAllTodosResponse{}
	err = json.NewDecoder(res.Body).Decode(&response)
	s.Require().NoError(err)

	s.Assert().Equal(0, len(response.Todos), "expected 0 todos")
}
