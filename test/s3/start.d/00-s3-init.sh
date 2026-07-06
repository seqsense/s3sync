#!/bin/sh

set -eu

awslocal s3 mb s3://example-bucket
awslocal s3 cp /fixture/README.md s3://example-bucket
awslocal s3 cp /fixture/README.md s3://example-bucket/foo/
awslocal s3 cp /fixture/README.md s3://example-bucket/bar/baz/

awslocal s3 mb s3://s3-source
awslocal s3 cp /fixture/README.md s3://s3-source
awslocal s3 cp /fixture/README.md s3://s3-source/foo/
awslocal s3 cp /fixture/README.md s3://s3-source/bar/baz/

awslocal s3 mb s3://s3-destination

awslocal s3 mb s3://s3-destination2

awslocal s3 mb s3://example-bucket-escaped

awslocal s3 mb s3://example-bucket-upload
awslocal s3 cp /fixture/README.md s3://example-bucket-upload/dest_only_file

awslocal s3 mb s3://example-bucket-upload-file

awslocal s3 mb s3://example-bucket-delete
awslocal s3 cp /fixture/README.md s3://example-bucket-delete/dest_only_file

awslocal s3 mb s3://example-bucket-delete-file
awslocal s3 cp /fixture/README.md s3://example-bucket-delete-file
awslocal s3 cp /fixture/README.md s3://example-bucket-delete-file/dest_only_file

awslocal s3 mb s3://example-bucket-dryrun
awslocal s3 cp /fixture/README.md s3://example-bucket-dryrun/dest_only_file

awslocal s3 mb s3://example-bucket-directory
awslocal s3api put-object --bucket example-bucket-directory --key test/

awslocal s3 mb s3://example-bucket-mime
