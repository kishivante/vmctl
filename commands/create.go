package commands

import (
	"fmt"
	"vmctl/config"
	"vmctl/proxmox"
)

func CreateVM() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Config okunamadı:", err)
		return
	}

	client := proxmox.Client{
		BaseURL: cfg.BaseURL,
		Token:   cfg.Token,
	}

	err = client.CreateVM(cfg.Node, 102, "test-vm")
	if err != nil {
		fmt.Println("VM oluşturma hatası:", err)
	}
}
