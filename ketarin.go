package main

import (
    "fmt"
    "os"
    "os/exec"
    "io"
    "io/ioutil"
    "path/filepath"
    "time"
    "errors"
    "encoding/xml"
    "regexp"
)

var KetarinAppDataDirName = "Ketarin"
var DatabaseFileName = "jobs.db"
var ChocopkgupConfigPath = `C:\tools\ChocolateyPackageUpdater\chocopkgup.exe.config`
var SettingFileName = "Ketarin.xml"
var PreUpdateCommand = `chocopkgup /p {appname} /v {version} /u "{preupdate-url}" /u64 "{url64}" /pp "{file}" /disablepush`

// 最初の行に改行を入れないこと!!
var CustomColumns = `System.String:<?xml version="1.0" encoding="utf-16"?>
<dictionary>
  <item>
    <key>
      <string>Version</string>
    </key>
    <value>
      <string>{version}</string>
    </value>
  </item>
</dictionary>
`

type AppSettings struct {
    XMLName xml.Name `xml:"appSettings"`
    Adds    []*Add   `xml:"add"`
}

type Add struct {
    XMLName xml.Name `xml:"add"`
    Key     string   `xml:"key,attr"`
    Value   string   `xml:"value,attr"`
}

type SupportedRuntime struct {
    XMLName xml.Name `xml:"supportedRuntime"`
    Version string   `xml:"version,attr"`
    Sku     string   `xml:"sku,attr"`
}

type Startup struct {
    XMLName           xml.Name          `xml:"startup"`
    SupportedRuntime  *SupportedRuntime
}

type Configuration struct {
    XMLName     xml.Name      `xml:"configuration"`
    AppSettings *AppSettings  `xml:"appSettings"`
    Startup     *Startup      `xml:"startup"`
}

func getKetarinDatabase() string {
    appdata := filepath.Join(os.Getenv("APPDATA"), KetarinAppDataDirName)
    dbPath := filepath.Join(appdata, DatabaseFileName)
    
    return dbPath
}

func getKetarinDatabaseBackup() string {
    appdata := filepath.Join(os.Getenv("APPDATA"), KetarinAppDataDirName)
    date := time.Now().Format("2006-01-02-150405")
    dbPath := filepath.Join(appdata, DatabaseFileName + "_" + date + ".noel.bak")
    
    return dbPath
}

func SwapKetarinDatabase() error {
    dbFile, err := os.Open(getKetarinDatabase())
    
    if err != nil {
        return err
    }
    
    defer dbFile.Close()
    
    destFile, err := os.Create(DatabaseFileName)
    
    if err != nil {
        return err
    }
    
    defer destFile.Close()
    
    
    bakPath := getKetarinDatabaseBackup()
    bakFile, err := os.Create(bakPath)
    
    if err != nil {
        return err
    }
    
    defer bakFile.Close()
    
    if _, err = io.Copy(destFile, dbFile); err != nil {
        return err
    }
    
    _, err = io.Copy(bakFile, dbFile)
    
    return err
}

func RestoreKetarinDatabase() error {
    dbFile, err := os.Create(getKetarinDatabase())
    
    if err != nil {
        return err
    }
    
    defer dbFile.Close()
    
    swapFile, err := os.Open(DatabaseFileName)
    
    if err != nil {
        return err
    }
    
    defer swapFile.Close()
    
    _, err = io.Copy(dbFile, swapFile)
    return err
}

func ClearKetarinDatabase() error {
    dbPath := getKetarinDatabase()
    return os.Remove(dbPath)
}


func SetChocopkgupPackageFolder(pkgDir string) error {
    fmt.Println("> Change chocopkgup.exe.config")
    
    path := ChocopkgupConfigPath
    xmlFile, err := ioutil.ReadFile(path)
    
    if err != nil {
        return err
    }
    
    var conf Configuration
    if err = xml.Unmarshal(xmlFile, &conf); err != nil {
        return err
    }
    
    appSettings := conf.AppSettings
    
    for _, add := range(appSettings.Adds) {
        if add.Key == "PackagesFolder" {
            // already same directory
            if add.Value == pkgDir {
                return nil
            }
            
            fmt.Println("PackagesFolder (old): " + add.Value)
            fmt.Println("PackagesFolder (new): " + pkgDir)
            
            add.Value = pkgDir
            newXmlFile, err := xml.Marshal(conf)
            
            if err != nil {
                return err
            }
            
            validXmlFile := xml.Header + string(newXmlFile)
            err = ioutil.WriteFile(path, []byte(validXmlFile), os.ModePerm)
            if err != nil {
                return err
            } else {
                return nil
            }
        }
    }
    
    return errors.New(`Can't find PackagesFolder setting`)
}

