package garkov

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fogleman/gg"
)

func Garkov() string {
	comic, err := getComic()
	if err != nil {
		panic(err)
	}

	base := initializeImage()

	cb := comic.Find(".commentblock")
	css, _ := cb.Attr("style")
	i := strings.Index(css, ".gif")
	if i < 0 {
		panic("Cannot find comic number")
	}
	number, err := strconv.Atoi(css[i-3 : i])
	if err != nil {
		panic(err)
	}
	drawStrip(base, number)

	text := cb.Find(".comment")
	ch := text.First().Children()

	topOffset := 33
	leftOffset := 3
	base.SetRGB(0, 0, 0)
	ch.Each(func(i int, el *goquery.Selection) {
		style, _ := el.Attr("style")
		positions := strings.Split(style, ";")[2:]
		top := 0
		left := 0
		for i, el := range positions {
			if i == 0 {
				el = strings.Trim(el, "top: x")
				top, _ = strconv.Atoi(el)
			} else if i == 1 {
				el = strings.Trim(el, "left: px")
				left, _ = strconv.Atoi(el)
			}
		}
		top += topOffset
		left += leftOffset
		letters := el.Children()
		letters.Each(func(j int, el *goquery.Selection) {
			attr, _ := el.Attr("src")
			path := fmt.Sprintf("resources/%s", strings.ReplaceAll(attr, ".gif", ".png"))
			letter, err := gg.LoadPNG(path)
			if err != nil {
				panic(err)
			}
			width := letter.Bounds().Size().X
			base.DrawImage(letter, left, top)
			left += width
		})
	})

	os.Mkdir("temp/", 0755)
	name := randomName(15)
	base.SavePNG("temp/" + name)
	return name
}

func initializeImage() *gg.Context {
	img := gg.NewContext(606, 211)
	img.SetRGB(1, 1, 1)
	img.DrawRectangle(0, 0, 606, 211)
	img.Fill()

	banner, err := gg.LoadPNG("resources/headmast.png")
	if err != nil {
		panic(err)
	}
	img.DrawImage(banner, 3, 0)

	return img
}

func drawStrip(img *gg.Context, number int) {
	bgPath := fmt.Sprintf("resources/%03d.png", number)
	strip, err := gg.LoadPNG(bgPath)
	if err != nil {
		panic(err)
	}
	img.DrawImage(strip, 3, 33)
}

func getComic() (*goquery.Selection, error) {
	res, err := http.Get("http://joshmillard.com/garkov/")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("non 200 status code: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	comic := doc.Find(".comicborder")
	return comic, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomName(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b) + ".png"
}
