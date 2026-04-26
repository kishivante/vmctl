package commands

import (
	"fmt"
	"vmctl/config"
	"vmctl/proxmox"
)

func MonitorVM() {
	cfg, _ := config.LoadConfig()
	client := proxmox.Client{
		BaseURL: cfg.BaseURL,
		Token:   cfg.Token,
	}

	stats, _ := client.GetVMStats(cfg.Node, 102)
	fmt.Println(stats)
}
