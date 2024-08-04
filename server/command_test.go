package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseCommand(t *testing.T) {
	type args struct {
		cmd string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		cmd     *Command
	}{
		{
			name: "valid SET",
			args: args{
				cmd: "*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue\r\n",
			},
			wantErr: false,
			cmd: &Command{
				Type:  SET,
				Key:   "mykey",
				Value: "myvalue",
			},
		},

		{
			name: "valid GET",
			args: args{
				cmd: "*2\r\n$3\r\nGET\r\n$5\r\nmykey\r\n",
			},
			wantErr: false,
			cmd: &Command{
				Type:  GET,
				Key:   "mykey",
				Value: "",
			},
		},

		{
			name: "valid PING",
			args: args{
				cmd: "*1\r\n$4\r\nPING\r\n",
			},
			wantErr: false,
			cmd: &Command{
				Type:  GET,
				Key:   "mykey",
				Value: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := parseCommand(tt.args.cmd)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseArray() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, res, tt.cmd)
		})
	}
}
