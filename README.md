# vmctl

Proxmox API üzerinden VM oluşturma, clone, backup ve monitoring yapan CLI aracı.

## Kurulum

```bash
go build -o vmctl.exe

vmctl.exe create
vmctl.exe clone
vmctl.exe backup
vmctl.exe monitor