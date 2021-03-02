package rest_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"

	"github.com/MarioCarrion/todo-api/internal"
	"github.com/MarioCarrion/todo-api/internal/rest"
	"github.com/MarioCarrion/todo-api/internal/rest/resttesting"
)

func TestTasks_Post(t *testing.T) {
	// XXX: Test "serviceArgs"

	t.Parallel()

	type output struct {
		expectedStatus int
		expected       interface{}
		target         interface{}
	}

	tests := []struct {
		name   string
		setup  func(*resttesting.FakeTaskService)
		input  []byte
		output output
	}{
		{
			"OK: 201",
			func(s *resttesting.FakeTaskService) {
				s.CreateReturns(
					internal.Task{
						ID:          "1-2-3",
						Description: "new task",
					},
					nil)
			},
			func() []byte {
				b, _ := json.Marshal(&rest.CreateTasksRequest{
					Description: "new task",
				})

				return b
			}(),
			output{
				http.StatusCreated,
				&rest.CreateTasksResponse{
					Task: rest.Task{
						ID:          "1-2-3",
						Description: "new task",
					},
				},
				&rest.CreateTasksResponse{},
			},
		},
		{
			"ERR: 400",
			func(*resttesting.FakeTaskService) {},
			[]byte(`{"invalid":"json`),
			output{
				http.StatusBadRequest,
				&rest.ErrorResponse{
					Error: "invalid request",
				},
				&rest.ErrorResponse{},
			},
		},
		{
			"ERR: 500",
			func(s *resttesting.FakeTaskService) {
				s.CreateReturns(internal.Task{},
					errors.New("service error"))
			},
			[]byte(`{}`),
			output{
				http.StatusInternalServerError,
				&rest.ErrorResponse{
					Error: "create failed",
				},
				&rest.ErrorResponse{},
			},
		},
	}

	//-

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router := mux.NewRouter()
			svc := &resttesting.FakeTaskService{}
			tt.setup(svc)

			rest.NewTaskHandler(svc).Register(router)

			//-

			res := doRequest(router,
				httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(tt.input)))

			//-

			assertResponse(t, res, test{tt.output.expected, tt.output.target})

			if tt.output.expectedStatus != res.StatusCode {
				t.Fatalf("expected code %d, actual %d", tt.output.expectedStatus, res.StatusCode)
			}
		})
	}
}

func TestTasks_Read(t *testing.T) {
	// XXX: Test "serviceArgs"

	t.Parallel()

	type output struct {
		expectedStatus int
		expected       interface{}
		target         interface{}
	}

	tests := []struct {
		name   string
		setup  func(*resttesting.FakeTaskService)
		output output
	}{
		{
			"OK: 200",
			func(s *resttesting.FakeTaskService) {
				s.TaskReturns(
					internal.Task{
						ID:          "a-b-c",
						Description: "existing task",
					},
					nil)
			},
			output{
				http.StatusOK,
				&rest.ReadTasksResponse{
					Task: rest.Task{
						ID:          "a-b-c",
						Description: "existing task",
					},
				},
				&rest.ReadTasksResponse{},
			},
		},
		{
			"ERR: 500",
			func(s *resttesting.FakeTaskService) {
				s.TaskReturns(internal.Task{},
					errors.New("service error"))
			},
			output{
				http.StatusInternalServerError,
				&rest.ErrorResponse{
					Error: "find failed",
				},
				&rest.ErrorResponse{},
			},
		},
	}

	//-

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router := mux.NewRouter()
			svc := &resttesting.FakeTaskService{}
			tt.setup(svc)

			rest.NewTaskHandler(svc).Register(router)

			//-

			res := doRequest(router,
				httptest.NewRequest(http.MethodGet, "/tasks/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", nil))

			//-

			assertResponse(t, res, test{tt.output.expected, tt.output.target})

			if tt.output.expectedStatus != res.StatusCode {
				t.Fatalf("expected code %d, actual %d", tt.output.expectedStatus, res.StatusCode)
			}
		})
	}
}

func TestTasks_Update(t *testing.T) {
	// XXX: Test "serviceArgs"

	t.Parallel()

	type output struct {
		expectedStatus int
		expected       interface{}
		target         interface{}
	}

	tests := []struct {
		name   string
		setup  func(*resttesting.FakeTaskService)
		input  []byte
		output output
	}{
		{
			"OK: 200",
			func(s *resttesting.FakeTaskService) {},
			func() []byte {
				b, _ := json.Marshal(&rest.UpdateTasksRequest{
					Description: "update task",
				})

				return b
			}(),
			output{
				http.StatusOK,
				&struct{}{},
				&struct{}{},
			},
		},
		{
			"ERR: 400",
			func(*resttesting.FakeTaskService) {},
			[]byte(`{"invalid":"json`),
			output{
				http.StatusBadRequest,
				&rest.ErrorResponse{
					Error: "invalid request",
				},
				&rest.ErrorResponse{},
			},
		},
		{
			"ERR: 500",
			func(s *resttesting.FakeTaskService) {
				s.UpdateReturns(errors.New("service error"))
			},
			[]byte(`{}`),
			output{
				http.StatusInternalServerError,
				&rest.ErrorResponse{
					Error: "update failed",
				},
				&rest.ErrorResponse{},
			},
		},
	}

	//-

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router := mux.NewRouter()
			svc := &resttesting.FakeTaskService{}
			tt.setup(svc)

			rest.NewTaskHandler(svc).Register(router)

			//-

			res := doRequest(router,
				httptest.NewRequest(http.MethodPut, "/tasks/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", bytes.NewReader(tt.input)))

			//-

			assertResponse(t, res, test{tt.output.expected, tt.output.target})

			if tt.output.expectedStatus != res.StatusCode {
				t.Fatalf("expected code %d, actual %d", tt.output.expectedStatus, res.StatusCode)
			}
		})
	}
}

type test struct {
	expected interface{}
	target   interface{}
}

func doRequest(router *mux.Router, req *http.Request) *http.Response {
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	return rr.Result()
}

func assertResponse(t *testing.T, res *http.Response, test test) {
	t.Helper()

	if err := json.NewDecoder(res.Body).Decode(test.target); err != nil {
		t.Fatalf("couldn't decode %s", err)
	}
	defer res.Body.Close()

	if !cmp.Equal(test.expected, test.target) {
		t.Fatalf("expected results don't match: %s", cmp.Diff(test.expected, test.target))
	}
}