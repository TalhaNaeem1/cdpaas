package authService

type Workspace struct {
	Id                 int    `json:"id"`
	CompanyName        string `json:"company_name"`
	Country            string `json:"country"`
	WorkspaceUrl       string `json:"workspace_url"`
	CustomerAdminRole  string `json:"customer_admin_role"`
	Industry           string `json:"industry"`
	TeamSize           string `json:"team_size"`
	PhoneNumber        string `json:"phone_number"`
	AirbyteWorkspaceId string `json:"airbyte_workspace_id"`
}
type Role struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type User struct {
	Id                 int           `json:"id"`
	IsActive           bool          `json:"is_active"`
	IsAnonymous        bool          `json:"is_anonymous"`
	IsTenantEndUser    bool          `json:"is_tenant_end_user"`
	IsTenantAdmin      bool          `json:"is_tenant_admin"`
	IsTenantSuperAdmin bool          `json:"is_tenant_super_admin"`
	IsZtnaSuperAdmin   bool          `json:"is_ztna_super_admin"`
	IsAuthenticated    bool          `json:"is_authenticated"`
	IsStaff            bool          `json:"is_staff"`
	IsSuperuser        bool          `json:"is_superuser"`
	Status             string        `json:"status"`
	Email              string        `json:"email"`
	FirstName          string        `json:"first_name"`
	LastName           string        `json:"last_name"`
	IsPasswordSet      bool          `json:"is_password_set"`
	Role               Role          `json:"role"`
	Workspace          Workspace     `json:"workspace"`
	UserPermissions    []interface{} `json:"user_permissions"`
	Groups             []interface{} `json:"groups"`
	SessionId          string        `json:"session_id"`
}

type UserInfo struct {
	Success bool `json:"success"`
	Payload struct {
		User User `json:"_user"`
	} `json:"payload"`
	Errors struct {
	} `json:"errors"`
	Description string `json:"description"`
}
