# Results worker for stränge.de riddle generation

This repository contains the **results worker** microservice for generating riddles for [strangui](https://github.com/polarfoxDev/strangui).

## Overview

The results worker is a Go-based service that retrieves finished job results from a Redis query and stores them as JSON files in a specified directory.

## Features

- **JSON output** from generated riddles
- **Queue management** via Redis
- **Logging** with logrus for easy debugging and monitoring

## How It Works

1. The worker monitors the Redis queue `generate-riddle-result`.
2. Take the result from the queue.
3. Save it in a JSON file.

## Project Structure

[`main.go`](./main.go) contains most of the logic: a worker loop, Redis integration, configuration, logging and storing files.

## Setup

### Prerequisites

- Go 1.20+
- Redis server

### Installation

Clone the repository and install dependencies:

```bash
git clone https://github.com/polarfoxDev/straenge-results-worker.git
cd straenge-results-worker
go mod tidy
```

### Configuration

Create a `.env` file in the root directory with the following variables:

```env
REDIS_URL=localhost:6379
LOG_LEVEL=info
BASE_FILE_PATH=/path/to/save/riddles
```

### Running the Worker

Start the worker:

```bash
go run main.go
```

## Contributing

Any contributions you make are greatly appreciated.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue.
Don't forget to give the project a star! Thanks!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the **MIT** License. See [`LICENSE`](./LICENSE) for more information.
