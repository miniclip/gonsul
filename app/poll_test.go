package app

import (
	"github.com/miniclip/gonsul/tests/mocks"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestPoll_RunPoll(t *testing.T) {
	RegisterTestingT(t)

	// Create our mocks and our Once mode
	cfg, log, _, _ := getCommonMocks()
	once := &mocks.Ionce{}

	// Create our assertions
	cfg.On("GetPollInterval").Return(1)
	log.On("PrintInfo", mock.Anything).Return()
	log.On("PrintDebug", mock.Anything).Return()
	once.On("RunOnce").Return()

	poll := getMockedPoll(cfg, log, once)

	// Run our application mode
	poll.RunPoll()

	// Create our expectations
	Expect(cfg.AssertExpectations(t)).To(BeTrue(), "Assert GetPollInterval")
	Expect(cfg.AssertNumberOfCalls(t, "GetPollInterval", 1))

	Expect(log.AssertExpectations(t)).To(BeTrue(), "Assert Logger")

	Expect(once.AssertExpectations(t)).To(BeTrue(), "Assert Once Run")
	Expect(once.AssertNumberOfCalls(t, "RunOnce", 1))
}
