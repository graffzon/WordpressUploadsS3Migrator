package main

import (
  "fmt"
  "os"
  "io/ioutil"
)

func main() {
  // wordpressDumpPath := os.Args[1]
  // wordpressDumpBytes, err := ioutil.ReadFile(wordpressDumpPath)
	// if err != nil {
	// 	fmt.Println(err)
	// }
  // wordpressDump := string(wordpressDumpBytes)
  // fmt.Println(wordpressDump(1)) // remove later
  //

  wordpressUploadsFolder := os.Args[2]
  names := uploadsNames(wordpressUploadsFolder)
  fmt.Println(len(names))
}

func uploadsNames(folder string) []string {
  var names []string

  files, err := ioutil.ReadDir(folder)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
    if file.IsDir() {
      names = append(names, uploadsNames(folder + "/" + file.Name())...)
    } else {
      names = append(names, file.Name())
    }
	}
  return names
}
