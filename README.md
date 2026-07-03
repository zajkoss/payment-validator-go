# payment-validator-go

A CLI tool for concurrent validation of ISO 20022 pain.001 payment files, written in Go.

Built as a learning project to explore Go's concurrency model (goroutines, channels, worker pool) — coming from a background in COBOL/mainframe payment processing and C#/.NET backend development.

## What it does

- Scans a directory for `.xml` payment files
- Validates each file in parallel using a worker pool
- Reports validation errors per file with a summary

## Validation rules

- IBAN must not be blank
- IBAN length validated per country (PL=28, DE=22, GB=22, NL=18)
- IBAN country code must be on whitelist (PL, DE, GB, NL)
- Amount must be a positive number
- Currency must be on whitelist (PLN, EUR, USD)

## Project structure

```
payment-validator-go/
  cmd/validator/
    main.go              # entry point, CLI args, report output
  internal/validator/
    worker.go            # worker pool, XML parsing, validation logic
  testdata/              # sample pain.001 XML files (valid and invalid)
```

## Run

```bash
go run cmd/validator/main.go ./testdata
```
## Run with Docker

```bash
docker compose up
```

Or build manually:
```bash
docker build -t payment-validator .
docker run payment-validator ./testdata
```

## Example output

```
=== Validator start ===
        Directory: ./testdata
WP: Starting worker 1
WP: Starting worker 2
WP: Starting worker 3
WP: Loading data to process...
WP: No more files to process, closing jobs channel...
WP: Grab all result in one array...
                                Error while deserializing file: EOF
                                Error while deserializing file: XML syntax error on line 9: element <Amt> closed by </CdtTrfTxInf>
=== Validation Report ===
        ./testdata/example.xml - FAIL
                - EOF
        ./testdata/invalid_malformed.xml - FAIL
                - XML syntax error on line 9: element <Amt> closed by </CdtTrfTxInf>
        ./testdata/invalid_all_errors.xml - FAIL
                - IBAN cannot be blank
                - Amount must be positive
                - Currency must be 3 characters
        ./testdata/invalid_iban_country.xml - FAIL
                - IBAN has to start with a letters
        ./testdata/invalid_wrong_currency.xml - OK
        ./testdata/invalid_missing_iban.xml - FAIL
                - IBAN cannot be blank
        ./testdata/invalid_negative_amount.xml - FAIL
                - Amount must be positive
        ./testdata/valid_payment.xml - OK
        ./testdata/valid_batch_payment.xml - OK
Total files: 9 Total valid: 3, Total invalid: 6
```

## Why Go

In mainframe/COBOL environments, parallel batch processing requires complex JCL orchestration and dedicated infrastructure. Go achieves the same with a lightweight worker pool — goroutines cost ~2KB vs ~1MB for OS threads, making it ideal for high-throughput payment processing pipelines.

This pattern maps directly to cloud-native microservices architecture (EKS, Kubernetes) where Go is the de facto standard.