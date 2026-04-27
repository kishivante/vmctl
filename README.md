# vmctl

Proxmox API üzerinden VM oluşturma, clone, backup ve monitoring yapan CLI aracı.

Şu an sadece create test ettim. Diğer araçları da daha sonra güncelleyip ekleyeceğim.

## Kurulum

```bash
go build -o vmctl.exe

vmctl.exe create
vmctl.exe clone
vmctl.exe backup
vmctl.exe monitor
