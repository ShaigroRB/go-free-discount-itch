package fditch

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// itch.io hostname
	hostname = "https://itch.io"

	// default parameters when doing API calls
	onsaleParams = "/on-sale?format=json&page"

	// max retries for API calls
	maxRetries = 5
)

// fetchDataWithRetry returns the response of an HTTP GET request for a given URL.
// It returns the response and a detailed error if any.
func fetchDataWithRetry(url string) (*http.Response, error) {
	for retries := 0; retries < maxRetries; retries++ {
		resp, err := http.Get(url)

		if err != nil {
			return nil, err
		}

		if resp.StatusCode == 429 && retries < maxRetries {
			fmt.Printf("Failed to fetch data from %s. Status code 429. Retrying in 2 seconds...\n", url)
			time.Sleep(2 * time.Second)
			continue
		}

		if resp.StatusCode != 200 && resp.StatusCode != 429 {
			return nil, fmt.Errorf("Failed to fetch data from %s. Status code %d, ", url, resp.StatusCode)
		}

		return resp, nil
	}
	return nil, fmt.Errorf("Failed to fetch data from %s after %d retries", url, maxRetries)
}

// getJSON returns the content of a page for a given category.
// It returns the JSON as a string and an error if any.
func getJSON(category string, page int) (string, error) {
	url := fmt.Sprintf("%s/%s%s=%d", hostname, category, onsaleParams, page)

	resp, err := fetchDataWithRetry(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

// getSales returns the content of a sales page and an error if any.
func getSales(link string) (string, error) {
	url := fmt.Sprintf("%s%s", hostname, link)
	resp, err := fetchDataWithRetry(url)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// getContent puts the content of a page in a list for a given category.
// It returns whether it was the last page and an error if any.
func getContent(category string, page int, list *[]Content) (isLastPage bool, err error) {
	json, err := getJSON(category, page)
	if err != nil {
		fmt.Printf(`
		Function: getContent::getJSON
		Context:
		- category: %s
		- page:     %d

		Error: %s
		`, category, page, err)
		return isLastPage, err
	}

	content := Content{}
	err = content.FromJSON(json)
	if err != nil {
		fmt.Printf(`
		Function: getContent::content.FromJSON
		Context:
		- category: %s
		- page:     %d
		- json:	    %s

		Error: %s
		`, category, page, json, err)
		return isLastPage, err
	}

	*list = append(*list, content)

	isLastPage = content.NumItems < 30

	return isLastPage, nil
}

// Type that represents a function to get a Content for a specific category.
type GetCategoryContentFn func(int, *[]Content) (bool, error)

// GetGameAssetsContent puts in a list the `game-assets` type content for a given page.
// It returns whether it was the last page and an error if any.
func GetGameAssetsContent(page int, list *[]Content) (isLastPage bool, err error) {
	return getContent("game-assets", page, list)
}

// GetBooksContent puts in a list the `books` type content for a given page.
// It returns whether it was the last page and an error if any.
func GetBooksContent(page int, list *[]Content) (isLastPage bool, err error) {
	return getContent("books", page, list)
}

// GetComicsContent puts in a list the `comics` type content for a given page.
// It returns whether it was the last page and an error if any.
func GetComicsContent(page int, list *[]Content) (isLastPage bool, err error) {
	return getContent("comics", page, list)
}

// GetToolsContent puts in a list the `tools` type content for a given page.
// It returns whether it was the last page and an error if any.
func GetToolsContent(page int, list *[]Content) (isLastPage bool, err error) {
	return getContent("tools", page, list)
}

// GetGamesContent puts in a list the `games` type content for a given page.
// It returns whether it was the last page and an error if any.
func GetGamesContent(page int, list *[]Content) (isLastPage bool, err error) {
	return getContent("games", page, list)
}

// GetPhysicalGamesContent puts in a list the `physical-games` type content for a given page.
// It returns whether it was the last page and an error if any.
func GetPhysicalGamesContent(page int, list *[]Content) (isLastPage bool, err error) {
	return getContent("physical-games", page, list)
}

// GetSoundstracksContent puts in a list the `soundtracks` type content for a given page.
// It returns whether it was the last page and an error if any.
func GetSoundtracksContent(page int, list *[]Content) (isLastPage bool, err error) {
	return getContent("soundtracks", page, list)
}

// GetGameModsContent puts in a list the `game-mods` type content for a given page.
// It returns whether it was the last page and an error if any.
func GetGameModsContent(page int, list *[]Content) (isLastPage bool, err error) {
	return getContent("game-mods", page, list)
}

// GetMiscContent puts in a list the `misc` type content for a given page.
// It returns whether it was the last page and an error if any.
func GetMiscContent(page int, list *[]Content) (isLastPage bool, err error) {
	return getContent("misc", page, list)
}
