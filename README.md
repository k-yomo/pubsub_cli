# pubsub_cli
Very Handy pubsu_cli for Pub/Sub Emulatro

## Installation
### Homebrew
```
$ brew tap k-yomo/pubsub_cli
$ brew install pubsub_cli 
```

### Go 
```
$ go get github.com/k-yomo/pubsub_cli
```

### Pub/Sub Emulator
- Make sure Pub/Sub Emulator is running before executing commands.
```
gcloud beta emulators pubsub start --host-port=0.0.0.0:8432
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
      --host string      emulator host (You can also set 'PUBSUB_EMULATOR_HOST' to env variable)
      --project string   gcp project id (You can also set 'GCP_PROJECT_ID' to env variable)
```

### Publish
```
  pubsub_cli publish test_topic '{"jsonKey":"value"}' --host localhost:8432 --project test_project
```

### Subscribe
```
  pubsub_cli subscribe test_topic --host localhost:8432 --project test_project
```

### Register Push Endpoint
```
  pubsub_cli register_push test_topic http://localhost:1323/subscribe --host localhost:8432 --project test_project
```
