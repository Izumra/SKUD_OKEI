package resp

type SuccessAuth struct {
	Username  string
	SessionId string `json:"-"`
}
