package main

import (
	"fmt"
	"os"
	"vmctl/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Kullanım: vmctl [create|backup|monitor]")
		return
	}

	switch os.Args[1] {

	case "create":
		commands.CreateVM()

	case "backup":
		commands.BackupVM()

	case "monitor":
		commands.MonitorVM()

	default:
		fmt.Println("Bilinmeyen komut")
	}
}
