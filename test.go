package main

import (
    "fmt"
    "os"
    "os/exec"
)

func testManual(name string, install bool) error {
    fmt.Println("> cd " + name)
    if err := os.Chdir(name); err != nil {
        return err
    }
    
    fmt.Println("> cd")
    wd, err := os.Getwd();
    if err != nil {
        return err
    } else {
        fmt.Println(wd)
    }
    
    fmt.Println("> choco pack")
    cmd := exec.Command("choco", "pack")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    if err := cmd.Run(); err != nil {
        return err
    }
    
    if install {
        fmt.Println("> choco install " + name)
        cmd = exec.Command("choco", "install", name, "-Force", "-Source", `"%cd%;http://chocolatey.org/api/v2/"`)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        
        if err := cmd.Start(); err != nil {
            return err
        }
        
        if err := cmd.Wait(); err != nil {
            return err
        }
    }
    
    fmt.Println("> cd ..")
    if err := os.Chdir(".."); err != nil {
        return err
    }
    
    fmt.Print("\n")
    
    return nil
}