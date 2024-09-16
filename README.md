# go-bcov

<p align="center">
  <img src="https://github.com/ALX99/go-bcov/blob/main/logo.png" />
</p>

## Description

`go-bcov` is a tool to calculate branch coverage from Go coverage reports.

It reads the coverage file from standard input and writes to standard output.

### Supported output formats

- [Sonarqube coverage report](https://docs.sonarsource.com/sonarqube/latest/analyzing-source-code/test-coverage/generic-test-data/) for SonarQube from Go coverage files.
  - Both line and branch coverage are supported.

## Installation

```bash
go install github.com/alx99/go-bcov@latest
```

## Usage

### Sonarqube coverage report

```bash
go test -coverprofile=coverage.out -covermode count ./...
go-bcov -format sonar-cover-report < coverage.out > coverage.xml

# And upload ...
sonar-scanner-cli \
    -Dsonar.sources=. -Dsonar.exclusions=**/*_test.go,**/*_mock.go \
    -Dsonar.tests=. -Dsonar.test.inclusions=**/*_test.go \
    -Dsonar.coverageReportPaths=coverage.xml -Dsonar.go.coverage.reportPaths=coverage.txt
```
