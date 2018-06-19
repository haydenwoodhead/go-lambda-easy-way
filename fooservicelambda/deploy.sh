#!/usr/bin/env bash
sam package --template-file template.yml --output-template-file output.yml --s3-bucket aklgolambdatalk
sam deploy --template-file output.yml --stack-name FooService --capabilities CAPABILITY_IAM