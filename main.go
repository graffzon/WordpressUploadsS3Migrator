package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

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
		Region: aws.String("eu-west-1")},
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	uploader := s3manager.NewUploader(sess)
	fmt.Println(uploader)

	bucket := "zonovme-assets"
	key := "foo/bar" + files[0]
	file, err := os.Open(files[0])
	if err != nil {
		log.Fatal(err)
	}
	upParams := &s3manager.UploadInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file,
	}

	// Perform an upload.
	result, err := uploader.Upload(upParams)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

	// Step 3. After file uploaded, find it in the dump and change the path
}

func uploadsNames(folder string) []string {
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
			names = append(names, folder+file.Name())
		}
	}
	return names
}
