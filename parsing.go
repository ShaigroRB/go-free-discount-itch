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

// isNodeAGame checks if a html node is a game node
func isNodeAGame(node *html.Node) bool {
	if len(node.Attr) == 0 {
		return false
	}
	if node.DataAtom == atom.Div {
		divPossibleAttrs := getDivPossibleAttrs(node)
		return divPossibleAttrs.dataGameID != ""
	}
	return false
}

type divPossibleAttrs struct {
	dataGameID string
	class      string
}

// getDivPossibleAttrs returns the possible attributes for a <div> node.
func getDivPossibleAttrs(node *html.Node) divPossibleAttrs {
	var divPossibleAttrs divPossibleAttrs
	for _, attr := range node.Attr {
		if attr.Key == "data-game_id" {
			divPossibleAttrs.dataGameID = attr.Val
		}
		if attr.Key == "class" {
			divPossibleAttrs.class = attr.Val
		}
	}
	return divPossibleAttrs
}

type aPossibleAttrs struct {
	class string
	href  string
}

// getAPossibleAttrs returns the possible attributes for an <a> node.
func getAPossibleAttrs(node *html.Node) aPossibleAttrs {
	var aPossibleAttrs aPossibleAttrs
	for _, attr := range node.Attr {
		if attr.Key == "class" {
			aPossibleAttrs.class = attr.Val
		}
		if attr.Key == "href" {
			aPossibleAttrs.href = attr.Val
		}
	}
	return aPossibleAttrs
}

type imgPossibleAttrs struct {
	dataLazySrc string
}

// getImgPossibleAttrs returns the possible attributes for an <img> node.
func getImgPossibleAttrs(node *html.Node) imgPossibleAttrs {
	var imgPossibleAttrs imgPossibleAttrs
	for _, attr := range node.Attr {
		if attr.Key == "data-lazy_src" {
			imgPossibleAttrs.dataLazySrc = attr.Val
		}
	}
	return imgPossibleAttrs
}

type spanPossibleAttrs struct {
	class string
}

// getSpanPossibleAttrs returns the possible attributes for an <span> node.
func getSpanPossibleAttrs(node *html.Node) spanPossibleAttrs {
	var spanPossibleAttrs spanPossibleAttrs
	for _, attr := range node.Attr {
		if attr.Key == "class" {
			spanPossibleAttrs.class = attr.Val
		}
	}
	return spanPossibleAttrs
}

func printNode(root *html.Node, depth int) {
	fmt.Printf("%*s%s\n", depth*2, "", root.Data)
	for _, attr := range root.Attr {
		fmt.Printf("%*s%s=\"%s\"\n", depth*2, "", attr.Key, attr.Val)
	}
	for child := root.FirstChild; child != nil; child = child.NextSibling {
		printNode(child, depth+1)
	}
}

func printNodes(nodes []*html.Node, depth int) {
	for _, node := range nodes {
		printNode(node, depth)
	}
}

// nodeToItemsWithoutEndDate returns a channel full of incomplete Item based on a node.
// Those Item are missing the end date of the sale.
func nodeToItemsWithoutEndDate(root *html.Node, maxItems int) chan Item {
	// no problem with channel even without coroutine because the size is specified
	cells := make(chan Item, maxItems)
	nodes := getPreOrderQueue(root)

	var cell Item

	continueTillNextGame := true

	for _, node := range nodes {
		if continueTillNextGame && !isNodeAGame(node) {
			continue
		} else {
			continueTillNextGame = false
		}

		switch node.DataAtom {
		case atom.Div:
			divPossibleAttrs := getDivPossibleAttrs(node)
			if divPossibleAttrs.dataGameID != "" {
				// consider that the item has been fully parsed
				if cell.ID != "" {
					cells <- cell
				}
				cell = Item{Platforms: []string{}}
				cell.ID = divPossibleAttrs.dataGameID
			} else if divPossibleAttrs.class != "" {
				switch divPossibleAttrs.class {
				case "game_genre":
					cell.Genre = node.FirstChild.Data
				case "game_author":
					cell.Author = node.FirstChild.FirstChild.Data
				case "game_text":
					cell.Description = node.FirstChild.Data
				case "sale_tag":
					if node.FirstChild.Data != "-100%" {
						continueTillNextGame = true
						cell.ID = ""
					}
				// cuz yes, reverse sales are a thing in itch.io
				case "sale_tag reverse_sale":
					if node.FirstChild.Data != "-100%" {
						continueTillNextGame = true
						cell.ID = ""
					}
				}
			}
		case atom.A:
			aPossibleAttrs := getAPossibleAttrs(node)
			if aPossibleAttrs.class != "" {
				switch aPossibleAttrs.class {
				case "title game_link":
					cell.Link = aPossibleAttrs.href
					cell.Title = node.FirstChild.Data
				case "price_tag meta_tag sale":
					cell.SalesLink = aPossibleAttrs.href
				}
			}
		case atom.Img:
			imgPossibleAttrs := getImgPossibleAttrs(node)
			if imgPossibleAttrs.dataLazySrc != "" {
				cell.ImgLink = imgPossibleAttrs.dataLazySrc
			}
		case atom.Span:
			spanPossibleAttrs := getSpanPossibleAttrs(node)
			switch spanPossibleAttrs.class {
			case "icon icon-windows8":
				cell.Platforms = append(cell.Platforms, "Windows")
			case "icon icon-apple":
				cell.Platforms = append(cell.Platforms, "macOS")
			case "icon icon-tux":
				cell.Platforms = append(cell.Platforms, "Linux")
			case "icon icon-android":
				cell.Platforms = append(cell.Platforms, "Android")
			case "web_flag":
				cell.Platforms = append(cell.Platforms, "Web")
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
	// if end_date is not at the end of the body
	regx := regexp.MustCompile(`end_date\".*\",`)
	matches := regx.FindStringSubmatch(body)
	// if end_date is at the end of the body
	if len(matches) == 0 {
		regx = regexp.MustCompile(`end_date\".*\"}`)
		matches = regx.FindStringSubmatch(body)
	}

	if len(matches) == 0 {
		fmt.Printf(`
		Function: parseEndDate
		Context:
		- body: %s

		Error: End date for the item was not found.\n`, body)

		return "End date for the item was not found. Please report this bug to https://github.com/ShaigroRB/go-free-discount-itch/issues"
	}

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
		fmt.Printf(`
		Function: ConvertContentToItems::html.Parse
		Context:
		- content: %s

		Error: %s\n`, content.Content, err)
		return nil, err
	}

	partialItems := nodeToItemsWithoutEndDate(node, content.NumItems)
	items := make(chan Item, len(partialItems))

	defer close(items)

	for partialItem := range partialItems {
		body, err := getSales(partialItem.SalesLink)
		if err != nil {
			fmt.Printf(`
			Function: ConvertContentToItems::getSales
			Context:
			- salesLink: %s

			Error: %s\n`, partialItem.SalesLink, err)
			return items, err
		}

		partialItem.EndDate = parseEndDate(body)

		items <- partialItem
	}

	return items, nil
}
