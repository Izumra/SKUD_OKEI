package integrserv

type ServiceError struct {
	ErrorCode             string
	Description           string
	InnerExceptionMessage string
}
