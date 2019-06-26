package entities

const OperationInsert = "INSERT"
const OperationUpdate = "UPDATE"
const OperationDelete = "DELETE"
const OperationAll = "ALL"

// Our single operation structure
type operation struct {
	opType string
	entry  Entry
}

// Our operations matrix
type OperationMatrix struct {
	total      int
	inserts    int
	updates    int
	deletes    int
	operations []operation
}

// GetType ...
func (op *operation) GetType() string {
	return op.opType
}

// GetVerb ...
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

// GetPath ...
func (op *operation) GetPath() string {
	return op.entry.KVPath
}

// GetValue ...
func (op *operation) GetValue() string {
	return op.entry.Value
}

// AddInsert ...
func (matrix *OperationMatrix) AddInsert(entry Entry) {
	// Increment our total number of operations
	matrix.total++
	matrix.inserts++
	matrix.operations = append(matrix.operations, operation{opType: OperationInsert, entry: entry})
}

// AddUpdate ...
func (matrix *OperationMatrix) AddUpdate(entry Entry) {
	// Increment our total number of operations
	matrix.total++
	matrix.updates++
	matrix.operations = append(matrix.operations, operation{opType: OperationUpdate, entry: entry})
}

// AddDelete ...
func (matrix *OperationMatrix) AddDelete(entry Entry) {
	// Increment our total number of operations
	matrix.total++
	matrix.deletes++
	matrix.operations = append(matrix.operations, operation{opType: OperationDelete, entry: entry})
}

// HasDeletes ...
func (matrix *OperationMatrix) HasDeletes() bool {
	return matrix.deletes > 0
}

// GetTotalOps ...
func (matrix *OperationMatrix) GetTotalOps() int {
	return matrix.total
}

// GetTotalInserts ...
func (matrix *OperationMatrix) GetTotalInserts() int {
	return matrix.inserts
}

// GetTotalUpdates ...
func (matrix *OperationMatrix) GetTotalUpdates() int {
	return matrix.updates
}

// GetTotalDeletes ...
func (matrix *OperationMatrix) GetTotalDeletes() int {
	return matrix.deletes
}

// GetOperations ...
func (matrix *OperationMatrix) GetOperations() []operation {
	return matrix.operations
}

// NewOperationsMatrix ...
func NewOperationsMatrix() OperationMatrix {
	return OperationMatrix{0, 0, 0, 0, nil}
}
