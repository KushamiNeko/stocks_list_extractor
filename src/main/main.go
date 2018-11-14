package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	lists = `/home/neko/programming_projects/golang/ibd_stocks_list_parser/stocks_list`
)

func main() {

	fmt.Printf("\n\n")

	files, err := ioutil.ReadDir(lists)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {

		if filepath.Ext(file.Name()) != ".xls" && filepath.Ext(file.Name()) != ".csv" {
			continue
		}

		if filepath.Ext(file.Name()) == ".csv" {
			//fmt.Printf("Removing File: %q\n", filepath.Join(lists, file.Name()))
			err = os.Remove(filepath.Join(lists, file.Name()))
			if err != nil {
				fmt.Println(err)
				return
			}

		}

		if filepath.Ext(file.Name()) == ".xls" {

			err = xlsToCsv(filepath.Join(lists, file.Name()), false)
			if err != nil {
				fmt.Println(err)
				return
			}

			csvFile := filepath.Join(lists, strings.Replace(file.Name(), filepath.Ext(file.Name()), ".csv", -1))
			//fmt.Println(csvFile)

			err = readCSV(csvFile)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	}

	fmt.Printf("\n\n")

}

func xlsToCsv(src string, debug bool) error {

	if debug {
		fmt.Printf("Converting File: %q\n", src)
	}

	var outb bytes.Buffer
	var errb bytes.Buffer

	cmd := exec.Command("libreoffice", "--headless", "--convert-to", "csv", src, "--outdir", lists)

	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()

	if err != nil {
		fmt.Printf("\n%v\n", errb.String())
		return err
	}

	if debug {
		fmt.Printf("\n%v\n", outb.String())
		//fmt.Printf("Removing File: %q\n", src)
	}

	//err = os.Remove(src)
	//if err != nil {
	//fmt.Printf("\n%v\n", err)
	//return err
	//}

	return nil
}

func readCSV(file string) error {

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	printSymbol(records)

	return nil
}

func printSymbol(records [][]string) {

	printing := false

	for _, record := range records {
		col := record[0]

		if col == "Symbol" {
			printing = true
			continue
		}

		if printing {
			if col == "" {
				printing = false
				break
			}

			fmt.Println(col)
		}
	}
}
