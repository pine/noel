package main

import (
    "fmt"
    "os/exec"
)

func checkEnv() error {
    fmt.Println("# checking choco")
    if _, err := exec.LookPath("choco"); err != nil {
        return err
    }
    
    fmt.Println("# checking ketarin")
    if _, err := exec.LookPath("ketarin"); err != nil {
        fmt.Print("\n")
        fmt.Println("Try to install used by chocolatey as following:")
        fmt.Println("$ choco install chocolateypackageupdater")
        return err
    }
    
    fmt.Println("# checking powershell")
    if _, err := exec.LookPath("powershell"); err != nil {
        return err
    }
    
    return nil
}
