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

func StartVM() {
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)

	vmid := startCmd.Int("vmid", 0, "VM ID")
	startCmd.Parse(os.Args[2:])

	reader := bufio.NewReader(os.Stdin)

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

	err = client.StartVM(cfg.Node, *vmid)
	if err != nil {
		fmt.Println("Hata:", err)
		return
	}

	fmt.Println("VM başlatıldı.")
}