func WaitKetarinProcess() error {
    cmd := exec.Command(
        "powershell",
        "-NoProfile", "-ExecutionPolicy", "unrestricted",
        "-Command", "Wait-Process -Name ketarin")
    
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    return cmd.Run()
}

func RunKetarin() error {
    tmpDir, err := ioutil.TempDir("", "nodel")
    
    if err != nil {
        return err
    }
    
    logPath := filepath.Join(tmpDir, "ketarin.log")
    cmd := exec.Command("ketarin", "/silent", "/log=" + logPath)
    
    if err := cmd.Run(); err != nil {
        return err
    }
    
    if err:= WaitKetarinProcess(); err != nil {
        return err
    }
    
    log, err := ioutil.ReadFile(logPath);
    
    if err != nil {
        return err
    }
    
    fmt.Println(string(log))
    
    return os.RemoveAll(tmpDir)
}

func InstallKetarinSetting(data TestData) error {
    settingPath := filepath.Join(data.Name, SettingFileName)
    fmt.Println(settingPath)
    
    if _, err := os.Stat(settingPath); err != nil {
        return errors.New("Setting file not found!\n" + settingPath)
    }
    
    cmd := exec.Command("ketarin", "/install=" + settingPath, "/exit")
    
    if err := cmd.Run(); err != nil {
        return err
    }
    
    return WaitKetarinProcess()
}

func UpdateKetarinSettings() error {
    db := NewKetarinDb(getKetarinDatabase())
    defer db.Close()
    
    if err := db.SetSetting("PreUpdateCommand", PreUpdateCommand); err != nil {
        return err
    }
    
    return db.SetSetting("CustomColumns", CustomColumns)
}

func FixPkgTargetPath(name string) error {
    filePath := filepath.Join(name, SettingFileName)
    xml, err := ioutil.ReadFile(filePath)
    
    if err != nil {
        return err
    }
    
    pattern, err := regexp.Compile(`<TargetPath>[^<]+</TargetPath>`)
    
    if err != nil {
        return err
    }
    
    tempdir, err := ioutil.TempDir("", name)
    
    if err != nil {
        return err
    }
    fmt.Println(tempdir)
    
    targetPath := fmt.Sprintf(`<TargetPath>%s</TargetPath>`, tempdir)
    replaced := pattern.ReplaceAllString(string(xml), targetPath)
    
    return ioutil.WriteFile(filePath, []byte(replaced), os.ModePerm)
}

func ClearOutputDir() error {
    return os.RemoveAll("_output")
}

func InstallKetarinPkg(data TestData) error {
    pkgDir := filepath.Join("_output", data.Name)
    infos, err := ioutil.ReadDir(pkgDir)
    
    if err != nil {
        return err
    }
    
    for _, info := range(infos) {
        name := info.Name()
        
        if info.IsDir() {
            if matched, _ := regexp.MatchString(`^[0-9\.]+$`, name); matched {
                dir := filepath.Join(pkgDir, name)
                
                if err := os.Chdir(dir); err != nil {
                    return err
                }
                
                return InstallChocoPkg(data)
            }
        }
    }
    
    return nil
}

func TestKetarinAutomatic(data TestData) error {
    fmt.Println("> cd")
    wd, err := os.Getwd();
    if err != nil {
        return err
    } else {
        fmt.Println(wd)
    }
    
    fmt.Println("> Clear output dir")
    if err := ClearOutputDir(); err != nil {
        return err
    }
    
    
    fmt.Println("> Swap Ketarin database")
    if err := SwapKetarinDatabase(); err != nil {
        return err
    }
    
    fmt.Println("> Clear ketarin database")
    if err := ClearKetarinDatabase(); err != nil {
        return err
    }
    
    defer func(){
        fmt.Println("> Restore ketarin database")
        RestoreKetarinDatabase()
    }()
    
    fmt.Println("> Set chocopkgup PackageFolder")
    if err := SetChocopkgupPackageFolder(wd); err != nil {
        return err
    }
    
    fmt.Println("> Fix Package TargetPath")
    if err := FixPkgTargetPath(data.Name); err != nil {
        fmt.Println(err)
    }
    
    fmt.Println("> Install ketarin settings")
    
    if err := InstallKetarinSetting(data); err != nil {
        fmt.Println(err)
        return err
    }
    
    fmt.Println("> Update ketarin settings")
    
    if err := UpdateKetarinSettings(); err != nil {
        fmt.Println(err)
        return err
    }
    
    fmt.Println("> Run ketarin")
    
    if err := RunKetarin(); err != nil {
        fmt.Println(err)
        return err
    }
    
    if data.Install {
        if err := InstallKetarinPkg(data); err != nil {
            return err
        }
    }
    
    fmt.Println("> cd")
    if err := os.Chdir(wd); err != nil {
        return err
    }
    
    return nil
}