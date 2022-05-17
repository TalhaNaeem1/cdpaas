package dataProduct_test

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
	"github.com/stretchr/testify/require"
	"pipelineService/clients/authService"
	mock_authservice "pipelineService/clients/authService/mocks"
	"pipelineService/handlers/v1/test"
	"pipelineService/models/v1"
	mockStore "pipelineService/services/db/mocks"
	"pipelineService/utils"
)

// TestCreateDataProduct tests all the scenarios while creating a Data Product.
func TestCreateDataProduct(t *testing.T) {
	mockDataProduct := createRandomDataProduct()

	testCaseSuite := []struct {
		testScenario  string
		body          models.DataProduct
		buildStubs    func(store *mockStore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		//{
		//	testScenario: "BadRequest",
		//
		//	getUserDetails: func(client *mock_authservice.MockHttpClient) {
		//		test.MockGetUserInfo(client)
		//	},
		//
		//	body: models.DataProduct{
		//		ProductID:             uuid.UUID{},
		//		Name:                  "",
		//		DataProductGovernance: mockDataProduct.DataProductGovernance,
		//		DataDomain:            mockDataProduct.DataDomain,
		//		Description:           mockDataProduct.Description,
		//		DataProductStatus:     mockDataProduct.DataProductStatus,
		//		LastUpdated:           mockDataProduct.LastUpdated,
		//		Owner:                 mockDataProduct.Owner,
		//		WorkspaceID:           mockDataProduct.WorkspaceID,
		//	},
		//
		//	buildStubs: func(store *mockStore.MockStore) {},
		//
		//	checkResponse: func(recorder *httptest.ResponseRecorder) {
		//		require.Equal(t, http.StatusBadRequest, recorder.Code)
		//	},
		//},
		{
			testScenario: "InternalServerError",

			body: mockDataProduct,

			buildStubs: func(store *mockStore.MockStore) {
				arg := models.DataProduct{
					ProductID:             mockDataProduct.ProductID,
					Name:                  mockDataProduct.Name,
					DataProductGovernance: mockDataProduct.DataProductGovernance,
					DataDomain:            mockDataProduct.DataDomain,
					Description:           mockDataProduct.Description,
					DataProductStatus:     mockDataProduct.DataProductStatus,
					LastUpdated:           mockDataProduct.LastUpdated,
					Owner:                 mockDataProduct.Owner,
					WorkspaceID:           mockDataProduct.WorkspaceID,
				}
				store.EXPECT().CreateDataProduct(arg).Times(1).Return(models.DataProduct{}, sql.ErrConnDone)
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

			body: mockDataProduct,

			buildStubs: func(store *mockStore.MockStore) {
				arg := models.DataProduct{
					ProductID:             mockDataProduct.ProductID,
					Name:                  mockDataProduct.Name,
					DataProductGovernance: mockDataProduct.DataProductGovernance,
					DataDomain:            mockDataProduct.DataDomain,
					Description:           mockDataProduct.Description,
					DataProductStatus:     mockDataProduct.DataProductStatus,
					LastUpdated:           mockDataProduct.LastUpdated,
					Owner:                 mockDataProduct.Owner,
					WorkspaceID:           mockDataProduct.WorkspaceID,
				}
				store.EXPECT().CreateDataProduct(arg).Times(1).Return(mockDataProduct, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   mockDataProduct,
				}
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

			server := test.NewTestServer(test.DATA_PRODUCT, store, nil, authServiceClient)

			url := test.BaseURL + "data-products/"
			expectedResp, err := test.MakeHttpRequest(server, http.MethodPost, url, nil, body)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// TestGetDataProduct tests all the scenarios while getting a Data Products.
func TestGetDataProduct(t *testing.T) {
	mockDataProductView := createRandomDataProductView()
	mockUser := createRandomUserDetails(1122, 1122)
	testCaseSuite := []struct {
		testScenario   string
		productID      string
		buildStubs     func(store *mockStore.MockStore)
		getUserDetails func(client *mock_authservice.MockHttpClient)
		checkResponse  func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "BadRequest_BadUUID",

			getUserDetails: func(client *mock_authservice.MockHttpClient) {},

			productID: "notUuid",

			buildStubs: func(store *mockStore.MockStore) {},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "InternalServerError",

			getUserDetails: func(client *mock_authservice.MockHttpClient) {},

			productID: mockDataProductView.ProductID.String(),

			buildStubs: func(store *mockStore.MockStore) {
				arg := mockDataProductView.ProductID
				store.EXPECT().GetDataProduct(arg).Times(1).Return(models.DataProductView{}, sql.ErrNoRows)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "Success",

			getUserDetails: func(client *mock_authservice.MockHttpClient) {
				test.MockGetUserByID(client, 1122, 1122)
			},

			productID: mockDataProductView.ProductID.String(),

			buildStubs: func(store *mockStore.MockStore) {
				arg := mockDataProductView.ProductID
				store.EXPECT().GetDataProduct(arg).Times(1).Return(mockDataProductView, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data: models.GetDataProductView{
						DataProduct: mockDataProductView,
						Owner:       mockUser.Payload.UserInfo,
					},
				}
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

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)
			testCase.getUserDetails(httpMockClient)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.DATA_PRODUCT, store, nil, authServiceClient)
			url := fmt.Sprintf("%sdata-products/%s/", test.BaseURL, testCase.productID)
			expectedResp, err := test.MakeHttpRequest(server, http.MethodGet, url, nil, nil)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// TestGetAllDataProducts tests all the scenarios while getting a list of Data Products.
func TestGetAllDataProducts(t *testing.T) {
	mockGetAllDataProductsView := []models.GetAllDataProductsView{createRandomGetAllDataProductsView(), createRandomGetAllDataProductsView(), createRandomGetAllDataProductsView()}
	testCaseSuite := []struct {
		testScenario  string
		buildStubs    func(store *mockStore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "InternalServerError",

			buildStubs: func(store *mockStore.MockStore) {
				workspaceID := 1122
				store.EXPECT().GetAllDataProducts(workspaceID).Times(1).Return([]models.GetAllDataProductsView{}, sql.ErrNoRows)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testScenario: "Success",

			buildStubs: func(store *mockStore.MockStore) {
				workspaceID := 1122
				store.EXPECT().GetAllDataProducts(workspaceID).Times(1).Return(mockGetAllDataProductsView, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   mockGetAllDataProductsView}
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

			httpMockClient := mock_authservice.NewMockHttpClient(ctrl)

			authServiceClient := authService.NewClient(httpMockClient)

			server := test.NewTestServer(test.DATA_PRODUCT, store, nil, authServiceClient)
			url := test.BaseURL + "data-products/"
			expectedResp, err := test.MakeHttpRequest(server, http.MethodGet, url, nil, nil)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// TestUpdateDataProduct test all the scenarios while updating a data product.
func TestUpdateDataProduct(t *testing.T) {
	mockDataProduct := createRandomDataProduct()

	testCaseSuite := []struct {
		testScenario  string
		productID     string
		body          models.DataProduct
		buildStubs    func(store *mockStore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "InternalServerError",

			productID: mockDataProduct.ProductID.String(),

			body: mockDataProduct,

			buildStubs: func(store *mockStore.MockStore) {
				arg := models.DataProduct{
					ProductID:             mockDataProduct.ProductID,
					Name:                  mockDataProduct.Name,
					DataProductGovernance: mockDataProduct.DataProductGovernance,
					DataDomain:            mockDataProduct.DataDomain,
					Description:           mockDataProduct.Description,
					DataProductStatus:     mockDataProduct.DataProductStatus,
					LastUpdated:           mockDataProduct.LastUpdated,
					Owner:                 mockDataProduct.Owner,
					WorkspaceID:           mockDataProduct.WorkspaceID,
				}
				store.EXPECT().UpdateDataProduct(arg).Times(1).Return(models.DataProduct{}, sql.ErrConnDone)
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

			productID: mockDataProduct.ProductID.String(),

			body: mockDataProduct,

			buildStubs: func(store *mockStore.MockStore) {
				arg := models.DataProduct{
					ProductID:             mockDataProduct.ProductID,
					Name:                  mockDataProduct.Name,
					DataProductGovernance: mockDataProduct.DataProductGovernance,
					DataDomain:            mockDataProduct.DataDomain,
					Description:           mockDataProduct.Description,
					DataProductStatus:     mockDataProduct.DataProductStatus,
					LastUpdated:           mockDataProduct.LastUpdated,
					Owner:                 mockDataProduct.Owner,
					WorkspaceID:           mockDataProduct.WorkspaceID,
				}
				store.EXPECT().UpdateDataProduct(arg).Times(1).Return(mockDataProduct, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   mockDataProduct,
				}
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

			server := test.NewTestServer(test.DATA_PRODUCT, store, nil, authServiceClient)

			url := fmt.Sprintf("%sdata-products/%s/", test.BaseURL, testCase.productID)
			expectedResp, err := test.MakeHttpRequest(server, http.MethodPut, url, nil, body)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// AddPipeline test all the scenarios while binding a list of pipelines to data-product.
func TestAddPipeline(t *testing.T) {
	mockProductsPipelines := []models.ProductsPipelines{createRandomProductsPipelines(), createRandomProductsPipelines()}
	mockInputPipelines := createRandomInputPipelines()

	testCaseSuite := []struct {
		testScenario  string
		productID     string
		body          models.InputProductsPipelines
		buildStubs    func(store *mockStore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			testScenario: "InternalServerError",

			productID: mockProductsPipelines[0].ProductID.String(),

			body: mockInputPipelines,

			buildStubs: func(store *mockStore.MockStore) {
				productID := mockProductsPipelines[0].ProductID
				arg := []models.ProductsPipelines{
					{
						ProductID:  mockProductsPipelines[0].ProductID,
						PipelineID: mockInputPipelines.Pipelines[0],
					},
					{
						ProductID:  mockProductsPipelines[0].ProductID,
						PipelineID: mockInputPipelines.Pipelines[1],
					},
				}
				store.EXPECT().AddPipeline(arg, productID).Times(1).Return([]models.ProductsPipelines{}, sql.ErrConnDone)
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

			productID: mockProductsPipelines[0].ProductID.String(),

			body: mockInputPipelines,
			buildStubs: func(store *mockStore.MockStore) {
				productID := mockProductsPipelines[0].ProductID
				arg := []models.ProductsPipelines{
					{
						ProductID:  mockProductsPipelines[0].ProductID,
						PipelineID: mockInputPipelines.Pipelines[0],
					},
					{
						ProductID:  mockProductsPipelines[0].ProductID,
						PipelineID: mockInputPipelines.Pipelines[1],
					},
				}

				mockProductsPipelines[0].PipelineID = mockInputPipelines.Pipelines[0]
				mockProductsPipelines[1].PipelineID = mockInputPipelines.Pipelines[1]

				store.EXPECT().AddPipeline(arg, productID).Times(1).Return(mockProductsPipelines, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				res := models.Response{
					Status: utils.SUCCESS,
					Errors: "",
					Data:   mockProductsPipelines,
				}
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

			server := test.NewTestServer(test.DATA_PRODUCT, store, nil, authServiceClient)

			url := fmt.Sprintf("%sdata-products/%s/add-pipeline/", test.BaseURL, testCase.productID)
			expectedResp, err := test.MakeHttpRequest(server, http.MethodPost, url, nil, body)
			require.NoError(t, err)

			testCase.checkResponse(expectedResp)
		})
	}
}

// createRandomUserDetails populates and return the UserDetails Model with random values.
func createRandomUserDetails(uID int, wsID int) models.UserDetails {
	user := models.UserDetails{
		Success: true,
		Payload: models.Payload{UserInfo: models.UserInfo{
			ID:          uID,
			Email:       "ztna-emumba@outlook.com",
			WorkspaceID: wsID,
			FirstName:   "Ztna",
			LastName:    "Admin",
		}},
		Errors:      struct{}{},
		Description: "Hi I am a new User",
	}

	return user
}

//createRandomGetAllDataProductsView populates and return the GetAllDataProductsView Model with random values.
func createRandomGetAllDataProductsView() models.GetAllDataProductsView {
	pID, _ := uuid.NewV1()
	GetAllDataProductsView := models.GetAllDataProductsView{
		ProductID: pID,
		Name:      utils.RandomString(5),
		DataProductGovernance: []string{
			utils.RandomString(5), utils.RandomString(5),
		},
		DataDomain:        utils.RandomString(5),
		Description:       utils.RandomString(10),
		DataProductStatus: "completed",
		LastUpdated:       1647606617639,
		Owner:             1122,
		WorkspaceID:       1122,
		PipelineCount:     1,
	}

	return GetAllDataProductsView
}

//createRandomDataProductView populates and return the DataProductView Model with random values.
func createRandomDataProductView() models.DataProductView {
	pID, _ := uuid.NewV1()
	DataProductView := models.DataProductView{
		ProductID: pID,
		Name:      utils.RandomString(5),
		DataProductGovernance: []string{
			utils.RandomString(5), utils.RandomString(5),
		},
		DataDomain:        utils.RandomString(5),
		Description:       utils.RandomString(10),
		DataProductStatus: "completed",
		LastUpdated:       1647606617639,
		Owner:             1122,
		WorkspaceID:       1122,
		Pipelines: []map[string]interface{}{
			{
				"name": utils.RandomString(5),
			},
			{
				"name": utils.RandomString(5),
			},
		},
	}

	return DataProductView
}

// createRandomDataProduct populates and return the DataProduct Model with random values.
func createRandomDataProduct() models.DataProduct {
	pID, _ := uuid.NewV1()
	randomName := utils.RandomString(5)
	randomDomain := utils.RandomString(5)
	randomDescription := utils.RandomString(10)
	randomDataProductStatus := "completed"

	Dp := models.DataProduct{
		ProductID: pID,
		//Name:      utils.RandomString(5),
		Name: &randomName,
		DataProductGovernance: []string{
			utils.RandomString(5), utils.RandomString(5),
		},
		DataDomain:        &randomDomain,
		Description:       &randomDescription,
		DataProductStatus: &randomDataProductStatus,
		LastUpdated:       1647606617639,
		Owner:             1122,
		WorkspaceID:       1122,
	}

	return Dp
}

// createRandomProductPipelines populates and return the ProductPipelines Model with random values.
func createRandomProductsPipelines() models.ProductsPipelines {
	//productPipelinesID, _ :=uuid.NewV1()
	pipelineID, _ := uuid.NewV1()
	productID, _ := uuid.NewV1()

	return models.ProductsPipelines{
		//ProductsPipelinesID: productPipelinesID,
		ProductID:  productID,
		PipelineID: pipelineID.String(),
	}
}

// createRandomInputPipelines populates and return the ProductPipelines Model with random values.
func createRandomInputPipelines() models.InputProductsPipelines {
	pipelineID0, _ := uuid.NewV1()
	pipelineID1, _ := uuid.NewV1()

	return models.InputProductsPipelines{
		Pipelines: []string{pipelineID0.String(), pipelineID1.String()},
	}
}

// TestMain runs the package level test in TestMode.
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
