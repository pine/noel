package main

import (
    "fmt"
    "os/exec"
)

func checkEnv() error {
    fmt.Println("# checking cpack")
    _, err := exec.LookPath("cpack")
    
    if err != nil {
        return err
    }
    
    return nil
}