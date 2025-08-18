package constant

const (
	ErrInvalidInput                 = "invalid input data"
	ErrEmailAlreadyExists           = "email already exists"
	ErrPasswordHash                 = "failed to hash password"
	ErrCreateAuthUser               = "failed to create auth user"
	ErrCreateUserProfile            = "failed to create user profile"
	ErrPublishEvent                 = "failed to publish verification event"
	ErrInternalServer               = "internal server error"
	ErrWeakPassword                 = "password must be at least 6 characters"
	ErrGenerateToken                = "failed to generate verification token"
	ErrMarshalRequest               = "unable to process request data"
	ErrCreateHTTPRequest            = "unable to create request to server"
	ErrSendHTTPRequest              = "unable to connect to the server"
	ErrUnmarshalResponse            = "unable to read response from the server"
	ErrUserNotFound                 = "User not found"
	ErrAlreadyVerified              = "Email already verified"
	ErrTokenRequired                = "error.token_required"
	ErrInvalidToken                 = "error.invalid_token"
	ErrGetUserFailed                = "error.get_user_failed"
	ErrUserAlreadyVerified          = "error.user_already_verified"
	ErrFailedToUpdateUser           = "error.failed_to_update_user"
	ErrInvalidCredentials           = "error.invalid_credentials"
	ErrGenerateTokenFailed          = "error.generate_token_failed"
	ErrInvalidRequest               = "error.invalid_request"
	ErrExpiredOrInvalidRefreshToken = "error.expired_or_invalid_refresh_token"
	ErrInvalidUserRefreshToken      = "error.invalid_user_refresh_token"
	ErrUserNotActive                = "error.user_not_active"
	ErrUserNotVerified              = "error.user_not_verified"
	ErrStrongPassword               = "Password must be at least 8 characters and contain upper, lower, number, and special character"
)

const (
	SuccessAccountVerified = "success.account_verified_successfully"
	SuccessSignUp          = "Sign up successful. Please verify your email"
	SuccessLogin           = "success.login"
	SuccessRefreshToken    = "success.refresh_token"
)

const (
	USER_ROLE     = "user"
	CreateUserUrl = "/api/v1/users/"
)

const (
	EventTypeVerifyEmail = "VERIFY_EMAIL"
)
