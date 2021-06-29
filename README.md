# sns-cwalarm

Simple program to send a test cloudwatch alarm to a sns topic.

## Usage

```bash
Usage of ./sns-cwalarm:
  -description string
        the alarm description to send (default "you can just chill!")
  -name string
        the alarm name to send (default "testing-alarm")
  -topic string
        arn to the sns topic to send the alarm to
```

## Build

Build with go using `go build`.
