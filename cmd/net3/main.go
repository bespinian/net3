package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/bespinian/net3/pkg/net3"
	"github.com/urfave/cli/v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	n3 := net3.New(clientset)

	app := &cli.App{
		Name:  "net3",
		Usage: "debug k8s workload network traffic",
		Commands: []*cli.Command{
			{
				Name:    "topo",
				Aliases: []string{"t"},
				Usage:   "display the network topology between a source and a destination",
				Action: func(c *cli.Context) error {
					args := c.Args()
					if args.Len() != 3 {
						return errors.New("needs more args")
					}
					n3.Topo(args.Get(0), args.Get(1))
					return nil
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
