package integrserv

type OperationResult struct {
	Success bool
	Result  any
	Error   ServiceError
}
