package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"sync"

	"github.com/gocolly/colly/v2"
)

type autoInc struct {
	sync.Mutex
	id int
}

func (a *autoInc) ID() (id int) {
	a.Lock()
	defer a.Unlock()
	id = a.id
	a.id++
	return
}

type Video struct {
	Id    int
	Titel string
	Url   string
	Dato  string
}

type Person struct {
	Id      int
	Navn    string
	Website string
	Videos  []Video
}

var ai autoInc

func NewPerson() *Person {
	return &Person{
		Id: ai.ID(),
	}
}

func (person *Person) AddVideo(item Video) []Video {
	person.Videos = append(person.Videos, item)
	return person.Videos
}

func main() {
	Videos := []Video{}
	person := Person{
		Id:      NewPerson().Id,
		Navn:    "Micl",
		Website: "test.dk",
		Videos:  Videos,
	}

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
		r_findDate, _ := regexp.Compile("([0-9]{2}.[0-9]{2}.[0-9]{4})|([0-9]{4}.[0-9]{2}.[0-9]{2})")
		OtherDatoField, _ := regexp.Compile("/^(?:(?:31(/|-|.)(?:0?[13578]|1[02]))1|(?:(?:29|30)(/|-|.)(?:0?[1,3-9]|1[0-2])2))(?:(?:1[6-9]|[2-9]d)?d{2})$|^(?:29(/|-|.)0?23(?:(?:(?:1[6-9]|[2-9]d)?(?:0[48]|[2468][048]|[13579][26])|(?:(?:16|[2468][048]|[3579][26])00))))$|^(?:0?[1-9]|1d|2[0-8])(/|-|.)(?:(?:0?[1-9])|(?:1[0-2]))4(?:(?:1[6-9]|[2-9]d)?d{2})$/g")

		if r_validVideoLink.MatchString(link) == true {
			// info.ID = int(ai.ID())
			videoTitel := e.Text
			videoDato := ""
			u, _ := url.Parse(link)
			videoUrl := ("http://micl-easj.Dk" + u.EscapedPath())
			if r_findDate.MatchString(link) == true {
				dato := r_findDate.Copy().Find([]byte(link))
				// Ved ikke hvorfor men nogle gange skriver han YYYY.MM.DD
				// Hvis den giver true, skal jeg lave om på formateringen af datoen.
				r_wongDateFormate, _ := regexp.Compile("([0-9]{4}.[0-9]{2}.[0-9]{2})")
				if r_wongDateFormate.MatchString(link) == true {
					year := fmt.Sprint(string(dato[:+4]))
					month := fmt.Sprint(string(dato[5:7]))
					date := string(dato[len(dato)-2:])
					videoDato = (date + "." + month + "." + year)
				} else {
					// Her har den fundet datoeen men den står ikke forkert.
					videoDato = string(dato)
				}
			} else if OtherDatoField.MatchString(link) == true {
				// Tester om datoen står på en anden måde?
				dato := OtherDatoField.Copy().Find([]byte(link))
				videoDato = string(dato)
			} else if r_validVideoLink.MatchString(e.Text) == true
		{
				dato := r_findDate.Copy().Find([]byte(link))
				// Ved ikke hvorfor men nogle gange skriver han YYYY.MM.DD
				// Hvis den giver true, skal jeg lave om på formateringen af datoen.
				r_wongDateFormate, _ := regexp.Compile("([0-9]{4}.[0-9]{2}.[0-9]{2})")
				if r_wongDateFormate.MatchString(link) == true {
					year := fmt.Sprint(string(dato[:+4]))
					month := fmt.Sprint(string(dato[5:7]))
					date := string(dato[len(dato)-2:])
					videoDato = (date + "." + month + "." + year)
				} else {
					// Her har den fundet datoeen men den står ikke forkert.
					videoDato = string(dato)
				}

			}else {
				// Her har den ikke fundet er dato. og sætter derfpr datoem til N/A
				videoDato = "N/A"
			}
			//fmt.Println(personData)
			//NewVideo := Video{Titel: videoTitel, Url: videoUrl, Dato: videoDato}
			//NewVideo.print()
			unescapedPath, err := url.PathUnescape(videoUrl)
			if err != nil {
				log.Fatal(err)
			}
			newVideo := Video{Id: NewPerson().Id, Titel: videoTitel, Url: unescapedPath, Dato: videoDato}
			person.AddVideo(newVideo)
			fmt.Printf("Video_ID: %d -> Dato: %s Titel: %s\n", newVideo.Id, newVideo.Dato, newVideo.Titel)

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
	fmt.Println("test!")
	mydata, _ := json.MarshalIndent(person, "", "")

	f, err := os.OpenFile("myfile.data", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString(string(mydata)); err != nil {
		panic(err)
	}
	fmt.Println(person)

}
