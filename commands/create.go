package commands

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"vmctl/config"
	"vmctl/proxmox"
)

func CreateVM() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)

	vmid := createCmd.Int("vmid", 0, "VM ID")
	name := createCmd.String("name", "", "VM name")
	memory := createCmd.Int("memory", 1024, "RAM MB")
	cores := createCmd.Int("cores", 1, "CPU core count")
	disk := createCmd.Int("disk", 10, "Disk size GB")

	createCmd.Parse(os.Args[2:])

	reader := bufio.NewReader(os.Stdin)

	// VMID sor
	if *vmid == 0 {
		fmt.Print("VM ID gir: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		id, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Geçersiz VM ID")
			return
		}
		*vmid = id
	}

	// NAME sor
	if *name == "" {
		fmt.Print("VM ismi gir: ")
		input, _ := reader.ReadString('\n')
		*name = strings.TrimSpace(input)
	}

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

	err = client.CreateVM(cfg.Node, *vmid, *name, *memory, *cores, *disk)
	if err != nil {
		fmt.Println("VM oluşturma hatası:", err)
		return
	}

	fmt.Println("VM oluşturma isteği gönderildi.")
}
