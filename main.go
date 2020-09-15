package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	fileUrls := readInFile()
	path := "/home/kneeru/Pictures/9gag/auto"
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			panic(errDir)
		}
	}
	var wg sync.WaitGroup
	for i := 0; i < len(fileUrls); i++ {
		wg.Add(1)
		err := DownloadFile(fmt.Sprintf("/home/kneeru/Pictures/9gag/auto/photo_%d.jpg", i+1), fileUrls[i], &wg)
		if err != nil {
			fmt.Println(err)
			fmt.Println(fileUrls[i])
			fmt.Printf("Error downloading %d ", i)
			return
		}
	}
	fmt.Println("Waiting for all the tasks to complete")
	wg.Wait()
	fmt.Println("Done...")
	return
}

// readInFile : this will read in the file for the contents line by line
func readInFile() []string {
	result := []string{}
	f, err := os.Open("./photos.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// TODO: link copy pasted from the browser has to be dressed up
		s := scanner.Text()
		// from the last index of . we trim off the extension ..
		// assumption here is to replace the webp extension with jpg
		// we are assuming we would be downloading only photos
		i := strings.LastIndex(s, ".")            //https://www.geeksforgeeks.org/how-to-find-the-last-index-value-of-specified-string-in-golang/
		url := fmt.Sprintf("%s.jpg", s[:i][:i-2]) // this is how you can reach the image cache of 9gag
		// removing the last 2 characters before the . can make you reach the image cache downloadable resource
		result = append(result, url)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return result
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string, wg *sync.WaitGroup) error {
	defer wg.Done()
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write the body to file
	_, err = io.Copy(f, resp.Body)
	return err
}
