package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"tutorial_template/steps"

	"github.com/buildium-org/codium-harness/env"
	"github.com/buildium-org/codium-harness/logger"
	"github.com/buildium-org/codium-harness/meta"
	testrunners "github.com/buildium-org/codium-harness/testrunners"
)

var cliSteps = []func(config *testrunners.CliTestConfig) error{
	steps.ExampleCliStep,
}

var serverSteps = []func(config *testrunners.ServerTestConfig) error{
	steps.ExampleServerStep,
}

func main() {
	m := meta.NewFromJson(env.GetMetaJson())
	executable := m.GetExecutablePath()

	fmt.Println("Running up to stage: ", m.Stage)

	ctx := context.Background()
	log := logger.NewLogger()

	serverRunner := testrunners.NewServerRunner(m, serverSteps, testrunners.NewTestServer(executable, log))
	testrunners.Run(ctx, serverRunner)

	cliRunner := testrunners.NewCliRunner(m, cliSteps, executable)
	testrunners.Run(ctx, cliRunner)
}

func resolveExecutable(m *meta.Meta) string {
	executableDir := m.ExecutableDir
	if executableDir == "" {
		executableDir = os.Getenv("CLIENT_DIR")
		if executableDir == "" {
			executableDir = "/app/bin"
		}
	}
	return filepath.Join(executableDir, m.Entrypoint)
}
