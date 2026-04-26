package commands

import (
	"vmctl/config"
	"vmctl/proxmox"
)

func BackupVM() {
	cfg, _ := config.LoadConfig()
	client := proxmox.Client{
		BaseURL: cfg.BaseURL,
		Token:   cfg.Token,
	}

	client.BackupVM(cfg.Node, 102)
}
