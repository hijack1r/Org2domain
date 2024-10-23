package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"https://github.com/hijack1r/Org2domain/v1/pkg/runner"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	fmt.Println(`
  ____            ___     _                       _       
 / __ \          |__ \   | |                     (_)      
| |  | |_ __ __ _   ) |__| | ___  _ __ ___   __ _ _ _ __  
| |  | | "__/ _" | / // _" |/ _ \| "_ " _ \ / _" | | "_ \ 
| |__| | | | (_| |/ /| (_| | (_) | | | | | | (_| | | | | |
 \____/|_|  \__  |____\__ _|\___/|_| |_| |_|\__ _|_|_| |_|
             __/ |                                        
            |___/
                      			Author: @h_ijack1r
                      			Github: https://github.com/hijack1r`)

	fileName := flag.String("f", "", "Enter the target file name.")
	proxyFlag := flag.String("p", "", "If need proxy to access Google, please set the proxy address.")
	outfileName := flag.String("o", "", "Exported file name.")
	flag.Parse()

	var err error
	orgFile, err := os.Open(*fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer orgFile.Close()

	if *proxyFlag != "" {
		proxyAddress := strings.TrimSpace(*proxyFlag)
		proxyURL, err := url.Parse(proxyAddress)
		if err != nil {
			fmt.Println("Error parsing proxy URL:", err)
			return
		}
		http.DefaultTransport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		fmt.Printf("Use Proxy ：%s\n", proxyAddress)
	} else {
		fmt.Println("No proxy used")
	}

	var outfile *os.File
	if *outfileName != "" {
		if strings.HasSuffix(*outfileName, ".csv") {
			outfile, err = os.Create(*outfileName)
			if err != nil {
				fmt.Println("Error opening file:", err)
				return
			}
			defer outfile.Close()
		} else {
			fmt.Printf("FileName %s invalid，Please use .csv suffix fileName\n", *outfileName)
			return
		}
	} else {
		fmt.Println("File name cannot be blank. Please enter file name", err)
		return
	}

	if outfile == nil {
		fmt.Println("The output file was not created correctly")
		return
	}

	_, err = outfile.Write([]byte{0xEF, 0xBB, 0xBF})
	if err != nil {
		fmt.Println("Error writing BOM to file:", err)
		return
	}

	Writer := csv.NewWriter(outfile)
	defer Writer.Flush()

	header := []string{"OrgName", "Domain"}
	err = Writer.Write(header)
	if err != nil {
		fmt.Println("Error writing header to CSV:", err)
		return
	}

	orgScanner := bufio.NewScanner(orgFile)
	content, err := ioutil.ReadFile(*fileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	lines := strings.Split(string(content), "\n")

	var nonEmptyLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}

	lineCount := len(nonEmptyLines)
	currentIndex := 0

	fmt.Println()
	fmt.Println("Exporting domain name. Please wait...")

	for orgScanner.Scan() {
		searchKeyword := orgScanner.Text()
		searchKeyword = strings.TrimSpace(searchKeyword)

		searchURL := fmt.Sprintf("https://www.google.com/search?q=%s&num=%d&sca_esv=d265e0c6ccb16526&ei=v4k5Zp_SJKzA1e8P-5i9kAs&ved=0ahUKEwjfrf-wt_qFAxUsYPUHHXtMD7IQ4dUDCBA&uact=5&oq=%s&gs_lp=Egxnd3Mtd2l6LXNlcnAiO3NpdGU6ZW4ud2lraXBlZGlhLm9yZyBBbWVyaWNhbiBBY2FkZW15IG9mIEZvcmVuc2ljIFNjaWVuY2VzSLYFUABYAHABeACQAQCYAQCgAQCqAQC4AQPIAQCYAgCgAgCYAwCIBgGSBwCgBwA&sclient=gws-wiz-serp", url.QueryEscape(searchKeyword+" wiki"), 20, url.QueryEscape(searchKeyword+" wiki"))

		response, err := http.Get(searchURL)
		if err != nil {
			fmt.Println("Error fetching search results:", err)
			return
		}
		defer response.Body.Close()

		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			fmt.Println("Error parsing HTML:", err)
			return
		}

		var urls []string
		urlMap := make(map[string]bool)
		doc.Find("a[href]").Each(func(index int, element *goquery.Selection) {
			href, exists := element.Attr("href")
			if exists && strings.HasPrefix(href, "/url?q=") {
				re := regexp.MustCompile(`^/url\?q=(.+?)&`)
				match := re.FindStringSubmatch(href)
				if len(match) == 2 {
					decodedURL, err := url.QueryUnescape(match[1])
					if err != nil {
						decodedURL = match[1]
					}
					if strings.HasPrefix(decodedURL, "http://") || strings.HasPrefix(decodedURL, "https://") {
						if !urlMap[decodedURL] {
							urlMap[decodedURL] = true
							urls = append(urls, decodedURL)
						}
					}
				}
			}
		})

		var found bool
		var filteredUrls []string
		for _, url := range urls {
			if !strings.HasPrefix(url, "https://maps.google.com/") {
				filteredUrls = append(filteredUrls, url)
			}
		}
		urls = filteredUrls

		if runner.CheckEn(searchKeyword) {
			for _, searchUrl := range urls {
				searchUrl = strings.TrimSpace(searchUrl)
				if strings.HasPrefix(searchUrl, "https://en.wikipedia.org") {
					found = true
					domain, err := runner.MatchDomain(searchUrl)
					if err != nil {
						fmt.Println("Error match:", err)
						return
					}

					data := []string{searchKeyword, domain}
					err = Writer.Write(data)
					if err != nil {
						fmt.Println("Error writing data to CSV:", err)
					}
					break
				} else {
					fmt.Println("Not find https://en.wikipedia.org begin")
					data := []string{searchKeyword, "No org was found."}
					err = Writer.Write(data)
					if err != nil {
						fmt.Println("Error writing data to CSV:", err)
					}
				}
			}

		} else {
			for _, searchUrl := range urls {
				if strings.HasPrefix(searchUrl, "https://zh.wikipedia.org") {
					found = true

					domain, err := runner.MatchDomain(searchUrl)
					if err != nil {
						fmt.Println("Error match:", err)
						return
					}

					data := []string{searchKeyword, domain}
					err = Writer.Write(data)
					if err != nil {
						fmt.Println("Error writing data to CSV:", err)
					}
					break
				}
			}

			if !found {
				data := []string{searchKeyword, "No org was found."}
				err = Writer.Write(data)
				if err != nil {
					fmt.Println("Error writing data to CSV:", err)
				}
			}
		}

		if currentIndex < lineCount {
			currentIndex++
		}
		runner.PrintProgress(currentIndex, lineCount)
	}
	fmt.Println()
	fmt.Printf("[+] The domain in %s has been matched,Output in %s.", *fileName, *outfileName)
	if err := orgScanner.Err(); err != nil {
		fmt.Printf("Error scanning %s file:", *fileName, err)
		return
	}
}
