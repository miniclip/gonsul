package structs

const OperationInsert	= "INSERT"
const OperationUpdate	= "UPDATE"
const OperationDelete	= "DELETE"
const OperationAll		= "ALL"

// A consul GET response
type ConsulResult struct {
	LockIndex		int			`json:"lockIndex"`
	Key				string		`json:"Key"`
	Flags			int			`json:"Flags"`
	Value			string		`json:"Value"`
	CreateIndex		int			`json:"CreateIndex"`
	ModifyIndex		int			`json:"ModifyIndex"`
}

// A consul Transaction payload
type ConsulTxn struct {
	KV				ConsulTxnKV	`json:"KV"`
}

// A consul KV Transaction payload
type ConsulTxnKV struct {
	Verb			*string		`json:"Verb"`
	Key				*string		`json:"Key"`
	Value			*string		`json:"Value,omitempty"`
}

// Our single operation structure
type operation struct {
	opType 	string
	entry   Entry
}

func (op *operation) GetType() string {
	return op.opType
}

func (op *operation) GetVerb() string {
	switch op.opType {
	case OperationInsert:
		return "set"
	case OperationUpdate:
		return "set"
	case OperationDelete:
		return "delete"
	}

	return "get"
}

func (op *operation) GetPath() string {
	return op.entry.KVPath
}

func (op *operation) GetValue() string {
	return op.entry.Value
}

// Our operations matrix
type OperationMatrix struct {
	total			int
	inserts			int
	updates			int
	deletes			int
	operations		[]operation
}

func (matrix *OperationMatrix) AddInsert(entry Entry) {
	// Increment our total number of operations
	matrix.total++
	matrix.inserts++
	matrix.operations = append(matrix.operations, operation{opType: OperationInsert, entry: entry})
}

func (matrix *OperationMatrix) AddUpdate(entry Entry) {
	// Increment our total number of operations
	matrix.total++
	matrix.updates++
	matrix.operations = append(matrix.operations, operation{opType: OperationUpdate, entry: entry})
}

func (matrix *OperationMatrix) AddDelete(entry Entry) {
	// Increment our total number of operations
	matrix.total++
	matrix.deletes++
	matrix.operations = append(matrix.operations, operation{opType: OperationDelete, entry: entry})
}

func (matrix *OperationMatrix) HasDeletes() bool {
	return matrix.deletes > 0
}

func (matrix *OperationMatrix) GetTotalOps() int {
	return matrix.total
}

func (matrix *OperationMatrix) GetTotalInserts() int {
	return matrix.inserts
}

func (matrix *OperationMatrix) GetTotalUpdates() int {
	return matrix.updates
}

func (matrix *OperationMatrix) GetTotalDeletes() int {
	return matrix.deletes
}

func (matrix *OperationMatrix) GetOperations() []operation {
	return matrix.operations
}

func NewOperationsMatrix() OperationMatrix {
	return OperationMatrix{0, 0, 0, 0, nil}
}