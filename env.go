package main

import (
    "fmt"
    "os/exec"
)

func checkEnv() error {
    fmt.Println("# checking choco")
    _, err := exec.LookPath("choco")
    
    if err != nil {
        return err
    }
    
    return nil
}