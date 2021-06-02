package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"testing"
)

func newTestRootCmd(t *testing.T) *cobra.Command {
	t.Helper()

	out := &bytes.Buffer{}
	rootCmd := newRootCmd(out)
	rootCmd.PersistentFlags().Set(projectFlagName, "test")
	rootCmd.PersistentFlags().Set(hostFlagName, "localhost:8085")
	rootCmd.PersistentFlags().Set(createTopicIfNotExistFlagName, "")
	return rootCmd
}
