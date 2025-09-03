package constant

const (
	ErrInvalidInput          = "invalid input data"
	ErrEmailAlreadyExists    = "email already exists"
	ErrDatabase              = "database operation failed"
	ErrInternalServer        = "internal server error"
	ErrCreateUser            = "failed to create user"
	ErrWeakPassword          = "password must be at least 8 characters"
	ErrUserNotFound          = "error.user_not_found"
	ErrInvalidRequest        = "error.invalid_request"
	ErrUpdateFailed          = "error.update_failed"
	ErrUnauthorized          = "error.unauthorized"
	ErrInvalidUserID         = "error.invalid_user_id"
	ErrNameRequired          = "error.name_required"
	ErrInvalidPhoneNumber    = "error.invalid_phone_number"
	ErrInvalidEmailType      = "error.invalid_email_type"
	ErrFailedToFetchUserList = "error.failed_to_fetch_user_list"
	ErrInvalidPageParameter  = "Invalid page parameter. Must be a positive integer"
	ErrInvalidLimitParameter = "Invalid limit parameter. Must be a positive integer between 1 and 100"
	ErrMarshalRequest        = "unable to process request data"
	ErrCreateHTTPRequest     = "unable to create request to server"
	ErrSendHTTPRequest       = "unable to connect to the server"
	ErrUnmarshalResponse     = "unable to read response from the server"
	ErrInvalidRole           = "invalid role"
)

const (
	UpdateAuthUserURL = "api/v1/auth/users"
)

const (
	RoleUser      = "user"
	RoleModerator = "moderator"
	RoleAdmin     = "admin"
)
