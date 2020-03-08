package cmd

import "testing"

func setTestRootVariables(t *testing.T) (clear func()) {
	t.Helper()
	projectID = "test"
	emulatorHost = "localhost:8085"
	return func() {
		projectID = ""
		emulatorHost = ""
	}
}
