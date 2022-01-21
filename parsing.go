package fditch

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// getPreOrderQueue traverses an *html.Node in preorder to fill in a queue.
func getPreOrderQueue(root *html.Node) []*html.Node {
	queue := make([]*html.Node, 0)
	res := make([]*html.Node, 0)

	queue = append(queue, root)

	for len(queue) > 0 {
		node := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		res = append(res, node)

		for child := node.LastChild; child != nil; child = child.PrevSibling {
			queue = append(queue, child)
		}
	}

	return res
}

// nodeToItemsWithoutEndDate returns a channel full of incomplete Item based on a node.
// Those Item are missing the end date of the sale.
func nodeToItemsWithoutEndDate(root *html.Node, maxItems int) chan Item {
	// no problem with channel even without coroutine because the size is specified
	cells := make(chan Item, maxItems)
	nodes := getPreOrderQueue(root)

	var cell Item

	isNodeAGame := func(n *html.Node) bool {
		if len(n.Attr) == 0 {
			return false
		}
		return n.DataAtom == atom.Div && n.Attr[0].Key == "data-game_id"
	}

	continueTillNextGame := true

	for _, node := range nodes {
		if continueTillNextGame && !isNodeAGame(node) {
			continue
		} else {
			continueTillNextGame = false
		}

		switch node.DataAtom {
		case atom.Div:
			if len(node.Attr) > 0 {
				if attr := node.Attr[0]; attr.Key == "data-game_id" {
					if cell.ID != "" {
						cells <- cell
					}
					cell = Item{}
					cell.ID = attr.Val
				} else if attr.Key == "class" && attr.Val == "game_author" {
					cell.Author = node.FirstChild.FirstChild.Data
				} else if attr.Key == "class" &&
					(attr.Val == "sale_tag" || attr.Val == "sale_tag reverse_sale") {
					// cause yes, reverse sales are a thing in itch.io
					if node.FirstChild.Data != "-100%" {
						continueTillNextGame = true
						cell.ID = ""
					}
				} else if len(node.Attr) > 1 {
					if attr := node.Attr[1]; attr.Key == "data-background_image" {
						cell.ImgLink = attr.Val
					} else if attr.Key == "class" && attr.Val == "game_text" {
						cell.Description = node.FirstChild.Data
					}
				}
			}
		case atom.A:
			if len(node.Attr) > 1 {
				if attr := node.Attr[0]; attr.Key == "class" && attr.Val == "title game_link" {
					cell.Link = node.Attr[1].Val
					cell.Title = node.FirstChild.Data
				} else if len(node.Attr) > 2 {
					if attr = node.Attr[2]; attr.Key == "class" && attr.Val == "price_tag meta_tag sale" {
						cell.SalesLink = node.Attr[0].Val
					}
				}
			}
		default:
			continue
		}
	}

	close(cells)
	return cells
}

// parseEndDate looks for the end date hidden in the body of a sales page.
func parseEndDate(body string) string {
	regx := regexp.MustCompile(`end_date\".*\",`)
	matches := regx.FindStringSubmatch(body)
	regx = regexp.MustCompile(`[0-9]+-[^\"]*`)
	matches = regx.FindStringSubmatch(matches[0])

	return matches[0]
}

// ConvertContentToItems converts a Content to a channel full of Item.
// Only the Items at -100% sales will be kept.
// It also does the needed API calls to get the end date for each Item.
// It may return an error if any arises.
func ConvertContentToItems(content Content) (chan Item, error) {
	reader := strings.NewReader(content.Content)
	node, err := html.Parse(reader)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	partialItems := nodeToItemsWithoutEndDate(node, content.NumItems)
	items := make(chan Item, len(partialItems))

	defer close(items)

	for partialItem := range partialItems {
		body, err := getSales(partialItem.SalesLink)
		if err != nil {
			fmt.Print(err)
			return items, err
		}

		partialItem.EndDate = parseEndDate(body)

		items <- partialItem
	}

	return items, nil
}
