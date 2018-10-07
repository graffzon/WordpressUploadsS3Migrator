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
	uploadedBasePath := ""
	acl := "public-read"

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
			ACL:    &acl,
		}
		result, err := uploader.Upload(upParams)
		if err != nil {
			fmt.Println(err)
		}
		if uploadedBasePath == "" {
			fmt.Println("key: " + strings.Replace(filename, wordpressUploadsFolder, "", 1))
			fmt.Println("location: " + result.Location)
			uploadedBasePath = strings.Replace(result.Location, strings.Replace(filename, wordpressUploadsFolder, "", 1), "", 1)
		}
	}
	fmt.Println(uploadedBasePath)

	// Perform an upload.
	fmt.Println("Alright!")

	// Step 3. After file uploaded, find it in the dump and change the path
	wordpressDumpPath := os.Args[1]
	wordpressDumpBytes, err := ioutil.ReadFile(wordpressDumpPath)
	if err != nil {
		fmt.Println(err)
	}
	wordpressDump := string(wordpressDumpBytes)

	wordpressUploadsWebPath := os.Args[3]
	for _, filename := range files {
		fullFileNameInDump := strings.Replace(filename, wordpressUploadsFolder, wordpressUploadsWebPath, 1)
		fullFileNameInS3 := strings.Replace(filename, wordpressUploadsFolder, uploadedBasePath, 1)
		fmt.Println("fullFileNameInDump: " + fullFileNameInDump + "        \n fullFileNameInS3:" + fullFileNameInS3)

		wordpressDump = strings.Replace(wordpressDump, fullFileNameInDump, fullFileNameInS3, -1)
	}
	ioutil.WriteFile(wordpressDumpPath, []byte(wordpressDump), 0644)
	fmt.Println(wordpressDump[1])
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
