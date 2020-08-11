package main

import (
	"flag"
	"os"

	"github.com/spf13/pflag"
	"k8s.io/klog"
)

type Options struct {
	ShowVersion bool

	APIServerHost  string
	KubeConfigFile string
}

func (options *Options) BindFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&options.ShowVersion, "version", false,
		`Show release information about the event-cwl-exporter and exit.`)
	fs.StringVar(&options.APIServerHost, "apiserver-host", "",
		`Address of the Kubernetes API server.
        Takes the form "protocol://address:port". If not specified, it is assumed the
        program runs inside a Kubernetes cluster and local discovery is attempted.`)
	fs.StringVar(&options.KubeConfigFile, "kubeconfig", "",
		`Path to a kubeconfig file containing authorization and API server information.`)

}

func getOptions() (*Options, error) {

	options := &Options{}
	fs := pflag.NewFlagSet("", pflag.ExitOnError)
	options.BindFlags(fs)

	_ = flag.Set("logtostderr", "true")
	_ = flag.Set("v", "2")
	fs.AddGoFlagSet(flag.CommandLine)

	klogFs := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFs)
	fs.AddGoFlagSet(klogFs)

	_ = fs.Parse(os.Args)

	return options, nil
}
