# checkupd

[![Go Report Card](https://goreportcard.com/badge/github.com/AdhityaRamadhanus/checkupd)](https://goreportcard.com/report/github.com/AdhityaRamadhanus/checkupd)  

self-hosted endpoint monitoring daemon and status pages based on https://github.com/sourcegraph/checkup

<p>
  <a href="#installation">Installation |</a>
  <a href="#setting-up-checkupd-with-fs">Setting FS |</a>
    <a href="#setting-up-checkupd-with-s3">Setting S3 |</a>
  <a href="#setting-endpoints">Setting Endpoints |</a>
  <a href="#notifier-slack">Notifier |</a>
  <a href="#licenses">License</a>
  <br><br>
  <blockquote>
	Checkupd is self-hosted health checks and status pages, written in Go using checkup (instead of using them as dependency i decide to copy the file to this project) and grpc as backend.

  There is much work to do for this project to be complete. Use it carefully.

  Checkupd currently supports:

  - Checking HTTP endpoints
  - Checking TCP endpoints (TLS supported)
  - Checking of DNS services & record existence  
  - Storing results on S3 and local filesystem
  - Add and delete endpoints on the fly, you don't need to change the config file everytime you decide to add/delete an endpoint
  - Easy to setup and deploy status page (100% static)
  - Get notified via slack and email (soon)
  </blockquote>
</p>

Installation
----------- 
* git clone
* make
```bash
NAME:
   checklist - Checkup server cli 

USAGE:
   checklist [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR:
   Adhitya Ramadhanus <adhitya.ramadhanus@gmail.com>

COMMANDS:
     add-http  Add endpoints to checkup
     add-tcp   Add tcp endpoints to checkup
     check     list and check endpoints
     list      list endpoint
     delete    delete endpoint
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

```bash
NAME:
   checkupd - Checkupd daemon 

USAGE:
   checkupd [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR:
   Adhitya Ramadhanus <adhitya.ramadhanus@gmail.com>

COMMANDS:
     setup-page    Setup statuspage
     setup-daemon  Setup daemon
     daemon        run daemon
     help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Setting up checkupd with FS
All you need to do is setup configuration for daemon and status page

1. **Daemon Setup**
  ```bash
  $ checkupd setup-daemon fs --log=<logs dir, by default its ./checkup_config/logs>
  ```

  It will create a directory called checkup_config and checkup_config/logs and checkup.json
  Example of fs storage checkup.json generated by cli 
  ```json
  {
    "checkers": [],
    "storage": {
      "provider": "fs",
      "dir": "./checkup_config/logs/"
    }
  }
  ```

2. **Status Page Setup**
  ```bash
  $ checkupd setup-page fs --log=<logs dir, by default its logs> --url=<url to serve the page, with port ex: localhost:80, mycheckup.com:80>
  ```

  It will create a directory called caddy_config and put config.js inside statuspage/js and index.html in statuspage
  Example of fs config.js generated by cli 
  ```js
  checkup.config = {
    "timeframe": 1 * time.Day,
    "refresh_interval": 60,
    "storage": {
      "url": "logs"
    },
    "status_text": {
      "healthy": "Situation Normal",
      "degraded": "Degraded Service",
      "down": "Service Disruption"
    }
  };
  ```

### Setting up checkupd with S3
First you need S3 bucket for this and setting privileges for this bucket, unfortunately checkupd doesn't support automatic provisioning like checkup so you have to do manual provision https://github.com/sourcegraph/checkup/wiki/Provisioning-S3-Manually

Then, just like with the FS, you need to setup configuration for daemon and status page

1. **Daemon Setup**
  ```bash
  $ checkupd setup-daemon s3 --i=<s3 AccessKeyID> --k=<s3 SecretKey> --r=<s3 Region> --b=<s3 Bucket Name>
  ```

  It will create a directory called caddy_config and put config.js inside statuspage/js and index.html in statuspage
  Example of fs config.js generated by cli 
  ```json
  {
    "checkers":[],
    "storage": {
      "provider": "s3",
      "access_key_id": "<yours>",
      "secret_access_key": "<yours>",
      "bucket": "<yours>",
      "region": "us-east-1"
    }
  }
  ```
  
2. **Status Page Setup**
  ```bash
  $ checkupd setup-daemon s3 --i=<s3 AccessKeyID> --k=<s3 SecretKey> --r=<s3 Region> --b=<s3 Bucket Name> --url=<url to serve the page, with port ex: localhost:80, mycheckup.com:80>
  ```

  It will create a directory called checkup_config and checkup_config/logs and checkup.json
  Example of fs storage checkup.json generated by cli 
  ```js
  checkup.config = {
    "timeframe": 1 * time.Day,
    "refresh_interval": 60,
    "storage": {
      "AccessKeyID": "{{.AccessKeyID}}",
      "SecretAccessKey": "{{.SecretAccessKey}}",
      "Region": "{{.Region}}",
      "BucketName": "{{.Bucket}}"
    },
    "status_text": {
      "healthy": "Situation Normal",
      "degraded": "Degraded Service",
      "down": "Service Disruption"
    }
  };
  ```

Running Daemon
---------------
```bash
$ checkupd daemon
```

Serve Status Page
---------------
```bash
docker-compose up
```
(Modified version, still in progress)
![checkup](https://cloud.githubusercontent.com/assets/5761975/25096466/888ca154-23ca-11e7-910d-59be4c610989.png)

Setting Endpoints
----------------
You can manually add endpoints to the generated checkup.json (see Setting up checkupd with FS or S3)

Example of such configuration

```json
{
	"checkers": [{
		"type": "http",
		"endpoint_name": "Example HTTP",
		"endpoint_url": "http://www.example.com",
		"attempts": 5
	},
	{
		"type": "tcp",
		"endpoint_name": "Example TCP",
		"endpoint_url": "example.com:80",
		"attempts": 5
	},
	{
		"type": "tcp",
		"endpoint_name": "Example TCP with TLS enabled and a valid certificate chain",
		"endpoint_url": "example.com:443",
		"attempts": 5,
		"tls": true
	},
	{
		"type": "tcp",
		"endpoint_name": "Example TCP with TLS enabled and a self-signed certificate chain",
		"endpoint_url": "example.com:8443",
		"attempts": 5,
		"timeout": "2s",
		"tls": true,
		"tls_ca_file": "testdata/ca.pem"
	},
	{
		"type": "tcp",
		"endpoint_name": "Example TCP with TLS enabled and verification disabled",
		"endpoint_url": "example.com:8443",
		"attempts": 5,
		"timeout": "2s",
		"tls": true,
		"tls_skip_verify": true
	},
	{
		"type": "dns",
		"endpoint_name": "Example DNS test of endpoint_url looking up host.example.com",
		"endpoint_url": "ns.example.com:53",
		"hostname_fqdn": "host.example.com",
		"timeout": "2s"
	}],
	"storage": {
		"provider": "s3",
		"access_key_id": "<yours>",
		"secret_access_key": "<yours>",
		"bucket": "<yours>",
		"region": "us-east-1"
	}
}
```

Or you can add them on the fly when the daemon run

1. Adding Tcp endpoint
```bash
NAME:
   checklist add-tcp - Add tcp endpoints to checkup

USAGE:
   checklist add-tcp [command options] [arguments...]

OPTIONS:
   --name value     Name of endpoint
   --address value  Address to check
   --tcp-tls        Is it tls endpoint?
   --host value     grpc server address (default: "/tmp/checkupd.sock")
```

2. Adding Http endpoint
```bash
NAME:
   checklist add-http - Add endpoints to checkup

USAGE:
   checklist add-http [command options] [arguments...]

OPTIONS:
   --name value  Name of endpoint
   --url value   URL to check
   --host value  grpc server address (default: "/tmp/checkupd.sock")
```

Just like adding endpoint, you can either modify the checkup.json or delete them through cli

1. Deleting Endpoint
```bash
NAME:
   checklist delete - delete endpoint

USAGE:
   checklist delete [command options] [arguments...]

OPTIONS:
   --tls         Name of endpoint
   --host value  grpc server address (default: "/tmp/checkupd.sock")
   --name value  endpoint name
```

List Endpoint
------------------------
```bash
NAME:
   checklist list - list endpoint

USAGE:
   checklist list [command options] [arguments...]

OPTIONS:
   --host value  grpc server address (default: "/tmp/checkupd.sock")
```

Check Endpoint
------------------------
```bash
NAME:
   checklist check - list and check endpoints

USAGE:
   checklist check [command options] [arguments...]

OPTIONS:
   --host value  grpc server address (default: "/tmp/checkupd.sock")
```

Notifier (Slack)
----------------
* To use this notifier you need bot integration in your team and channel id where this bot will notify you, refer to this link https://api.slack.com/bot-users
* Add notifier section to the generated checkup.json
* Example of checkup.json 
```json
{
	"checkers": [{
		"type": "tcp",
		"endpoint_name": "redis",
		"endpoint_url": "localhost:6379",
		"attempts": 5
	}],
	"storage": {
		"provider": "fs",
		"dir": "./checkup_config/logs"
	},
  "notifier": {
      "name": "slack",
      "token": "your token",
      "channel": "your channel id"
  }
}
```


Todo
-----------
* Email Notifier

License
----

MIT © [Adhitya Ramadhanus]

