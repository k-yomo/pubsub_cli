package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"os"
	"testing"
)

func Test_newRootCmd(t *testing.T) {
	tests := []struct {
		name   string
		before func()
		check  func(cmd *cobra.Command)
		after  func()
	}{
		{
			name: "project can be set from env variable",
			before: func() {
				os.Setenv("GCP_PROJECT_ID", "test")
			},
			check: func(gotCmd *cobra.Command) {
				gotProject := gotCmd.Flag("project").Value.String()
				if gotProject != "test" {
					t.Errorf("newRootCmd() got project = %v, want %v", gotProject, "test")
					return
				}
			},
			after: func() {
				os.Setenv("GCP_PROJECT_ID", "")
			},
		},
		{
			name: "pubsub emulator host can be set from env variable",
			before: func() {
				os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8080")
			},
			check: func(gotCmd *cobra.Command) {
				gotProject := gotCmd.Flag("host").Value.String()
				if gotProject != "localhost:8080" {
					t.Errorf("newRootCmd() got pubsub emulator host = %v, want %v", gotProject, "localhost:8080")
					return
				}
			},
			after: func() {
				os.Setenv("PUBSUB_EMULATOR_HOST", "")
			},
		},
		{
			name: "gcp credential file path can be set from env variable",
			before: func() {
				os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "cred.json")
			},
			check: func(gotCmd *cobra.Command) {
				gotProject := gotCmd.Flag("cred-file").Value.String()
				if gotProject != "cred.json" {
					t.Errorf("newRootCmd() got gcp credential file path = %v, want %v", gotProject, "cred.json")
					return
				}
			},
			after: func() {
				os.Setenv("GCP_CREDENTIAL_FILE_PATH", "")
			},
		},
		// TODO: test variables can be set from flags
		{
			name:   "publish command is registered",
			before: func() {},
			check: func(gotCmd *cobra.Command) {
				short := newPublishCmd(&bytes.Buffer{}).Short
				for _, cmd := range gotCmd.Commands() {
					if cmd.Short == short {
						return
					}
				}
				t.Errorf("newRootCmd() want %v", short)
			},
			after: func() {},
		},
		{
			name:   "subscribe command is registered",
			before: func() {},
			check: func(gotCmd *cobra.Command) {
				short := newSubscribeCmd(&bytes.Buffer{}).Short
				for _, cmd := range gotCmd.Commands() {
					if cmd.Short == short {
						return
					}
				}
				t.Errorf("newRootCmd() want %v", short)
			},
			after: func() {},
		},
		{
			name:   "register_push command is registered",
			before: func() {},
			check: func(gotCmd *cobra.Command) {
				short := newRegisterPushCmd(&bytes.Buffer{}).Short
				for _, cmd := range gotCmd.Commands() {
					if cmd.Short == short {
						return
					}
				}
				t.Errorf("newRootCmd() want %v", short)
			},
			after: func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
			tt.check(newRootCmd())
			tt.after()
		})
	}
}
