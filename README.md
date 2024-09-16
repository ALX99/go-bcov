# go-sonarcover

<p align="center">
  <img src="https://github.com/ALX99/go-sonarcover/blob/main/logo.png" />
</p>

## Description

`go-sonarcover` is a tool to generate [generic coverage reports](https://docs.sonarsource.com/sonarqube/latest/analyzing-source-code/test-coverage/generic-test-data/) for SonarQube from Go coverage files.

It reads the coverage file from standard input and writes the generic coverage report to standard output
which includes both line and branch coverage.

## Installation

```bash
go install github.com/alx99/go-sonarcover@latest
```

## Usage

```bash
go test -coverprofile=coverage.out -covermode count ./...
go-sonarcover < coverage.out > coverage.xml

# And upload ...
sonar-scanner-cli \
    -Dsonar.sources=. -Dsonar.exclusions=**/*_test.go,**/*_mock.go \
    -Dsonar.tests=. -Dsonar.test.inclusions=**/*_test.go \
    -Dsonar.coverageReportPaths=coverage.xml -Dsonar.go.coverage.reportPaths=coverage.txt
```
