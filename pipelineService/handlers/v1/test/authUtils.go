package test

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	mock_authservice "pipelineService/clients/authService/mocks"
	"pipelineService/env"
	"pipelineService/models/v1"
)

const SessionIdKey = "sessionid"
const SessionIdValue = "6jgxlgfdxrm8z4vnf9evqpqhj6bi14as"
const UserID = "1122"
const WorkspaceID = "1122"
const AirByteWorkspaceID = "a152379e-01a1-11ec-82d6-a312edcd9c7b"

func MockAddAuthorization(request *http.Request) {
	cookie := http.Cookie{
		Name:   SessionIdKey,
		Value:  SessionIdValue,
		Path:   "/",
		Domain: "HttpOnly",
	}
	request.AddCookie(&cookie)
	request.Header.Set("userID", UserID)
	request.Header.Set("workspaceID", WorkspaceID)
	request.Header.Set("airbyteWorkspaceID", AirByteWorkspaceID)
}

func MockGetUserByID(client *mock_authservice.MockHttpClient, userID int, workspaceID int) {
	authServiceURL := fmt.Sprintf("%s/auth-service/api/v1/accounts/internal/user-from-id?user_id=%d", env.Env.AuthServiceAddress, userID)

	res := http.Response{
		Status:     "",
		StatusCode: http.StatusOK,
		Proto:      "",
		ProtoMajor: 0,
		ProtoMinor: 0,
		Header:     nil,
		Body: io.NopCloser(strings.NewReader(fmt.Sprintf(`{
    "success": true,
    "payload": {
        "user": {
            "id": %d,
            "email": "ztna-emumba@outlook.com",
            "workspace_id": %d,
            "first_name": "Ztna",
            "last_name": "Admin"
        }
    },
    "errors": {},
    "description": "User object with given user_id."
}`, userID, workspaceID))),
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}

	client.EXPECT().Get(authServiceURL).Times(1).Return(&res, nil)
}

func CreateRandomUserDetails(uID int, wsID int) models.UserDetails {
	u := models.UserDetails{
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

	return u
}
