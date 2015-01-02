package main

import (
    "fmt"
    "os"
    "os/exec"
)

func testManual(name string) error {
    fmt.Println("> cd " + name)
    if err := os.Chdir(name); err != nil {
        return err
    }
    
    fmt.Println("> cd")
    if dir, err := os.Getwd(); err != nil {
        return err
    } else {
        fmt.Println(dir)
    }
    
    fmt.Println("> cpack")
    cmd := exec.Command("cpack")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    if err := cmd.Run(); err != nil {
        return err
    }
    
    fmt.Println("> cd ..")
    if err := os.Chdir(".."); err != nil {
        return err
    }
    
    fmt.Print("\n")
    
    return nil
}