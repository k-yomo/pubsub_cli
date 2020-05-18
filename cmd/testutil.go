package cmd

import (
	"github.com/spf13/cobra"
	"testing"
)

func newTestRootCmd(t *testing.T) *cobra.Command {
	t.Helper()

	rootCmd := newRootCmd()
	rootCmd.PersistentFlags().Set("project", "test")
	rootCmd.PersistentFlags().Set("host", "localhost:8085")
	return rootCmd
}
