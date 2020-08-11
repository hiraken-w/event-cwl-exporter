package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hiraken-w/event-cwl-exporter/internal/controller"
	cloudwatchlogs "github.com/hiraken-w/event-cwl-exporter/internal/output"
	"github.com/hiraken-w/event-cwl-exporter/version"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"k8s.io/sample-controller/pkg/signals"
)

const (
	// High enough QPS to fit all expected use cases. QPS=0 is not set here, because
	// client code is overriding it.
	defaultQPS = 1e6
	// High enough Burst to fit all expected use cases. Burst=0 is not set here, because
	// client code is overriding it.
	defaultBurst = 1e6
)

func main() {
	klog.InitFlags(nil)
	fmt.Println(version.String())
	stopCh := signals.SetupSignalHandler()
	options, err := getOptions()
	if err != nil {
		klog.Fatal(err)
	}
	if options.ShowVersion {
		os.Exit(0)
	}
	config, err := buildRestConfig(options)
	if err != nil {
		klog.Fatal(err)
	}
	clientset, _ := kubernetes.NewForConfig(config)
	informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)
	cwlclient := cloudwatchlogs.NewCloudWatchLogs(getLogGroupName(), getLogStreamName(), getRegionName())
	controller := controller.NewController(clientset,
		informerFactory.Core().V1().Events(),
		cwlclient,
	)
	informerFactory.Start(stopCh)
	if err := controller.Run(1, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}

// buildRestConfig creates a new Kubernetes REST configuration. apiserverHost is
// the URL of the API server in the format protocol://address:port/pathPrefix,
// kubeConfig is the location of a kubeconfig file. If defined, the kubeconfig
// file is loaded first, the URL of the API server read from the file is then
// optionally overridden by the value of apiserverHost.
// If neither apiserverHost nor kubeConfig are passed in, we assume the
// controller runs inside Kubernetes and fallback to the in-cluster config. If
// the in-cluster config is missing or fails, we fallback to the default config.
func buildRestConfig(options *Options) (*rest.Config, error) {
	restCfg, err := clientcmd.BuildConfigFromFlags(options.APIServerHost, options.KubeConfigFile)
	if err != nil {
		return nil, err
	}
	restCfg.QPS = defaultQPS
	restCfg.Burst = defaultBurst
	return restCfg, nil
}

func getLogGroupName() string {
	logGroupName, found := os.LookupEnv("CW_LOG_GROUP_NAME")
	if !found {
		logGroupName = "my-log-group"
	}
	return logGroupName
}

func getLogStreamName() string {
	logStreamName, found := os.LookupEnv("CW_LOG_STREAM_NAME")
	if !found {
		logStreamName = "my-log-stream"
	}
	return logStreamName
}

func getRegionName() string {
	regionName, found := os.LookupEnv("AWS_REGION")
	if !found {
		regionName = "us-west-2"
	}
	return regionName
}
