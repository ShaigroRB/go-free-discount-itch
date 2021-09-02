package fditch

import (
	"encoding/json"
	"fmt"
)

// arrayToJSON converts an array to a JSON string.
func arrayToJSON(items []Item) (string, error) {
	result, err := json.Marshal(items)
	return string(result), err
}

// getCategoryAllContents returns a list containing all the Content for a category.
// It returns an error with it if any.
func getCategoryAllContents(getCategoryContentFn GetCategoryContentFn) ([]Content, error) {
	var err error

	list := make([]Content, 0)
	page := 0
	for {
		page++
		isLastPage, err := getCategoryContentFn(page, &list)
		if err != nil || isLastPage {
			break
		}
	}

	return list, err
}

// errorToJSON takes an error an returns a json like { 'error': 'this is the error' }.
func errorToJSON(err error) string {
	return fmt.Sprintf("{\"error\": \"%s\"}", err.Error())
}

// getCategoryItems returns all items using a function to get the page contents for a category.
func getCategoryItems(getCategoryContentFn GetCategoryContentFn) ([]Item, error) {
	// initialize to always return a list, even when there is no item
	items := make([]Item, 0)

	// get a list with all the Content for a category
	pageContentList, err := getCategoryAllContents(getCategoryContentFn)

	if err != nil {
		return items, err
	}

	// loop on each element in the list
	for _, pageContent := range pageContentList {

		// for each Content, parse the items in it
		itemsForContent, err := ConvertContentToItems(pageContent)

		if err != nil {
			return items, err
		}

		// add those items to the list of all the items
		for item := range itemsForContent {
			items = append(items, item)
		}
	}

	return items, nil
}

// getCategoryItemsAsJSON returns a JSON string containing all items using a function to get the page contents for a category.
func getCategoryItemsAsJSON(getCategoryContentFn GetCategoryContentFn) string {
	items, err := getCategoryItems(getCategoryContentFn)

	if err != nil {
		return errorToJSON(err)
	}

	// convert all the items to a single json
	result, err := arrayToJSON(items)

	if err != nil {
		return errorToJSON(err)
	}
	return result
}

// setResultForCategory calls a function with a GetCategoryContentFn as parameter depending on a given category.
func setResultForCategory(category Category, setResult func(fn GetCategoryContentFn)) {
	switch category {
	case GameAssets:
		setResult(GetGameAssetsContent)
	case Books:
		setResult(GetBooksContent)
	case Comics:
		setResult(GetComicsContent)
	case Tools:
		setResult(GetToolsContent)
	case Games:
		setResult(GetGamesContent)
	case PhysicalGames:
		setResult(GetPhysicalGamesContent)
	case Soundtracks:
		setResult(GetSoundtracksContent)
	case GameMods:
		setResult(GetGameModsContent)
	case Misc:
		setResult(GetMiscContent)
	}
}

// GetCategoryItems returns a list containing all items for a given category.
func GetCategoryItems(category Category) ([]Item, error) {
	var items []Item
	var err error

	// inner function that gets the items based for the category
	setResult := func(fn GetCategoryContentFn) {
		items, err = getCategoryItems(fn)
	}

	// set the resulting items and error for the given category
	setResultForCategory(category, setResult)

	return items, err
}

// GetCategoryItemsAsJSON returns a JSON string containing all items for a given category.
func GetCategoryItemsAsJSON(category Category) string {
	var itemsAsJSON string

	// inner function that gets the items as json for a given category
	setResult := func(fn GetCategoryContentFn) {
		itemsAsJSON = getCategoryItemsAsJSON(fn)
	}

	// set the resulting json for the given category
	setResultForCategory(category, setResult)

	return itemsAsJSON
}
