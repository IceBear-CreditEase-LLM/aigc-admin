package dockerapi

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"path/filepath"
)

//go:generate gowrap gen -g -p ./ -i Service -bt "ce_log:logging.go ce_trace:tracing.go"
type Service interface {
	Update(ctx context.Context, name string, id string, config Config) (newId string, err error)
	Create(ctx context.Context, name string, config Config) (id string, err error)
	Stop(ctx context.Context, id string) (err error)
	Restart(ctx context.Context, id string) (err error)
	Remove(ctx context.Context, id string) (err error)
	Logs(ctx context.Context, id string) (log string, err error)
	Status(ctx context.Context, id string) (status string, err error)
}
type Middleware func(service Service) Service

type service struct {
	workspace string
	dockerCli *client.Client
}

func (s *service) getCli(ctx context.Context) *client.Client {
	return s.dockerCli
}

func (s *service) Update(ctx context.Context, name string, id string, config Config) (newId string, err error) {
	err = s.Stop(ctx, id)
	if err != nil {
		err = fmt.Errorf("stop err: %s", err.Error())
		return
	}
	err = s.Remove(ctx, id)
	if err != nil {
		err = fmt.Errorf("remove err: %s", err.Error())
		return
	}

	return s.Create(ctx, name, config)
}

func (s *service) Create(ctx context.Context, name string, config Config) (id string, err error) {

	exposedPorts := make(nat.PortSet, 0)
	hostPortBindings := make(nat.PortMap, 0)
	hostBinds := make([]string, 0, 0)

	if len(config.ConfigData) > 0 {
		err = config.saveConfigToLocal(name, s.workspace)
		if err != nil {
			err = fmt.Errorf("config.saveConfigToLocal err: %s", err.Error())
			return
		}
	}

	for k, v := range config.Ports {
		exposedPort := fmt.Sprintf("%s/tcp", v)
		exposedPorts[nat.Port(exposedPort)] = struct{}{}
		hostPortBindings[nat.Port(exposedPort)] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: k,
			},
		}
	}

	for _, v := range config.Volumes {
		if config.HasConfigData(v.Key) {
			hostBinds = append(hostBinds, fmt.Sprintf("%s:%s", filepath.Join(s.workspace, name, v.Key), v.Value))
			//hostMount = append(hostMount, mount.Mount{
			//	Type:   mount.TypeBind,
			//	Source: filepath.Join(s.workspace, name, v.Key),
			//	Target: v.Value,
			//})
		} else {
			hostBinds = append(hostBinds, fmt.Sprintf("%s:%s", v.Key, v.Value))

			//hostMount = append(hostMount, mount.Mount{
			//	Type:   mount.TypeBind,
			//	Source: v.Key,
			//	Target: v.Value,
			//})
		}
	}

	ccf := &container.Config{
		Image:        config.Image,
		Cmd:          config.Command,
		Env:          config.EnvVars,
		ExposedPorts: exposedPorts,
	}

	hcf := &container.HostConfig{
		PortBindings: hostPortBindings,
		Resources: container.Resources{
			CPUCount: config.CPU,
			Memory:   config.Memory,
		},
		//Mounts: hostMount,
		Binds: hostBinds,
	}
	var dr []container.DeviceRequest
	if config.GPU > 0 {
		dr = append(dr, container.DeviceRequest{
			Driver:       "nvidia",
			Count:        config.GPU,
			Capabilities: [][]string{{"gpu"}},
		})
	}
	hcf.Resources.DeviceRequests = dr

	resp, err := s.getCli(ctx).ContainerCreate(ctx, ccf, hcf, nil, nil, name)
	if err != nil {
		err = fmt.Errorf("ContainerCreate err: %s", err.Error())
		return
	}

	err = s.getCli(ctx).ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		err = fmt.Errorf("ContainerStart err: %s", err.Error())
		return
	}

	//if len(config.ConfigData) > 0 {

	//var myConfigmapDir = "/myconfigmap/"
	//
	//execConfig := types.ExecConfig{
	//	Cmd:          []string{"/bin/bash", "-c", fmt.Sprintf("mkdir -p %s", myConfigmapDir)},
	//	AttachStdout: true,
	//	AttachStderr: true,
	//}

	//var execIDResp types.IDResponse
	//execIDResp, err = s.getCli(ctx).ContainerExecCreate(ctx, resp.ID, execConfig)
	//if err != nil {
	//	err = fmt.Errorf("%s ContainerExecCreate err: %s", "mkdir", err.Error())
	//	return
	//}

	//execStartCheck := types.ExecStartCheck{}
	//var attachResp types.HijackedResponse
	//attachResp, err = s.getCli(ctx).ContainerExecAttach(ctx, execIDResp.ID, execStartCheck)
	////stdout, stderr := attachResp.Reader, attachResp.Reader
	////outputStdout, outputStderr := io.NopCloser(stdout), io.NopCloser(stderr)
	//b, _ := io.ReadAll(attachResp.Reader)
	//fmt.Println("attachresp: ", string(b))
	//if err != nil {
	//	err = fmt.Errorf("ContainerExecAttach err: %s", err.Error())
	//	return
	//}
	//defer attachResp.Close()

	//var t *bytes.Buffer
	//t, err = config.Tar(name, s.workspace)
	//if err != nil {
	//	err = fmt.Errorf("config.Tar err: %s", err.Error())
	//	return
	//}

	//err = s.getCli(ctx).CopyToContainer(ctx, resp.ID, filepath.Join(myConfigmapDir), t, types.CopyToContainerOptions{
	//	AllowOverwriteDirWithFile: true,
	//	CopyUIDGID:                false,
	//})
	//if err != nil {
	//	err = fmt.Errorf("CopyToContainer err: %s", err.Error())
	//	return
	//}
	//
	//var lnCmd []string
	//
	//for _, v := range config.Volumes {
	//	if config.HasConfigData(v.Key) {
	//		lnCmd = append(lnCmd, fmt.Sprintf(`ln -sf %s %s`, filepath.Join(myConfigmapDir, v.Key), v.Value))
	//	}
	//}
	//
	//cmd := strings.Join(lnCmd, ";")
	//if len(lnCmd) > 0 {
	//	execConfig = types.ExecConfig{
	//		Cmd:          []string{"/bin/bash", "-c", cmd},
	//		AttachStdout: true,
	//		AttachStderr: true,
	//	}
	//
	//	execIDResp, err = s.getCli(ctx).ContainerExecCreate(ctx, resp.ID, execConfig)
	//	if err != nil {
	//		err = fmt.Errorf("%s ContainerExecCreate err: %s", cmd, err.Error())
	//		return
	//	}
	//
	//	execStartCheck = types.ExecStartCheck{}
	//	attachResp, err = s.getCli(ctx).ContainerExecAttach(ctx, execIDResp.ID, execStartCheck)
	//	if err != nil {
	//		err = fmt.Errorf("ContainerExecAttach err: %s", err.Error())
	//		return
	//	}
	//
	//	b, _ = io.ReadAll(attachResp.Reader)
	//	fmt.Println("attachresp: ", string(b))
	//	if err != nil {
	//		err = fmt.Errorf("ContainerExecAttach err: %s", err.Error())
	//		return
	//	}
	//
	//	defer attachResp.Close()
	//}

	//}

	return resp.ID, nil
}

