package scraper

import (
	"bufio"
	"fmt"
	"net/http"
)

func Search(title string) error {
	url := "https://clipzag.com/search?q=" + title
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return nil
}
