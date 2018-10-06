package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

	// Step 1. Create a list of files
	wordpressUploadsFolder := os.Args[2]
	files := uploadsNames(wordpressUploadsFolder)
	fmt.Println(len(files))

	// Step 2. Instead of creating a list, upload each file into S3
	//  If smth is not uploaded - print the error.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	uploader := s3manager.NewUploader(sess)
	fmt.Println(uploader)

	bucket := "zonovme-assets"
	for _, filename := range files {
		key := "uploads/" + strings.Replace(filename, wordpressUploadsFolder, "", 1)
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
		}
		upParams := &s3manager.UploadInput{
			Bucket: &bucket,
			Key:    &key,
			Body:   file,
		}
		_, err = uploader.Upload(upParams)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Perform an upload.
	fmt.Println("Alright!")

	// Step 3. After file uploaded, find it in the dump and change the path
}

func uploadsNames(folder string) []string {
	fmt.Println("Looking into: " + folder)
	// var files []os.FileInfo
	var names []string

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		if file.IsDir() {
			names = append(names, uploadsNames(folder+"/"+file.Name())...)
		} else {
			names = append(names, folder+"/"+file.Name())
		}
	}
	return names
}
