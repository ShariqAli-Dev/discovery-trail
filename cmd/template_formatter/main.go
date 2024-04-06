package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	routesPath, err := filepath.Abs("ui/html")
	if err != nil {
		log.Fatal(err)
	}
	err = filepath.Walk(routesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != routesPath && info.Name() == "partials" {
			searchFileName := getBuiltJS("index")
			baseTemplatePath := routesPath + "/partials/base_templ.go"
			baseTemplate, err := os.ReadFile(baseTemplatePath)
			if err != nil {
				return err
			}
			baseTemplateContent := string(baseTemplate)
			templateContent := strings.Replace(baseTemplateContent, "index.js", searchFileName+".js", 1)
			if err := os.WriteFile(baseTemplatePath, []byte(templateContent), 0644); err != nil {
				return err
			}
			return filepath.SkipDir
		}

		if strings.HasSuffix(info.Name(), ".go") {
			fileName := info.Name()
			fileNameWithoutSuffix := fileName[:len(fileName)-9]
			searchFileName := getBuiltJS(fileNameWithoutSuffix + "-")
			if searchFileName == "" {
				return nil
			}

			template, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			templateContent := string(template)
			searchString := fmt.Sprintf(`BaseHTML("%s",`, fileNameWithoutSuffix)
			replacementString := fmt.Sprintf(`BaseHTML("%s",`, searchFileName)

			templateContent = strings.Replace(templateContent, searchString, replacementString, 1)
			if err := os.WriteFile(path, []byte(templateContent), 0644); err != nil {
				log.Fatal(err)
			}

			return nil
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)

	}
}

func getBuiltJS(fileBase string) string {
	assetsDirectoryPath, err := filepath.Abs("ui/static/dist/assets")
	if err != nil {
		log.Fatal(err)
	}
	dir, err := os.Open(assetsDirectoryPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()

	files, err := dir.Readdirnames(-1)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if strings.Contains(file, fileBase) {
			return file[:len(file)-3]
		}
	}
	return ""
}
