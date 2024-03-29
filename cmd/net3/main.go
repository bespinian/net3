package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bespinian/net3/internal/net3"
	"github.com/urfave/cli/v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	topoArgsCount      = 2
	proxyAddArgsCount  = 2
	proxyListArgsCount = 0
)

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
				Name:    "proxy",
				Aliases: []string{"p"},
				Usage:   "manage logging proxies",
				Subcommands: []*cli.Command{
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "list logging proxies",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "namespace",
								Aliases: []string{"n"},
								Value:   "default",
								Usage:   "the target namespace",
							},
							&cli.BoolFlag{
								Name:    "all-namespaces",
								Aliases: []string{"i"},
								Value:   false,
								Usage:   "if enabled, all namespaces are considered",
							},
							&cli.StringFlag{
								Name:    "service",
								Aliases: []string{"s"},
								Usage:   "name of service (defaults to all)",
							},
						},
						Action: func(c *cli.Context) error {
							args := c.Args()

							if args.Len() != proxyListArgsCount {
								return errors.New("usage: net3 proxy list") //nolint:goerr113
							}
							err = n3.ListProxies(c.String("namespace"), c.Bool("all-namespaces"), c.String("service"))
							if err != nil {
								return fmt.Errorf("error listing proxies: %w", err)
							}

							return nil
						},
					},
					{
						Name:    "add",
						Aliases: []string{"a"},
						Usage:   "add a logging proxy to the pods of an existing service",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "namespace",
								Aliases: []string{"n"},
								Value:   "default",
								Usage:   "the target namespace",
							},
							&cli.StringFlag{
								Name:    "container-name",
								Aliases: []string{"c"},
								Value:   "net3-proxy",
								Usage:   "name of proxy container",
							},
							&cli.StringFlag{
								Name:    "image",
								Aliases: []string{"i"},
								Value:   "bespinian/net3-proxy",
								Usage:   "image for proxy container",
							},
						},
						Action: func(c *cli.Context) error {
							args := c.Args()

							if args.Len() != proxyAddArgsCount {
								return errors.New("usage: net3 proxy add DESTINATION PORT") //nolint:goerr113
							}
							portInt, convErr := (strconv.Atoi(args.Get(1)))
							if convErr != nil {
								return fmt.Errorf("error converting argument %q to a port number: %w", args.Get(1), err)
							}
							port := int32(portInt)
							err = n3.AddProxy(c.String("namespace"), args.Get(0), c.String("container-name"), c.String("image"), port)
							if err != nil {
								return fmt.Errorf("error adding proxy: %w", err)
							}

							return nil
						},
					},
					{
						Name:    "remove",
						Aliases: []string{"r"},
						Usage:   "remove a logging proxy from the pods of an existing service",
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

							if args.Len() != proxyAddArgsCount {
								return errors.New("usage: net3 proxy remove DESTINATION PORT") //nolint:goerr113
							}
							portInt, convErr := (strconv.Atoi(args.Get(1)))
							if convErr != nil {
								return fmt.Errorf("error converting argument %q to a port number: %w", args.Get(1), err)
							}
							port := int32(portInt)
							err = n3.RemoveProxy(c.String("namespace"), args.Get(0), port)
							if err != nil {
								return fmt.Errorf("error removing proxy: %w", err)
							}

							return nil
						},
					},
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
