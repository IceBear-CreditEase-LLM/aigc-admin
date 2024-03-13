package runtime

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type k8s struct {
	k8sClient     *kubernetes.Clientset
	k8sConfig     *rest.Config
	createOptions CreationOptions
}

// NewK8s 创建k8s
func NewK8s(opts ...CreationOption) (Service, error) {
	createOptions := CreationOptions{
		namespace: "default",
		shmSize:   "16Gi",
	}

	for _, opt := range opts {
		opt(&createOptions)
	}

	var k8sClient *kubernetes.Clientset
	var k8sConfig *rest.Config
	var err error

	if createOptions.k8sConfigPath != "" {

		configLoad := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: createOptions.k8sConfigPath},
			&clientcmd.ConfigOverrides{},
		)

		k8sConfig, err := configLoad.ClientConfig()
		if err != nil {
			return nil, err
		}

		k8sClient, err = kubernetes.NewForConfig(k8sConfig)
		if err != nil {
			return nil, err
		}
	}

	if createOptions.k8sTokenModel.host != "" && createOptions.k8sTokenModel.token != "" {
		k8sConfig = &rest.Config{
			Host:        createOptions.k8sTokenModel.host,
			BearerToken: createOptions.k8sTokenModel.token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: createOptions.k8sTokenModel.insecure,
			},
		}

		k8sClient, err = kubernetes.NewForConfig(k8sConfig)
		if err != nil {
			return nil, err
		}
	}

	return &k8s{
		k8sClient:     k8sClient,
		k8sConfig:     k8sConfig,
		createOptions: createOptions,
	}, nil
}

// CreateJob 创建job
func (s *k8s) CreateJob(ctx context.Context, config Config) (jobName string, err error) {
	if config.ShmSize == "" {
		config.ShmSize = s.createOptions.shmSize
	}

	config.namespace = s.createOptions.namespace

	if len(config.ConfigData) > 0 {
		configmap := config.GenConfigMap()

		_, err = s.k8sClient.CoreV1().ConfigMaps(config.namespace).Create(ctx, &configmap, v1.CreateOptions{})
		if err != nil {
			return "", err
		}
	}

	job, err := config.GenJob()
	if err != nil {
		return "", err
	}

	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&job)
	if err != nil {
		panic(err)
	}
	jobB, _ := yaml.Marshal(unstructuredObj)
	fmt.Println(string(jobB))

	_, err = s.k8sClient.BatchV1().Jobs(config.namespace).Create(ctx, &job, v1.CreateOptions{})
	if err != nil {
		return "", err
	}
	return config.ServiceName, nil
}

// CreateDeployment 创建deployment
func (s *k8s) CreateDeployment(ctx context.Context, config Config) (deploymentName string, err error) {
	if config.ShmSize == "" {
		config.ShmSize = s.createOptions.shmSize
	}

	config.namespace = s.createOptions.namespace

	if len(config.ConfigData) > 0 {
		configmap := config.GenConfigMap()

		_, err = s.k8sClient.CoreV1().ConfigMaps(config.namespace).Create(ctx, &configmap, v1.CreateOptions{})
		if err != nil {
			return "", err
		}
	}

	deployment, err := config.GenDeployment()
	if err != nil {
		return "", err
	}

	_, err = s.k8sClient.AppsV1().Deployments(config.namespace).Create(ctx, &deployment, v1.CreateOptions{})
	if err != nil {
		return "", err
	}

	return config.ServiceName, nil
}

