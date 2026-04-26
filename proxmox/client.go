package proxmox

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL string
	Token   string
}

func (c *Client) CreateVM(node string, vmid int, name string) error {
	url := fmt.Sprintf("%s/api2/json/nodes/%s/qemu", c.BaseURL, node)
	fmt.Println("Request URL:", url)

	data := map[string]interface{}{
		// önemöli burası vm ayarları var
		"vmid":   vmid,
		"name":   name,
		"memory": 1024, // ram
		"cores":  1,    // cpu için ayar vs

		"scsihw": "virtio-scsi-pci",
		"scsi0":  "local-lvm:10", // disk boyutu 10
		"net0":   "virtio,bridge=vmbr0",
	}

	jsonData, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "PVEAPIToken="+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	fmt.Println("Create status:", resp.Status)
	return nil
}

func (c *Client) BackupVM(node string, vmid int) error {
	url := fmt.Sprintf("%s/api2/json/nodes/%s/qemu", c.BaseURL, node)

	data := map[string]interface{}{
		"vmid": vmid,
		"mode": "snapshot",
	}

	jsonData, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "PVEAPIToken="+c.Token)

	_, err := http.DefaultClient.Do(req)
	return err
}

func (c *Client) GetVMStats(node string, vmid int) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api2/json/nodes/%s/qemu/%d/status/current", c.BaseURL, node, vmid)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "PVEAPIToken="+c.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result, nil
}
