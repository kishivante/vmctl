package commands

import (
	"fmt"
	"vmctl/config"
	"vmctl/proxmox"
)

func CloneVM() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Config hatası:", err)
		return
	}

	client := proxmox.Client{
		BaseURL: cfg.BaseURL,
		Token:   cfg.Token,
		CACert:  cfg.CACert,
	}

	err = client.CloneVM(cfg.Node, 9000, 200, "cloned-vm")
	if err != nil {
		fmt.Println("Clone hatası:", err)
	}
}
