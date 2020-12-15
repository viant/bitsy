package resource

type Operation int

const (
	OperationUndefined = -1

	OperationAdded = Operation(iota)
	OperationModified
	OperationDeleted
)
