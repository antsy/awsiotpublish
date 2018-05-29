# AWS IoT Publish

Serverless Golang lambda function which posts a message to a AWS IoT topic.

## Prerequisites

 * AWS account
   * IoT topic configured
 * [Go](https://golang.org/)
 * [Serverless](https://github.com/serverless/serverless)

## Deployment

To build run `make`.

To deploy to the cloud, configure your profile and region and stuff in `serverless.yml`, then run `sls deploy`.

## Try it yourself!

1. Open AWS console and search for 'iot core'.
2. Select 'test' and create test topic for yourself.
3. Send `POST {url}/{stage}/iot/{topic}/{message}` (Use the address you got from the deployment).
