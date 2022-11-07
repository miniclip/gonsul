package app

import (
	"github.com/miniclip/gonsul/internal/config"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"testing"
)

func TestOnce_RunRead(t *testing.T) {
	RegisterTestingT(t)

	// Create our mocks and our Once mode
	cfg, log, exp, imp := getCommonMocks()
	read := NewRead(cfg, log, imp)
	mode := config.StrategyRead
	// Create our transitive variable (between exporter and importer)
	transitive := map[string]string{"test": "stuff"}

	// Create our assertions
	cfg.On("GetStrategy").Return(mode)
	log.On("PrintInfo", mock.Anything).Return()
	log.On("PrintDebug", mock.Anything).Return()
	exp.On("Start").Return(transitive)
	imp.On("Start", transitive).Return()

	// Run our application mode
	read.RunRead()

	// Create our expectations
	Expect(cfg.AssertExpectations(t)).To(BeTrue(), "Assert GetStrategy")
	Expect(cfg.AssertNumberOfCalls(t, "GetStrategy", 1))

	Expect(exp.AssertExpectations(t)).To(BeTrue(), "Assert Exporter Start")
	Expect(exp.AssertNumberOfCalls(t, "Start", 1))

	Expect(imp.AssertExpectations(t)).To(BeTrue(), "Assert Importer Start")
	Expect(imp.AssertNumberOfCalls(t, "Start", 1))
}
