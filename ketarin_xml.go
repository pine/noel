package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "path/filepath"
    "regexp"
)

var SettingFileName = "Ketarin.xml"

func UpdateKetarinXml(name string) error {
    filePath := filepath.Join(name, SettingFileName)
    xmlBytes, err := ioutil.ReadFile(filePath)
    
    if err != nil {
        return err
    }
    
    xml := string(xmlBytes)
    
    if xml, err = replaceTargetPath(xml, name); err != nil {
        return err
    }
    
    if xml, err = replacePreviousLocation(xml); err != nil {
        return err
    }
    
    return ioutil.WriteFile(filePath, []byte(xml), os.ModePerm)
}

func replaceTargetPath(xml string, name string) (string, error) {
    pattern, err := regexp.Compile(`<TargetPath>[^<]+</TargetPath>`)
    
    if err != nil {
        return "", err
    }
    
    tempdir, err := ioutil.TempDir("", name)
    
    if err != nil {
        return "", err
    }
    fmt.Println("TargetPath = " + tempdir)
    
    targetPath := fmt.Sprintf(`<TargetPath>%s</TargetPath>`, tempdir)
    return pattern.ReplaceAllString(xml, targetPath), nil
}

func replacePreviousLocation(xml string) (string, error) {
    pattern, err := regexp.Compile(`<PreviousLocation>[^<]+</PreviousLocation>`)
    
    if err != nil {
        return "", err
    }
    
    fmt.Println(`PreviousLocation = ""`)
    
    empty := `<PreviousLocation />`
    return pattern.ReplaceAllString(xml, empty), nil
}
