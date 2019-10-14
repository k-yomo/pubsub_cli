# pubsub_cli

## Installation
### Homebrew
```
$ brew tap k-yomo/pubsub_cli
$ brew install pubsub_cli 
```

### Homebrew
```
$ go get github.com/k-yomo/pubsub_cli
```

## Usage
```
Usage:
  pubsub_cli [command]

Available Commands:
  help        Help about any command
  publish     publish Pub/Sub message
  subscribe   subscribe Pub/Sub topic

Flags:
  -h, --help             help for pubsub_cli
      --host string      emulator host (You can also set 'PUBSUB_EMULATOR_HOST' to env variable) (default "localhost:8432")
      --project string   gcp project id (You can also set 'GCP_PROJECT_ID' to env variable) (default "dev")
```

### Publish
```
  pubsub_cli publish test_topic '{"jsonKey":"value"}' --host localhost:8432 --project test_project
```

### Subscribe
```
  pubsub_cli subscribe test_topic --host localhost:8432 --project test_project
```
