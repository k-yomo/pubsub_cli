package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"testing"
)

func Test_subscribe(t *testing.T) {
	rootCmd := newTestRootCmd(t)

	type args struct {
		rootCmd *cobra.Command
		args    []string
	}

	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name:    "subscribe topic with invalid name causes error",
			args:    args{rootCmd: rootCmd, args: []string{"subscribe", "1"}},
			wantErr: true,
		},
		{
			name: "parent cmd without projectFlag causes error",
			args: args{rootCmd: func() *cobra.Command {
				cmd := &cobra.Command{}
				cmd.PersistentFlags().String(hostFlagName, "host", "")
				cmd.PersistentFlags().String(credFileFlagName, "cred.json", "")
				return cmd
			}(), args: []string{"publish", "test_topic", "hello"}},
			wantErr: true,
		},
		{
			name: "parent cmd without hostFlag causes error",
			args: args{rootCmd: func() *cobra.Command {
				cmd := &cobra.Command{}
				cmd.PersistentFlags().String(projectFlagName, "project", "")
				cmd.PersistentFlags().String(credFileFlagName, "cred.json", "")
				return cmd
			}(), args: []string{"publish", "test_topic", "hello"}},
			wantErr: true,
		},
		{
			name: "parent cmd without credFileFlag causes error",
			args: args{rootCmd: func() *cobra.Command {
				cmd := &cobra.Command{}
				cmd.PersistentFlags().String(projectFlagName, "project", "")
				cmd.PersistentFlags().String(hostFlagName, "host", "")
				return cmd
			}(), args: []string{"publish", "test_topic", "hello"}},
			wantErr: true,
		},
		// TODO: test regular cases
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cmd := newSubscribeCmd(out)
			tt.args.rootCmd.SetArgs(tt.args.args)
			tt.args.rootCmd.AddCommand(cmd)

			err := tt.args.rootCmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("subscribe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("subscribe() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
