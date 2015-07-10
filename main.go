package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"

	"github.com/trayio/bunny/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/trayio/bunny/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/trayio/bunny/nodes"
	"github.com/trayio/bunny/temporary"
)

const (
	configTemplate = `[{rabbit, [{cluster_nodes, {[{{range $index, $element := .}}{{if $index}}, {{end}}'rabbit@{{.Host}}'{{end}}], disc}}]}].`
)

func main() {
	var (
		rabbits  []*nodes.Node
		region   string
		config   string
		cfg      *aws.Config
		sigChan  = make(chan os.Signal, 1)
		tmpl     = template.Must(template.New("config").Parse(configTemplate))
		tickInit *time.Ticker
		tick     *time.Ticker
	)

	flag.StringVar(&region, "region", "us-west-1", "AWS region")
	flag.StringVar(&config, "destination", "/tmp/rabbitmq-cluster.conf", "Destination for generated configuration")
	flag.Parse()

	signal.Notify(
		sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	cfg = auth(region)

	tickInit = time.NewTicker(time.Second * 2)
	tick = time.NewTicker(time.Minute * 5)

	for {
		select {
		case <-tickInit.C:
			rabbits = nodes.Collect(cfg)
			if len(rabbits) < 2 {
				continue
			}
			if config == "-" {
				tmpl.Execute(os.Stdout, rabbits)
			} else {
				tmpFile, _ := temporary.NewFile()
				tmpl.Execute(tmpFile.File, rabbits)
				tmpFile.Move(config)
				tickInit.Stop()
			}

		case <-tick.C:
			if config == "-" {
				tmpl.Execute(os.Stdout, rabbits)
			} else {
				rabbits = nodes.Collect(cfg)
				tmpFile, _ := temporary.NewFile()
				tmpl.Execute(tmpFile.File, rabbits)
				tmpFile.Move(config)
			}

		case <-sigChan:
			close(sigChan)
			tickInit.Stop()
			tick.Stop()
			return
		}
	}
}

func auth(region string) *aws.Config {
	credentials := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.SharedCredentialsProvider{},
			&credentials.EnvProvider{},
			&credentials.EC2RoleProvider{},
		},
	)

	return &aws.Config{
		Credentials: credentials,
		Region:      region,
	}
}
