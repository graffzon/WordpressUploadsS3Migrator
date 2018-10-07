# Migrator of Wordpress Uploads to S3
Once I had a necessity to migrate my Wordpress blog from a plain hosting to the AWS. And one of the first steps I had to do is to migrate my uploads. I wanted to be safe and first upload them to the S3, then change the links in the hosted blog, and only then, after making sure it works, migrate the blog itself.
Back then I haven't find any tool to help me with this, so I created this small script. Feel free to contribute.

## Steps to migrate your uploads
1. You have to get your SQL dump. Usually you can get it from the management panel of your hosting (in my case it was ISP Manager). Also find out in advance, how can upload the SQL dump back again when you modify it.
2. As the next step, you need to download an archive with your files from hosting. There are a handful of those tools. I used BackupGuard for this, you can use whichever you want.
3. Setup your S3 bucket in the AWS. You can do it either manually, or using the Terraform code, which you can find in the `infrastructure/` directory.
4. Clone this repo and run the script with the following params:

```bin/migrator /Users/username/Downloads/u8526307_default.sql /Users/username/Downloads/sg_backup_opt_20181006084001/wp-content/uploads http://zonov.me/wp-content/uploads bucket-for-assets```

Where the first parameter is your SQL dump file, the second one is your downloaded uploads directory, the third param is your current website `uploads` folder, the fourth one is the name of your bucket in the S3.

Here I'm assuming that your AWS S3 region is `eu-central-1`, if it isn't - feel free to change it in the code and then just run it with `go run`. (Before that you will need to install `go` and `dep`, f.e. using Homebrew).

## Contribution
I would really appreciate your contribution. Here are some things you could do:
1. Improve logging
2. Parallelise the files upload
3. Make the script even more generic, probably extract the region into a param, if so - make it optional, etc.
4. Add tests
5. Improve the documentation

If I haven't mentioned something, but you still want to contribute - feel free to submit an issue or a PR.
