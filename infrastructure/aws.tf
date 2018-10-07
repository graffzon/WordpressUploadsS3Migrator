provider "aws" {
  region = "eu-central-1"
}

resource "aws_s3_bucket" "blog-uploads" {
  bucket = "zonovme-assets"
  acl    = "public-read"

  versioning {
    enabled = true
  }
}