// GetDeploymentLogs 获取部署的日志
func (s *k8s) GetDeploymentLogs(ctx context.Context, deploymentName string) (log string, err error) {
	config := Config{
		ServiceName: deploymentName,
		namespace:   s.createOptions.namespace,
	}

	pods, err := s.k8sClient.CoreV1().Pods(s.createOptions.namespace).List(ctx, v1.ListOptions{
		LabelSelector: v1.FormatLabelSelector(&v1.LabelSelector{
			MatchLabels: config.GenDeploymentLabels(),
		}),
	})

	if err != nil {
		err = errors.Wrap(err, "get pods error")
		return "", err
	}

	if len(pods.Items) == 0 {
		return "", errors.New("no pod found")
	}

	podLogOpts := corev1.PodLogOptions{
		Container:  pods.Items[0].Spec.Containers[0].Name,
		Follow:     false,
		Previous:   false,
		Timestamps: true,
		//LimitBytes: &byteReadLimit,
		//TailLines:  &lineReadLimit,
	}
	logs := s.k8sClient.CoreV1().Pods(config.namespace).GetLogs(pods.Items[0].Name, &podLogOpts)
	podLogs, err := logs.Stream(ctx)
	if err != nil {
		err = errors.Wrap(err, "Stream")
		return
	}
	defer podLogs.Close()

	result, err := io.ReadAll(podLogs)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// GetJobLogs 获取job的日志
func (s *k8s) GetJobLogs(ctx context.Context, jobName string) (log string, err error) {
	config := Config{
		ServiceName: jobName,
		namespace:   s.createOptions.namespace,
	}
	pods, err := s.k8sClient.CoreV1().Pods(s.createOptions.namespace).List(ctx, v1.ListOptions{
		LabelSelector: v1.FormatLabelSelector(&v1.LabelSelector{
			MatchLabels: config.GenJobLabels(),
		}),
	})

	if err != nil {
		err = errors.Wrap(err, "GetPods")
		return
	}
	if len(pods.Items) == 0 {
		err = fmt.Errorf("pod not found")
		return
	}

	//var lineReadLimit int64 = 5000
	//var byteReadLimit int64 = 500000
	podLogOpts := corev1.PodLogOptions{
		Container:  pods.Items[0].Spec.Containers[0].Name,
		Follow:     false,
		Previous:   false,
		Timestamps: true,
		//LimitBytes: &byteReadLimit,
		//TailLines:  &lineReadLimit,
	}
	logs := s.k8sClient.CoreV1().Pods(config.namespace).GetLogs(pods.Items[0].Name, &podLogOpts)
	podLogs, err := logs.Stream(ctx)
	if err != nil {
		err = errors.Wrap(err, "Stream")
		return
	}
	defer podLogs.Close()

	result, err := io.ReadAll(podLogs)
	if err != nil {
		return "", err
	}
	return string(result), nil

}

// GetJobStatus 获取job的状态
func (s *k8s) GetJobStatus(ctx context.Context, jobName string) (status string, err error) {
	config := Config{
		ServiceName: jobName,
		namespace:   s.createOptions.namespace,
	}

	pods, err := s.k8sClient.CoreV1().Pods(config.namespace).List(ctx, v1.ListOptions{
		LabelSelector: v1.FormatLabelSelector(&v1.LabelSelector{
			MatchLabels: config.GenJobLabels(),
		}),
	})

	if err != nil {
		err = errors.Wrap(err, "GetPods")
		return
	}
	if len(pods.Items) == 0 {
		err = fmt.Errorf("pod not found")
		return
	}

	status = string(pods.Items[0].Status.Phase)
	return
}

// GetDeploymentStatus 获取部署的状态
func (s *k8s) GetDeploymentStatus(ctx context.Context, deploymentName string) (status string, err error) {
	config := Config{
		ServiceName: deploymentName,
		namespace:   s.createOptions.namespace,
	}

	pods, err := s.k8sClient.CoreV1().Pods(config.namespace).List(ctx, v1.ListOptions{
		LabelSelector: v1.FormatLabelSelector(&v1.LabelSelector{
			MatchLabels: config.GenDeploymentLabels(),
		}),
	})

	if err != nil {
		err = errors.Wrap(err, "GetPods")
		return
	}
	if len(pods.Items) == 0 {
		err = fmt.Errorf("pod not found")
		return
	}

	status = string(pods.Items[0].Status.Phase)
	return
}

// RemoveJob 删除job
func (s *k8s) RemoveJob(ctx context.Context, jobName string) (err error) {
	deletePolicy := v1.DeletePropagationForeground
	err = s.k8sClient.BatchV1().Jobs(s.createOptions.namespace).Delete(ctx, jobName, v1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		err = errors.Wrap(err, "Job Delete")
		return
	}

	configmap, err := s.k8sClient.CoreV1().ConfigMaps(s.createOptions.namespace).Get(ctx, jobName, v1.GetOptions{})
	if err != nil {
		if k8serr.IsNotFound(err) {
			return
		} else {
			return err
		}
	}

	err = s.k8sClient.CoreV1().ConfigMaps(s.createOptions.namespace).Delete(ctx, configmap.Name, v1.DeleteOptions{})
	if err != nil {
		err = errors.Wrap(err, "ConfigMap Delete")
		return

	}

	return
}

// RemoveDeployment 删除部署
func (s *k8s) RemoveDeployment(ctx context.Context, deploymentName string) (err error) {
	err = s.k8sClient.AppsV1().Deployments(s.createOptions.namespace).Delete(ctx, deploymentName, v1.DeleteOptions{})
	if err != nil {
		err = errors.Wrap(err, "delete deployment")
		return
	}

	configmap, err := s.k8sClient.CoreV1().ConfigMaps(s.createOptions.namespace).Get(ctx, deploymentName, v1.GetOptions{})
	if err != nil {
		if k8serr.IsNotFound(err) {
			return
		} else {
			return err
		}
	}

	err = s.k8sClient.CoreV1().ConfigMaps(s.createOptions.namespace).Delete(ctx, configmap.Name, v1.DeleteOptions{})
	if err != nil {
		err = errors.Wrap(err, "ConfigMap Delete")
		return

	}

	return
}