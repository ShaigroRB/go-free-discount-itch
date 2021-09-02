# go-free-discount-itch
[itch.io](https://itch.io) is a marketplace for games, game assets, ...  
The items can be on discount and even be 100% on discount, essentially making the items free.

The go-free-discount-itch package (*fditch for short*) is a package that exposes methods to get all items on 100% discount of itch.io.

## How to use
### Get the package
`go get github.com/ShaigroRB/go-free-discount-itch`

### Examples
```go
func main() {
    // get all items that are on 100% discount as json
    for _, category := range fditch.Categories {
        jsonString := fditch.GetCategoryItemsAsJSON(category)
        fmt.Println(jsonString)
    }

    // get all the games on 100% discount
    games, err := fditch.GetCategoryItems(fditch.Games)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, game := range games {
            // print only the link of the game
			fmt.Println(game.Link)
		}
	}
}
```


## Structure of one cell (*item*)
How I represent it:  
html_element interesting_attribute (why) [text] (??)

- div data-game_id
    - a href (game link)
        - div data-background_image (image for the cell)
            - div ??
            - div data-gif (gif for the cell) ??
    - div
        - a
            - span
    - div
        - div
            - a href (game link) [game title]
            - a href (sales link)
                - div [price value]
                - div [sale tag]
        - div title (game description) [game description]
        - div
            - a href (game author) [game author]
        - div ?
            - span title (windows) ??
            - span title (linux) ??
            - span title (mac) ??

## TODO list
1. [ ] Get max amount of results
    1. [ ] Get items/on-sale
    2. [ ] Parse to get the number of results (`curl "https://itch.io/items/on-sale" | grep -i "<nobr class=\"game\_count\".*</nobr>"`)
2. [x] Get the json of all pages
    1. [x] Get category/on-sale?format=json&page=.. for each json
    2. [x] Put each one in a struct (PageContent)
3. [x] Parse the content as incomplete items (*missing the end date for sales*)
    1. [x] Read the content as html nodes
    2. [x] Construct items based on the nodes
        1. Split the nodes to keep only the "*game_cells*" nodes (*Spoiler: it's not worth it*)
        2. [x] Create items from those "*game_cells*" nodes. Don't forget to only keep the **100%** on sales items. (*careful of +100% sales*)
4. [x] Get the end date for each item
    1. [x] Get html content from the sales link
    2. [x] Parse it to get the end date for the sales (`{"start_date":"2021-05-28T10:00:35Z","id":50563,"end_date":"2021-05-30T10:02:59Z","can_be_bought":true,"actual_price":398}`)
5. [x] Create a JSON string out of all the items
6. [x] Keep it simple by removing all concurrency. Any concurrency should be done by the user of the package.
7. [x] Method to get all items for a category as a list of Items
8. [ ] Tests

## License
This project is under the MIT license.