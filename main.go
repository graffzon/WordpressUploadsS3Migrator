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
const bucketName = "zonovme-assets"

func main() {
	// Step 1. Create a list of files
	wordpressUploadsFolder := os.Args[2]
	files := getUploadsNames(wordpressUploadsFolder)
	fmt.Println(len(files))

	// Step 2. Upload each file from the list into S3.
	//  If smth is not uploaded - print the error.
	err := uploadFilesToS3(files, wordpressUploadsFolder)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Step 3. After file uploaded, find it in the dump and change the path
	wordpressDumpPath := os.Args[1]
	wordpressDumpBytes, err := ioutil.ReadFile(wordpressDumpPath)
	if err != nil {
		fmt.Println(err)
	}
	wordpressDump := string(wordpressDumpBytes)

	wordpressUploadsWebPath := os.Args[3]
	replaceOldURLsWithNew(files, wordpressDump, wordpressDumpPath, wordpressUploadsFolder, wordpressUploadsWebPath)
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

func uploadFilesToS3(files []string, wordpressUploadsFolder string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(bucketRegion)},
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	uploader := s3manager.NewUploader(sess)

	bucket := bucketName
	acl := "public-read"

	// Will be taken from the first upload
	objects := []s3manager.BatchUploadObject{}
	// S3BasePath := ""

	for _, filename := range files {
		key := "uploads/" + strings.Replace(filename, wordpressUploadsFolder, "", 1)
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
		}
		upParams := []s3manager.BatchUploadObject{
			{
				Object: &s3manager.UploadInput{
					Bucket: &bucket,
					Key:    &key,
					Body:   file,
					ACL:    &acl,
				},
			},
		}

		objects = append(objects, upParams...)

		// result, err := uploader.Upload(upParams)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// if S3BasePath == "" {
		// 	fmt.Println("key: " + strings.Replace(filename, wordpressUploadsFolder, "", 1))
		// 	fmt.Println("location: " + result.Location)
		// 	S3BasePath = strings.Replace(result.Location, strings.Replace(filename, wordpressUploadsFolder, "", 1), "", 1)
		// }
	}
	// fmt.Println(S3BasePath)

	iter := &s3manager.UploadObjectsIterator{Objects: objects}
	if err := uploader.UploadWithIterator(aws.BackgroundContext(), iter); err != nil {
		return err
	}
	return nil
}

func replaceOldURLsWithNew(files []string, wordpressDump string, wordpressDumpPath string, wordpressUploadsFolder string, wordpressUploadsWebPath string) {
	// replace `uploads` with a variable
	bucketUrl := "https://" + bucketName + ".s3." + bucketRegion + ".amazonaws.com/uploads"
	for _, filename := range files {
		fullFileNameInDump := strings.Replace(filename, wordpressUploadsFolder, wordpressUploadsWebPath, 1)
		fullFileNameInS3 := strings.Replace(filename, wordpressUploadsFolder, bucketUrl, 1)
		fmt.Println("fullFileNameInDump: " + fullFileNameInDump + "        \n fullFileNameInS3:" + fullFileNameInS3)

		// wordpressDump = strings.Replace(wordpressDump, fullFileNameInDump, fullFileNameInS3, -1)
	}
	ioutil.WriteFile(wordpressDumpPath, []byte(wordpressDump), 0644)
	fmt.Println(wordpressDump[1])
}
