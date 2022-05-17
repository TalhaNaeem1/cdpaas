package pipeline_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	mock_airbyte "pipelineService/clients/airbyte/mocks"
	"pipelineService/clients/authService"
	mock_authservice "pipelineService/clients/authService/mocks"
	"pipelineService/env"
	"pipelineService/handlers/v1/test"
	"pipelineService/models/v1"
	mockStore "pipelineService/services/db/mocks"
	"pipelineService/utils"
)

//TestCreatePipeline tests all the scenarios while creating a pipeline.
func TestCreatePipeline(t *testing.T) {
	mockPipeline := createRandomPipeline()

	testCaseSuite := []struct {
		testScenario  string
		body          models.Pipeline
		productID     string
		buildStubs    func(store *mockStore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "BadRequest_BadModel",

			body: models.Pipeline{
				PipelineID:         uuid.UUID{},
				Name:               "",
				PipelineGovernance: mockPipeline.PipelineGovernance,
				CreatedAt:          mockPipeline.CreatedAt,
			},

			buildStubs: func(store *mockStore.MockStore) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "InternalServerError_FailedCreation",

			body: mockPipeline,

			buildStubs: func(store *mockStore.MockStore) {
				arg1 := mockPipeline
				store.EXPECT().CreatePipeline(arg1).Times(1).Return(models.Pipeline{}, sql.ErrConnDone)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: "Something went wrong",
					Data:   nil,
				}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "Success",

			body: mockPipeline,

			buildStubs: func(store *mockStore.MockStore) {
				arg1 := mockPipeline
				store.EXPECT().CreatePipeline(arg1).Times(1).Return(mockPipeline, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   mockPipeline}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
	}

	for i := range testCaseSuite {
		testCase := testCaseSuite[i]

		t.Run(testCase.testScenario, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockStore.NewMockStore(ctrl)
			testCase.buildStubs(store)

			body, e := json.Marshal(testCase.body)
			require.NoError(t, e)

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.PIPELINE, store, nil, authServiceClient)
			url := fmt.Sprintf("%spipelines/", test.BaseURL)
			expectedResp, err := test.MakeHttpRequest(server, http.MethodPost, url, nil, body)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

//TestUpdatePipeline tests all the scenarios while updating a pipeline.
func TestUpdatePipeline(t *testing.T) {
	mockPipeline := createRandomPipeline()

	mockUpdatePipeline := models.UpdatePipeline{
		PipelineID:         mockPipeline.PipelineID,
		Name:               mockPipeline.Name,
		PipelineGovernance: mockPipeline.PipelineGovernance,
	}

	testCaseSuite := []struct {
		testScenario  string
		body          models.UpdatePipeline
		pipelineID    string
		buildStubs    func(store *mockStore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "BadRequest_BadModel",

			pipelineID: mockPipeline.PipelineID.String(),

			body: models.UpdatePipeline{
				PipelineID:         uuid.UUID{},
				Name:               "",
				PipelineGovernance: nil,
			},

			buildStubs: func(store *mockStore.MockStore) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "InternalServerError_FailedCreation",

			pipelineID: mockPipeline.PipelineID.String(),

			body: mockUpdatePipeline,

			buildStubs: func(store *mockStore.MockStore) {
				arg1 := mockUpdatePipeline
				store.EXPECT().UpdatePipeline(arg1).Times(1).Return(models.Pipeline{}, sql.ErrConnDone)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: "Something went wrong",
					Data:   nil,
				}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "Success",

			pipelineID: mockPipeline.PipelineID.String(),

			body: mockUpdatePipeline,

			buildStubs: func(store *mockStore.MockStore) {
				arg1 := mockUpdatePipeline
				store.EXPECT().UpdatePipeline(arg1).Times(1).Return(mockPipeline, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   mockPipeline}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
	}

	for i := range testCaseSuite {
		testCase := testCaseSuite[i]

		t.Run(testCase.testScenario, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockStore.NewMockStore(ctrl)
			testCase.buildStubs(store)

			body, e := json.Marshal(testCase.body)
			require.NoError(t, e)

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.PIPELINE, store, nil, authServiceClient)
			url := fmt.Sprintf("%spipelines/%s/", test.BaseURL, testCase.pipelineID)
			expectedResp, err := test.MakeHttpRequest(server, http.MethodPut, url, nil, body)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

//TestCreatePipelineConnectionOnAirByte tests all the scenarios while creating a pipeline connection on AirByte.
func TestCreatePipelineConnectionOnAirByte(t *testing.T) {
	mockCreatePipelineReq := createRandomPipelineReq()
	mockWorkSpaceID := test.AirByteWorkspaceID
	mockConnectionID, _ := uuid.NewV1()
	mockAirByteSourceID, _ := uuid.NewV1()
	mockAirByteDestID, _ := uuid.NewV1()
	testCaseSuite := []struct {
		testScenario  string
		body          models.CreatePipelineRequest
		queryAirByte  func(querier *mock_airbyte.MockAirByteQuerier)
		buildStubs    func(store *mockStore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "BadRequest_BadModel",

			body: models.CreatePipelineRequest{
				SourceID:      "",
				DestinationID: "",
				Schedule:      mockCreatePipelineReq.Schedule,
				SyncCatalog:   mockCreatePipelineReq.SyncCatalog,
			},

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {},

			buildStubs: func(store *mockStore.MockStore) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "InternalServerError_BadAirByteInfo",

			body: mockCreatePipelineReq,

			buildStubs: func(store *mockStore.MockStore) {
				arg0 := mockCreatePipelineReq.SourceID
				arg1 := mockCreatePipelineReq.DestinationID
				store.EXPECT().GetSourceAndDestinationAirbyteInfo(arg0, arg1).Times(1).Return(models.AirbyteSourceAndDestinations{}, sql.ErrNoRows)
			},

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: "Something went wrong",
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "InternalServerError_BadAirByteCall",

			body: mockCreatePipelineReq,

			buildStubs: func(store *mockStore.MockStore) {
				arg0 := mockCreatePipelineReq.SourceID
				arg1 := mockCreatePipelineReq.DestinationID
				store.EXPECT().GetSourceAndDestinationAirbyteInfo(arg0, arg1).Times(1).
					Return(models.AirbyteSourceAndDestinations{
						ConnectionID:         mockConnectionID.String(),
						SourceID:             mockCreatePipelineReq.SourceID,
						DestinationID:        mockCreatePipelineReq.DestinationID,
						AirbyteSourceID:      mockAirByteSourceID.String(),
						AirbyteDestinationID: mockAirByteDestID.String(),
					}, nil)
			},

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {
				mockCreatePipelineReq.Operations[0].WorkspaceId = mockWorkSpaceID

				arg := models.CreatePipelineAirbyteRequest{
					DestinationId:       mockAirByteDestID.String(),
					SourceId:            mockAirByteSourceID.String(),
					NamespaceDefinition: utils.AIRBYTE_DEFAULT_NAMESPACE_DEFINITION,
					NamespaceFormat:     utils.AIRBYTE_DEFAULT_NAMESPACE_FORMAT,
					Prefix:              *mockCreatePipelineReq.Prefix,
					Status:              utils.AIRBYTE_DEFAULT_STATUS,
					Schedule:            mockCreatePipelineReq.Schedule,
					SyncCatalog:         mockCreatePipelineReq.SyncCatalog,
					Operations:          mockCreatePipelineReq.Operations,
				}

				querier.EXPECT().CreateConnection(arg).Times(1).Return(models.CreatePipelineAirbyteResponse{}, errors.New("bad airByte response"))
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: errors.New("bad airByte response").Error(),
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "InternalServerError_BadUpdateConnInfo",

			body: mockCreatePipelineReq,

			buildStubs: func(store *mockStore.MockStore) {
				arg0 := mockCreatePipelineReq.SourceID
				arg1 := mockCreatePipelineReq.DestinationID
				store.EXPECT().GetSourceAndDestinationAirbyteInfo(arg0, arg1).Times(1).
					Return(models.AirbyteSourceAndDestinations{
						ConnectionID:         mockConnectionID.String(),
						SourceID:             mockCreatePipelineReq.SourceID,
						DestinationID:        mockCreatePipelineReq.DestinationID,
						AirbyteSourceID:      mockAirByteSourceID.String(),
						AirbyteDestinationID: mockAirByteDestID.String(),
					}, nil)

				arg2 := models.Connection{
					ConnectionID:          mockConnectionID.String(),
					AirbyteConnectionID:   mockConnectionID.String(),
					AirbyteStatus:         "UP",
					AirbyteFrequencyUnits: mockCreatePipelineReq.Schedule.Units,
					AirbyteTimeUnit:       mockCreatePipelineReq.Schedule.TimeUnit,
					Owner:                 1122,
					WorkspaceID:           1122,
				}
				arg3 := mockCreatePipelineReq.DestinationID
				store.EXPECT().UpdateConnectionInfo(arg2, arg3).Times(1).
					Return(sql.ErrConnDone)
			},

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {
				mockCreatePipelineReq.Operations[0].WorkspaceId = mockWorkSpaceID

				arg := models.CreatePipelineAirbyteRequest{
					DestinationId:       mockAirByteDestID.String(),
					SourceId:            mockAirByteSourceID.String(),
					NamespaceDefinition: utils.AIRBYTE_DEFAULT_NAMESPACE_DEFINITION,
					NamespaceFormat:     utils.AIRBYTE_DEFAULT_NAMESPACE_FORMAT,
					Prefix:              *mockCreatePipelineReq.Prefix,
					Status:              utils.AIRBYTE_DEFAULT_STATUS,
					Schedule:            mockCreatePipelineReq.Schedule,
					SyncCatalog:         mockCreatePipelineReq.SyncCatalog,
					Operations:          mockCreatePipelineReq.Operations,
				}

				querier.EXPECT().CreateConnection(arg).Times(1).
					Return(models.CreatePipelineAirbyteResponse{
						ConnectionId: mockConnectionID.String(),
						Status:       "UP",
						Schedule:     mockCreatePipelineReq.Schedule}, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: "Something went wrong",
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "Success",

			body: mockCreatePipelineReq,

			buildStubs: func(store *mockStore.MockStore) {
				arg0 := mockCreatePipelineReq.SourceID
				arg1 := mockCreatePipelineReq.DestinationID
				store.EXPECT().GetSourceAndDestinationAirbyteInfo(arg0, arg1).Times(1).
					Return(models.AirbyteSourceAndDestinations{
						ConnectionID:         mockConnectionID.String(),
						SourceID:             mockCreatePipelineReq.SourceID,
						DestinationID:        mockCreatePipelineReq.DestinationID,
						AirbyteSourceID:      mockAirByteSourceID.String(),
						AirbyteDestinationID: mockAirByteDestID.String(),
					}, nil)

				arg2 := models.Connection{
					ConnectionID:          mockConnectionID.String(),
					AirbyteConnectionID:   mockConnectionID.String(),
					AirbyteStatus:         "UP",
					AirbyteFrequencyUnits: mockCreatePipelineReq.Schedule.Units,
					AirbyteTimeUnit:       mockCreatePipelineReq.Schedule.TimeUnit,
					Owner:                 1122,
					WorkspaceID:           1122,
				}
				arg3 := mockCreatePipelineReq.DestinationID
				store.EXPECT().UpdateConnectionInfo(arg2, arg3).Times(1).
					Return(nil)
			},

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {
				mockCreatePipelineReq.Operations[0].WorkspaceId = mockWorkSpaceID

				arg := models.CreatePipelineAirbyteRequest{
					DestinationId:       mockAirByteDestID.String(),
					SourceId:            mockAirByteSourceID.String(),
					NamespaceDefinition: utils.AIRBYTE_DEFAULT_NAMESPACE_DEFINITION,
					NamespaceFormat:     utils.AIRBYTE_DEFAULT_NAMESPACE_FORMAT,
					Prefix:              *mockCreatePipelineReq.Prefix,
					Status:              utils.AIRBYTE_DEFAULT_STATUS,
					Schedule:            mockCreatePipelineReq.Schedule,
					SyncCatalog:         mockCreatePipelineReq.SyncCatalog,
					Operations:          mockCreatePipelineReq.Operations,
				}

				querier.EXPECT().CreateConnection(arg).Times(1).
					Return(models.CreatePipelineAirbyteResponse{
						ConnectionId: mockConnectionID.String(),
						Status:       "UP",
						Schedule:     mockCreatePipelineReq.Schedule}, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   "pipeline created successfully"}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
	}

	for i := range testCaseSuite {
		testCase := testCaseSuite[i]

		t.Run(testCase.testScenario, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockStore.NewMockStore(ctrl)
			testCase.buildStubs(store)

			querier := mock_airbyte.NewMockAirByteQuerier(ctrl)
			testCase.queryAirByte(querier)

			body, e := json.Marshal(testCase.body)
			require.NoError(t, e)

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.PIPELINE, store, querier, authServiceClient)
			url := fmt.Sprintf("%spipelines/connections/", test.BaseURL)
			expectedResp, err := test.MakeHttpRequest(server, http.MethodPost, url, nil, body)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

//TestUpdatePipelineConnectionOnAirByte tests all the scenarios while updating a pipeline connection on AirByte.
func TestUpdatePipelineConnectionOnAirByte(t *testing.T) {
	mockCreatePipelineReq := createRandomPipelineReq()
	mockConnectionID, _ := uuid.NewV1()
	mockConnection := createRandomConnection(mockConnectionID.String())

	mockUpdatePipelineReq := models.UpdatePipelineAirByteRequest{
		ConnectionId: mockConnection.AirbyteConnectionID,
		Prefix:       mockCreatePipelineReq.Prefix,
		SyncCatalog:  mockCreatePipelineReq.SyncCatalog,
		Schedule:     mockCreatePipelineReq.Schedule,
		Status:       "active",
		Operations:   mockCreatePipelineReq.Operations,
	}

	testCaseSuite := []struct {
		testScenario  string
		body          models.UpdatePipelineAirByteRequest
		connectionID  string
		queryAirByte  func(querier *mock_airbyte.MockAirByteQuerier)
		buildStubs    func(store *mockStore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "InternalServerError",

			connectionID: mockConnectionID.String(),

			body: mockUpdatePipelineReq,

			buildStubs: func(store *mockStore.MockStore) {
				store.EXPECT().GetConnection(mockConnectionID.String()).Times(1).Return(models.Connection{}, sql.ErrConnDone)
			},

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: "Something went wrong",
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "InternalServerError_BadAirByte",

			connectionID: mockConnectionID.String(),

			body: mockUpdatePipelineReq,

			buildStubs: func(store *mockStore.MockStore) {
				store.EXPECT().GetConnection(mockConnectionID.String()).Times(1).Return(mockConnection, nil)
			},

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {
				querier.EXPECT().UpdateConnection(mockUpdatePipelineReq).Times(1).
					Return(models.CreatePipelineAirbyteResponse{}, errors.New("Bad AirByte response"))
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: "Bad AirByte response",
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "InternalServerError",

			connectionID: mockConnectionID.String(),

			body: mockUpdatePipelineReq,

			buildStubs: func(store *mockStore.MockStore) {
				store.EXPECT().GetConnection(mockConnectionID.String()).Times(1).Return(mockConnection, nil)

				arg := models.Connection{
					ConnectionID:          mockConnectionID.String(),
					AirbyteFrequencyUnits: mockUpdatePipelineReq.Schedule.Units,
					AirbyteTimeUnit:       mockUpdatePipelineReq.Schedule.TimeUnit,
				}
				store.EXPECT().UpdateConnectionSchedule(arg).Times(1).Return(sql.ErrConnDone)
			},

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {
				querier.EXPECT().UpdateConnection(mockUpdatePipelineReq).Times(1).
					Return(models.CreatePipelineAirbyteResponse{
						Schedule: mockUpdatePipelineReq.Schedule,
					}, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: "Something went wrong",
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "Success",

			connectionID: mockConnectionID.String(),

			body: mockUpdatePipelineReq,

			buildStubs: func(store *mockStore.MockStore) {
				store.EXPECT().GetConnection(mockConnectionID.String()).Times(1).Return(mockConnection, nil)

				arg := models.Connection{
					ConnectionID:          mockConnectionID.String(),
					AirbyteFrequencyUnits: mockUpdatePipelineReq.Schedule.Units,
					AirbyteTimeUnit:       mockUpdatePipelineReq.Schedule.TimeUnit,
				}
				store.EXPECT().UpdateConnectionSchedule(arg).Times(1).Return(nil)
			},

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {
				querier.EXPECT().UpdateConnection(mockUpdatePipelineReq).Times(1).
					Return(models.CreatePipelineAirbyteResponse{Schedule: mockUpdatePipelineReq.Schedule}, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   "pipeline updated successfully"}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
	}

	for i := range testCaseSuite {
		testCase := testCaseSuite[i]

		t.Run(testCase.testScenario, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockStore.NewMockStore(ctrl)
			testCase.buildStubs(store)

			querier := mock_airbyte.NewMockAirByteQuerier(ctrl)
			testCase.queryAirByte(querier)

			body, e := json.Marshal(testCase.body)
			require.NoError(t, e)

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.PIPELINE, store, querier, authServiceClient)
			url := fmt.Sprintf("%spipelines/connections/%s/", test.BaseURL, testCase.connectionID)
			expectedResp, err := test.MakeHttpRequest(server, http.MethodPut, url, nil, body)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

func TestRunManualSyncOnAirByte(t *testing.T) {
	mockManualConnectionSyncResponse := createRandomManualConnectionSyncResponse()
	cid, _ := uuid.NewV1()
	mockConnectionID := cid.String()

	testCaseSuite := []struct {
		testScenario  string
		productID     string
		connectionID  string
		queryAirByte  func(querier *mock_airbyte.MockAirByteQuerier)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "BadRequest",

			connectionID: "",

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {
				arg := make(map[string]interface{})
				arg["connectionId"] = ""

				querier.EXPECT().SyncConnectionManually(arg).Times(1).Return(models.ManualConnectionSyncResponse{}, errors.New("bad connection ID"))
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: errors.New("bad connection ID").Error(),
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "Success",

			connectionID: mockConnectionID,

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {
				arg := make(map[string]interface{})
				arg["connectionId"] = mockConnectionID

				querier.EXPECT().SyncConnectionManually(arg).Times(1).
					Return(mockManualConnectionSyncResponse, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   mockManualConnectionSyncResponse}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
	}

	for i := range testCaseSuite {
		testCase := testCaseSuite[i]

		t.Run(testCase.testScenario, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			querier := mock_airbyte.NewMockAirByteQuerier(ctrl)
			testCase.queryAirByte(querier)

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.PIPELINE, nil, querier, authServiceClient)
			url := fmt.Sprintf("%spipelines/connections/%s/sync/", test.BaseURL, testCase.connectionID)
			expectedResp, err := test.MakeHttpRequest(server, http.MethodPost, url, nil, nil)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// TestGetAllPipelines tests all the scenarios while getting all the pipelines of the specific product.
func TestGetAllPipelines(t *testing.T) {
	mockPipeline := createRandomPipeline()
	mockAirByteConnectID0, _ := uuid.NewV1()
	mockAirByteConnectID1, _ := uuid.NewV1()
	mockConnectionMeta := createRandomConnectionMeta()

	mockPipelineMetaData := []models.PipelinesMetaData{createRandomPipelineMeta(mockAirByteConnectID0.String()), createRandomPipelineMeta(mockAirByteConnectID1.String())}

	testCaseSuite := []struct {
		testScenario   string
		getUserDetails func(client *mock_authservice.MockHttpClient)
		buildStubs     func(store *mockStore.MockStore)
		queryAirByte   func(querier *mock_airbyte.MockAirByteQuerier)
		checkResponse  func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "InternalServerError",

			getUserDetails: func(client *mock_authservice.MockHttpClient) {
			},

			buildStubs: func(store *mockStore.MockStore) {
				arg2 := mockPipeline.WorkspaceID
				store.EXPECT().GetAllPipelines(arg2).Times(1).Return([]models.PipelinesMetaData{}, sql.ErrConnDone)
			},

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: "Something went wrong",
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "Success",

			getUserDetails: func(client *mock_authservice.MockHttpClient) {
				test.MockGetUserByID(client, 0, 1122)
				test.MockGetUserByID(client, 0, 1122)
			},

			buildStubs: func(store *mockStore.MockStore) {
				arg2 := mockPipeline.WorkspaceID
				store.EXPECT().GetAllPipelines(arg2).Times(1).Return(mockPipelineMetaData, nil)
			},

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {
				arg := make(map[string]interface{})
				arg["connectionId"] = mockAirByteConnectID0
				arg["withRefreshedCatalog"] = false
				querier.EXPECT().GetConnectionDetails(arg).Times(1).Return(mockConnectionMeta, nil)

				arg1 := make(map[string]interface{})
				arg1["connectionId"] = mockAirByteConnectID1
				arg1["withRefreshedCatalog"] = false
				querier.EXPECT().GetConnectionDetails(arg1).Times(1).Return(mockConnectionMeta, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   mockPipelineMetaData}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
	}

	for i := range testCaseSuite {
		testCase := testCaseSuite[i]

		t.Run(testCase.testScenario, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockStore.NewMockStore(ctrl)
			testCase.buildStubs(store)

			airbyteQuerier := mock_airbyte.NewMockAirByteQuerier(ctrl)
			testCase.queryAirByte(airbyteQuerier)

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)
			testCase.getUserDetails(httpMockClient)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.PIPELINE, store, airbyteQuerier, authServiceClient)
			url := fmt.Sprintf("%spipelines/", test.BaseURL)
			expectedResp, err := test.MakeHttpRequest(server, http.MethodGet, url, nil, nil)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// TestGetPipeline tests all the scenarios while getting the pipeline of the specific product.
func TestGetPipeline(t *testing.T) {
	abConnectionID, _ := uuid.NewV1()
	mockPipelineView := createRandomPipelineView(abConnectionID.String())
	mockUser := test.CreateRandomUserDetails(1, 1122)
	mockConnectionMeta := createRandomConnectionMeta()

	testCaseSuite := []struct {
		testScenario   string
		productID      string
		pipelineID     string
		getUserDetails func(client *mock_authservice.MockHttpClient)
		buildStubs     func(store *mockStore.MockStore)
		queryAirByte   func(querier *mock_airbyte.MockAirByteQuerier)
		checkResponse  func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "BadPipelineID",

			getUserDetails: func(client *mock_authservice.MockHttpClient) {},

			productID: "BadUUID",

			pipelineID: "BadPipelineUUID",

			buildStubs: func(store *mockStore.MockStore) {},

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "FailedGetPipelineReq",

			getUserDetails: func(client *mock_authservice.MockHttpClient) {},

			productID: mockPipelineView.ProductID.String(),

			pipelineID: mockPipelineView.PipelineID.String(),

			buildStubs: func(store *mockStore.MockStore) {
				arg := mockPipelineView.PipelineID
				store.EXPECT().GetPipeline(arg).Times(1).Return(models.PipelineView{}, sql.ErrConnDone)
			},

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: "Something went wrong",
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "FailedGetUserByID",

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {
				requestBody := make(map[string]interface{})
				requestBody["connectionId"] = abConnectionID
				requestBody["withRefreshedCatalog"] = false

				querier.EXPECT().GetConnectionDetails(requestBody).Times(1).Return(mockConnectionMeta, nil)
			},

			getUserDetails: func(client *mock_authservice.MockHttpClient) {
				ownerID := mockPipelineView.Owner
				arg := fmt.Sprintf("%s/auth-service/api/v1/accounts/internal/user-from-id?user_id=%d", env.Env.AuthServiceAddress, ownerID)
				res := http.Response{}
				client.EXPECT().Get(arg).Times(1).Return(&res, errors.New("User Not Found"))
			},

			productID: mockPipelineView.ProductID.String(),

			pipelineID: mockPipelineView.PipelineID.String(),

			buildStubs: func(store *mockStore.MockStore) {
				arg := mockPipelineView.PipelineID
				store.EXPECT().GetPipeline(arg).Times(1).Return(mockPipelineView, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: "User Not Found",
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "Success",

			getUserDetails: func(client *mock_authservice.MockHttpClient) {
				test.MockGetUserByID(client, 1, 1122)
			},

			productID: mockPipelineView.ProductID.String(),

			pipelineID: mockPipelineView.PipelineID.String(),

			buildStubs: func(store *mockStore.MockStore) {
				arg := mockPipelineView.PipelineID
				store.EXPECT().GetPipeline(arg).Times(1).Return(mockPipelineView, nil)
			},

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {
				requestBody := make(map[string]interface{})
				requestBody["connectionId"] = abConnectionID
				requestBody["withRefreshedCatalog"] = false

				querier.EXPECT().GetConnectionDetails(requestBody).Times(1).Return(mockConnectionMeta, nil)
				mockPipelineView.AirbyteStatus = mockConnectionMeta.LatestSyncJobStatus
				mockPipelineView.AirbyteLastRun = mockConnectionMeta.LatestSyncJobCreatedAt
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data: models.GetPipelineDetails{
						Pipeline: mockPipelineView,
						Owner:    mockUser.Payload.UserInfo,
					}}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
	}
	for i := range testCaseSuite {
		testCase := testCaseSuite[i]

		t.Run(testCase.testScenario, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockStore.NewMockStore(ctrl)
			testCase.buildStubs(store)

			airbyteQuerier := mock_airbyte.NewMockAirByteQuerier(ctrl)
			testCase.queryAirByte(airbyteQuerier)

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)
			testCase.getUserDetails(httpMockClient)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.PIPELINE, store, airbyteQuerier, authServiceClient)
			url := fmt.Sprintf("%spipelines/%s/", test.BaseURL, testCase.pipelineID)
			expectedResp, err := test.MakeHttpRequest(server, http.MethodGet, url, nil, nil)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// TestGetSourceSchemaFromAirByteConnection tests all the scenarios while getting the Schema of the specific connection.
func TestGetSourceSchemaFromAirByteConnection(t *testing.T) {
	cid, _ := uuid.NewV1()
	mockConnectionID := cid.String()

	testCaseSuite := []struct {
		testScenario  string
		connectionID  string
		queryAirByte  func(querier *mock_airbyte.MockAirByteQuerier)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "BadRequest",

			connectionID: mockConnectionID,

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {
				arg := mockConnectionID
				querier.EXPECT().GetConnectionSchema(arg).Times(1).
					Return(models.ConnectionSourceSchema{}, errors.New("bad connection ID"))
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: errors.New("bad connection ID").Error(),
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "Success",

			connectionID: mockConnectionID,

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {
				arg := mockConnectionID
				querier.EXPECT().GetConnectionSchema(arg).Times(1).
					Return(models.ConnectionSourceSchema{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   models.ConnectionSourceSchema{}}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
	}

	for i := range testCaseSuite {
		testCase := testCaseSuite[i]

		t.Run(testCase.testScenario, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			querier := mock_airbyte.NewMockAirByteQuerier(ctrl)
			testCase.queryAirByte(querier)

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.PIPELINE, nil, querier, authServiceClient)
			url := fmt.Sprintf("%spipelines/connections/%s/schema/", test.BaseURL, testCase.connectionID)
			expectedResp, err := test.MakeHttpRequest(server, http.MethodGet, url, nil, nil)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// TestFetchSyncHistoryFromAirByte tests all the scenarios while getting the Sync history of the specific connection.
func TestFetchSyncHistoryFromAirByte(t *testing.T) {
	cid, _ := uuid.NewV1()
	mockConnectionID := cid.String()

	testCaseSuite := []struct {
		testScenario  string
		connectionID  string
		queryAirByte  func(querier *mock_airbyte.MockAirByteQuerier)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "BadRequest",

			connectionID: mockConnectionID,

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {
				arg := models.SyncHistoryRequest{
					ConfigTypes: []string{
						utils.SYNC,
						utils.RESET_CONNECTION,
					},
					ConfigId: mockConnectionID,
				}

				querier.EXPECT().FetchSyncHistory(arg).Times(1).
					Return(models.SyncHistoryResponse{}, errors.New("bad connection ID"))
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: errors.New("bad connection ID").Error(),
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "Success",

			connectionID: mockConnectionID,

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {
				arg := models.SyncHistoryRequest{
					ConfigTypes: []string{
						utils.SYNC,
						utils.RESET_CONNECTION,
					},
					ConfigId: mockConnectionID,
				}

				querier.EXPECT().FetchSyncHistory(arg).Times(1).
					Return(models.SyncHistoryResponse{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   models.SyncHistoryResponse{}}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
	}

	for i := range testCaseSuite {
		testCase := testCaseSuite[i]

		t.Run(testCase.testScenario, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			querier := mock_airbyte.NewMockAirByteQuerier(ctrl)
			testCase.queryAirByte(querier)

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.PIPELINE, nil, querier, authServiceClient)
			url := fmt.Sprintf("%spipelines/connections/%s/sync/history/", test.BaseURL, testCase.connectionID)
			expectedResp, err := test.MakeHttpRequest(server, http.MethodGet, url, nil, nil)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// TestGetJobLogsFromAirByte tests all the scenarios while getting the job logs of the specific connection job.
func TestGetJobLogsFromAirByte(t *testing.T) {
	mockJobID := strconv.Itoa(int(utils.RandomInt(1, 10)))

	testCaseSuite := []struct {
		testScenario  string
		jobID         string
		queryAirByte  func(querier *mock_airbyte.MockAirByteQuerier)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "BadRequest",

			jobID: mockJobID,

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {
				arg, _ := strconv.Atoi(mockJobID)
				querier.EXPECT().GetJobLogs(arg).Times(1).
					Return(models.JobLogs{}, errors.New("bad job ID"))
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: errors.New("bad job ID").Error(),
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "Success",

			jobID: mockJobID,

			queryAirByte: func(querier *mock_airbyte.MockAirByteQuerier) {
				arg, _ := strconv.Atoi(mockJobID)
				querier.EXPECT().GetJobLogs(arg).Times(1).
					Return(models.JobLogs{}, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   models.JobLogs{}}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
	}

	for i := range testCaseSuite {
		testCase := testCaseSuite[i]

		t.Run(testCase.testScenario, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			querier := mock_airbyte.NewMockAirByteQuerier(ctrl)
			testCase.queryAirByte(querier)

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.PIPELINE, nil, querier, authServiceClient)
			url := fmt.Sprintf("%spipelines/connections/sync/logs/%s/", test.BaseURL, testCase.jobID)
			expectedResp, err := test.MakeHttpRequest(server, http.MethodGet, url, nil, nil)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// createRandomManualConnectionSyncResponse populates and returns the ManualConnectionSyncResponse Model with random values.
func createRandomManualConnectionSyncResponse() models.ManualConnectionSyncResponse {
	job := models.Job{
		ID:         int(utils.RandomInt(1, 10)),
		ConfigType: utils.RandomString(5),
		ConfigID:   utils.RandomString(5),
		Status:     utils.RandomString(5),
	}

	r := models.ManualConnectionSyncResponse{
		Job:      job,
		Attempts: nil,
	}

	return r
}

//createRandomPipelineReq populates and return the CreatePipelineRequest Model with random values.
func createRandomPipelineReq() models.CreatePipelineRequest {
	sID, _ := uuid.NewV1()
	dID, _ := uuid.NewV1()

	stream := models.Stream{
		Name:                    utils.RandomString(5),
		JsonSchema:              nil,
		SupportedSyncModes:      []string{utils.RandomString(5), utils.RandomString(5)},
		SourceDefinedCursor:     false,
		DefaultCursorField:      []string{utils.RandomString(5), utils.RandomString(5)},
		SourceDefinedPrimaryKey: [][]string{{utils.RandomString(5), utils.RandomString(5)}, {utils.RandomString(5), utils.RandomString(5)}},
		Namespace:               utils.RandomString(5),
	}
	streams := []models.Streams{
		{Stream: stream,
			Config: models.Config{
				SyncMode:            utils.RandomString(5),
				CursorField:         []string{utils.RandomString(5), utils.RandomString(5)},
				DestinationSyncMode: utils.RandomString(5),
				PrimaryKey:          [][]string{{utils.RandomString(5), utils.RandomString(5)}, {utils.RandomString(5), utils.RandomString(5)}},
				AliasName:           utils.RandomString(5),
				Selected:            false,
			}},
	}

	normalization := &models.Normalization{
		Option: utils.RandomString(5),
	}

	dbt := &models.Dbt{
		GitRepoUrl:    utils.RandomString(5),
		GitRepoBranch: utils.RandomString(5),
		DockerImage:   utils.RandomString(5),
		DbtArguments:  utils.RandomString(5),
	}

	operatorConfiguration := models.OperatorConfiguration{
		OperatorType:  utils.RandomString(5),
		Normalization: normalization,
		Dbt:           dbt,
	}

	operation := models.Operations{
		WorkspaceId:           test.AirByteWorkspaceID,
		Name:                  utils.RandomString(5),
		OperatorConfiguration: operatorConfiguration,
	}

	p := utils.RandomString(5)
	pr := models.CreatePipelineRequest{
		SourceID:      sID.String(),
		DestinationID: dID.String(),
		Schedule: &models.Schedule{
			Units:    1,
			TimeUnit: "minutes",
		},
		SyncCatalog: models.SyncCatalog{Streams: streams},
		Operations:  []models.Operations{operation},
		Prefix:      &p,
	}

	return pr
}

//createRandomConnectionMeta populates and return the ConnectionMeta Model with random values.
func createRandomConnectionMeta() models.ConnectionMeta {
	cm := models.ConnectionMeta{
		LatestSyncJobCreatedAt: int(utils.RandomInt(1, 10)),
		LatestSyncJobStatus:    utils.RandomString(5),
	}

	return cm
}

//createRandomConnection populates and return the Connection Model with random values.
func createRandomConnection(connectionID string) models.Connection {
	aBConnID, _ := uuid.NewV1()
	cm := models.Connection{
		ConnectionID:        connectionID,
		AirbyteConnectionID: aBConnID.String(),
	}

	return cm
}

//createRandomPipelineMeta populates and return the PipelineMeta Model with random values.
func createRandomPipelineMeta(AbConnecID string) models.PipelinesMetaData {
	pm := models.PipelinesMetaData{
		PipelineName:        utils.RandomString(5),
		SourceName:          utils.RandomString(5),
		DestinationName:     utils.RandomString(5),
		AirbyteStatus:       "",
		AirbyteLastRun:      0,
		AirbyteConnectionID: AbConnecID,
		Owner:               "{user:admin}",
	}

	return pm
}

//createRandomPipelineView populates and return the PipelineView Model with random values.
func createRandomPipelineView(abConnID string) models.PipelineView {
	pID, _ := uuid.NewV1()

	p := models.PipelineView{
		PipelineID:         pID,
		Name:               utils.RandomString(5),
		PipelineGovernance: []string{utils.RandomString(5), utils.RandomString(5)},
		CreatedAt:          1647606617639,
		Product: []map[string]interface{}{
			{
				"Name": utils.RandomString(5),
			},
			{"Name": utils.RandomString(5)},
		},
		SourceName:          utils.RandomString(5),
		DestinationName:     utils.RandomString(5),
		AirbyteStatus:       "",
		AirbyteLastRun:      0,
		AirbyteConnectionID: abConnID,
		Owner:               1,
	}

	return p
}

//createRandomPipeline populates and return the Pipeline Model with random values.
func createRandomPipeline() models.Pipeline {
	pID, _ := uuid.NewV1()
	p := models.Pipeline{
		PipelineID:         pID,
		Name:               utils.RandomString(5),
		PipelineGovernance: []string{utils.RandomString(5), utils.RandomString(5)},
		CreatedAt:          1647606617639,
		Owner:              1122,
		WorkspaceID:        1122,
	}

	return p
}

// createRandomDataProduct populates and return the DataProduct Model with random values.
//func createRandomDataProduct() models.DataProduct {
//	pID, _ := uuid.NewV1()
//	Dp := models.DataProduct{
//		ProductID: pID,
//		Name:      utils.RandomString(5),
//		DataProductGovernance: []string{
//			utils.RandomString(5), utils.RandomString(5),
//		},
//		DataDomain:  utils.RandomString(5),
//		Description: utils.RandomString(10),
//		DataProductStatus: "completed",
//		LastUpdated: 1647606617639,
//		Owner:       1122,
//		WorkspaceID: 1122,
//	}
//
//	return Dp
//}

// TestMain runs the package level test in TestMode.
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
