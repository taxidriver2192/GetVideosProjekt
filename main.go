package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"regexp"

	"github.com/gocolly/colly/v2"
)

type Video struct {
	Id        string            `json:"id,omitempty"`
	Url       string            `json:"Video,omitempty"`
	InfoVideo map[string]string `json:"artist,omitempty"`
}

func main() {

	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("micl-easj.dk"),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		link := e.Attr("href")

		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		r_validVideoLink, _ := regexp.Compile("(m+p+4)")
		if r_validVideoLink.MatchString(link) == true {
			//matched, noDateError := regexp.Match("/^(?:(?:31(/|-|.)(?:0?[13578]|1[02]))1|(?:(?:29|30)(/|-|.)(?:0?[1,3-9]|1[0-2])2))(?:(?:1[6-9]|[2-9]d)?d{2})$|^(?:29(/|-|.)0?23(?:(?:(?:1[6-9]|[2-9]d)?(?:0[48]|[2468][048]|[13579][26])|(?:(?:16|[2468][048]|[3579][26])00))))$|^(?:0?[1-9]|1d|2[0-8])(/|-|.)(?:(?:0?[1-9])|(?:1[0-2]))4(?:(?:1[6-9]|[2-9]d)?d{2})$/g", []byte(link))
			var info Video
			info.Id = "1"
			u, TestError := url.Parse(link)
			if TestError != nil {
				log.Fatal(TestError)
			}
			info.Url = ("http://micl-easj.dk" + u.EscapedPath())
			m := make(map[string]string)
			info.InfoVideo = m
			info.InfoVideo["firstName"] = "Michael"
			info.InfoVideo["lastName"] = "Claudius"
			info.InfoVideo["email"] = "test@test.dk"
			info.InfoVideo["website"] = "http://micl-easj.dk/"

			r_findDate, _ := regexp.Compile("([0-9]{2}.[0-9]{2}.[0-9]{4})|([0-9]{4}.[0-9]{2}.[0-9]{2})")
			if r_findDate.MatchString(link) == true {
				dato := r_findDate.Copy().Find([]byte(link))
				// Ved ikke hvorfor men nogle gange skriver han YYYY.MM.DD
				// Hvis den giver true, skal jeg lave om på formateringen af datoen.
				r_wongDateFormate, _ := regexp.Compile("([0-9]{4}.[0-9]{2}.[0-9]{2})")
				if r_wongDateFormate.MatchString(link) == true {
					year := fmt.Sprint(string(dato[0:4]))
					month := fmt.Sprint(string(dato[5:7]))
					date := string(dato[len(dato)-2:])
					newDate := date + "." + month + "." + year
					info.InfoVideo["date"] = string(newDate)
				} else {
					// Her har den fundet datoeen men den står ikke forkert.
					info.InfoVideo["date"] = string(dato)
				}
			} else {
				// Her har den ikke fundet er dato. og sætter derfpr datoem til N/A
				info.InfoVideo["date"] = "N/A"
			}

			var data []byte
			data, _ = json.MarshalIndent(info, "", "    ")

			fmt.Println(string(data))
			// d1 := []byte(link + "\ngo\n")
			// err := ioutil.WriteFile("/links-txt/"+e.Text+".txt", d1, 0644)
			// fmt.Printf("wrote %d bytes\n", d1)
			// check(err)
			// f, err := os.Create("/links-txt/" + e.Text + ".txt")
			// check(err)

			// defer f.Close()
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
