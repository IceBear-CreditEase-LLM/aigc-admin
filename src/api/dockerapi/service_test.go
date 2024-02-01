package dockerapi

import (
	"context"
	"fmt"
	"testing"
)

var workspace = `/mnt/c/Users/37160/Documents/docker/workspace`
var testId = "49252d3057c11817d37e3a9c23d2e2c5a5c753ea712c3f735f267d6596125a78"
var svc = New(workspace)

var configMap = map[string]string{
	"./index.html": "say hello",
}

func Test_service_Create(t *testing.T) {
	type args struct {
		ctx           context.Context
		name          string
		composeConfig Config
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "create",
			args: args{
				ctx:  context.Background(),
				name: "nginx",
				composeConfig: Config{
					Image: "nginx:latest",
					//Command: `/bin/sh -c "sleep 3600"`,
					Ports: map[string]string{
						"80": "80",
					},
					CPU:    1,
					Memory: 500 * 1024 * 1024,
					Volumes: []Volume{
						{
							Key:   "./index.html",
							Value: "/usr/share/nginx/html/index.html",
						},
					},
					EnvVars:    []string{"TZ=Asia/Shanghai"},
					ConfigData: configMap,
				},
			},
			wantOut: "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := svc.Create(tt.args.ctx, tt.args.name, tt.args.composeConfig)
			fmt.Println("create id", id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_service_Down(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "down",
			args: args{
				ctx:  context.Background(),
				name: testId,
			},
			wantOut: "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Remove(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Down() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_service_Logs(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "logs",
			args: args{
				ctx:  context.Background(),
				name: testId,
			},
			wantOut: "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := svc.Logs(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Logs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut != tt.wantOut {
				t.Errorf("Logs() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_service_Restart(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "restart",
			args: args{
				ctx:  context.Background(),
				name: testId,
			},
			wantOut: "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Restart(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Restart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_service_Status(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "status",
			args: args{
				ctx:  context.Background(),
				name: testId,
			},
			wantOut: "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := svc.Status(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Status() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut != tt.wantOut {
				t.Errorf("Status() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_service_Stop(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "stop",
			args: args{
				ctx:  context.Background(),
				name: testId,
			},
			wantOut: "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Stop(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
