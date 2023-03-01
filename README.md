# go-free-discount-itch

[itch.io](https://itch.io) is a marketplace for games, game assets, ...  
The items can be on discount and even be 100% on discount, essentially making the items free.

The go-free-discount-itch package (_fditch for short_) is a package that exposes methods to get all items on 100% discount of itch.io.

## How to use

### Get the package

`go get -u github.com/ShaigroRB/go-free-discount-itch`

### Examples

```go
package main

import (
    "fmt"
    fditch "github.com/ShaigroRB/go-free-discount-itch"
)

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

## Structure of one cell (_item_)

How I represent it:  
html_element interesting_attribute (why) [text\] (??)

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

<details><summary>Example of an item as HTML</summary>

```html
<div data-game_id="1325209" class="game_cell has_cover lazy_images" dir="auto">
  <a
    tabindex="-1"
    class="thumb_link game_link"
    href="https://raidgames-studios.itch.io/sinsfromgod2"
    data-label="game:1325209:thumb"
    data-action="game_grid"
  >
    <div
      class="game_thumb"
      data-background_image="https://img.itch.zone/aW1nLzc4MTA1NzEuanBn/315x250%23c/%2BDiYI6.jpg"
      style="background-color: #000;"
    ></div>
  </a>
  <div class="game_cell_tools">
    <a
      data-register_action="add_to_collection"
      href="/g/raidgames-studios/sinsfromgod2/add-to-collection?source=browse"
      class="action_btn add_to_collection_btn"
    >
      <span class="icon icon-playlist_add"></span>
      Add to collection
    </a>
  </div>
  <div class="game_cell_data">
    <div class="game_title">
      <a
        class="title game_link"
        href="https://raidgames-studios.itch.io/sinsfromgod2"
        data-label="game:1325209:title"
        data-action="game_grid"
      >
        SinsFromGod 2
      </a>
      <a
        href="/s/63424/all-games-100-off"
        title="Pay $0 or more for this Game"
        class="price_tag meta_tag sale"
      >
        <div class="price_value">$0</div>
        <div class="sale_tag">-100%</div>
      </a>
    </div>
    <div title="Horror Game Based In A Hotel." class="game_text">
      Horror Game Based In A Hotel.
    </div>
    <div class="game_author">
      <a
        href="https://raidgames-studios.itch.io"
        data-label="user:4785070"
        data-action="game_grid"
      >
        RaidGames Studios
      </a>
    </div>
    <div class="game_genre">Survival</div>
    <div class="game_platform">
      <span title="Download for Windows" class="icon icon-windows8"></span>
      <span title="Download for Linux" class="icon icon-tux"></span>
    </div>
  </div>
</div>
```

</details>

## TODO list

1. [ ] Get max amount of results
   1. [ ] Get items/on-sale
   2. [ ] Parse to get the number of results (`curl "https://itch.io/items/on-sale" | grep -i "<nobr class=\"game\_count\".*</nobr>"`)
2. [x] Get the json of all pages
   1. [x] Get category/on-sale?format=json&page=.. for each json
   2. [x] Put each one in a struct (PageContent)
3. [x] Parse the content as incomplete items (_missing the end date for sales_)
   1. [x] Read the content as html nodes
   2. [x] Construct items based on the nodes
      1. Split the nodes to keep only the "_game_cells_" nodes (_Spoiler: it's not worth it_)
      2. [x] Create items from those "_game_cells_" nodes. Don't forget to only keep the **100%** on sales items. (_careful of +100% sales_)
4. [x] Get the end date for each item
   1. [x] Get html content from the sales link
   2. [x] Parse it to get the end date for the sales (`{"start_date":"2021-05-28T10:00:35Z","id":50563,"end_date":"2021-05-30T10:02:59Z","can_be_bought":true,"actual_price":398}`)
5. [x] Create a JSON string out of all the items
6. [x] Keep it simple by removing all concurrency. Any concurrency should be done by the user of the package.
7. [x] Method to get all items for a category as a list of Items
8. [ ] Tests

## License

This project is under the MIT license.
