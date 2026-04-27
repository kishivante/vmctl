package proxmox

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Client struct {
	BaseURL string
	Token   string
	CACert  string
}

type VM struct {
	VMID int    `json:"vmid"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type VMListResponse struct {
	Data []VM `json:"data"`
}

func (c *Client) getHTTPClient() (*http.Client, error) {
	if c.CACert == "" {
		return &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}, nil
	}

	cert, err := os.ReadFile(c.CACert)
	if err != nil {
		return nil, err
	}

	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(cert); !ok {
		return nil, fmt.Errorf("sertifika okunamadı")
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caPool,
			},
		},
	}, nil
}

func (c *Client) checkResponse(resp *http.Response, action string) error {
	body, _ := io.ReadAll(resp.Body)

	fmt.Println(action+" status:", resp.Status)

	if len(body) > 0 {
		fmt.Println("Response:", string(body))
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("%s başarısız: %s", action, resp.Status)
	}

	return nil
}

func (c *Client) VMExists(node string, vmid int) (bool, error) {
	vms, err := c.ListVMs(node)
	if err != nil {
		return false, err
	}

	for _, vm := range vms {
		if vm.VMID == vmid {
			return true, nil
		}
	}

	return false, nil
}

func (c *Client) CreateVM(node string, vmid int, name string, memory int, cores int, disk int) error {
	exists, err := c.VMExists(node, vmid)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("VMID %d zaten mevcut", vmid)
	}

	url := fmt.Sprintf("%s/api2/json/nodes/%s/qemu", c.BaseURL, node)

	data := map[string]interface{}{
		"vmid":   vmid,
		"name":   name,
		"memory": memory,
		"cores":  cores,
		"scsihw": "virtio-scsi-pci",
		"scsi0":  fmt.Sprintf("local-lvm:%d", disk),
		"net0":   "virtio,bridge=vmbr0",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "PVEAPIToken="+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client, err := c.getHTTPClient()
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.checkResponse(resp, "Create")
}

func (c *Client) ListVMs(node string) ([]VM, error) {
	url := fmt.Sprintf("%s/api2/json/nodes/%s/qemu", c.BaseURL, node)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "PVEAPIToken="+c.Token)

	client, err := c.getHTTPClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("VM listesi alınamadı: %s - %s", resp.Status, string(body))
	}

	var result VMListResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) StartVM(node string, vmid int) error {
	url := fmt.Sprintf("%s/api2/json/nodes/%s/qemu/%d/status/start", c.BaseURL, node, vmid)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "PVEAPIToken="+c.Token)

	client, err := c.getHTTPClient()
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.checkResponse(resp, "Start")
}

func (c *Client) CloneVM(node string, templateID int, newID int, name string) error {
	url := fmt.Sprintf("%s/api2/json/nodes/%s/qemu/%d/clone", c.BaseURL, node, templateID)

	data := map[string]interface{}{
		"newid": newID,
		"name":  name,
		"full":  1,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "PVEAPIToken="+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client, err := c.getHTTPClient()
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.checkResponse(resp, "Clone")
}

func (c *Client) SuspendVM(node string, vmid int) error {
	url := fmt.Sprintf("%s/api2/json/nodes/%s/qemu/%d/status/suspend", c.BaseURL, node, vmid)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "PVEAPIToken="+c.Token)

	client, err := c.getHTTPClient()
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.checkResponse(resp, "Suspend")
}

func (c *Client) ResumeVM(node string, vmid int) error {
	url := fmt.Sprintf("%s/api2/json/nodes/%s/qemu/%d/status/resume", c.BaseURL, node, vmid)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "PVEAPIToken="+c.Token)

	client, err := c.getHTTPClient()
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.checkResponse(resp, "Resume")
}

func (c *Client) UpdateVM(node string, vmid int, updates map[string]interface{}) error {
	url := fmt.Sprintf("%s/api2/json/nodes/%s/qemu/%d/config", c.BaseURL, node, vmid)

	jsonData, err := json.Marshal(updates)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "PVEAPIToken="+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client, err := c.getHTTPClient()
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.checkResponse(resp, "Update")
}

func (c *Client) ResizeDisk(node string, vmid int, disk string, size string) error {
	url := fmt.Sprintf("%s/api2/json/nodes/%s/qemu/%d/resize", c.BaseURL, node, vmid)

	data := map[string]interface{}{
		"disk": disk,
		"size": size,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "PVEAPIToken="+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client, err := c.getHTTPClient()
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.checkResponse(resp, "Disk Resize")
}

func (c *Client) BackupVM(node string, vmid int) error {
	url := fmt.Sprintf("%s/api2/json/nodes/%s/vzdump", c.BaseURL, node)

	data := map[string]interface{}{
		"vmid": vmid,
		"mode": "snapshot",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "PVEAPIToken="+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client, err := c.getHTTPClient()
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.checkResponse(resp, "Backup")
}

func (c *Client) GetVMStats(node string, vmid int) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api2/json/nodes/%s/qemu/%d/status/current", c.BaseURL, node, vmid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "PVEAPIToken="+c.Token)

	client, err := c.getHTTPClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("istatistik alınamadı: %s - %s", resp.Status, string(body))
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
