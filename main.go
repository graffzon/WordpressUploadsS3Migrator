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

const bucketRegion = "eu-central-1"

func main() {
	// Reading command line params
	wordpressDumpPath := os.Args[1]
	wordpressUploadsFolder := os.Args[2]
	wordpressUploadsWebPath := os.Args[3]
	bucketName := os.Args[4]

	// Step 1. Create a list of files
	files := getUploadsNames(wordpressUploadsFolder)
	fmt.Println("Number of files to be sent: " + string(len(files)))

	// Step 2. Upload each file from the list into S3.
	//  If smth is not uploaded - print the error.
	err := uploadFilesToS3(files, wordpressUploadsFolder, bucketName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Step 3. After files uploaded, find them in the dump and change the path
	wordpressDumpBytes, err := ioutil.ReadFile(wordpressDumpPath)
	if err != nil {
		fmt.Println(err)
	}
	wordpressDump := string(wordpressDumpBytes)

	replaceOldURLsWithNew(files, wordpressDump, wordpressDumpPath, wordpressUploadsFolder, wordpressUploadsWebPath, bucketName)
}

func getUploadsNames(folder string) []string {
	fmt.Println("Looking into: " + folder)
	var names []string

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		if file.IsDir() {
			names = append(names, getUploadsNames(folder+"/"+file.Name())...)
		} else {
			names = append(names, folder+"/"+file.Name())
		}
	}
	return names
}

func uploadFilesToS3(files []string, wordpressUploadsFolder string, bucketName string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(bucketRegion)},
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	acl := "public-read"
	uploader := s3manager.NewUploader(sess)
	objects := []s3manager.BatchUploadObject{}

	for _, filename := range files {
		key := "uploads/" + strings.Replace(filename, wordpressUploadsFolder, "", 1)
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
		}
		upParams := []s3manager.BatchUploadObject{
			{
				Object: &s3manager.UploadInput{
					Bucket: &bucketName,
					Key:    &key,
					Body:   file,
					ACL:    &acl,
				},
			},
		}

		objects = append(objects, upParams...)
	}

	iter := &s3manager.UploadObjectsIterator{Objects: objects}
	if err := uploader.UploadWithIterator(aws.BackgroundContext(), iter); err != nil {
		return err
	}
	return nil
}

func replaceOldURLsWithNew(files []string, wordpressDump string, wordpressDumpPath string, wordpressUploadsFolder string, wordpressUploadsWebPath string, bucketName string) {
	// TODO: replace `uploads` with a variable
	bucketUrl := "https://" + bucketName + ".s3." + bucketRegion + ".amazonaws.com/uploads"
	for _, filename := range files {
		fullFileNameInDump := strings.Replace(filename, wordpressUploadsFolder, wordpressUploadsWebPath, 1)
		fullFileNameInS3 := strings.Replace(filename, wordpressUploadsFolder, bucketUrl, 1)
		// Uncomment if you want to see all the files' names
		// fmt.Println("fullFileNameInDump: " + fullFileNameInDump + "\n fullFileNameInS3:" + fullFileNameInS3)

		wordpressDump = strings.Replace(wordpressDump, fullFileNameInDump, fullFileNameInS3, -1)
	}
	ioutil.WriteFile(wordpressDumpPath, []byte(wordpressDump), 0644)
	fmt.Println(wordpressDump[1])
}
