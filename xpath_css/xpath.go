package xpath_css

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func Xpath() {
	url := "http://quotes.toscrape.com/"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3776.0 Safari/537.36")
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	doc, _ := htmlquery.Parse(resp.Body)
	list := htmlquery.Find(doc, "//div[@class=\"quote\"]")

	for _, n := range list {
		content := htmlquery.FindOne(n, ".//span[1]")
		author := htmlquery.FindOne(n, "/span[2]//small")

		fmt.Printf("%s-%s\n", htmlquery.InnerText(author), htmlquery.InnerText(content))
	}
}

func XpathIPToAddr(ip string) {
	//url := "https://www.ipshudi.com/66.249.65.123.htm"
	//url := "https://www.ipshudi.com/119.131.198.248.htm"
	url := fmt.Sprintf("https://www.ipshudi.com/%s.htm", ip)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3776.0 Safari/537.36")
	client := &http.Client{Timeout: time.Second * 3}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
	//fmt.Println(222)
	doc, err := htmlquery.Parse(strings.NewReader(string(b)))
	if err != nil {
		panic(err)
	}
	span, err := htmlquery.Query(doc, "/html/body/div[1]/div[1]/div[1]/div/div[3]/table/tbody/tr[1]/td[2]/span")
	if err != nil {
		panic(err)
	}
	if span == nil || span.FirstChild == nil {
		fmt.Println("span is nil")
		return
	}
	td, _ := htmlquery.Query(doc, "/html/body/div[1]/div[1]/div[1]/div/div[3]/table/tbody/tr[2]/td[2]/span[1]")
	isp := ""
	if td != nil && td.FirstChild != nil {
		isp = td.FirstChild.Data // 国外isp可能查不到
	}
	/*
		Addr= 美国 俄克拉荷马州 普赖尔  isp= 谷歌云
		Addr= 中国 安徽省   isp= 电信
		Addr= 俄罗斯联邦 莫斯科市   isp=
	*/
	fmt.Println("Addr=", span.FirstChild.Data, "isp=", isp)
}

type bdCloudApiRet struct {
	Code string `json:"code"` // Success
	Data struct {
		Country  string `json:"country"`
		Prov     string `json:"prov"`
		City     string `json:"city"`
		District string `json:"district"`
		Isp      string `json:"isp"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// QueryAddrByIP_BaiduAPI
// 使用这个页面的API：https://qifu.baidu.com/?activeKey=SEARCH_IP&trace=apistore_ip_aladdin&activeId=SEARCH_IP_ADDRESS&ip=
func QueryAddrByIP_BaiduAPI(ip string) {
	url := fmt.Sprintf("https://gwgp-cekvddtwkob.n.bdcloudapi.com/ip/geo/v1/district?ip=%s", ip)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3776.0 Safari/537.36")
	client := &http.Client{Timeout: time.Second * 3}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	ret := new(bdCloudApiRet)
	_ = json.Unmarshal(b, ret)

	fmt.Printf("%s %+v \n", ip, ret)
}

func QueryAddrByIP_IPCn(ip string) {
	url := fmt.Sprintf("https://www.ip.cn/ip/%s.html", ip)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3776.0 Safari/537.36")
	client := &http.Client{Timeout: time.Second * 3}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	//fmt.Println(string(b))
	//fmt.Println(222)
	doc, err := htmlquery.Parse(strings.NewReader(string(b)))
	if err != nil {
		panic(err)
	}
	span, err := htmlquery.Query(doc, "//*[@id=\"tab0_address\"]")
	if err != nil {
		panic(err)
	}
	if span == nil || span.FirstChild == nil {
		fmt.Println("span is nil")
		return
	}
	fmt.Println("Addr=", span.FirstChild.Data) // 英国  伦敦 伦敦 (国外IP的查询结果的最后一位不是ISP)
	d := strings.ReplaceAll(span.FirstChild.Data, "  ", " ")
	ss := strings.Split(d, " ")
	for _, i := range ss {
		println(i)
	}
}

func QueryAddrByIP_Chinaz(ip string) {
	url := fmt.Sprintf("https://ip.tool.chinaz.com/%s", ip)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3776.0 Safari/537.36")
	client := &http.Client{Timeout: time.Second * 3}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	//t.Log(string(b)) // 直接打印 不出全部内容，只能strings.Contains判断
	if strings.Contains(string(b), "香港") {
		println(111)
	}
	doc, err := htmlquery.Parse(strings.NewReader(string(b)))
	if err != nil {
		panic(err)
	}
	span, err := htmlquery.Query(doc, "//*[@id=\"infoLocation\"]")
	if err != nil {
		panic(err)
	}
	if span == nil || span.FirstChild == nil {
		fmt.Println("span is nil")
		return
	}
	fmt.Println("Addr=", span.FirstChild.Data)
}

// 有JS反扒，暂无法破解
func QueryAddrByIP_Ping0(ip string) {
	url := fmt.Sprintf("https://ping0.cc/ip/%s", ip)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3776.0 Safari/537.36")
	req.Header.Add("cookie", "jskey=1139d09f4785e2e01cb81caaf7228b0e")
	client := &http.Client{Timeout: time.Second * 3}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b)) // 直接打印 不出全部内容，只能strings.Contains判断
	if strings.Contains(string(b), "香港") {
		println(111)
	}
	doc, err := htmlquery.Parse(strings.NewReader(string(b)))
	if err != nil {
		panic(err)
	}
	span, err := htmlquery.Query(doc, "//*[@id=\"check\"]/div[2]/div[1]/div[2]/div[2]/div[2]")
	if err != nil {
		panic(err)
	}
	if span == nil || span.FirstChild == nil {
		fmt.Println("span is nil")
		return
	}
	fmt.Println("Addr=", span.FirstChild.Data)
}

/*
eg.
Addr=&{Status:success Country:中国 CountryCode:CN Region:GD 	  RegionName:广东 City:广州市 Zip: Lat:23.1181 Lon:113.2539 Timezone:Asia/Shanghai      Isp:Chinanet Org:Chinanet GD As:AS4134 CHINANET-BACKBONE Query:119.131.198.248}
Addr=&{Status:success Country:香港 CountryCode:HK Region:KYT   RegionName:油尖旺區 City:旺角 Zip:96521 Lat:22.316 Lon:114.172 Timezone:Asia/Hon      g_Kong Isp:Nearoute Limited Org:Kidc Limited As:AS134972 IKUUU NETWORK LTD Query:103.151.172.30}
*/
type Ip125Resp struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

func QueryAddrByIP_Ip125(ip string) {
	url := fmt.Sprintf("https://ip125.com/api/%s?lang=zh-CN", ip)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3776.0 Safari/537.36")
	req.Header.Add("cookie", "jskey=1139d09f4785e2e01cb81caaf7228b0e")
	client := &http.Client{Timeout: time.Second * 3}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	r := new(Ip125Resp)
	_ = json.NewDecoder(resp.Body).Decode(r)

	fmt.Printf("Addr=%+v\n", r)
	if r.Status != "success" {
		return
	}
}
