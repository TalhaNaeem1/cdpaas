package models

type Workspace struct {
	WorkspaceId             string        `json:"workspaceId"`
	CustomerId              string        `json:"customerId"`
	Email                   string        `json:"email"`
	Name                    string        `json:"name"`
	Slug                    string        `json:"slug"`
	InitialSetupComplete    bool          `json:"initialSetupComplete"`
	DisplaySetupWizard      bool          `json:"displaySetupWizard"`
	AnonymousDataCollection bool          `json:"anonymousDataCollection"`
	News                    bool          `json:"news"`
	SecurityUpdates         bool          `json:"securityUpdates"`
	Notifications           []interface{} `json:"notifications"`
	FirstCompletedSync      interface{}   `json:"firstCompletedSync"`
	FeedbackDone            interface{}   `json:"feedbackDone"`
}

type Workspaces struct {
	Workspaces []Workspace `json:"workspaces"`
}

type WorkspaceRequest struct {
	Email                   string        `json:"email"`
	Name                    string        `json:"name"`
}

type WorkspaceAPIResponse struct {
	WorkspaceId             string        `json:"workspaceId"`
}

type WorkspaceResponse struct {
	Status string `json:"status" example:"success"`
	Errors string `json:"errors" example:""`
	Data   struct {
		WorkspaceId string `json:"workspaceId"`
	} `json:"data"`
}