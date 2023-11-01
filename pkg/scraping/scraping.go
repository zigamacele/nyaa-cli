package scraping

import (
	"fmt"
	"github.com/gocolly/colly"
	"nyaa-cli/pkg/scraping/types"
	"strings"
)

//WEB SCRAPING

var Torrents []types.Torrent

func Run() {
	c := colly.NewCollector()

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		var torrent []string
		e.ForEach("td", func(_ int, td *colly.HTMLElement) {
			td.ForEach("a", func(_ int, a *colly.HTMLElement) {
				link := a.Attr("href")
				if strings.Contains(link, "/?c=") {
					link = strings.ReplaceAll(link, "/?c=", "")
				}
				if strings.Contains(link, "#comments") {
					return
				}

				//fmt.Printf("Link: %s\n", link)
				torrent = append(torrent, link)
			})

			var tdText = strings.ReplaceAll(td.Text, "\n", "")
			tdText = strings.ReplaceAll(tdText, "\t", "")
			if tdText != "" {
				torrent = append(torrent, tdText)
				//fmt.Printf("%q\n", tdText)
			}

		})

		Torrents = append(Torrents, types.Torrent{
			Category:  torrent[0],
			View:      torrent[1],
			Name:      torrent[2],
			Download:  torrent[3],
			Magnet:    torrent[4],
			Size:      torrent[5],
			Date:      torrent[6],
			Seeders:   torrent[7],
			Leechers:  torrent[8],
			Downloads: torrent[9],
		})

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://nyaa.si")
}
