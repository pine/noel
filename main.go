package main

import (
    "os"
    "fmt"
    "flag"
    
    "github.com/wsxiaoys/terminal"
    "github.com/shiena/ansicolor"
    "gopkg.in/fatih/set.v0"
)

var Version = "1.0.0-beta"
var SettingFile = "noel.json"

func getChangedPackages() ([]string, error) {
    return getChangedRootDirs(1)
}

func contains(col []string, val string) bool {
    for _, x := range(col) {
        if x == val {
            return true
        }
    }
    
    return false
}

func printOk(w terminal.TerminalWriter, msg string) {
    w.Color("g").
        Print(msg).
        Print("\n").
        Reset()
}

func printSkip(w terminal.TerminalWriter, msg string) {
    w.Color("b").
        Print(msg).
        Print("\n").
        Reset()
}

func printError(w terminal.TerminalWriter, err error, msg string) {
    w.Color("r").
        Print("\n").
        Print(err).
        Print("\n\n").
        Print(msg).
        Print("\n").
        Reset()
}

func getAllPackages(pkgs *pkgs) []string {
    allPkgs := pkgs.Manual
    
    for _, pkg := range(pkgs.Automatic) {
        allPkgs = append(allPkgs, pkg)
    }
    
    return allPkgs
}

func getTestPkgs(hackPkgs []string, changedPkgs []string, pkgs *pkgs) []string {
    testPkgs := set.New()
    
    for _, pkg := range(changedPkgs) {
        testPkgs.Add(pkg)
    }
    
    for _, pkg := range(hackPkgs) {
        if pkg == "<all>" {
           allPkgs := getAllPackages(pkgs)
           
           for _, pkgByConf := range(allPkgs) {
               testPkgs.Add(pkgByConf)
           }
           
           continue
        }
        
        testPkgs.Add(pkg)
    }
    
    return set.StringSlice(testPkgs)
}

func main() {
    var install bool
    flag.BoolVar(&install, "install", false, "Install package")
    flag.Parse()
    
    stdout := terminal.TerminalWriter { ansicolor.NewAnsiColorWriter(os.Stdout) }
    
    fmt.Println("Start Noel v" + Version + " [Chocolatey Packages Test Runner]\n")
    
    // ----------------------------------------------------
    
    fmt.Println("Check environment:")
    if err := checkEnv(); err != nil {
        printError(stdout, err, "Failed")
        return
    } else {
        printOk(stdout, "\nSucceeded\n")
    }
    
    // ----------------------------------------------------
    
    stdout.Print("Load settings: ")
    pkgs, err := loadPkgs("noel.json")
    
    if err != nil {
        printError(stdout, err, "Failed")
        return
    } else {
        printOk(stdout, "Succeeded")
    }
    
    // ----------------------------------------------------
    
    fmt.Print("Detect hacks of commit messages: ")
    
    hackPkgs, _ := getPackageNamesOfCommitMessage()
    
    if err != nil {
        printError(stdout, err, "Failed")
        return
    } else {
        printOk(stdout, "Succeeded")
    }
    
    // ----------------------------------------------------
    
    fmt.Println("\nHack packages:")
    
    for _, pkg := range(hackPkgs) {
        fmt.Println("    " + pkg)
    }
    
    fmt.Printf("\n    %d packages\n\n", len(hackPkgs))
    
    // ----------------------------------------------------
    
    fmt.Print("Detect package changes: ")
    changedPkgs, err := getChangedPackages()
    
    if err != nil {
        printError(stdout, err, "Failed")
        return
    } else {
        printOk(stdout, "Succeeded")
    }
    
    // ----------------------------------------------------
    
    fmt.Println("\nChanged packages:")
    
    for _, pkg := range(changedPkgs) {
        fmt.Println("    " + pkg)
    }
    
    fmt.Printf("\n    %d packages", len(changedPkgs))
    
    if len(changedPkgs) == 0 && len(hackPkgs) == 0 {
        printSkip(stdout, "\n\nNo changed\n\n")
        return
    }
    
    // ----------------------------------------------------
    
    fmt.Println("\n\nStart tests:\n")
    
    testPkgs := getTestPkgs(hackPkgs, changedPkgs, pkgs)
    
    for _, pkg := range(testPkgs) {
        fmt.Print("Test for [" + pkg + "]: ")
        
        data := TestData {
            Name: pkg,
            Install: install,
            Timeout: 60,
        }
        
        if contains(pkgs.Manual, pkg) {
            fmt.Println("Manual tests")
            
            if err := TestManual(data); err != nil {
                printError(stdout, err, "Failed")
                return
            
            } else {
                printOk(stdout, "Succeeded\n")
            }
            continue
        }
        
        if contains(pkgs.Automatic, pkg) {
            fmt.Println("Automatic tests")
            
            if err := TestAutomatic(data); err != nil {
                printError(stdout, err, "Failed")
                return
            
            } else {
                printOk(stdout, "Succeeded\n")
            }
            continue
        }
        
        printSkip(stdout, "Skip")
    }
    
    printOk(stdout, "\nAll test succeeded")
}