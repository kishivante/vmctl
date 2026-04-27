package commands

import (
	"fmt"
	"vmctl/config"
	"vmctl/proxmox"
)

func ResumeVM() {
	cfg, _ := config.LoadConfig()

	client := proxmox.Client{
		BaseURL: cfg.BaseURL,
		Token:   cfg.Token,
	}

	err := client.ResumeVM(cfg.Node, 102)
	if err != nil {
		fmt.Println("Hata:", err)
	}
}
