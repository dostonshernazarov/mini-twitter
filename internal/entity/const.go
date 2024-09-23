package entity

// role for users
const (
	RoleAdmin        = "admin"
	RoleUser         = "user"
	RoleUnauthorized = "unauthorized"
	RoleUnknown      = "unknown"
)

// error messages
const (
	NoAccess           string = "You have no access"
	RequireRefresh     string = "Require refresh"
	NotFoundData       string = "Data was not found"
	ServerError        string = "Server error"
	IncorrectData      string = "Data was incorrect"
	EmailUsed          string = "Email already used"
	UsernameTaken      string = "Username already taken"
	SendOTP            string = "We have send verification code your email"
	OTPExpired         string = "OTP Expired"
	TokenExpired       string = "Token Expired"
	WrongLoginOrPasswd string = "Wrong login or password"
	UploadingError     string = "Error happened while upload files"
)
