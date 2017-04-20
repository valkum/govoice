package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/flosch/pongo2.v3"
	"gopkg.in/yaml.v2"
)

type item struct {
	ID           string  `yaml:"id"`
	Description  string  `yaml:"description"`
	Quantity     float64 `yaml:"quantity"`
	PricePerUnit float64 `yaml:"pricePerUnit"`
	Total        float64 `yaml:"total"`
}
type contact struct {
	Name   string `yaml:"name"`
	Street string `yaml:"street"`
	Mobile string `yaml:"mobile"`
	ZIP    string `yaml:"zip"`
	City   string `yaml:"city"`
}

type invoice struct {
	Date  string  `yaml:"date"`
	To    contact `yaml:"to"`
	Items []item  `yaml:"items"`
}

// type invoice struct {
// 	Date  string              `yaml:"date"`
// 	to    map[string]string   `yaml:"to"`
// 	items []map[string]string `yaml:"items"`
// }
var config *viper.Viper

func main() {
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName("Govoice")        // name of config file (without extension)
	config.AddConfigPath("/etc/govoice/")  // path to look for the config file in
	config.AddConfigPath("$HOME/.govoice") // call multiple times to add many search paths
	config.AddConfigPath(".")              // optionally look for config in the working directory
	err := config.ReadInConfig()           // Find and read the config file
	if err != nil {                        // Handle errors reading the config file
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		usage()
	}

	out := config.GetString("out_dir")
	from := contact{}
	err = config.UnmarshalKey("from", &from)
	if err != nil {
		log.Fatal(err)
	}

	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		inFileName := arg

		t := invoice{}
		inFile, err := ioutil.ReadFile(inFileName)
		if err != nil {
			log.Fatal(fmt.Errorf("error: %v", err))
		}

		err = yaml.Unmarshal(inFile, &t)
		if err != nil {
			log.Fatal(fmt.Errorf("error: %v", err))
		}

		var total float64

		for _, item := range t.Items {
			total += item.Total
		}
		fName := filepath.Base(inFileName)
		extName := filepath.Ext(inFileName)
		bName := fName[:len(fName)-len(extName)]
		err = os.MkdirAll(out, 0770)
		if err != nil {
			log.Panic(err)
		}

		renderHTML(from, t, bName, total, out+"/"+bName+".html")
		if config.GetBool("pdf") {
			renderPDF(out+"/"+bName+".html", out+"/"+bName+".pdf")
		}
	}
}

func renderHTML(from contact, t invoice, invoiceName string, total float64, outFileString string) {
	tpl := pongo2.Must(pongo2.FromFile(config.GetString("template")))
	log.Printf("Using template %s", config.GetString("template"))

	log.Printf("Write %s.html", outFileString)

	outFile, err := os.OpenFile(outFileString, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	err = tpl.ExecuteWriter(pongo2.Context{"from": from, "to": t.To, "date": t.Date, "items": t.Items, "name": invoiceName, "total": total}, outFile)
}

func renderPDF(inFileString string, outFileString string) {
	cmd := exec.Command("electron-pdf", inFileString, outFileString)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Command finished with error: %v", err)
	}
	log.Printf("Command finished")
}
func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s inputfile [inputfile]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}
