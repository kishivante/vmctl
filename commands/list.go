package commands

import (
	"fmt"
	"vmctl/config"
	"vmctl/proxmox"
)

func ListVMs() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Config okunamadı:", err)
		return
	}

	client := proxmox.Client{
		BaseURL: cfg.BaseURL,
		Token:   cfg.Token,
		CACert:  cfg.CACert,
	}

	vms, err := client.ListVMs(cfg.Node)
	if err != nil {
		fmt.Println("VM listesi alınamadı:", err)
		return
	}

	fmt.Println("Mevcut VM'ler:")
	fmt.Println("-------------------------")

	for _, vm := range vms {
		fmt.Printf("ID: %-5d Name: %s\n", vm.VMID, vm.Name)
	}
}
