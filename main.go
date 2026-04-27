package main

import (
	"fmt"
	"os"
	"vmctl/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Kullanım: vmctl [create|backup|monitor|clone|list]")
		return
	}

	switch os.Args[1] {

	case "create":
		commands.CreateVM()

	case "start":
		commands.StartVM()

	case "backup":
		commands.BackupVM()

	case "monitor":
		commands.MonitorVM()

	case "clone":
		commands.CloneVM()

	case "list":
		commands.ListVMs()

	case "update":
		commands.UpdateVM()

	case "suspend":
		commands.SuspendVM()

	case "resume":
		commands.ResumeVM()

	default:
		fmt.Println("Bilinmeyen komut")
	}
}
