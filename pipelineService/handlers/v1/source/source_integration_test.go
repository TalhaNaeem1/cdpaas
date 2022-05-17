package source_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	mockairbyte "pipelineService/clients/airbyte/mocks"
	"pipelineService/clients/authService"
	mock_authservice "pipelineService/clients/authService/mocks"
	"pipelineService/env"
	"pipelineService/handlers/v1/test"
	"pipelineService/models/v1"
	mockStore "pipelineService/services/db/mocks"
	"pipelineService/utils"
)

// TestConfigureSourceOnAirByte tests all the scenarios while configuring the sources on airByte.
func TestConfigureSourceOnAirByte(t *testing.T) {
	mockSourceDefID, _ := uuid.NewV1()
	mockAirByteSourceID, _ := uuid.NewV1()
	mockSourceConnectorReq := createRandomSourceConnectorRequestAPI(mockSourceDefID.String())
	mockSourceConnectorRes := models.CreateSourceConnectorResponseAirbyte{
		AirbyteSourceId: mockAirByteSourceID.String(),
		SourceName:      mockSourceConnectorReq.Name,
		CreateSourceConnectorRequest: models.CreateSourceConnectorRequest{
			AirbyteSourceDefinitionId: mockSourceDefID.String(),
			ConnectionConfiguration:   nil,
			Name:                      utils.RandomString(5),
		},
	}

	testCaseSuite := []struct {
		testScenario  string
		body          models.CreateSourceConnectorRequestAPI
		queryAirByte  func(querier *mockairbyte.MockAirByteQuerier)
		buildStubs    func(store *mockStore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "BadRequest",

			body: models.CreateSourceConnectorRequestAPI{
				CreateSourceConnectorRequest: models.CreateSourceConnectorRequest{},
				Pipeline:                     "",
			},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {},

			buildStubs: func(store *mockStore.MockStore) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "InternalServerError",

			body: mockSourceConnectorReq,

			buildStubs: func(store *mockStore.MockStore) {},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				arg := map[string]interface{}{
					"sourceDefinitionId":      mockSourceDefID.String(),
					"connectionConfiguration": mockSourceConnectorReq.ConnectionConfiguration,
				}
				querier.EXPECT().CheckSourceConnection(arg).Times(1).Return(errors.New("Bad Source Connector"))
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: "Bad Source Connector",
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "InternalServerError",

			body: mockSourceConnectorReq,

			buildStubs: func(store *mockStore.MockStore) {},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				arg0 := map[string]interface{}{
					"sourceDefinitionId":      mockSourceDefID.String(),
					"connectionConfiguration": mockSourceConnectorReq.ConnectionConfiguration,
				}
				querier.EXPECT().CheckSourceConnection(arg0).Times(1).Return(nil)

				arg := models.CreateSourceConnectorRequestAirbyte{
					CreateSourceConnectorRequest: mockSourceConnectorReq.CreateSourceConnectorRequest,
					WorkspaceId:                  test.AirByteWorkspaceID,
				}

				querier.EXPECT().CreateSourceConnectorOnAirByte(arg).
					Times(1).Return(models.CreateSourceConnectorResponseAirbyte{}, errors.New("can't create source on AirByte"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: errors.New("can't create source on AirByte").Error(),
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "Internal Server Error",

			body: mockSourceConnectorReq,

			buildStubs: func(store *mockStore.MockStore) {
				arg0 := models.Source{
					SourceName:                mockSourceConnectorRes.SourceName,
					AirbyteSourceID:           mockSourceConnectorRes.AirbyteSourceId,
					AirbyteSourceDefinitionID: mockSourceConnectorRes.AirbyteSourceDefinitionId,
					Owner:                     1122,
					WorkspaceID:               1122,
				}
				arg1 := models.Connection{
					PipelineID: mockSourceConnectorReq.Pipeline,
				}

				store.EXPECT().CreateConnectionAndSourceAgainstAPipeline(arg0, arg1).Times(1).
					Return(models.Source{
						SourceID:                  mockSourceDefID.String(),
						SourceName:                arg0.SourceName,
						AirbyteSourceID:           arg0.AirbyteSourceID,
						AirbyteSourceDefinitionID: arg0.AirbyteSourceDefinitionID,
						ConnectionID:              mockSourceDefID.String(),
						Owner:                     1122,
						WorkspaceID:               1122,
					}, models.Connection{
						ConnectionID: mockSourceDefID.String(),
						PipelineID:   arg1.PipelineID,
						CreatedAt:    0,
					}, sql.ErrConnDone)
			},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				arg0 := map[string]interface{}{
					"sourceDefinitionId":      mockSourceDefID.String(),
					"connectionConfiguration": mockSourceConnectorReq.ConnectionConfiguration,
				}
				querier.EXPECT().CheckSourceConnection(arg0).Times(1).Return(nil)

				arg := models.CreateSourceConnectorRequestAirbyte{
					CreateSourceConnectorRequest: mockSourceConnectorReq.CreateSourceConnectorRequest,
					WorkspaceId:                  test.AirByteWorkspaceID,
				}

				querier.EXPECT().CreateSourceConnectorOnAirByte(arg).
					Times(1).Return(mockSourceConnectorRes, nil)
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

			body: mockSourceConnectorReq,

			buildStubs: func(store *mockStore.MockStore) {
				arg0 := models.Source{
					SourceName:                mockSourceConnectorRes.SourceName,
					AirbyteSourceID:           mockSourceConnectorRes.AirbyteSourceId,
					AirbyteSourceDefinitionID: mockSourceConnectorRes.AirbyteSourceDefinitionId,
					Owner:                     1122,
					WorkspaceID:               1122,
				}
				arg1 := models.Connection{
					PipelineID: mockSourceConnectorReq.Pipeline,
				}

				store.EXPECT().CreateConnectionAndSourceAgainstAPipeline(arg0, arg1).Times(1).
					Return(models.Source{
						SourceID:                  mockSourceDefID.String(),
						SourceName:                arg0.SourceName,
						AirbyteSourceID:           arg0.AirbyteSourceID,
						AirbyteSourceDefinitionID: arg0.AirbyteSourceDefinitionID,
						ConnectionID:              mockSourceDefID.String(),
						Owner:                     1122,
						WorkspaceID:               1122,
					}, models.Connection{
						ConnectionID: mockSourceDefID.String(),
						PipelineID:   arg1.PipelineID,
						CreatedAt:    0,
					}, nil)
			},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				arg0 := map[string]interface{}{
					"sourceDefinitionId":      mockSourceDefID.String(),
					"connectionConfiguration": mockSourceConnectorReq.ConnectionConfiguration,
				}
				querier.EXPECT().CheckSourceConnection(arg0).Times(1).Return(nil)

				arg := models.CreateSourceConnectorRequestAirbyte{
					CreateSourceConnectorRequest: mockSourceConnectorReq.CreateSourceConnectorRequest,
					WorkspaceId:                  test.AirByteWorkspaceID,
				}

				querier.EXPECT().CreateSourceConnectorOnAirByte(arg).
					Times(1).Return(mockSourceConnectorRes, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				mockAPIRes := models.CreateSourceConnectorResposneData{
					SourceID:     mockSourceDefID.String(),
					SourceName:   mockSourceConnectorRes.SourceName,
					ConnectionID: mockSourceDefID.String(),
					PipelineID:   mockSourceConnectorReq.Pipeline,
				}

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   mockAPIRes}
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

			airByte := mockairbyte.NewMockAirByteQuerier(ctrl)
			testCase.queryAirByte(airByte)

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.SOURCE, store, airByte, authServiceClient)
			url := test.BaseURL + "sources/"
			expectedResp, err := test.MakeHttpRequest(server, http.MethodPost, url, nil, body)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// TestGetConnectionSummary tests all the scenarios while getting the Connection Summary.
func TestGetConnectionSummary(t *testing.T) {
	mockSourceID, _ := uuid.NewV1()
	mockConnectionSummary := createRandomConnectionSummary()
	mockConnectionSummaryResponseAirByte := createRandomConnectionSummaryResponseAB()
	mockUser := test.CreateRandomUserDetails(1122, 1122)
	testCaseSuite := []struct {
		testScenario   string
		sourceID       string
		buildStubs     func(store *mockStore.MockStore)
		queryAirByte   func(querier *mockairbyte.MockAirByteQuerier)
		getUserDetails func(client *mock_authservice.MockHttpClient)
		checkResponse  func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "Bad Response from DB",

			sourceID: mockSourceID.String(),

			buildStubs: func(store *mockStore.MockStore) {
				store.EXPECT().GetSourceAndConnectionDetails(mockSourceID.String()).Times(1).Return(models.ConnectionSummary{}, sql.ErrConnDone)
			},

			getUserDetails: func(client *mock_authservice.MockHttpClient) {},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {},

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
			testScenario: "Bad Response from AirByte",

			sourceID: mockSourceID.String(),

			buildStubs: func(store *mockStore.MockStore) {
				store.EXPECT().GetSourceAndConnectionDetails(mockSourceID.String()).Times(1).Return(mockConnectionSummary, nil)
			},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				querier.EXPECT().GetConnectionSummary(mockConnectionSummary.AirbyteConnectionID).Times(1).
					Return(models.ConnectionSummaryAirByte{}, errors.New("failed to request AirByte"))
			},

			getUserDetails: func(client *mock_authservice.MockHttpClient) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: "failed to request AirByte",
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "Bad Response from Auth",

			sourceID: mockSourceID.String(),

			buildStubs: func(store *mockStore.MockStore) {
				store.EXPECT().GetSourceAndConnectionDetails(mockSourceID.String()).Times(1).Return(mockConnectionSummary, nil)
			},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				querier.EXPECT().GetConnectionSummary(mockConnectionSummary.AirbyteConnectionID).Times(1).
					Return(mockConnectionSummaryResponseAirByte, nil)
			},

			getUserDetails: func(client *mock_authservice.MockHttpClient) {
				ownerID := mockConnectionSummary.Owner
				arg := fmt.Sprintf("%s/auth-service/api/v1/accounts/internal/user-from-id?user_id=%d", env.Env.AuthServiceAddress, ownerID)
				res := http.Response{}
				client.EXPECT().Get(arg).Times(1).Return(&res, errors.New("User Not Found"))
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

			sourceID: mockSourceID.String(),

			buildStubs: func(store *mockStore.MockStore) {
				store.EXPECT().GetSourceAndConnectionDetails(mockSourceID.String()).Times(1).Return(mockConnectionSummary, nil)
			},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				querier.EXPECT().GetConnectionSummary(mockConnectionSummary.AirbyteConnectionID).Times(1).
					Return(mockConnectionSummaryResponseAirByte, nil)
			},

			getUserDetails: func(client *mock_authservice.MockHttpClient) {
				test.MockGetUserByID(client, 1122, 1122)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data: models.ConnectionSummaryResponse{
						AirByteSummary:    mockConnectionSummaryResponseAirByte,
						SourceName:        mockConnectionSummary.SourceName,
						Owner:             mockUser.Payload.UserInfo,
						ConfigurationDate: mockConnectionSummary.CreatedAt,
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

			airByte := mockairbyte.NewMockAirByteQuerier(ctrl)
			testCase.queryAirByte(airByte)

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)
			testCase.getUserDetails(httpMockClient)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.SOURCE, store, airByte, authServiceClient)
			url := fmt.Sprintf("%ssources/%s/summary/", test.BaseURL, testCase.sourceID)
			expectedResp, err := test.MakeHttpRequest(server, http.MethodGet, url, nil, nil)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// TestEditSourceOnAirByte tests all the scenarios while Editing the configured sources on airByte.
func TestEditSourceOnAirByte(t *testing.T) {
	mockSourceID, _ := uuid.NewV1()
	mockSource := createRandomSource(mockSourceID.String())
	mockEditConnectorReq := models.EditSourceConnectorRequest{ConnectionConfiguration: "{port:5432}"}
	testCaseSuite := []struct {
		testScenario  string
		body          models.EditSourceConnectorRequest
		sourceID      string
		queryAirByte  func(querier *mockairbyte.MockAirByteQuerier)
		buildStubs    func(store *mockStore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "BadRequest",

			sourceID: "Bad Source ID",

			body: models.EditSourceConnectorRequest{},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {},

			buildStubs: func(store *mockStore.MockStore) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "Internal Server Error",

			sourceID: mockSourceID.String(),

			body: models.EditSourceConnectorRequest{},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {},

			buildStubs: func(store *mockStore.MockStore) {
				arg := mockSourceID.String()
				store.EXPECT().GetSource(arg).Times(1).Return(models.Source{}, sql.ErrConnDone)
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
			testScenario: "Bad Request",

			sourceID: mockSourceID.String(),

			body: models.EditSourceConnectorRequest{},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {},

			buildStubs: func(store *mockStore.MockStore) {
				arg := mockSourceID.String()
				store.EXPECT().GetSource(arg).Times(1).Return(mockSource, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "Bad Request to AirByte",

			sourceID: mockSourceID.String(),

			body: mockEditConnectorReq,

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				arg := map[string]interface{}{
					"sourceDefinitionId":      mockSource.AirbyteSourceDefinitionID,
					"connectionConfiguration": mockEditConnectorReq.ConnectionConfiguration,
				}
				querier.EXPECT().CheckSourceConnection(arg).Times(1).Return(errors.New("Bad Source Connector"))
			},

			buildStubs: func(store *mockStore.MockStore) {
				arg := mockSourceID.String()
				store.EXPECT().GetSource(arg).Times(1).Return(mockSource, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: errors.New("Bad Source Connector").Error(),
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "Bad Request to AirByte",

			sourceID: mockSourceID.String(),

			body: mockEditConnectorReq,

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				arg := map[string]interface{}{
					"sourceDefinitionId":      mockSource.AirbyteSourceDefinitionID,
					"connectionConfiguration": mockEditConnectorReq.ConnectionConfiguration,
				}
				querier.EXPECT().CheckSourceConnection(arg).Times(1).Return(nil)

				arg1 := models.EditSourceConnectorRequestAirByte{
					AirByteSourceID:         mockSource.AirbyteSourceID,
					ConnectionConfiguration: mockEditConnectorReq.ConnectionConfiguration,
					Name:                    mockSource.SourceName,
				}
				querier.EXPECT().EditSourceConnectorOnAirByte(arg1).Times(1).
					Return(models.CreateSourceConnectorResponseAirbyte{}, errors.New("bad request to AirByte"))
			},

			buildStubs: func(store *mockStore.MockStore) {
				arg := mockSourceID.String()
				store.EXPECT().GetSource(arg).Times(1).Return(mockSource, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: errors.New("bad request to AirByte").Error(),
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "Success",

			sourceID: mockSourceID.String(),

			body: mockEditConnectorReq,

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				arg := map[string]interface{}{
					"sourceDefinitionId":      mockSource.AirbyteSourceDefinitionID,
					"connectionConfiguration": mockEditConnectorReq.ConnectionConfiguration,
				}
				querier.EXPECT().CheckSourceConnection(arg).Times(1).Return(nil)

				arg1 := models.EditSourceConnectorRequestAirByte{
					AirByteSourceID:         mockSource.AirbyteSourceID,
					ConnectionConfiguration: mockEditConnectorReq.ConnectionConfiguration,
					Name:                    mockSource.SourceName,
				}
				querier.EXPECT().EditSourceConnectorOnAirByte(arg1).Times(1).
					Return(models.CreateSourceConnectorResponseAirbyte{}, nil)
			},

			buildStubs: func(store *mockStore.MockStore) {
				arg := mockSourceID.String()
				store.EXPECT().GetSource(arg).Times(1).Return(mockSource, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   "Source Edited Successfully"}
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

			airByte := mockairbyte.NewMockAirByteQuerier(ctrl)
			testCase.queryAirByte(airByte)

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.SOURCE, store, airByte, authServiceClient)
			url := test.BaseURL + "sources/" + testCase.sourceID + "/"
			expectedResp, err := test.MakeHttpRequest(server, http.MethodPut, url, nil, body)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// TestGetConfiguredSource tests all the scenarios while getting the configured sources on airByte.
func TestGetConfiguredSource(t *testing.T) {
	mockSourceID, _ := uuid.NewV1()
	mockSource := createRandomSource(mockSourceID.String())
	mockConfiguredSource := createRandomConfiguredSource(mockSourceID.String())
	testCaseSuite := []struct {
		testScenario  string
		sourceID      string
		queryAirByte  func(querier *mockairbyte.MockAirByteQuerier)
		buildStubs    func(store *mockStore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{

		{
			testScenario: "Internal Server Error",

			sourceID: mockSourceID.String(),

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {},

			buildStubs: func(store *mockStore.MockStore) {
				arg := mockSourceID.String()
				store.EXPECT().GetSource(arg).Times(1).Return(models.Source{}, sql.ErrConnDone)
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
			testScenario: "Bad Request to AirByte",

			sourceID: mockSourceID.String(),

			buildStubs: func(store *mockStore.MockStore) {
				arg := mockSourceID.String()
				store.EXPECT().GetSource(arg).Times(1).Return(mockSource, nil)
			},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				querier.EXPECT().GetConfiguredSource(mockSource.AirbyteSourceID).Times(1).
					Return(models.ConfiguredSource{}, errors.New("bad request to AirByte"))
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: errors.New("bad request to AirByte").Error(),
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "Success",

			sourceID: mockSourceID.String(),

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				querier.EXPECT().GetConfiguredSource(mockSource.AirbyteSourceID).Times(1).
					Return(mockConfiguredSource, nil)
			},

			buildStubs: func(store *mockStore.MockStore) {
				arg := mockSourceID.String()
				store.EXPECT().GetSource(arg).Times(1).Return(mockSource, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   mockConfiguredSource}
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

			airByte := mockairbyte.NewMockAirByteQuerier(ctrl)
			testCase.queryAirByte(airByte)

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.SOURCE, store, airByte, authServiceClient)
			url := test.BaseURL + "sources/" + testCase.sourceID + "/"
			expectedResp, err := test.MakeHttpRequest(server, http.MethodGet, url, nil, nil)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// TestGetSupportedSources tests all the scenarios while getting the supported sources.
func TestGetSupportedSources(t *testing.T) {
	mockSupportedSources := []models.SupportedSources{
		createRandomSupportedSource(),
		createRandomSupportedSource(),
		createRandomSupportedSource()}

	testCaseSuite := []struct {
		testScenario  string
		buildStubs    func(store *mockStore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "InternalServerError",

			buildStubs: func(store *mockStore.MockStore) {
				store.EXPECT().GetSupportedSources().Times(1).Return([]models.SupportedSources{}, sql.ErrNoRows)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "Success",

			buildStubs: func(store *mockStore.MockStore) {
				store.EXPECT().GetSupportedSources().Times(1).Return(mockSupportedSources, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   mockSupportedSources}
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

			airByte := mockairbyte.NewMockAirByteQuerier(ctrl)

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.SOURCE, store, airByte, authServiceClient)
			url := test.BaseURL + "sources/"
			expectedResp, err := test.MakeHttpRequest(server, http.MethodGet, url, nil, nil)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// TestGetSourceSpecifications tests all the scenarios while getting the source specifications.
func TestGetSourceSpecifications(t *testing.T) {
	mockSourceSpecification := createRandomSourceSpecification()
	name := "randomDefinition"
	SourceDefinitions := []models.SourceDefinition{createRandomSourceDefinition(name), createRandomSourceDefinition("random")}
	mockSourceDefinitions := models.SourceDefinitions{SourceDefinitions: SourceDefinitions}

	testCaseSuite := []struct {
		testScenario  string
		query         string
		queryAirByte  func(querier *mockairbyte.MockAirByteQuerier)
		buildStubs    func(store *mockStore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "BadGetSourceDefinitionsCall",

			query: mockSourceDefinitions.SourceDefinitions[0].Name,

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				querier.EXPECT().GetSourceDefinitions().Times(1).Return(models.SourceDefinitions{}, errors.New("Bad Request"))
			},

			buildStubs: func(store *mockStore.MockStore) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "BadGetSourceSpecificationsCall",

			query: mockSourceDefinitions.SourceDefinitions[0].Name,

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				querier.EXPECT().GetSourceDefinitions().Times(1).Return(mockSourceDefinitions, nil)

				arg := mockSourceDefinitions.SourceDefinitions[0].SourceDefinitionID

				querier.EXPECT().GetSourceSpecification(arg).Times(1).Return(models.SourceSpecification{}, errors.New("Bad Request"))
			},

			buildStubs: func(store *mockStore.MockStore) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "Success",

			query: mockSourceDefinitions.SourceDefinitions[0].Name,

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				querier.EXPECT().GetSourceDefinitions().Times(1).Return(mockSourceDefinitions, nil)

				arg := mockSourceDefinitions.SourceDefinitions[0].SourceDefinitionID

				querier.EXPECT().GetSourceSpecification(arg).Times(1).Return(mockSourceSpecification, nil)
			},

			buildStubs: func(store *mockStore.MockStore) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   mockSourceSpecification}
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

			airByte := mockairbyte.NewMockAirByteQuerier(ctrl)
			testCase.queryAirByte(airByte)

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.SOURCE, store, airByte, authServiceClient)

			query := make(map[string]string)
			query["source"] = testCase.query

			url := test.BaseURL + "sources/specification/"
			expectedResp, err := test.MakeHttpRequest(server, http.MethodGet, url, query, nil)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// TestDiscoverSourceSchema tests all the scenarios while discovering the source schema.
func TestDiscoverSourceSchema(t *testing.T) {
	randomSourceID := "randomSourceID"
	mockRandomSource := createRandomSource(randomSourceID)
	testCaseSuite := []struct {
		testScenario  string
		query         string
		queryAirByte  func(querier *mockairbyte.MockAirByteQuerier)
		buildStubs    func(store *mockStore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "Couldn't get discovered source schema",

			query: randomSourceID,

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {},

			buildStubs: func(store *mockStore.MockStore) {
				arg := randomSourceID
				store.EXPECT().GetSource(arg).Times(1).Return(models.Source{}, sql.ErrNoRows)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "Bad Air Byte Response",

			query: randomSourceID,

			buildStubs: func(store *mockStore.MockStore) {
				arg := randomSourceID
				store.EXPECT().GetSource(arg).Times(1).Return(mockRandomSource, nil)
			},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				arg := mockRandomSource.AirbyteSourceID
				querier.EXPECT().DiscoverSourceSchema(arg).Times(1).Return(models.SourceSchema{}, errors.New("AirByte Server is Down"))
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				res := models.Response{
					Status: utils.ERROR,
					Errors: errors.New("couldn't get discover source schema").Error(),
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "Success",

			query: randomSourceID,

			buildStubs: func(store *mockStore.MockStore) {
				arg := randomSourceID
				store.EXPECT().GetSource(arg).Times(1).Return(mockRandomSource, nil)
			},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				arg := mockRandomSource.AirbyteSourceID
				querier.EXPECT().DiscoverSourceSchema(arg).Times(1).Return(models.SourceSchema{}, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   models.SourceSchema{}}
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

			airByte := mockairbyte.NewMockAirByteQuerier(ctrl)
			testCase.queryAirByte(airByte)

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.SOURCE, store, airByte, authServiceClient)

			query := make(map[string]string)
			query["source_id"] = testCase.query

			url := test.BaseURL + "sources/discover/schema/"
			expectedResp, err := test.MakeHttpRequest(server, http.MethodGet, url, query, nil)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

//createRandomConnectionSummary populates and return the ConnectionSummary model with random values.
func createRandomConnectionSummary() models.ConnectionSummary {
	abID, _ := uuid.NewV1()

	cs := models.ConnectionSummary{
		AirbyteConnectionID: abID.String(),
		SourceName:          utils.RandomString(5),
		Owner:               1122,
		CreatedAt:           utils.RandomInt(1, 10),
	}

	return cs
}

//createRandomConnectionSummaryResponseAB populates and return the ConnectionSummaryAirByte model with random values.
func createRandomConnectionSummaryResponseAB() models.ConnectionSummaryAirByte {
	cs := models.ConnectionSummaryAirByte{
		Schedule: struct {
			Units    int    `json:"units"`
			TimeUnit string `json:"timeUnit"`
		}{Units: int(utils.RandomInt(1, 10)),
			TimeUnit: "minutes"},
		ConnectionStatus: "active",
	}

	return cs
}

// createRandomSourceConnectorRequestAPI populates and return the CreateSourceConnectorRequest Model with random values.
func createRandomSourceConnectorRequestAPI(sourceDefID string) models.CreateSourceConnectorRequestAPI {
	sID, _ := uuid.NewV1()
	Scr := models.CreateSourceConnectorRequestAPI{
		CreateSourceConnectorRequest: models.CreateSourceConnectorRequest{
			AirbyteSourceDefinitionId: sourceDefID,
			ConnectionConfiguration:   "{This will vary from source to source}",
			Name:                      utils.RandomString(5),
		},
		Pipeline: sID.String(),
	}

	return Scr
}

// createRandomSupportedSource populates and return the Supported Source Model with random values.
func createRandomSupportedSource() models.SupportedSources {
	sID, _ := uuid.NewV1()
	Ss := models.SupportedSources{
		ID:   sID,
		Name: utils.RandomString(5),
		Type: utils.RandomString(5),
	}

	return Ss
}

// createRandomSource populates and return the SourceModel with Random values.
func createRandomSource(sid string) models.Source {
	s := models.Source{
		SourceID:                  sid,
		SourceName:                utils.RandomString(5),
		AirbyteSourceID:           utils.RandomString(5),
		AirbyteSourceDefinitionID: utils.RandomString(5),
		ConnectionID:              utils.RandomString(5),
	}

	return s
}

// createRandomConfiguredSource populates and return the ConfiguredSource with Random values.
func createRandomConfiguredSource(sid string) models.ConfiguredSource {
	cs := models.ConfiguredSource{
		SourceDefinitionId:      sid,
		SourceId:                sid,
		WorkspaceId:             test.AirByteWorkspaceID,
		ConnectionConfiguration: "{}",
		Name:                    utils.RandomString(5),
		SourceName:              utils.RandomString(5),
	}

	return cs
}

// createRandomSourceDefinition populates and return the Source Definition Model with random values.
func createRandomSourceDefinition(name string) models.SourceDefinition {
	Sd := models.SourceDefinition{
		SourceDefinitionID: utils.RandomString(5),
		Name:               name,
		DockerRepository:   utils.RandomString(5),
		DockerImageTag:     utils.RandomString(5),
		DocumentationURL:   utils.RandomString(5),
		Icon:               utils.RandomString(5),
	}

	return Sd
}

// createRandomSourceSpecification populates and return the Source Specification Model with random values.
func createRandomSourceSpecification() models.SourceSpecification {
	Ss := models.SourceSpecification{
		SourceDefinitionID:      utils.RandomString(5),
		DocumentationURL:        utils.RandomString(10),
		ConnectionSpecification: nil,
		AuthSpecification:       nil,
		AdvancedAuth:            nil,
		JobInfo:                 nil,
	}

	return Ss
}

// TestMain runs the package level test in TestMode.
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
