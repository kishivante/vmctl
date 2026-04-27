package commands

import (
	"flag"
	"fmt"
	"os"
	"vmctl/config"
	"vmctl/proxmox"
)

func UpdateVM() {
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)

	vmid := updateCmd.Int("vmid", 0, "VM ID")
	ram := updateCmd.Int("ram", 0, "RAM (MB)")
	cpu := updateCmd.Int("cpu", 0, "CPU cores")
	disk := updateCmd.String("disk", "", "Disk resize (+10G gibi)")

	updateCmd.Parse(os.Args[2:])

	if *vmid == 0 {
		fmt.Println("VM ID gerekli: --vmid")
		return
	}

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

	updates := make(map[string]interface{})

	if *ram > 0 {
		updates["memory"] = *ram
	}

	if *cpu > 0 {
		updates["cores"] = *cpu
	}

	if len(updates) == 0 && *disk == "" {
		fmt.Println("Değişiklik gerekli: --ram, --cpu veya --disk")
		return
	}

	if len(updates) > 0 {
		err := client.UpdateVM(cfg.Node, *vmid, updates)
		if err != nil {
			fmt.Println("Update hatası:", err)
			return
		}
	}

	if *disk != "" {
		err := client.ResizeDisk(cfg.Node, *vmid, "scsi0", *disk)
		if err != nil {
			fmt.Println("Disk resize hatası:", err)
			return
		}
	}

	fmt.Println("İşlem tamamlandı")
}
