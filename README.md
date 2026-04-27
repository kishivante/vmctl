# vmctl

CLI tool for managing Proxmox VE virtual machines via API.

Proxmox API kullanarak sanal makineleri yönetmek için geliştirilmiş bir CLI aracıdır.

---

## 🚀 Features / Özellikler

- VM oluşturma (Create VM)
- VM klonlama (Clone VM)
- VM başlatma / durdurma / suspend / resume
- Kaynak güncelleme (RAM, CPU)
- Disk büyütme (resize - only grow)
- Backup alma
- Canlı izleme (monitor)

---

## ⚙️ Requirements / Gereksinimler

- Go (1.20+)
- Proxmox VE
- API Token

---

## 🔧 Setup / Kurulum

```bash
go build -o vmctl.exe
