# pubsub_cli
![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)
![Main Workflow](https://github.com/k-yomo/pubsub_cli/workflows/Main%20Workflow/badge.svg)
[![codecov](https://codecov.io/gh/k-yomo/pubsub_cli/branch/master/graph/badge.svg)](https://codecov.io/gh/k-yomo/pubsub_cli)
[![Go Report Card](https://goreportcard.com/badge/go-kit/kit)](https://goreportcard.com/report/k-yomo/pubsub_cli)

pubsub_cli is a super handy Pub/Sub CLI which lets you publish / subscribe Pub/Sub message right away!

## Installation
### CLI
#### Homebrew
```
$ brew tap k-yomo/pubsub_cli
$ brew install pubsub_cli 
```

#### Go 
```
$ go get github.com/k-yomo/pubsub_cli
```

### Pub/Sub Emulator
- Make sure Pub/Sub Emulator is running before executing commands.
```
$ gcloud beta emulators pubsub start --host-port=0.0.0.0:8085
```
 
## Usage

```
Usage:
  pubsub_cli [command]

Available Commands:
  help          Help about any command
  publish       publish Pub/Sub message
  register_push register Pub/Sub push endpoint
  subscribe     subscribe Pub/Sub topics

Flags:
  -c, --cred-file string   gcp credential file path (You can also set 'GCP_CREDENTIAL_FILE_PATH' to env variable)
      --help               help for pubsub_cli
  -h, --host string        emulator host (You can also set 'PUBSUB_EMULATOR_HOST' to env variable)
  -p, --project string     gcp project id (You can also set 'GCP_PROJECT_ID' to env variable)
```
※ When both of --host and --cred-file are set, emulator host will be prioritised for safety purpose.

## Examples
### Publish
```
$ pubsub_cli publish test_topic '{"key":"value"}' --host=localhost:8085 --project=test_project
```

### Subscribe
```
$ pubsub_cli subscribe test_topic another_topic --cred-file=credentials.json -p=test_project
```

### Register Push Endpoint
```
$ pubsub_cli register_push test_topic http://localhost:1323/subscribe -h=localhost:8085 -p=test_project
```

## Note
- Created topic won't be deleted automatically. 
- Unused subscription will be deleted in 24 hours.
