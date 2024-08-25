# goes-api-go

## Description
A lightweight and fast API for retrieving GOES satellite data stored in an S3 bucket, written in Go. This service is designed to be used in conjunction with goes-web, a lightweight web frontend for viewing GOES satellite data.

## Requirements
This service assumes that you have an S3 bucket that is filled with GOES satellite data, likely from [goesrecv](https://pietern.github.io/goestools/commands/goesrecv.html), a system for capturing and processing the satellite data locally. It assumes that you have the same directory structure as the data folder used by `goesproc` as it's writing data. Basically, this implies using some cron script to simply copy the data from the data folder used by `goesproc` (e.g. goes16) to an S3 bucket.

This project was developed against minio and digital ocean spaces, but should work with any S3-compatible storage.

## Configuration
The following environment variables are required:
- `ACCESS_KEY_ID` - The access key for the S3 bucket.
- `SECRET_ACCESS_KEY` - The secret access key for the S3 bucket.
- `S3_ENDPOINT` - The endpoint for the S3 API.
- `BUCKET_NAME` - The name of the S3 bucket.
- `USE_SSL_FOR_S3` - Whether to use SSL for S3. Set to `false` to use with simple local minio setups.
- `TRUSTED_PROXIES`- A comma-separated list of trusted proxies, in case you want to set them. This is useful if you are running behind a reverse proxy and want to get the real IP address of the client.

## API Usage
The following API endpoints are available:
- `/latest` - Returns the latest image available. Example:
```shell
> curl http://localhost:3000/latest
{ "latest": "2024-08-23" }
```
- `/available-dates` - Returns a list of dates with available images. Example:
```shell
> curl http://localhost:3000/available-dates
{ "availableDates": ["2024-08-23", "2024-08-22", "2024-08-21"] }
```
- `/archive/{date}` - Returns the image urls for a given date.Example:
```shell
> curl http://localhost:3000/archive/2024-05-16
{
  "imageUrls": [
    "http://localhost:9000/goes-16/false-color/fd/2024-05-16/OR_ABI-L2-CMIPF-M6CFC_G16_s20241371130206_e20241371139514_c20241371139566.png",
    "http://localhost:9000/goes-16/false-color/fd/2024-05-16/OR_ABI-L2-CMIPF-M6CFC_G16_s20241371200206_e20241371209514_c20241371209572.png",
    "http://localhost:9000/goes-16/false-color/fd/2024-05-16/OR_ABI-L2-CMIPF-M6CFC_G16_s20241371230206_e20241371239514_c20241371239566.png",
    "http://localhost:9000/goes-16/false-color/fd/2024-05-16/OR_ABI-L2-CMIPF-M6CFC_G16_s20241372330204_e20241372339512_c20241372339564.png"
  ]
}
```