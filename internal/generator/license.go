package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const rawLicensesRepo = "https://raw.githubusercontent.com/github/choosealicense.com/gh-pages/_licenses/"

var mapLicenses = map[string]string{
	`MIT License`: "mit.txt",
	`BSD 3-Clause "New" or "Revised" License`: "bsd-3-clause.txt",
	`ISC License`:                   "isc.txt",
	`Apache Software License 2.0`:   "apache-2.0.txt",
	`GNU General Public License v3`: "gpl-3.0.txt",
}

func GenerateLicense(idx, owner, destinationPath string) error {

	lf, ok := mapLicenses[idx]
	if !ok {
		return fmt.Errorf("%s unsupported license", idx)
	}

	resp, err := http.Get(rawLicensesRepo + lf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	s := strings.Split(string(data), "---\n")
	if len(s) < 3 {
		return errors.New("failed getting license")
	}
	licenseText := s[2]
	licenseText = strings.ReplaceAll(licenseText, "[year]", fmt.Sprintf("%d", time.Now().Year()))
	licenseText = strings.ReplaceAll(licenseText, "[fullname]", owner)

	var f *os.File
	f, err = os.Create("LICENSE")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(licenseText)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := GenerateLicense("MIT License", "igkx", "."); err != nil {
		panic(err)
	}
}
