package model

type OperationType struct {
	OperationTypeID int  `json:"operation_type_id"`
	IsCredit        bool `json:"is_credit"`
}
