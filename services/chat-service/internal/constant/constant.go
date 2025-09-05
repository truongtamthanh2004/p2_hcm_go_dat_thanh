package constant

const (
	ErrFailedToSaveMessage     = "error.failed_to_save_message"
	ErrFailedToGetConversation = "error.failed_to_get_conversation"
	ErrInvalidUserID           = "error.invalid_user_id"
	ErrUserIDRequired          = "error.user_id_required"
	ErrUpgradeFailed           = "error.failed_to_upgrade_to_websocket"
	ErrSameUserConversation    = "cannot get conversation between the same user"
	ErrMarshalRequest          = "failed to marshal request body"
	ErrCreateHTTPRequest       = "failed to create HTTP request"
	ErrSendHTTPRequest         = "failed to send HTTP request"
	ErrUnmarshalResponse       = "failed to unmarshal response body"
	ErrInternalServer          = "internal server error"
	ErrUserNotFound            = "user not found"
	ErrUnauthorized            = "unauthorized"
)

const (
	SuccessGetConversation = "success.get_conversation"
)

const (
	GetUserUrl = "/api/v1/users/%d"
)
