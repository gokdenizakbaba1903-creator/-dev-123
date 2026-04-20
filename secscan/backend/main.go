package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

type ScanRequest struct {
	URL string `json:"url"`
}

type ScanResult struct {
	ID        string         `json:"id"`
	URL       string         `json:"url"`
	Progress  int            `json:"progress"`
	Status    string         `json:"status"`
	Findings  []Finding      `json:"findings"`
	Score     string         `json:"score"`
	RadarData map[string]int `json:"radar_data"`
}

type Finding struct {
	Module   string `json:"module"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

var (
	scans   = make(map[string]*ScanResult)
	mu      sync.RWMutex
	clients = make(map[string]chan string)
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/api/scan", startScan)
	r.GET("/api/scan/:id", getScan)
	r.GET("/api/report/:id/pdf", downloadPDF)
	r.GET("/stream/:id", streamScan)

	fmt.Println("SecScan Backend running on :8080")
	r.Run(":8080")
}

func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}
	privateIPBlocks := []*net.IPNet{
		{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(8, 32)},
		{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(12, 32)},
		{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)},
		{IP: net.ParseIP("127.0.0.0"), Mask: net.CIDRMask(8, 32)},
		{IP: net.ParseIP("169.254.0.0"), Mask: net.CIDRMask(16, 32)},
	}
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func ssrfMiddleware(targetURL string) error {
	host := targetURL
	if strings.Contains(host, "://") {
		parts := strings.Split(host, "://")
		host = strings.Split(parts[1], "/")[0]
	} else {
		host = strings.Split(host, "/")[0]
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return err
	}
	for _, ip := range ips {
		if isPrivateIP(ip) {
			return fmt.Errorf("SSRF Protection: Private IP detected (%s)", ip.String())
		}
	}
	return nil
}

func startScan(c *gin.Context) {
	var req ScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := ssrfMiddleware(req.URL); err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	mu.Lock()
	scans[id] = &ScanResult{
		ID:        id,
		URL:       req.URL,
		Progress:  0,
		Status:    "Running",
		Findings:  []Finding{},
		RadarData: make(map[string]int),
	}
	mu.Unlock()

	go runScanner(id, req.URL)

	c.JSON(200, gin.H{"id": id})
}

func getScan(c *gin.Context) {
	id := c.Param("id")
	mu.RLock()
	scan, ok := scans[id]
	mu.RUnlock()
	if !ok {
		c.JSON(404, gin.H{"error": "Scan not found"})
		return
	}
	c.JSON(200, scan)
}

func streamScan(c *gin.Context) {
	id := c.Param("id")
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	clientChan := make(chan string)
	mu.Lock()
	clients[id] = clientChan
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(clients, id)
		mu.Unlock()
		close(clientChan)
	}()

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-clientChan; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})
}

func downloadPDF(c *gin.Context) {
	id := c.Param("id")
	mu.RLock()
	scan, ok := scans[id]
	mu.RUnlock()
	if !ok {
		c.JSON(404, gin.H{"error": "Scan not found"})
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	
	// HEADER
	pdf.SetFillColor(41, 128, 185)
	pdf.Rect(0, 0, 210, 40, "F")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 24)
	pdf.CellFormat(0, 20, "SecScan Security Report", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 10, fmt.Sprintf("Report Generated: %s", time.Now().Format("2006-01-02 15:04:05")), "", 1, "C", false, 0, "")
	pdf.Ln(15)

	// TARGET INFO
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Target Summary")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Target URL:")
	pdf.Cell(0, 10, scan.URL)
	pdf.Ln(8)
	pdf.Cell(40, 10, "Final Score:")
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(39, 174, 96)
	pdf.Cell(0, 10, scan.Score)
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(15)

	// FINDINGS TABLE HEADER
	pdf.SetFillColor(236, 240, 241)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(40, 10, "Module", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Severity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(120, 10, "Observation", "1", 1, "C", true, 0, "")

	// FINDINGS DATA
	pdf.SetFont("Arial", "", 10)
	for _, f := range scan.Findings {
		if f.Severity == "High" {
			pdf.SetTextColor(192, 57, 43)
		} else {
			pdf.SetTextColor(41, 128, 185)
		}
		
		pdf.CellFormat(40, 10, f.Module, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 10, f.Severity, "1", 0, "C", false, 0, "")
		pdf.CellFormat(120, 10, f.Message, "1", 1, "L", false, 0, "")
	}

	pdf.Ln(20)
	pdf.SetTextColor(127, 140, 141)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(0, 10, "Confidential - For Educational Purposes Only")

	c.Header("Content-Disposition", "attachment; filename=SecScan_Report.pdf")
	c.Header("Content-Type", "application/pdf")
	pdf.Output(c.Writer)
}

func runScanner(id string, targetURL string) {
	modules := []string{"Port Scanner", "Security Headers", "TLS Analyzer", "Directory Fuzzer", "XSS Test", "SQLi Test", "CVE Checker"}
	
	for i, mod := range modules {
		var findings []Finding
		score := 100

		switch mod {
		case "Port Scanner":
			findings, score = scanPorts(targetURL)
		case "Security Headers":
			findings, score = checkHeaders(targetURL)
		case "TLS Analyzer":
			findings, score = checkTLS(targetURL)
		case "Directory Fuzzer":
			findings, score = fuzzDirectories(targetURL)
		case "XSS Test":
			findings, score = scanXSS(targetURL)
		case "SQLi Test":
			findings, score = scanSQLi(targetURL)
		case "CVE Checker":
			findings, score = scanCVE(targetURL)
		default:
			findings = append(findings, Finding{Module: mod, Severity: "Low", Message: mod + " verification complete."})
			score = 100
		}

		mu.Lock()
		scan := scans[id]
		scan.Findings = append(scan.Findings, findings...)
		scan.RadarData[mod] = score
		scan.Progress = (i + 1) * 100 / len(modules)
		
		if scan.Progress >= 100 {
			scan.Status = "Completed"
			highCount := 0
			for _, f := range scan.Findings {
				if f.Severity == "High" { highCount++ }
			}
			if highCount == 0 { scan.Score = "A+" } else if highCount < 2 { scan.Score = "B" } else { scan.Score = "F" }
		}
		
		msg, _ := json.Marshal(scan)
		if ch, ok := clients[id]; ok {
			ch <- string(msg)
		}
		mu.Unlock()
	}
}

// REAL SCANNER LOGIC

func scanPorts(target string) ([]Finding, int) {
	host := getHost(target)
	ports := []int{21, 22, 23, 25, 53, 80, 110, 443, 3306, 8080}
	var findings []Finding
	openPorts := 0

	for _, port := range ports {
		address := fmt.Sprintf("%s:%d", host, port)
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err == nil {
			findings = append(findings, Finding{Module: "Port Scanner", Severity: "High", Message: fmt.Sprintf("Open port found: %d", port)})
			openPorts++
			conn.Close()
		}
	}

	score := 100 - (openPorts * 15)
	if score < 0 { score = 0 }
	if len(findings) == 0 {
		findings = append(findings, Finding{Module: "Port Scanner", Severity: "Low", Message: "No common sensitive ports are exposed."})
	}
	return findings, score
}

func checkHeaders(url string) ([]Finding, int) {
	if !strings.HasPrefix(url, "http") { url = "http://" + url }
	resp, err := http.Head(url)
	if err != nil {
		return []Finding{{Module: "Security Headers", Severity: "High", Message: "Connection failed"}}, 0
	}
	defer resp.Body.Close()

	var findings []Finding
	missing := 0
	required := []string{"Content-Security-Policy", "X-Frame-Options", "X-Content-Type-Options", "Strict-Transport-Security"}
	for _, h := range required {
		if resp.Header.Get(h) == "" {
			findings = append(findings, Finding{Module: "Security Headers", Severity: "High", Message: "Missing: " + h})
			missing++
		}
	}
	
	score := 100 - (missing * 25)
	if len(findings) == 0 {
		findings = append(findings, Finding{Module: "Security Headers", Severity: "Low", Message: "Essential security headers are present."})
	}
	return findings, score
}

func checkTLS(target string) ([]Finding, int) {
	host := getHost(target)
	conf := &tls.Config{InsecureSkipVerify: false}
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 3 * time.Second}, "tcp", host+":443", conf)
	
	if err != nil {
		return []Finding{{Module: "TLS Analyzer", Severity: "High", Message: "Insecure or missing TLS/SSL"}}, 20
	}
	defer conn.Close()
	return []Finding{{Module: "TLS Analyzer", Severity: "Low", Message: "Valid TLS configuration found."}}, 100
}

func fuzzDirectories(target string) ([]Finding, int) {
	if !strings.HasPrefix(target, "http") { target = "http://" + target }
	paths := []string{"/admin", "/.git", "/.env", "/config.php", "/api/v1"}
	var findings []Finding
	found := 0

	for _, path := range paths {
		resp, err := http.Get(target + path)
		if err == nil && resp.StatusCode == 200 {
			findings = append(findings, Finding{Module: "Directory Fuzzer", Severity: "High", Message: "Sensitive path found: " + path})
			found++
		}
		if resp != nil { resp.Body.Close() }
	}

	score := 100 - (found * 30)
	if score < 0 { score = 0 }
	if len(findings) == 0 {
		findings = append(findings, Finding{Module: "Directory Fuzzer", Severity: "Low", Message: "No common sensitive directories found."})
	}
	return findings, score
}

func getHost(target string) string {
	host := target
	if strings.Contains(host, "://") {
		parts := strings.Split(host, "://")
		host = strings.Split(parts[1], "/")[0]
	} else {
		host = strings.Split(host, "/")[0]
	}
	return host
}

func scanXSS(target string) ([]Finding, int) {
	if !strings.HasPrefix(target, "http") { target = "http://" + target }
	payload := "<svg/onload=alert(1)>"
	params := []string{"q", "s", "search", "id", "query"}
	var findings []Finding
	found := 0

	for _, p := range params {
		testURL := fmt.Sprintf("%s?%s=%s", target, p, payload)
		resp, err := http.Get(testURL)
		if err == nil {
			body, _ := io.ReadAll(resp.Body)
			if strings.Contains(string(body), payload) {
				findings = append(findings, Finding{Module: "XSS Test", Severity: "High", Message: fmt.Sprintf("Reflected XSS found in parameter: %s", p)})
				found++
			}
			resp.Body.Close()
		}
	}

	score := 100 - (found * 40)
	if score < 0 { score = 0 }
	if len(findings) == 0 {
		findings = append(findings, Finding{Module: "XSS Test", Severity: "Low", Message: "No obvious reflected XSS vulnerabilities found."})
	}
	return findings, score
}

func scanSQLi(target string) ([]Finding, int) {
	if !strings.HasPrefix(target, "http") { target = "http://" + target }
	payloads := []string{"'", "\""}
	params := []string{"id", "user", "item", "search"}
	errorPatterns := []string{"sql syntax", "mysql", "postgresql", "sqlite", "oracle", "native client"}
	var findings []Finding
	found := 0

	for _, p := range params {
		for _, pay := range payloads {
			testURL := fmt.Sprintf("%s?%s=%s", target, p, pay)
			resp, err := http.Get(testURL)
			if err == nil {
				if resp.StatusCode == 500 {
					findings = append(findings, Finding{Module: "SQLi Test", Severity: "High", Message: fmt.Sprintf("Potential SQLi (500 Error) in parameter: %s", p)})
					found++
				} else {
					body, _ := io.ReadAll(resp.Body)
					lowerBody := strings.ToLower(string(body))
					for _, pattern := range errorPatterns {
						if strings.Contains(lowerBody, pattern) {
							findings = append(findings, Finding{Module: "SQLi Test", Severity: "High", Message: fmt.Sprintf("SQL error pattern '%s' found in parameter: %s", pattern, p)})
							found++
							break
						}
					}
				}
				resp.Body.Close()
			}
		}
	}

	score := 100 - (found * 40)
	if score < 0 { score = 0 }
	if len(findings) == 0 {
		findings = append(findings, Finding{Module: "SQLi Test", Severity: "Low", Message: "No common SQL injection patterns detected."})
	}
	return findings, score
}

func scanCVE(target string) ([]Finding, int) {
	if !strings.HasPrefix(target, "http") { target = "http://" + target }
	resp, err := http.Head(target)
	if err != nil {
		return []Finding{{Module: "CVE Checker", Severity: "Low", Message: "Could not retrieve headers for CVE analysis."}}, 100
	}
	defer resp.Body.Close()

	server := resp.Header.Get("Server")
	poweredBy := resp.Header.Get("X-Powered-By")
	var findings []Finding
	score := 100

	vulnerableVersions := map[string]string{
		"Apache/2.4.49": "CVE-2021-41773 (Path Traversal)",
		"nginx/1.17.7":  "CVE-2019-20372 (Request Smuggling)",
		"PHP/5.4":       "End of Life version - high risk",
	}

	for ver, cve := range vulnerableVersions {
		if strings.Contains(server, ver) || strings.Contains(poweredBy, ver) {
			findings = append(findings, Finding{Module: "CVE Checker", Severity: "High", Message: fmt.Sprintf("Known vulnerability: %s detected.", cve)})
			score = 40
		}
	}

	if len(findings) == 0 {
		msg := "No known CVEs found based on fingerprinting."
		if server != "" { msg = fmt.Sprintf("Server fingerprinted as %s - no known critical CVEs.", server) }
		findings = append(findings, Finding{Module: "CVE Checker", Severity: "Low", Message: msg})
	}
	return findings, score
}