func (s *service) Stop(ctx context.Context, id string) (err error) {
	return s.getCli(ctx).ContainerStop(ctx, id, container.StopOptions{})
}

func (s *service) Restart(ctx context.Context, id string) (err error) {
	return s.getCli(ctx).ContainerRestart(ctx, id, container.StopOptions{})
}

func (s *service) Remove(ctx context.Context, id string) (err error) {
	fmt.Println("stop id", id)
	err = s.Stop(ctx, id)
	if err != nil {
		return
	}
	fmt.Println("remove id", id)
	return s.getCli(ctx).ContainerRemove(ctx, id, types.ContainerRemoveOptions{})
}

func (s *service) Logs(ctx context.Context, id string) (log string, err error) {
	out, err := s.getCli(ctx).ContainerLogs(ctx, id, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
	})
	if err != nil {
		err = fmt.Errorf("ContainerLogs err: %s", err.Error())
		return
	}

	b, err := io.ReadAll(out)
	if err != nil {
		err = fmt.Errorf("io.ReadAll err: %s", err.Error())
		return
	}
	return string(b), nil

}

func (s *service) Status(ctx context.Context, id string) (status string, err error) {
	cJson, err := s.getCli(ctx).ContainerInspect(ctx, id)
	if err != nil {
		err = fmt.Errorf("ContainerInspect err: %s", err.Error())
		return
	}

	return cJson.State.Status, nil
}

func New(workspace string) Service {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return &service{
		workspace: workspace,
		dockerCli: cli,
	}
}
