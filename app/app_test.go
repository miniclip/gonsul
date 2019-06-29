package app

import (
	"github.com/miniclip/gonsul/internal/config"
	"github.com/miniclip/gonsul/tests/mocks"

	. "github.com/onsi/gomega"
	"os"
	"testing"
)

func getCommonMocks() (cfg *mocks.IConfig, log *mocks.ILogger, exp *mocks.IExporter, imp *mocks.IImporter) {
	cfg = &mocks.IConfig{}
	log = &mocks.ILogger{}
	exp = &mocks.IExporter{}
	imp = &mocks.IImporter{}

	return
}

func getMockedOnce() Ionce {
	cfg, log, exp, imp := getCommonMocks()

	return NewOnce(cfg, log, exp, imp)
}

func getMockedHook() Ihook {
	http := &mocks.IHookHttp{}
	cfg, log, _, _ := getCommonMocks()

	return NewHook(http, cfg, log, getMockedOnce())
}

func getMockedPoll() Ipoll {
	cfg, log, _, _ := getCommonMocks()

	return NewPoll(cfg, log, getMockedOnce(), 1)
}

func TestApplication_Start(t *testing.T) {
	RegisterTestingT(t)

	// Create our table tests
	tests := []struct {Strategy string}{
		{Strategy: "ONCE"},
		{Strategy: "DRYRUN"},
		{Strategy: "POLL"},
		{Strategy: "HOOK"},
		{Strategy: "FAKE"},
	}

	// Create our required channel
	sigChan := make(chan os.Signal)

	for _, test := range tests {
		// Create our mocks for current test
		cfg, _, _, _ := getCommonMocks()
		once := &mocks.Ionce{}
		hook := &mocks.Ihook{}
		poll := &mocks.Ipoll{}

		// Create our application
		application := NewApplication(cfg, once, hook, poll, sigChan)

		// Always assert config GetStrategy
		cfg.On("GetStrategy").Return(test.Strategy)

		// Check current strategy
		switch test.Strategy {
		case config.StrategyDry, config.StrategyOnce:
			// Assert RunOnce
			once.On("RunOnce").Return()
			// Start application
			application.Start()
			// Validate expectations
			Expect(cfg.AssertExpectations(t)).To(BeTrue(), "Assert GetStrategy")
			Expect(cfg.AssertNumberOfCalls(t, "GetStrategy", 1))
			Expect(once.AssertExpectations(t)).To(BeTrue(), "Assert RunOnce")
			Expect(once.AssertNumberOfCalls(t, "RunOnce", 1))
		case config.StrategyHook:
			// Assert RunOnce
			hook.On("RunHook").Return()
			// Start application
			application.Start()
			// Validate expectations
			Expect(cfg.AssertExpectations(t)).To(BeTrue(), "Assert GetStrategy")
			Expect(cfg.AssertNumberOfCalls(t, "GetStrategy", 1))
			Expect(hook.AssertExpectations(t)).To(BeTrue(), "Assert RunHook")
			Expect(hook.AssertNumberOfCalls(t, "RunHook", 1))
		case config.StrategyPoll:
			// Assert RunOnce
			poll.On("RunPoll").Return()
			// Start application
			application.Start()
			// Validate expectations
			Expect(cfg.AssertExpectations(t)).To(BeTrue(), "Assert GetStrategy")
			Expect(cfg.AssertNumberOfCalls(t, "GetStrategy", 1))
			Expect(poll.AssertExpectations(t)).To(BeTrue(), "Assert RunPoll")
			Expect(poll.AssertNumberOfCalls(t, "RunPoll", 1))
		default:
			// Start application (On this test case, we need to make sure none of the applications run)
			application.Start()
			// Validate expectations
			Expect(cfg.AssertExpectations(t)).To(BeTrue(), "Assert GetStrategy")
			Expect(cfg.AssertNumberOfCalls(t, "GetStrategy", 1))
			Expect(once.AssertNumberOfCalls(t, "RunOnce", 0))
			Expect(hook.AssertNumberOfCalls(t, "RunHook", 0))
			Expect(poll.AssertNumberOfCalls(t, "RunPoll", 0))
		}
	}

}
