# go-bcov

<p align="center">
  <img src="https://github.com/ALX99/go-bcov/blob/main/logo.png" />
</p>

`go-bcov` is a tool to calculate branch coverage from Go coverage reports.

It reads the coverage file from standard input and writes to standard output.

## Installation

```bash
go install github.com/alx99/go-bcov@v1
```

## Usage

```bash
$ go-bcov -h

Usage of go-bcov:
  -format string
    	output format (default "reserved")
```

### Sonarqube coverage report

```bash
go test -coverprofile=coverage.out -covermode count ./...
go-bcov -format sonar-cover-report < coverage.out > coverage.xml

# and upload...
sonar-scanner-cli \
    -Dsonar.sources=. -Dsonar.exclusions=**/*_test.go,**/*_mock.go \
    -Dsonar.tests=. -Dsonar.test.inclusions=**/*_test.go \
    -Dsonar.coverageReportPaths=coverage.xml -Dsonar.go.coverage.reportPaths=coverage.txt
```

## Supported output formats

- [Sonarqube generic coverage report](https://docs.sonarsource.com/sonarqube/latest/analyzing-source-code/test-coverage/generic-test-data/)
  - Both line and branch coverage are supported.

