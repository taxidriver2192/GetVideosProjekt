package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/gocolly/colly/v2"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("micl-easj.dk"),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)

		r, _ := regexp.Compile("(m+p+4)")

		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		if r.MatchString(link) == true {
			d1 := []byte(link + "\ngo\n")
			err := ioutil.WriteFile("/links-txt/"+e.Text+".txt", d1, 0644)
			fmt.Printf("wrote %d bytes\n", d1)
			check(err)
			f, err := os.Create("/links-txt/" + e.Text + ".txt")
			check(err)

			defer f.Close()
		} else {
			c.Visit(e.Request.AbsoluteURL("http://micl-easj.dk/" + link))
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("http://micl-easj.dk/")
}
