package pkg

import (
	"flag"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

type config struct {
	homeDir    string
	kubeConfig *string
	frequency  int
	setup      bool
	AskVersion bool
	AskHelp    bool
}

// Check if configuration is initialized
func (c *config) Initialized() bool {
	return c != nil && c.setup
}

var CONFIG config

// Get configuration
func GetConfig() config {
	if CONFIG.Initialized() {
		return CONFIG
	}

	c := config{}

	// Home directory
	// FIXME replace by HomeDir() // k8s.io/client-go/kubernetes/util
	//c.homeDir = os.Getenv("USERPROFILE")
	//if h := os.Getenv("HOME"); h != "" {
	//	c.homeDir = h
	//}

	c.homeDir = homedir.HomeDir()

	// Kubernetes configuration

	//if h := c.homeDir; h != "" {
	//	flag.StringVar(c.kubeConfig,"kubeconfig", filepath.Join(h, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	//} else {
	//	flag.StringVar(c.kubeConfig,"kubeconfig", "", "absolute path to the kubeconfig file")
	//}

	if h := c.homeDir; h != "" {
		c.kubeConfig = flag.String("kubeconfig", filepath.Join(h, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		c.kubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	// Refreshing frequency
	flag.IntVar(&c.frequency, "frequency", 3, "refreshing frequency in seconds (default: 5)")

	// CLI Asks
	flag.BoolVar(&c.AskVersion, "version", false, "get Pody version")
	flag.BoolVar(&c.AskHelp, "help", false, "get Pody help")

	flag.Parse()

	c.setup = true
	CONFIG = c

	return c
}
