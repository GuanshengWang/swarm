package cli

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/docker/swarm/scheduler/filter"
	"github.com/docker/swarm/scheduler/strategy"
)

func homepath(p string) string {
	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
	}
	return filepath.Join(home, p)
}

func getDiscovery(c *cli.Context) string {
	if len(c.Args()) == 1 {
		return c.Args()[0]
	}
	return os.Getenv("SWARM_DISCOVERY")
}

var (
	flStore = cli.StringFlag{
		Name:  "rootdir",
		Value: homepath(".swarm"),
		Usage: "",
	}
	flAddr = cli.StringFlag{
		Name:   "addr",
		Value:  "127.0.0.1:2375",
		Usage:  "ip to advertise",
		EnvVar: "SWARM_ADDR",
	}

	// hack for go vet
	flHostsValue = cli.StringSlice([]string{"tcp://127.0.0.1:2375"})

	flHosts = cli.StringSliceFlag{
		Name:   "host, H",
		Value:  &flHostsValue,
		Usage:  "ip/socket to listen on",
		EnvVar: "SWARM_HOST",
	}
	flHeartBeat = cli.IntFlag{
		Name:  "heartbeat, hb",
		Value: 25,
		Usage: "time in second between each heartbeat",
	}
	flEnableCors = cli.BoolFlag{
		Name:  "api-enable-cors, cors",
		Usage: "enable CORS headers in the remote API",
	}
	flTLS = cli.BoolFlag{
		Name:  "tls",
		Usage: "use TLS; implied by --tlsverify=true",
	}
	flTLSCaCert = cli.StringFlag{
		Name:  "tlscacert",
		Usage: "trust only remotes providing a certificate signed by the CA given here",
	}
	flTLSCert = cli.StringFlag{
		Name:  "tlscert",
		Usage: "path to TLS certificate file",
	}
	flTLSKey = cli.StringFlag{
		Name:  "tlskey",
		Usage: "path to TLS key file",
	}
	flTLSVerify = cli.BoolFlag{
		Name:  "tlsverify",
		Usage: "use TLS and verify the remote",
	}
	flStrategy = cli.StringFlag{
		Name:  "strategy",
		Usage: "placement strategy to use [" + strings.Join(strategy.List(), ", ") + "]",
		Value: strategy.List()[0],
	}

	// hack for go vet
	flFilterValue = cli.StringSlice(filter.List())
	// DefaultFilterNumber is exported
	DefaultFilterNumber = len(flFilterValue)

	flFilter = cli.StringSliceFlag{
		Name:  "filter, f",
		Usage: "filter to use [" + strings.Join(filter.List(), ", ") + "]",
		Value: &flFilterValue,
	}

	flCluster = cli.StringFlag{
		Name:  "cluster-driver, c",
		Usage: "cluster driver to use [swarm]",
		Value: "swarm",
	}
	flClusterOpt = cli.StringSliceFlag{
		Name:  "cluster-opt",
		Usage: "cluster driver options",
		Value: &cli.StringSlice{},
	}
)
