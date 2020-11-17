package config

import (
	"fmt"
	"github.com/hpcng/warewulf/internal/pkg/util"
	"github.com/hpcng/warewulf/internal/pkg/wwlog"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)


type Config struct {
	Port            int    `yaml:"warewulfd port", envconfig:"WAREWULFD_PORT"`
	Ipaddr          string `yaml:"warewulfd ipaddr", envconfig:"WAREWULFD_IPADDR"`
	InsecureRuntime bool   `yaml:"insecure runtime"`
	Debug           bool   `yaml:"debug"`
	SysConfDir      string `yaml:"system config dir"`
	LocalStateDir   string `yaml:"local state dir"`
	Editor			string `yaml:"default editor", envconfig:"EDITOR"`
}

var c Config

func init() {
	fd, err := ioutil.ReadFile("/etc/warewulf/warewulf.conf")
	if err != nil {
		wwlog.Printf(wwlog.ERROR, "Could not read config file: %s\n", err)
		os.Exit(255)
	}

	err = yaml.Unmarshal(fd, &c)
	if err != nil {
		wwlog.Printf(wwlog.ERROR, "Could not unmarshal config file: %s\n", err)
		os.Exit(255)
	}

	err = envconfig.Process("", &c)
	if err != nil {
		wwlog.Printf(wwlog.ERROR, "Could not obtain environment configuration: %s\n", err)
		os.Exit(255)
	}

	if c.Ipaddr == "" {
		fmt.Printf("ERROR: 'warewulf ipaddr' has not been set in /etc/warewulf/warewulf.conf\n")
	}

	if c.SysConfDir == "" {
		c.SysConfDir = "/etc/warewulf"
	}
	if c.LocalStateDir == "" {
		c.LocalStateDir = "/var/warewulf"
	}

	util.ValidateOrDie("warewulf.conf", "warewulfd ipaddr", c.Ipaddr, "^[0-9]+.[0-9]+.[0-9]+.[0-9]+$")
	util.ValidateOrDie("warewulf.conf", "system config dir", c.SysConfDir, "^[a-zA-Z0-9-._:/]+$")
	util.ValidateOrDie("warewulf.conf", "local state dir", c.LocalStateDir, "^[a-zA-Z0-9-._:/]+$")
	util.ValidateOrDie("warewulf.conf", "default editor", c.LocalStateDir, "^[a-zA-Z0-9-._:/]+$")

}

func New() (Config) {
	return c
}

func (self *Config) NodeConfig() string {
	return fmt.Sprintf("%s/nodes.conf", self.LocalStateDir)
}

func (self *Config) OverlayDir() string {
	return fmt.Sprintf("%s/overlays/", self.LocalStateDir)
}

func (self *Config) SystemOverlayDir() string {
	return path.Join(self.OverlayDir(), "/system")
}

func (self *Config) RuntimeOverlayDir() string {
	return path.Join(self.OverlayDir(), "/runtime")
}

func (self *Config) SystemOverlaySource(overlayName string) string {
	if util.TaintCheck(overlayName, "^[a-zA-Z0-9-._]+$") == false {
		wwlog.Printf(wwlog.ERROR, "System overlay name contains illegal characters: %s\n", overlayName)
		os.Exit(1)
	}

	return path.Join(self.SystemOverlayDir(), overlayName)
}


func (self *Config) RuntimeOverlaySource(overlayName string) string {
	if util.TaintCheck(overlayName, "^[a-zA-Z0-9-._]+$") == false {
		wwlog.Printf(wwlog.ERROR, "Runtime overlay name contains illegal characters: %s\n", overlayName)
		os.Exit(1)
	}

	return path.Join(self.RuntimeOverlayDir(), overlayName)
}

func (self *Config) KernelImage(kernelVersion string) string {
	if util.TaintCheck(kernelVersion, "^[a-zA-Z0-9-._]+$") == false {
		wwlog.Printf(wwlog.ERROR, "Runtime overlay name contains illegal characters: %s\n", kernelVersion)
		os.Exit(1)
	}

	return fmt.Sprintf("%s/provision/kernel/vmlinuz-%s", self.LocalStateDir, kernelVersion)
}

func (self *Config) KmodsImage(kernelVersion string) string {
	if util.TaintCheck(kernelVersion, "^[a-zA-Z0-9-._]+$") == false {
		wwlog.Printf(wwlog.ERROR, "Runtime overlay name contains illegal characters: %s\n", kernelVersion)
		os.Exit(1)
	}

	return fmt.Sprintf("%s/provision/kernel/kmods-%s.img", self.LocalStateDir, kernelVersion)
}

func (self *Config) VnfsImage(vnfsNameClean string) string {
	if util.TaintCheck(vnfsNameClean, "^[a-zA-Z0-9-._:]+$") == false {
		wwlog.Printf(wwlog.ERROR, "Runtime overlay name contains illegal characters: %s\n", vnfsNameClean)
		os.Exit(1)
	}

	return fmt.Sprintf("%s/provision/vnfs/%s.img.gz", self.LocalStateDir, vnfsNameClean)
}

func (self *Config) SystemOverlayImage(nodeName string) string {
	if util.TaintCheck(nodeName, "^[a-zA-Z0-9-._:]+$") == false {
		wwlog.Printf(wwlog.ERROR, "System overlay name contains illegal characters: %s\n", nodeName)
		os.Exit(1)
	}

	return fmt.Sprintf("%s/provision/overlays/system/%s.img", self.LocalStateDir, nodeName)
}

func (self *Config) RuntimeOverlayImage(nodeName string) string {
	if util.TaintCheck(nodeName, "^[a-zA-Z0-9-._:]+$") == false {
		wwlog.Printf(wwlog.ERROR, "System overlay name contains illegal characters: %s\n", nodeName)
		os.Exit(1)
	}

	return fmt.Sprintf("%s/provision/overlays/runtime/%s.img", self.LocalStateDir, nodeName)
}

func (self *Config) VnfsChroot(vnfsNameClean string) string {
	if util.TaintCheck(vnfsNameClean, "^[a-zA-Z0-9-._:]+$") == false {
		wwlog.Printf(wwlog.ERROR, "Runtime overlay name contains illegal characters: %s\n", vnfsNameClean)
		os.Exit(1)
	}

	return fmt.Sprintf("%s/chroot/%s.img", self.LocalStateDir, vnfsNameClean)
}

