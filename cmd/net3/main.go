package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bespinian/net3/pkg/net3"
	"github.com/urfave/cli/v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const topoArgsCount = 2
const logArgsCount = 1

func main() {
	kubeconfig := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	if kubeconfig == "" {
		home := homeDir()
		kubeconfig = filepath.Join(home, clientcmd.RecommendedHomeDir, clientcmd.RecommendedFileName)
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(fmt.Errorf("error building k8s config from flags: %w", err))
	}

	k8s, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(fmt.Errorf("error creating k8s client set: %w", err))
	}

	n3 := net3.New(k8s)

	app := &cli.App{
		Name:  "net3",
		Usage: "debug k8s workload network traffic",

		Commands: []*cli.Command{
			{
				Name:    "topo",
				Aliases: []string{"t"},
				Usage:   "display the network topology between a source and a destination",

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "namespace",
						Aliases: []string{"n"},
						Value:   "default",
						Usage:   "the source namespace",
					},
				},

				Action: func(c *cli.Context) error {
					args := c.Args()

					if args.Len() != topoArgsCount {
						return errors.New("usage: net3 topo SOURCE DESTINATION") //nolint:goerr113
					}
					err = n3.Topo(c.String("namespace"), args.Get(0), args.Get(1))
					if err != nil {
						return fmt.Errorf("error creating topo: %w", err)
					}

					return nil
				},
			},
			{
				Name:    "log",
				Aliases: []string{"l"},
				Usage:   "add a request logging proxy to the pods of a service",

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "namespace",
						Aliases: []string{"n"},
						Value:   "default",
						Usage:   "the target namespace",
					},
				},

				Action: func(c *cli.Context) error {
					args := c.Args()

					if args.Len() != logArgsCount {
						return errors.New("usage: net3 log DESTINATION") //nolint:goerr113
					}
					err = n3.Log(c.String("namespace"), args.Get(0))
					if err != nil {
						return fmt.Errorf("error executing log command: %w", err)
					}

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
	h := os.Getenv("HOME")
	if h == "" {
		h = os.Getenv("USERPROFILE") // windows
	}
	return h
}
