package main

import (
	"io/ioutil"

	"github.com/codegangsta/cli"
	"github.com/kubernetes/deployment-manager/common"
	"gopkg.in/yaml.v2"
)

func init() {
	addCommands(deployCmd())
}

func deployCmd() cli.Command {
	return cli.Command{
		Name:      "deploy",
		Aliases:   []string{"install"},
		Usage:     "Deploy a chart into the cluster.",
		ArgsUsage: "[CHART]",
		Action:    func(c *cli.Context) { run(c, deploy) },
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "config,c",
				Usage: "The configuration YAML file for this deployment.",
			},
			cli.StringFlag{
				Name:  "name",
				Usage: "Name of deployment, used for deploy and update commands (defaults to template name)",
			},
			// TODO: I think there is a Generic flag type that we can implement parsing with.
			cli.StringFlag{
				Name:  "properties,p",
				Usage: "A comma-separated list of key=value pairs: 'foo=bar,foo2=baz'.",
			},
		},
	}
}

func deploy(c *cli.Context) error {

	// If there is a configuration file, use it.
	cfg := &common.Configuration{}
	if c.String("config") != "" {
		if err := loadConfig(cfg, c.String("config")); err != nil {
			return err
		}
	} else {
		cfg.Resources = []*common.Resource{
			{
				Properties: map[string]interface{}{},
			},
		}
	}

	// If there is a chart specified on the commandline, override the config
	// file with it.
	args := c.Args()
	if len(args) > 0 {
		cfg.Resources[0].Type = args[0]
	}

	// Override the name if one is passed in.
	if name := c.String("name"); len(name) > 0 {
		cfg.Resources[0].Name = name
	}

	if props, err := parseProperties(c.String("properties")); err != nil {
		return err
	} else if len(props) > 0 {
		// Coalesce the properties into the first props. We have no way of
		// knowing which resource the properties are supposed to be part
		// of.
		for n, v := range props {
			cfg.Resources[0].Properties[n] = v
		}
	}

	return client(c).PostDeployment(cfg)
}

// loadConfig loads a file into a common.Configuration.
func loadConfig(c *common.Configuration, filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, c)
}
