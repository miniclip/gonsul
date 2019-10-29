package app

import (
	"github.com/miniclip/gonsul/tests/mocks"

	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	
	"testing"
)

func TestHook_RunHook(t *testing.T) {
	RegisterTestingT(t)

	// Create our mocks and our Once mode
	cfg, log, _, _ := getCommonMocks()
	http := &mocks.IHookHttp{}
	once := &mocks.Ionce{}
	hook := getMockedHook(http, cfg, log, once)

	// Create our assertions
	http.On("Start", mock.Anything, mock.Anything).Return()
	log.On("PrintInfo", mock.Anything).Return()

	// Run our application mode
	hook.RunHook()

	// Create our expectations
	Expect(http.AssertExpectations(t)).To(BeTrue(), "Assert Http.Start")
	Expect(http.AssertNumberOfCalls(t, "Start", 1))
	Expect(log.AssertExpectations(t)).To(BeTrue(), "Assert Logger")
}
