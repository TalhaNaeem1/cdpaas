package destination_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gorm.io/datatypes" //nolint
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	// nolint
	"gorm.io/datatypes" 
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

// TestConfigureDestinationOnAirByte tests all the scenarios while configuring the destination on AirByte.
func TestConfigureDestinationOnAirByte(t *testing.T) {
	mockDestinationType := utils.RandomString(5)
	mockWorkSpaceID := test.AirByteWorkspaceID
	mockCreateDestinationConnectorRequest := createRandomDestinationConnectorRequest()
	testCaseSuite := []struct {
		testScenario  string
		body          models.CreateDestinationConnectorRequestAPI
		queryAirByte  func(querier *mockairbyte.MockAirByteQuerier)
		buildStubs    func(store *mockStore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "Bad Request",

			body: models.CreateDestinationConnectorRequestAPI{
				CreateDestinationConnectorRequest: mockCreateDestinationConnectorRequest,
				DestinationType:                   "",
			},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {},

			buildStubs: func(store *mockStore.MockStore) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "Internal Server Error",

			body: models.CreateDestinationConnectorRequestAPI{
				CreateDestinationConnectorRequest: mockCreateDestinationConnectorRequest,
				DestinationType:                   mockDestinationType,
			},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				arg := map[string]interface{}{
					"destinationDefinitionId": mockCreateDestinationConnectorRequest.AirbyteDestinationDefinitionId,
					"connectionConfiguration": mockCreateDestinationConnectorRequest.ConnectionConfiguration,
				}
				querier.EXPECT().CheckDestinationConnection(arg).Times(1).Return(errors.New("Bad Destination Connector"))
			},

			buildStubs: func(store *mockStore.MockStore) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: "Bad Destination Connector",
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "Internal Server Error",

			body: models.CreateDestinationConnectorRequestAPI{
				CreateDestinationConnectorRequest: mockCreateDestinationConnectorRequest,
				DestinationType:                   mockDestinationType,
			},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				arg0 := map[string]interface{}{
					"destinationDefinitionId": mockCreateDestinationConnectorRequest.AirbyteDestinationDefinitionId,
					"connectionConfiguration": mockCreateDestinationConnectorRequest.ConnectionConfiguration,
				}
				querier.EXPECT().CheckDestinationConnection(arg0).Times(1).Return(nil)

				arg := models.CreateDestinationConnectorRequestAirbyte{
					CreateDestinationConnectorRequest: mockCreateDestinationConnectorRequest,
					WorkspaceId:                       mockWorkSpaceID,
				}

				querier.EXPECT().CreateDestinationConnectorOnAirByte(arg).Times(1).
					Return(models.CreateDestinationConnectorResponseAirbyte{}, errors.New("unable to create Destination on AirByte"))
			},

			buildStubs: func(store *mockStore.MockStore) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: errors.New("unable to create Destination on AirByte").Error(),
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "Internal Server Error",

			body: models.CreateDestinationConnectorRequestAPI{
				CreateDestinationConnectorRequest: mockCreateDestinationConnectorRequest,
				DestinationType:                   mockDestinationType,
			},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				arg0 := map[string]interface{}{
					"destinationDefinitionId": mockCreateDestinationConnectorRequest.AirbyteDestinationDefinitionId,
					"connectionConfiguration": mockCreateDestinationConnectorRequest.ConnectionConfiguration,
				}
				querier.EXPECT().CheckDestinationConnection(arg0).Times(1).Return(nil)

				arg := models.CreateDestinationConnectorRequestAirbyte{
					CreateDestinationConnectorRequest: mockCreateDestinationConnectorRequest,
					WorkspaceId:                       mockWorkSpaceID,
				}

				querier.EXPECT().CreateDestinationConnectorOnAirByte(arg).Times(1).
					Return(models.CreateDestinationConnectorResponseAirbyte{
						AirbyteDestinationId: mockCreateDestinationConnectorRequest.AirbyteDestinationDefinitionId,
						DestinationName:      mockCreateDestinationConnectorRequest.Name,
						CreateDestinationConnectorRequestAirbyte: models.CreateDestinationConnectorRequestAirbyte{
							CreateDestinationConnectorRequest: mockCreateDestinationConnectorRequest,
							WorkspaceId:                       mockWorkSpaceID,
						},
					}, nil)
			},

			buildStubs: func(store *mockStore.MockStore) {
				arg := models.Destination{
					DestinationName:         mockCreateDestinationConnectorRequest.Name,
					AirbyteDestinationID:    mockCreateDestinationConnectorRequest.AirbyteDestinationDefinitionId,
					AirbyteDestDefinitionID: mockCreateDestinationConnectorRequest.AirbyteDestinationDefinitionId,
					DestinationType:         mockDestinationType,
					ConfigurationDetails:    mockCreateDestinationConnectorRequest.ConnectionConfiguration,
					Owner:                   1122,
					WorkspaceID:             1122,
				}
				store.EXPECT().CreateDestination(arg).Times(1).Return(models.Destination{}, sql.ErrConnDone)
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

			body: models.CreateDestinationConnectorRequestAPI{
				CreateDestinationConnectorRequest: mockCreateDestinationConnectorRequest,
				DestinationType:                   mockDestinationType,
			},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				arg0 := map[string]interface{}{
					"destinationDefinitionId": mockCreateDestinationConnectorRequest.AirbyteDestinationDefinitionId,
					"connectionConfiguration": mockCreateDestinationConnectorRequest.ConnectionConfiguration,
				}
				querier.EXPECT().CheckDestinationConnection(arg0).Times(1).Return(nil)

				arg := models.CreateDestinationConnectorRequestAirbyte{
					CreateDestinationConnectorRequest: mockCreateDestinationConnectorRequest,
					WorkspaceId:                       mockWorkSpaceID,
				}

				querier.EXPECT().CreateDestinationConnectorOnAirByte(arg).Times(1).
					Return(models.CreateDestinationConnectorResponseAirbyte{
						AirbyteDestinationId: mockCreateDestinationConnectorRequest.AirbyteDestinationDefinitionId,
						DestinationName:      mockCreateDestinationConnectorRequest.Name,
						CreateDestinationConnectorRequestAirbyte: models.CreateDestinationConnectorRequestAirbyte{
							CreateDestinationConnectorRequest: mockCreateDestinationConnectorRequest,
							WorkspaceId:                       mockWorkSpaceID,
						},
					}, nil)
			},

			buildStubs: func(store *mockStore.MockStore) {
				arg := models.Destination{
					DestinationName:         mockCreateDestinationConnectorRequest.Name,
					AirbyteDestinationID:    mockCreateDestinationConnectorRequest.AirbyteDestinationDefinitionId,
					AirbyteDestDefinitionID: mockCreateDestinationConnectorRequest.AirbyteDestinationDefinitionId,
					DestinationType:         mockDestinationType,
					ConfigurationDetails:    mockCreateDestinationConnectorRequest.ConnectionConfiguration,
					Owner:                   1122,
					WorkspaceID:             1122,
				}
				store.EXPECT().CreateDestination(arg).Times(1).Return(models.Destination{
					DestinationID:           mockWorkSpaceID,
					DestinationName:         mockCreateDestinationConnectorRequest.Name,
					AirbyteDestinationID:    mockCreateDestinationConnectorRequest.AirbyteDestinationDefinitionId,
					AirbyteDestDefinitionID: mockCreateDestinationConnectorRequest.AirbyteDestinationDefinitionId,
					DestinationType:         mockDestinationType,
					ConfigurationDetails:    mockCreateDestinationConnectorRequest.ConnectionConfiguration,
					Owner:                   1122,
					WorkspaceID:             1122,
				}, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				mockAPIRes := models.CreateDestinationConnectorResponseData{
					DestinationID:   mockWorkSpaceID,
					DestinationName: mockCreateDestinationConnectorRequest.Name,
					DestinationType: mockDestinationType,
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

			server := test.NewTestServer(test.DESTINATION, store, airByte, authServiceClient)
			url := test.BaseURL + "destinations/"
			expectedResp, err := test.MakeHttpRequest(server, http.MethodPost, url, nil, body)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// TestGetSupportedDestinations tests all the scenarios while getting the supported destinations.
func TestGetSupportedDestinations(t *testing.T) {
	mockSupportedDestinations := []models.SupportedDestinations{createRandomSupportedDestination(), createRandomSupportedDestination(), createRandomSupportedDestination()}

	testCaseSuite := []struct {
		testScenario  string
		buildStubs    func(store *mockStore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "InternalServerError",

			//getUserDetails: func(client *mock_authservice.MockHttpClient) {
			//	test.MockGetUserInfo(client)
			//},

			buildStubs: func(store *mockStore.MockStore) {
				store.EXPECT().GetSupportedDestinations().Times(1).Return([]models.SupportedDestinations{}, sql.ErrNoRows)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "Success",

			//getUserDetails: func(client *mock_authservice.MockHttpClient) {
			//	test.MockGetUserInfo(client)
			//},

			buildStubs: func(store *mockStore.MockStore) {
				store.EXPECT().GetSupportedDestinations().Times(1).Return(mockSupportedDestinations, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   mockSupportedDestinations}
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

			server := test.NewTestServer(test.DESTINATION, store, airByte, authServiceClient)
			url := test.BaseURL + "destinations/"
			expectedResp, err := test.MakeHttpRequest(server, http.MethodGet, url, nil, nil)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// TestGetConfiguredDestinations tests all the scenarios while getting the configured destinations.
func TestGetConfiguredDestinations(t *testing.T) {
	mockConfiguredDestinations := []models.ConfiguredDestination{createRandomConfiguredDestination(), createRandomConfiguredDestination(), createRandomConfiguredDestination()}
	mockWorkSpaceID := "1122"

	testCaseSuite := []struct {
		testScenario  string
		buildStubs    func(store *mockStore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "InternalServerError",

			buildStubs: func(store *mockStore.MockStore) {
				arg, _ := strconv.Atoi(mockWorkSpaceID)
				store.EXPECT().GetConfiguredDestination(arg).Times(1).Return([]models.ConfiguredDestination{}, sql.ErrNoRows)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "Success",

			buildStubs: func(store *mockStore.MockStore) {
				arg, _ := strconv.Atoi(mockWorkSpaceID)
				store.EXPECT().GetConfiguredDestination(arg).Times(1).Return(mockConfiguredDestinations, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   mockConfiguredDestinations}
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

			server := test.NewTestServer(test.DESTINATION, store, airByte, authServiceClient)
			url := test.BaseURL + "destinations/configured/"
			expectedResp, err := test.MakeHttpRequest(server, http.MethodGet, url, nil, nil)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// TestGetDestinationSummary tests all the scenarios while getting the destination Summary.
func TestGetDestinationSummary(t *testing.T) {
	mockDestinationID, _ := uuid.NewV1()
	mockDestinationSummary := createRandomDestinationSummary()
	mockUser := test.CreateRandomUserDetails(1122, 1122)
	testCaseSuite := []struct {
		testScenario   string
		destinationID  string
		buildStubs     func(store *mockStore.MockStore)
		getUserDetails func(client *mock_authservice.MockHttpClient)
		checkResponse  func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "Bad Destination ID",

			destinationID: "not a uuid",

			buildStubs: func(store *mockStore.MockStore) {},

			getUserDetails: func(client *mock_authservice.MockHttpClient) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "Bad Response from DB",

			destinationID: mockDestinationID.String(),

			buildStubs: func(store *mockStore.MockStore) {
				store.EXPECT().GetDestinationSummary(mockDestinationID).Times(1).Return(models.DestinationSummary{}, sql.ErrConnDone)
			},

			getUserDetails: func(client *mock_authservice.MockHttpClient) {},

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
			testScenario: "Bad Response from Auth",

			destinationID: mockDestinationID.String(),

			buildStubs: func(store *mockStore.MockStore) {
				store.EXPECT().GetDestinationSummary(mockDestinationID).Times(1).Return(mockDestinationSummary, nil)
			},

			getUserDetails: func(client *mock_authservice.MockHttpClient) {
				ownerID := mockDestinationSummary.Owner
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

			destinationID: mockDestinationID.String(),

			buildStubs: func(store *mockStore.MockStore) {
				store.EXPECT().GetDestinationSummary(mockDestinationID).Times(1).Return(mockDestinationSummary, nil)
			},

			getUserDetails: func(client *mock_authservice.MockHttpClient) {
				test.MockGetUserByID(client, 1122, 1122)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data: models.DestinationSummaryResponse{
						DestinationName: mockDestinationSummary.DestinationName,
						Owner:           mockUser.Payload.UserInfo,
						CreatedAt:       mockDestinationSummary.CreatedAt,
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

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)
			testCase.getUserDetails(httpMockClient)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.DESTINATION, store, airByte, authServiceClient)
			url := fmt.Sprintf("%sdestinations/%s/summary/", test.BaseURL, testCase.destinationID)
			expectedResp, err := test.MakeHttpRequest(server, http.MethodGet, url, nil, nil)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// TestGetDestinationSpecification tests all the scenarios while getting the destination Specifications.
func TestGetDestinationSpecification(t *testing.T) {
	mockDestinationSpecification := createRandomDestinationSpecification()
	mockDestinationDefName := "randomDestinationDefinition"
	mockDestinationDefID := utils.RandomString(10)
	destinationDefinitions := []models.DestinationDefinition{createRandomDestinationDefinition(mockDestinationDefID, mockDestinationDefName), createRandomDestinationDefinition("1", "ONE")}
	mockDestinationDefinitions := models.DestinationDefinitions{DestinationDefinitions: destinationDefinitions}

	testCaseSuite := []struct {
		testScenario  string
		query         string
		queryAirByte  func(querier *mockairbyte.MockAirByteQuerier)
		buildStubs    func(store *mockStore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "BadGetDestinationDefinitionsCall",

			//getUserDetails: func(client *mock_authservice.MockHttpClient) {
			//	test.MockGetUserInfo(client)
			//},

			query: "",

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {},

			buildStubs: func(store *mockStore.MockStore) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "BadGettingDestDefinitionsAB",

			//getUserDetails: func(client *mock_authservice.MockHttpClient) {
			//	test.MockGetUserInfo(client)
			//},

			query: mockDestinationDefName,

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				querier.EXPECT().GetDestinationDefinitions().Times(1).Return(models.DestinationDefinitions{},
					errors.New("can't get Destination Defs"))
			},

			buildStubs: func(store *mockStore.MockStore) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: errors.New("can't get Destination Defs").Error(),
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},

		{
			testScenario: "BadGetDestSpecificationsCall",

			//getUserDetails: func(client *mock_authservice.MockHttpClient) {
			//	test.MockGetUserInfo(client)
			//},

			query: mockDestinationDefName,

			buildStubs: func(store *mockStore.MockStore) {},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				querier.EXPECT().GetDestinationDefinitions().Times(1).Return(mockDestinationDefinitions, nil)

				arg := mockDestinationDefID

				querier.EXPECT().GetDestinationSpecification(arg).Times(1).Return(models.DestinationSpecification{}, errors.New("Bad Request"))
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

				res := models.Response{
					Status: utils.ERROR,
					Errors: errors.New("Bad Request").Error(),
					Data:   nil}
				actual, e := json.Marshal(res)
				require.NoError(t, e)
				test.ReqResBodyMatcher(t, recorder.Body, actual)
			},
		},
		{
			testScenario: "BadGetDestSpecificationsCall",

			//getUserDetails: func(client *mock_authservice.MockHttpClient) {
			//	test.MockGetUserInfo(client)
			//},

			query: "random",

			buildStubs: func(store *mockStore.MockStore) {},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				querier.EXPECT().GetDestinationDefinitions().Times(1).Return(mockDestinationDefinitions, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "Success",

			//getUserDetails: func(client *mock_authservice.MockHttpClient) {
			//	test.MockGetUserInfo(client)
			//},

			query: mockDestinationDefName,

			buildStubs: func(store *mockStore.MockStore) {},

			queryAirByte: func(querier *mockairbyte.MockAirByteQuerier) {
				querier.EXPECT().GetDestinationDefinitions().Times(1).Return(mockDestinationDefinitions, nil)

				arg := mockDestinationDefID

				querier.EXPECT().GetDestinationSpecification(arg).Times(1).Return(mockDestinationSpecification, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   mockDestinationSpecification}
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

			server := test.NewTestServer(test.DESTINATION, store, airByte, authServiceClient)

			query := make(map[string]string)
			query["destination"] = testCase.query

			url := test.BaseURL + "destinations/specification/"
			expectedResp, err := test.MakeHttpRequest(server, http.MethodGet, url, query, nil)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

//createRandomDestinationSummary populates and return the DestinationSummary model with random values.
func createRandomDestinationSummary() models.DestinationSummary {
	ds := models.DestinationSummary{
		DestinationName: utils.RandomString(5),
		Owner:           1122,
		CreatedAt:       utils.RandomInt(1, 10),
	}

	return ds
}

//createRandomDestinationDefinition populates and return the DestinationDefinition model with random values.
func createRandomDestinationDefinition(ddID string, ddName string) models.DestinationDefinition {
	dd := models.DestinationDefinition{
		DestinationDefinitionID: ddID,
		Name:                    ddName,
		DockerRepository:        utils.RandomString(5),
		DockerImageTag:          utils.RandomString(5),
		DocumentationURL:        utils.RandomString(5),
		Icon:                    utils.RandomString(5),
	}

	return dd
}

//createRandomDestinationSpecification populates and return the DestinationSpecification model with random values.
func createRandomDestinationSpecification() models.DestinationSpecification {
	ds := models.DestinationSpecification{
		DestinationDefinitionID:       utils.RandomString(5),
		DocumentationURL:              utils.RandomString(5),
		ConnectionSpecification:       nil,
		AuthSpecification:             nil,
		AdvancedAuth:                  nil,
		JobInfo:                       nil,
		SupportedDestinationSyncModes: nil,
		SupportsDbt:                   false,
		SupportsNormalization:         false,
	}

	return ds
}

func createRandomConfiguredDestination() models.ConfiguredDestination {
	dID, _ := uuid.NewV1()
	abID, _ := uuid.NewV1()
	d := models.ConfiguredDestination{
		DestinationID:           dID.String(),
		DestinationName:         utils.RandomString(5),
		AirbyteDestinationID:    abID.String(),
		AirbyteDestDefinitionID: abID.String(),
		DestinationType:         utils.RandomString(5),
	}

	return d
}

//createRandomSupportedDestination populates and return the SupportedDestinations model with random values.
func createRandomSupportedDestination() models.SupportedDestinations {
	sdID, _ := uuid.NewV1()
	sd := models.SupportedDestinations{
		ID:   sdID,
		Name: utils.RandomString(5),
		Type: utils.RandomString(5),
	}

	return sd
}

//createRandomDestinationConnectorRequest populates and return the CreateDestinationConnectorRequest model with random values.
func createRandomDestinationConnectorRequest() models.CreateDestinationConnectorRequest {
	scr := models.CreateDestinationConnectorRequest{
		AirbyteDestinationDefinitionId: utils.RandomString(10),
		ConnectionConfiguration:        datatypes.JSON("{Random Configuration}"),
		Name:                           utils.RandomString(5),
	}

	return scr
}

// TestMain runs the package level test in TestMode.
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
