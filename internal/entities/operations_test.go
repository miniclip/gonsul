package entities

import (
	. "github.com/onsi/gomega"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestNewOperationsMatrix(t *testing.T) {
	RegisterTestingT(t)

	operation := NewOperationsMatrix()
	rand.Seed(time.Now().UnixNano())
	totalIterations := rand.Intn(20-0) + 0
	inserts := 0
	deletes := 0
	updates := 0

	for i := 0; i < totalIterations; i++ {
		currOp := rand.Intn(3-1) + 1
		switch currOp {
		case 1:
			// inserts
			inserts++
			operation.AddInsert(Entry{
				KVPath: strconv.Itoa(i),
				Value: strconv.Itoa(i),
			})
		case 2:
			// deletes
			deletes++
			operation.AddDelete(Entry{
				KVPath: strconv.Itoa(i),
				Value: strconv.Itoa(i),
			})
		case 3:
			// deletes
			updates++
			operation.AddUpdate(Entry{
				KVPath: strconv.Itoa(i),
				Value: strconv.Itoa(i),
			})
		}
	}

	// Always check
	Expect(operation.GetTotalOps()).To(Equal(totalIterations), "Assert total operations")
	Expect(operation.GetTotalInserts()).To(Equal(inserts), "Assert total deletes")
	Expect(operation.GetTotalDeletes()).To(Equal(deletes), "Assert total deletes")
	Expect(operation.GetTotalUpdates()).To(Equal(updates), "Assert total deletes")

	if totalIterations > 0 {
		Expect(operation.GetOperations()).To(Not(BeNil()), "Assert there are operations")
	} else {
		Expect(operation.GetOperations()).To(BeNil(), "Assert there are no operations")
	}

	if deletes > 0 {
		Expect(operation.HasDeletes()).To(BeTrue(), "Assert there are deletes")
	} else {
		Expect(operation.HasDeletes()).To(BeFalse(), "Assert there are no deletes")
	}
}
