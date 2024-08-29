# wikitracker

This project is a real-time Wikipedia edit tracking system built using Golang, Apache Kafka, and KSQL. It monitors the Wikipedia edit stream, processes the data, and provides statistics on the most frequently edited pages within a given time window.

## Features

- Real-time streams tracking of Wikipedia edits
- Processing and aggregation of edit data using KSQL
- Show the Top-K of most edited pages
- RESTful API for retrieving statistics (unimplemented)
- Dockerized environment for easy setup and deployment

## Prerequisites

- Docker version 27.1.1
- Docker Compose 2.29.1
- Go 1.21

## Setup

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/wikitracker.git
   cd wikitracker
   ```

2. Start the Kafka ecosystem using Docker Compose:
   ```
   docker-compose up -d
   ```

3. Wait for all services to be up and running. You can check the status using:
   ```
   docker-compose ps
   ```

4. Create the necessary Kafka topics and KSQL streams and tables:
   ```
   docker-compose exec ksqldb-cli ksql http://ksqldb-server:8088
   ```
   Then paste the following KSQL commands:
   ```sql
   CREATE STREAM IF NOT EXISTS WIKI_EDITS_STREAM (
       TYPE STRING,
       `NAMESPACE` INTEGER,
       TITLE STRING,
       TITLE_URL STRING,
       USER STRING,
       BOT BOOLEAN
   ) WITH (
       KAFKA_TOPIC='wiki_events',
       VALUE_FORMAT='JSON'
   );

   CREATE TABLE IF NOT EXISTS WIKI_EDITS_COUNT_HW WITH (
       KAFKA_TOPIC = 'WIKI_EDITS_COUNT_HW',
       VALUE_FORMAT = 'JSON',
       KEY_FORMAT = 'JSON',
       PARTITIONS = 3,
       REPLICAS = 1
   ) AS
   SELECT
       TITLE,
       TITLE_URL AS URL,
       WINDOWSTART AS WINDOW_START,
       WINDOWEND AS WINDOW_END,
       COUNT(*) AS EDITS_COUNT
   FROM WIKI_EDITS_STREAM 
   WINDOW HOPPING (SIZE 30 MINUTES, ADVANCE BY 10 MINUTE)
   WHERE TYPE = 'edit' AND BOT = false
   GROUP BY TITLE, TITLE_URL
   EMIT CHANGES;
   ```

5. Run ./cmd/producer.go and ./cmd/webserver.go as a deamon:
   ```
    // build
    go build -o producer ./cmd/producer.go
    go build -o webserver ./cmd/webserver.go
   
   // run
    nohup ./webserver > webserver.log 2>&1 &
    nohup ./webserver > webserver.log 2>&1 &
   ```

## Usage

The application consists of two main components:

1. A streaming service(producer) that consumes Wikipedia edit events and produces them to a Kafka topic.
2. A consumer service(webserver) that processes the aggregated data and exposes it via a RESTful API.

### API Endpoints

- `/health`: Check the health status of the application
- `/wikitrack`: Get the current top-K most edited Wikipedia pages

## Configuration

The application uses environment variables for configuration. You can set the following variables:

- `KAFKA_SERVER`: Comma-separated list of Kafka broker addresses (default: "localhost:9092")
- `CONSUMER_GROUP`: Kafka consumer group ID (default: "consumer-group-01")
- `TOPIC`: Kafka topic to consume from (default: "WIKI_EDITS_COUNT_HW")

## Architecture

The system architecture consists of the following components:

1. Wikipedia SSE (Server-Sent Events) stream
2. Golang application for consuming and producing events
3. Apache Kafka for message queuing
4. KSQL for stream processing and aggregation

## License
This project is licensed under the MIT License.