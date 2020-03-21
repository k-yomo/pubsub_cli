package cmd

import (
	"bytes"
	"github.com/k-yomo/pubsub_cli/pkg"
	"github.com/spf13/cobra"
	"testing"
)

func Test_subscribe(t *testing.T) {
	pubsubClient, err := pkg.NewTestPubSubClient(t)
	if err != nil {
		t.Fatal(err)
	}
	clear := setTestRootVariables(t)
	defer clear()
	type args struct {
		in0          *cobra.Command
		pubsubClient *pkg.PubSubClient
		args         []string
	}

	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name:    "subscribe topic with invalid name causes error",
			args:    args{pubsubClient: pubsubClient, args: []string{"1"}},
			wantErr: true,
		},
		// TODO: test regular cases
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			err := newSubscribeCmd(out).RunE(tt.args.in0, tt.args.args)
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