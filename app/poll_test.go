package app

import (
	"github.com/miniclip/gonsul/tests/mocks"

	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"testing"
)

func TestPoll_RunPoll(t *testing.T) {
	RegisterTestingT(t)

	// Create our mocks, our Once mode and our application
	cfg, log, _, _ := getCommonMocks()
	once := &mocks.Ionce{}
	poll := getMockedPoll(cfg, log, once)

	// Create our assertions
	cfg.On("GetPollInterval").Return(1)
	log.On("PrintInfo", mock.Anything).Return()
	log.On("PrintDebug", mock.Anything).Return()
	once.On("RunOnce").Return()

	// Run our application mode
	poll.RunPoll()

	// Create our expectations
	Expect(cfg.AssertExpectations(t)).To(BeTrue(), "Assert GetPollInterval")
	Expect(cfg.AssertNumberOfCalls(t, "GetPollInterval", 1))

	Expect(log.AssertExpectations(t)).To(BeTrue(), "Assert Logger")

	Expect(once.AssertExpectations(t)).To(BeTrue(), "Assert Once Run")
	Expect(once.AssertNumberOfCalls(t, "RunOnce", 1))
}
