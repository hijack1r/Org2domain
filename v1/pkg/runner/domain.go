package runner

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func MatchDomain(url1 string) (string, error) {

	decodedURL, err := url.QueryUnescape(url1)
	if err != nil {
		return "", fmt.Errorf("Error decoding URL: %v", err)
	}

	if !(strings.HasPrefix(decodedURL, "http")) {
		return "", fmt.Errorf("Domain Not Found!")
	}

	resp, err := http.Get(decodedURL)
	if err != nil {
		return "", fmt.Errorf("Error sending GET request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %v", err)
	}

	responseBody := string(body)

	re := regexp.MustCompile(`<td class="infobox-data"><span class="url"><a rel="nofollow" class="external text" href="https?://([^"]+)` +
		`|><a rel="nofollow" class="external text" href="https?://([^"]+)`)

	matches := re.FindStringSubmatch(responseBody)

	if len(matches) > 0 {
		if matches[1] != "" {
			domain := matches[1]
			if strings.Contains(domain, "/") {
				domain = strings.Split(domain, "/")[0]
			}

			trimmedDomain := strings.TrimRight(domain, "/")
			return trimmedDomain, nil

		} else if matches[2] != "" {
			domain := matches[2]
			if strings.Contains(domain, "/") {
				domain = strings.Split(domain, "/")[0]
			}

			trimmedDomain := strings.TrimRight(domain, "/")
			return trimmedDomain, nil

		} else {
			fmt.Println("Domain name not matched in body")
		}
	} else {
		fmt.Println("No relevant field found")
	}

	return "", fmt.Errorf("No domain found in the response")
}
