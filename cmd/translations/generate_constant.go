package main

import (
	"crypto/sha256"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

func main() {
	bytes, err := os.ReadFile("assets/lang/en.yml")
	if err != nil {
		panic(err)
	}

	var translations map[string]interface{}

	fmt.Println("Generating translations.go...")
	if err := yaml.Unmarshal(bytes, &translations); err != nil {
		panic(err)
	}

	contents := "package translations\n\nconst (\n"

	//convertMapping := map[string]string{}

	for key := range translations {
		parts := strings.Split(key, ".")
		key2 := ""

		for _, part := range parts {
			key2 += strings.Title(part)
		}

		//convertMapping[key] = key2

		contents += "\t" + key2 + " = \"" + key + "\"\n"
	}

	contents += ")\n"

	header := "// Code generated by cmd/translations/generate_constant.go DO NOT EDIT.\n// Hash: " + fmt.Sprintf("%x", sha256.Sum256([]byte(contents))) + "\n\n"

	if err := os.WriteFile("translations/translations.go", []byte(header+contents), 0644); err != nil {
		panic(err)
	}

	fmt.Println("Done!")

	//files, err := readDirRecursive("practice")
	//if err != nil {
	//	panic(err)
	//}
	//
	//for _, file := range files {
	//	if strings.HasSuffix(file, ".go") {
	//		fmt.Println("Processing", file)
	//
	//		bytes2, err := os.ReadFile(file)
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		contents2 := string(bytes2)
	//		for key := range translations {
	//			contents2 = strings.ReplaceAll(contents2, "\""+key+"\"", "translations."+convertMapping[key])
	//		}
	//
	//		if err := os.WriteFile(file, []byte(contents2), 0644); err != nil {
	//			panic(err)
	//		}
	//	}
	//}
}

//func readDirRecursive(dir string) ([]string, error) {
//	var files []string
//
//	entries, err := os.ReadDir(dir)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, entry := range entries {
//		if entry.IsDir() {
//			subFiles, err := readDirRecursive(dir + "/" + entry.Name())
//			if err != nil {
//				return nil, err
//			}
//			files = append(files, subFiles...)
//		} else {
//			files = append(files, dir+"/"+entry.Name())
//		}
//	}
//
//	return files, nil
//}
