package validator

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
)

var validCountries = map[string]bool{
	"PL": true,
	"DE": true,
	"GB": true,
	"NL": true,
}

var ibanLengths = map[string]int{
	"PL": 28,
	"DE": 22,
	"GB": 22,
	"NL": 18,
}

type FileStatus struct {
	Filename string
	IsValid  bool
	Error    []string
}

type Document struct {
	XMLName          xml.Name         `xml:"Document"`
	CstmrCdtTrfInitn CstmrCdtTrfInitn `xml:"CstmrCdtTrfInitn"`
}

type CstmrCdtTrfInitn struct {
	PmtInf PmtInf `xml:"PmtInf"`
}

type PmtInf struct {
	CdtTrfTxInf []CdtTrfTxInf `xml:"CdtTrfTxInf"`
}

type CdtTrfTxInf struct {
	Amt      Amt      `xml:"Amt"`
	CdtrAcct CdtrAcct `xml:"CdtrAcct"`
}

type Amt struct {
	InstdAmt InstdAmt `xml:"InstdAmt"`
}

type InstdAmt struct {
	Ccy   string `xml:"Ccy,attr"`  // atrybut XML: <InstdAmt Ccy="PLN">
	Value string `xml:",chardata"` // wartość: 100.00
}

type CdtrAcct struct {
	Id Id `xml:"Id"`
}

type Id struct {
	IBAN string `xml:"IBAN"`
}

// Private method
func validateFile(path string) (FileStatus, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		fmt.Printf("\t\t\t\tError reading file: %s\n", err)
		return FileStatus{}, err
	}

	var doc Document
	err = xml.Unmarshal(data, &doc)
	if err != nil {
		fmt.Printf("\t\t\t\tError while deserializing file: %s\n", err)
		return FileStatus{}, err
	}

	result := FileStatus{Filename: path}
	for _, tx := range doc.CstmrCdtTrfInitn.PmtInf.CdtTrfTxInf {
		iban := tx.CdtrAcct.Id.IBAN
		currency := tx.Amt.InstdAmt.Ccy
		amount := tx.Amt.InstdAmt.Value

		if iban == "" {
			result.Error = append(result.Error, "IBAN cannot be blank")
		} else if !validCountries[iban[:2]] {
			result.Error = append(result.Error, "IBAN has to start with a letters")
		} else if len(iban) != ibanLengths[iban[:2]] {
			result.Error = append(result.Error, fmt.Sprintf("IBAN must be %d characters", ibanLengths[iban[:2]]))
		}

		value, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			result.Error = append(result.Error, err.Error())
		}

		if value <= 0 {
			result.Error = append(result.Error, "Amount must be positive")
		}

		if currency == "" {
			result.Error = append(result.Error, "Currency cannot be blank")
		} else if len(currency) != 3 {
			result.Error = append(result.Error, "Currency must be 3 characters")
		}
	}
	result.IsValid = len(result.Error) == 0
	return result, nil
}

// Public method
// Worker pool.
// Arg: paths - files to process, workers - jobs to run the same time
func RunWorkerPool(paths []string, workers int) []FileStatus {
	//input channel, string as input, and buffor allows to decide how many files shoule be loaded to memory
	jobs := make(chan string, len(paths))
	//output channel, how many result we can keep in memory before receive
	results := make(chan FileStatus, len(paths))

	//Start workers async
	for w := 1; w <= workers; w++ {
		fmt.Printf("WP: Starting worker %d\n", w)
		go func() {
			for path := range jobs {
				result, err := validateFile(path)
				//if err == nil {
				//	results <- result
				//}
				if err != nil {
					results <- FileStatus{Filename: path, IsValid: false, Error: []string{err.Error()}}
				} else {
					results <- result
				}
			}
		}()
	}

	fmt.Printf("WP: Loading data to process...\n")
	for _, path := range paths {
		jobs <- path
	}

	fmt.Printf("WP: No more files to process, closing jobs channel...\n")
	close(jobs)

	fmt.Printf("WP: Grab all result in one array...\n")
	var allResults []FileStatus
	for range paths {
		allResults = append(allResults, <-results)
	}
	return allResults
}
