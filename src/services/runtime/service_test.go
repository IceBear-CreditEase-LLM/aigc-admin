package runtime

import (
	"context"
	"fmt"
	"testing"
)

// var s Service = NewDocker()
// var id string = ""

// var s, err = NewK8s(WithK8sConfigPath("./k8sconfig.yaml"), WithNamespace("dev"))
var token string = `eyJhbGciOiJSUzI1NiIsImtpZCI6InBHSUo1VXhzUzI0VFM1OHY1TFZtQVlKVWdsbUNxQmp5R2pTZ3NNQ2pjaEEifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiQ.rQ-a28Qz1HxPZYMxYUjB51zqAK8M3hs8hvIKbh30r3Z2FpRjIBCqY1I1EKyHXsaTHs4qMD87QviT9v_Ffz8X-DM7xDNw`
var s, err = NewK8s(WithK8sToken("https://localhost:6443", token, true), WithNamespace("dev"))
var id string = "dataset-similar-task-1"
var image = "nginx"

func TestService_CreateJob(t *testing.T) {
	if err != nil {
		t.Error(err.Error())
		return
	}

	jobName, err := s.CreateDeployment(context.Background(), Config{
		ServiceName: fmt.Sprintf("dataset-similar-task-%d", 1),
		Image:       image,
		Cpu:         0,
		Memory:      0,
		GPU:         0,
		// Command: []string{
		// "/bin/bash",
		// "/app/dataset_analyze_similar.sh",
		// },
		Command: []string{
			"/bin/bash",
			"-c",
			"echo start sleep ;sleep 10000000;",
		},
		EnvVars: []string{
			"TZ=Asia/Shanghai",
		},
		ConfigData: map[string]string{
			"/app/dataset.json": "{\"instruction\":\"asdfasdfasdf\",\"input\":\"\",\"output\":\"ffff\",\"intent\":\"ffff\",\"document\":\"\",\"question\":\"fadgawevansdf\"}\n",
			"/app/hello":        "say hello",
		},
		Replicas: 1,
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(jobName)
	id = jobName
}

func TestService_GetJobStatus(t *testing.T) {
	status, err := s.GetDeploymentStatus(context.Background(), id)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(status)
}

func TestService_GetJobLog(t *testing.T) {
	log, err := s.GetDeploymentLogs(context.Background(), id)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(log)
}

// func TestService_RemoveJob(t *testing.T) {
// 	err := s.RemoveJob(context.Background(), id)
// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// }
