package commands

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"vmctl/config"
	"vmctl/proxmox"
)

const (
	green  = "\033[32m"
	red    = "\033[31m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	reset  = "\033[0m"
	clear  = "\033[H\033[2J"
)

func MonitorVM() {
	monitorCmd := flag.NewFlagSet("monitor", flag.ExitOnError)

	vmid := monitorCmd.Int("vmid", 0, "VM ID")
	interval := monitorCmd.Int("interval", 2, "Refresh interval seconds")

	monitorCmd.Parse(os.Args[2:])

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
		fmt.Println("Config okunamadı:", err)
		return
	}

	client := proxmox.Client{
		BaseURL: cfg.BaseURL,
		Token:   cfg.Token,
		CACert:  cfg.CACert,
	}

	for {
		stats, err := client.GetVMStats(cfg.Node, *vmid)

		fmt.Print(clear)

		if err != nil {
			fmt.Println("Monitor hatası:", err)
		} else {
			printMonitorDashboard(stats)
		}

		fmt.Println()
		fmt.Println("CTRL + C ile çıkış yapabilirsin.")

		time.Sleep(time.Duration(*interval) * time.Second)
	}
}

func printMonitorDashboard(stats map[string]interface{}) {
	data, ok := stats["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Geçersiz veri formatı")
		return
	}

	vmid := getFloat(data, "vmid")
	name := getString(data, "name")
	status := getString(data, "status")
	qmpstatus := getString(data, "qmpstatus")

	cpu := getFloat(data, "cpu") * 100
	cpus := getFloat(data, "cpus")

	mem := getFloat(data, "mem")
	maxmem := getFloat(data, "maxmem")

	disk := getFloat(data, "disk")
	maxdisk := getFloat(data, "maxdisk")

	netin := getFloat(data, "netin")
	netout := getFloat(data, "netout")

	uptime := getFloat(data, "uptime")

	state := statusIcon(status, qmpstatus)

	fmt.Println(blue + "VM MONITOR" + reset)
	fmt.Println()

	fmt.Printf("%-6s %-16s %-8s %-8s %-7s %-18s %-18s %-12s %-12s %-12s\n",
		"VMID",
		"NAME",
		"STATE",
		"CPU",
		"CORES",
		"MEMORY",
		"DISK",
		"NET-IN",
		"NET-OUT",
		"UPTIME",
	)

	fmt.Printf("%-6s %-16s %-8s %-8s %-7s %-18s %-18s %-12s %-12s %-12s\n",
		"----",
		"----",
		"-----",
		"---",
		"-----",
		"------",
		"----",
		"------",
		"-------",
		"------",
	)

	fmt.Printf("%-6.0f %-16s %-8s %-8s %-7.0f %-18s %-18s %-12s %-12s %-12s\n",
		vmid,
		cutText(name, 16),
		state,
		fmt.Sprintf("%.1f%%", cpu),
		cpus,
		cutText(formatUsage(mem, maxmem), 18),
		cutText(formatUsage(disk, maxdisk), 18),
		formatBytesShort(netin),
		formatBytesShort(netout),
		formatUptime(uptime),
	)

	fmt.Println()
	fmt.Println(green + "UP" + reset + " = running   " + red + "DOWN" + reset + " = stopped   " + yellow + "SUSP" + reset + " = suspended")
}

func statusIcon(status string, qmpstatus string) string {
	if status == "running" {
		return "UP"
	}

	if status == "stopped" {
		return "DOWN"
	}

	if qmpstatus == "paused" || status == "suspended" {
		return "SUSP"
	}

	return "UNK"
}

func formatUsage(used float64, max float64) string {
	if max <= 0 {
		return "-"
	}

	percent := (used / max) * 100
	return fmt.Sprintf("%s/%s %.0f%%", formatBytesShort(used), formatBytesShort(max), percent)
}

func formatBytesShort(bytes float64) string {
	gb := bytes / 1024 / 1024 / 1024
	if gb >= 1 {
		return fmt.Sprintf("%.1fgb", gb)
	}

	mb := bytes / 1024 / 1024
	return fmt.Sprintf("%.1fmb", mb)
}

func cutText(text string, max int) string {
	if len(text) <= max {
		return text
	}

	if max <= 3 {
		return text[:max]
	}

	return text[:max-3] + "..."
}

func getString(data map[string]interface{}, key string) string {
	value, ok := data[key].(string)
	if !ok {
		return "-"
	}
	return value
}

func getFloat(data map[string]interface{}, key string) float64 {
	value, ok := data[key].(float64)
	if !ok {
		return 0
	}
	return value
}

func formatUptime(seconds float64) string {
	if seconds <= 0 {
		return "0s"
	}

	h := int(seconds) / 3600
	m := (int(seconds) % 3600) / 60
	s := int(seconds) % 60

	return fmt.Sprintf("%02dh%02dm%02ds", h, m, s)
}
