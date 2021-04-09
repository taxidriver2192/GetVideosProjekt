package main

func main() {

	// DOWNLOADE FILE IN MAIN FUNKTION!
	// fileUrl := "http://micl-easj.dk/Machine%20Learning/Noter%20and%20Books/Machine%20Learning%20Landscape%20Ch.%201/ML%20Landscape%20No.%201%20Overview%20%202021.01.28.mp4"
	// err := DownloadFile("Test.mp4", fileUrl)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Downloaded: " + fileUrl)
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
// func DownloadFile(filepath string, url string) error {

// 	// Get the data
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	// Create the file
// 	out, err := os.Create(filepath)
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()

// 	// Write the body to file
// 	_, err = io.Copy(out, resp.Body)
// 	return err
// }
