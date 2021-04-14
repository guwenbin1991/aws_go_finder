package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"unicode"
	// "gopkg.in/yaml.v2"
)

var (
	FileName string
)

func main() {
	http.HandleFunc("/", goServer)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	//json
	fmt.Scanln(&FileName)
	GetFile(FileName)

}

func GetFile(filename string) {
	if path.Ext(FileName) == ".json" {
		ReadJson(FileName)
		fmt.Println("this is a json file")
	} else if path.Ext(FileName) == ".yaml" {
		fmt.Println("this is a yaml file")
		//	t := map[string]interface{}{}
		buffer, yaml_err := ioutil.ReadFile(FileName)
		//	yaml_err = yaml.Unmarshal(buffer, &t)
		println("buffer %s", buffer)
		if yaml_err != nil {
			log.Fatalf(yaml_err.Error())
		}
		// ReadYaml(buffer)
	} else if path.Ext(FileName) == ".backbone" || path.Ext(FileName) == ".mobile" || path.Ext(FileName) == ".system" {
		DecodeErl(FileName)
	} else {
		fmt.Println("file type do not exsit")
	}
}

func goServer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintln(w, r.Form)
	if r.Method == "GET" {
		fmt.Println("Geoge")
		fmt.Println(r.URL.MarshalBinary())
	}
}

func ReadJson(jsonFileName string) {
	jsonFile, err := os.Open(jsonFileName)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	fmt.Println(string(byteValue))
}

// func ReadYaml(yamlFileName []byte) {
// 	t2, err := yaml.Unmarshal(yamlFileName)
// 	if err != nil {
// 		fmt.Printf("err: %v\n", err)
// 		return
// 	}
// 	fmt.Println(string(t2))
// }

func DecodeErl(erlFileName string) {
	file, err := os.Open(erlFileName)
	if err != nil {
		log.Fatalf("Error when opening files: %s", err)
	}
	fileScanner := bufio.NewScanner(file)
	// read line by line
	fmt.Print("{")
	for fileScanner.Scan() {
		if fileScanner.Text() == "" || fileScanner.Text()[0] == '%' {
			continue
		} else {
			newLine := ""
			oldLine := fileScanner.Text()
			for i := 0; i < len(oldLine); i++ {
				if oldLine[i] == ',' && unicode.IsLetter(rune(oldLine[i-1])) && !unicode.IsPunct(rune(oldLine[i+1])) {
					newLine += `":"`
				} else if oldLine[i] == ',' && unicode.IsDigit(rune(oldLine[i-1])) {
					newLine += `",`
				} else if oldLine[i] == ',' && !unicode.IsPunct(rune(oldLine[i+1])) {
					newLine += `"`
				} else if oldLine[i] == '.' {
					continue
				} else if oldLine[i] == '{' {
					newLine += `{"`
				} else if oldLine[i] == '}' && !unicode.IsPunct(rune(oldLine[i-1])) {
					newLine += `"}`
				} else {
					newLine += string(oldLine[i])
				}
			}

			fmt.Print(strings.Replace(newLine, "undefined", "null", -1))
		}
	}
	fmt.Print("}")
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}
	file.Close()
}
