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

type Any interface{}
type autoInc_person struct {
	sync.Mutex
	id int
}

type autoInc_video struct {
	sync.Mutex
	id int
}

func (a *autoInc_person) ID_P() (id int) {
	a.Lock()
	defer a.Unlock()
	id = a.id
	a.id++
	return
}

func (b *autoInc_video) ID_V() (id int) {
	b.Lock()
	defer b.Unlock()
	id = b.id
	b.id++
	return
}

var ai_p autoInc_person
var ai_v autoInc_video

type Person struct {
	ID_Person   int
	PrimitiveID string
	Navn        string
	Website     string
	Videos      []Video
}
type Video struct {
	ID_Video int
	Titel    string
	Url      string
	Dato     string
}

type Videos []*Video

func (vs Videos) Process(f func(video *Video)) {
	for _, video := range vs {
		f(video)
	}
}

func (vs Videos) FindAll(f func(video *Video) bool) Videos {
	videos := make([]*Video, 0)
	vs.Process(func(v *Video) {
		if f(v) {
			videos = append(videos, v)
		}
	})
	return videos
}

func (vs Videos) Map(f func(video *Video) Any) []Any {
	result := make([]Any, 0)
	ix := 0
	vs.Process(func(v *Video) {
		result[ix] = f(v)
		ix++
	})
	return result
}

func MakeSortedAppender(manufacture []string) (func(video *Video), map[string]Videos) {
	sortedVideos := make(map[string]Videos)

	for _, m := range manufacture {
		sortedVideos[m] = make([]*Video, 0)
	}

	sortedVideos["Default"] = make([]*Video, 0)

	appender := func(v *Video) {
		if _, ok := sortedVideos[v.Dato]; ok {
			sortedVideos[v.Dato] = append(sortedVideos[v.Dato], v)
		} else {
			sortedVideos["Default"] = append(sortedVideos["Default"], v)
		}
	}

	return appender, sortedVideos
}

var PersonData Person
var VideoData Video

func NewPerson() *Person {
	return &Person{
		ID_Person: ai_p.ID_P(),
	}
}
func NewVideo() *Video {
	return &Video{
		ID_Video: ai_v.ID_V(),
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

		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		r_validVideoLink, _ := regexp.Compile("(m+p+4)")
		r_findDate, _ := regexp.Compile("([0-9]{2}.[0-9]{2}.[0-9]{4})|([0-9]{4}.[0-9]{2}.[0-9]{2})")
		OtherDatoField, _ := regexp.Compile("/^(?:(?:31(/|-|.)(?:0?[13578]|1[02]))1|(?:(?:29|30)(/|-|.)(?:0?[1,3-9]|1[0-2])2))(?:(?:1[6-9]|[2-9]d)?d{2})$|^(?:29(/|-|.)0?23(?:(?:(?:1[6-9]|[2-9]d)?(?:0[48]|[2468][048]|[13579][26])|(?:(?:16|[2468][048]|[3579][26])00))))$|^(?:0?[1-9]|1d|2[0-8])(/|-|.)(?:(?:0?[1-9])|(?:1[0-2]))4(?:(?:1[6-9]|[2-9]d)?d{2})$/g")

		if r_validVideoLink.MatchString(link) == true {
			// info.ID = int(ai.ID())
			tempVideoTitel := e.Text
			tempVideoDato := ""
			// Sådan at jeg kan se hvor mange den har fundet i terminalen.
			u, TestError := url.Parse(link)
			if TestError != nil {
				log.Fatal(TestError)
			}
			tempVideoUrl := ("http://micl-easj.Dk" + u.EscapedPath())
			if r_findDate.MatchString(link) == true {
				dato := r_findDate.Copy().Find([]byte(link))
				// Ved ikke hvorfor men nogle gange skriver han YYYY.MM.DD
				// Hvis den giver true, skal jeg lave om på formateringen af datoen.
				r_wongDateFormate, _ := regexp.Compile("([0-9]{4}.[0-9]{2}.[0-9]{2})")
				if r_wongDateFormate.MatchString(link) == true {
					year := fmt.Sprint(string(dato[:+4]))
					month := fmt.Sprint(string(dato[5:7]))
					date := string(dato[len(dato)-2:])
					tempVideoDato = (date + "." + month + "." + year)
				} else {
					// Her har den fundet datoeen men den står ikke forkert.
					tempVideoDato = string(dato)
				}
			} else if OtherDatoField.MatchString(link) == true {
				// Tester om datoen står på en anden måde?
				dato := OtherDatoField.Copy().Find([]byte(link))
				tempVideoDato = string(dato)
			} else {
				// Her har den ikke fundet er dato. og sætter derfpr datoem til N/A
				tempVideoDato = "N/A"
			}
			test123 := ai_v.id
			micl := &Video{test123, VideoData.Titel, tempVideoUrl, tempVideoDato}
			allVideos := Videos([]*Video{micl})
			allNewVideos := allVideos.FindAll(func(video *Video) bool {
				return (video.ID_Video == 1) && (video.ID_Video > 2010)
			})
			fmt.Println("All Videos: ", allVideos)
			fmt.Println("New Videos: ", allNewVideos)

			copyright := []string{"mich", "Person2", "Person3"}
			sortedAppender, sortedVideos := MakeSortedAppender(copyright)
			allVideos.Process(sortedAppender)
			fmt.Println("Map sortedVideos: ", sortedVideos)
			MichCount := len(sortedVideos["mich"])
			fmt.Println("Mic har lavet ", MichCount, "Videoer")
			VideoData.Titel = tempVideoTitel
			VideoData.Url = tempVideoUrl
			VideoData.Dato = tempVideoDato
			//fmt.Printf("Video_ID: %d -> Dato: %s Titel: %s\n", test123, VideoData.Dato, VideoData.Titel)

			mydata, _ := json.MarshalIndent(NewVideo(), PersonData.PrimitiveID, PersonData.Navn)

			f, err := os.OpenFile("myfile.data", os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			if _, err = f.WriteString(string(mydata)); err != nil {
				panic(err)
			}

		} else {
			c.Visit(e.Request.AbsoluteURL("http://micl-easj.dk/" + link))
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	NewPerson()
	PersonData.Navn = "test"
	PersonData.Website = "http://micl-easj.dk/"

	c.Visit("http://micl-easj.dk/")

}
