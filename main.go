package main

import (
    "os"
    "fmt"
    
    "github.com/wsxiaoys/terminal"
    "github.com/shiena/ansicolor"
)

var Version = "1.0.0"
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

func main() {
    stdout := terminal.TerminalWriter { ansicolor.NewAnsiColorWriter(os.Stdout) }
    
    fmt.Println("Start Noel v" + Version + " [Chocolatey Packages Test Runner]\n")
    
    fmt.Println("Check environment:")
    if err := checkEnv(); err != nil {
        printError(stdout, err, "Failed")
        return
    } else {
        printOk(stdout, "\nSucceeded\n")
    }
    
    stdout.Print("Load settings: ")
    pkgs, err := loadPkgs("noel.json")
    
    if err != nil {
        printError(stdout, err, "Failed")
        return
    } else {
        printOk(stdout, "Succeeded")
    }
    
    fmt.Print("Detect package changes: ")
    changedPkgs, err := getChangedPackages()
    
    if err != nil {
        printError(stdout, err, "Failed")
        return
    } else {
        printOk(stdout, "Succeeded")
    }
    
    fmt.Println("\nChanged packages:")
    
    for _, pkg := range(changedPkgs) {
        fmt.Println("    " + pkg)
    }
    
    fmt.Println("\nStart tests:\n")
    
    
    for _, pkg := range(changedPkgs) {
        fmt.Print("Test for [" + pkg + "]: ")
        
        if contains(pkgs.Manual, pkg) {
            fmt.Println("Manual tests")
            
            if err := testManual(pkg); err != nil {
                printError(stdout, err, "Failed")
                break
            
            } else {
                printOk(stdout, "Succeeded\n")
            }
            continue
        }
        
        if contains(pkgs.Automatic, pkg) {
            fmt.Println("Automatic tests")
            printSkip(stdout, "Skip")
            continue
        }
        
        printSkip(stdout, "Skip")
    }
    
    printOk(stdout, "\nAll test succeeded")
}