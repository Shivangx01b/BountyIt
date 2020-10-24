package main

import (
	"fmt"
	"net/http"
	"sync"
	"io/ioutil"
	"time"
	"net"
	"crypto/tls"
	"github.com/fatih/color"
	"flag"
	"bufio"
	"os"
	"log"
	"strings"

)

var Threads int
var recheck_url string
var method string
var body string
var payload string
var base_size int
var matcher string
var payloads []string 
var confirm []string
var verify bool
var signatures = []string{"Program Files", "Windows", "[boot loader]", "[drivers]", "HTTP /1.1", "HTTP /1.0", "About php.ini", "root:x:", "root:*"} 

func getClient() *http.Client {          
	tr := &http.Transport{
		MaxIdleConns:    30,
		IdleConnTimeout: time.Second,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   time.Second * 10,
			KeepAlive: time.Second,
		}).DialContext,
	}

	re := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &http.Client{
		Transport:     tr,
		CheckRedirect: re,
		Timeout:       time.Second * 10,
	}
}

func base_request(c *http.Client, u string, method string, matcher string) (int, string) {
		req, _ := http.NewRequest(method, u, nil)
		if req != nil {
			resp, _ := c.Do(req)
			if resp != nil {
				contents, _ := ioutil.ReadAll(resp.Body)
				if matcher == "check" {
					body = string(contents)
				}
				base_size = len(contents)
				resp.Body.Close()
			}
		}

	return base_size, body
}


func requester(c *http.Client,  u string, method string, list []string , verify bool, matcher string) {
	req_base, _ := base_request(c, u, method, matcher)
	for _, test := range list {
		url := strings.Replace(u, "FUZZ", test, -1)
		req_test, _ := base_request(c, url, method , matcher)
		if req_test != req_base {
			if verify != true {
				fmt.Printf("%v %s\n", color.RedString("[!] Potential vulnerability found at:..🛠") , url)
				fmt.Printf("%v\n", color.CyanString("[~] Storing for confirmation..✒"))
			}
			confirm = append(confirm, url)      
		}
	}
	if verify != true {
		fmt.Printf("%v\n",color.YellowString("[>] Staring confirmation tests..🔍"))
	}
	matcher = "check"
	for _, recheck_url = range confirm {
		_, checkbody := base_request(c, recheck_url, method, matcher)
		for _, query :=  range signatures {
			if strings.Contains(checkbody, query)  {
				fmt.Printf("%v %s\n", color.GreenString("[+] POC:..✨"), recheck_url)
			} 
		}
	}
}

func payloadlist(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		payloads = append(payloads, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return payloads
}

func Banner() {
	color.HiGreen(`
.____     _____.__   _____          
|    |  _/ ____\__| /     \   ____  
|    |  \   __\|  |/  \ /  \_/ __ \ 
|    |___|  |  |  /    Y    \  ___/ 
|_______ \__|  |__\____|__  /\___  >
        \/                \/     \/  v1.0  
										   `)
	color.HiRed("                 " + "Made with <3 by @shivangx01b")
	
}

func ParseArguments() {
	flag.IntVar(&Threads, "t", 40, "Number of workers to use..default 40. Ex: -t 50")
	flag.StringVar(&payload, "p", "", "Feed the list of payloads to fuzz. Ex: -p ~/wordlists/lfi.txt")  
	flag.StringVar(&method, "method",  "GET", "Add method name if required. Ex: -method PUT. Default \"GET\"")
	flag.BoolVar(&verify, "verify",  false, "Only prints confirmed results. Ex -verify ")
	flag.Parse()
}


func main() {
	ParseArguments()
	Banner()
	checkin, _ := os.Stdin.Stat()
	if checkin.Mode() & os.ModeNamedPipe > 0 {
		if payload != "" {
			list := payloadlist(payload)
			matcher = "nocheck"
			urls := make(chan string, Threads)
			processGroup := new(sync.WaitGroup)
			processGroup.Add(Threads)

			for i := 0; i < Threads; i++ {
				c := getClient()
				go func() {
					defer processGroup.Done()
					for u := range urls {
						requester(c, u, method, list, verify, matcher)
					}
				}()
			}

			sc := bufio.NewScanner(os.Stdin)

			for sc.Scan() {
				urls <- sc.Text()
			}
			close(urls)
			processGroup.Wait()
		} else {
			color.HiRed("\n[!] Must give payload list")
		}
	} else {
		color.HiRed("\n[!] Check: LfiMe -h for arguments")
	}
}
