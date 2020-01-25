# aws-ip-scanner

Gets a list of the public IPs on your AWS account and scans them for open ports

## How to use

```
go install github.com/zebroc/aws-ip-scanner

export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=
export AWS_REGION=eu-west-1

$GOPATH/bin/aws-ip-scanner
```

## Example output

```
1.2.3.4	22
4.3.2.1	80	HTTP (200): <html><header><title>This is title</title></header><body>Hello world</body></html>
5.4.3.2	443	HTTPS (200): {  "name" : "",  "cluster_name" : "ElasticSearch",  "cluster_uuid" : "abcdefg",  "version" : {    "number" : "7.1.0",    "build_flavor" : "oss",    "build_type" : "tar",    "build_hash" : "abcdefg",    "build_date" : "2000-01-01T01:00:00.000000Z",    "build_snapshot" : false,    "lucene_version" : "8.0.0",    "minimum_wire_compatibility_version" : "6.8.0",    "minimum_index_compatibility_version" : "6.0.0-beta1"  },  "tagline" : "You Know, for Search"}
6.5.4.3	80, 443	HTTP (200): Go away	HTTPS (200): Go away
7.6.5.4	6379
```

We found an open ElasticSearch instance here which should not be exposed…
7.6.5.4 seems to be running an open Redis server, also a bad idea

## What to keep in mind

Run this from an IP that is not whitelisted to see if your setup is open.

## Contribute

Fork, code and create a PR --> I'm happy for additions!
