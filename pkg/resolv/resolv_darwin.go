// +build darwin

package resolv

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"

	"github.com/kayrus/gof5/pkg/config"
)

var re = regexp.MustCompile(`\(\d+\)\s(.*)`)

func getInterfaces() ([]string, error) {
	v, err := exec.Command("networksetup", "-listnetworkserviceorder").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get a list of interfaces: %s: %s", v, err)
	}

	matches := re.FindAllStringSubmatch(string(v), -1)

	if len(matches) == 0 {
		return nil, fmt.Errorf("cannot find interfaces list")
	}

	var res []string
	for _, v := range matches {
		if len(v) == 2 {
			res = append(res, v[1])
		}
	}

	return res, nil
}

func ConfigureDNS(cfg *config.Config) error {
	ifaces, err := getInterfaces()
	if err != nil {
		return err
	}

	for _, iface := range ifaces {
		if len(cfg.DNS) == 0 {
			args := []string{
				"-setdnsservers",
				iface,
			}
			for _, v := range cfg.F5Config.Object.DNS {
				args = append(args, v.String())
			}
			v, err := exec.Command("networksetup", args...).Output()
			if err != nil {
				return fmt.Errorf("failed to set %q DNS servers on %q: %s: %s", args, iface, v, err)
			}
		} else {
			v, err := exec.Command("networksetup", "-setdnsservers", iface, cfg.ListenDNS.String()).Output()
			if err != nil {
				return fmt.Errorf("failed to set %q DNS server on %q: %s: %s", cfg.ListenDNS.String(), iface, v, err)
			}
		}

		if len(cfg.F5Config.Object.DNSSuffix) > 0 {
			args := []string{"-setsearchdomains", iface}
			args = append(args, cfg.F5Config.Object.DNSSuffix...)
			v, err := exec.Command("networksetup", args...).Output()
			if err != nil {
				return fmt.Errorf("failed to set %q DNS search prefix on %q: %s: %s", cfg.F5Config.Object.DNSSuffix, iface, v, err)
			}
		}
	}

	return nil
}

func RestoreDNS(cfg *config.Config) {
	ifaces, err := getInterfaces()
	if err != nil {
		log.Printf("Failed to restore DNS settings: %s", err)
		return
	}

	for _, iface := range ifaces {
		v, err := exec.Command("networksetup", "-setdnsservers", iface, "empty").Output()
		if err != nil {
			log.Printf("Failed to restore DNS servers on %q: %s: %s", iface, v, err)
		}

		if len(cfg.F5Config.Object.DNSSuffix) > 0 {
			v, err := exec.Command("networksetup", "-setsearchdomains", iface, "empty").Output()
			if err != nil {
				log.Printf("failed to restore DNS search prefix on %q: %s: %s", iface, v, err)
			}
		}
	}
}
