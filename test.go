package main

import (
    "fmt"
    "os"
    "os/exec"
    "time"
    "errors"
)


type TestData struct {
    Name    string
    Install bool
    Timeout int
}

func InstallChocoPkg(data TestData) error {
    fmt.Println("> choco install " + data.Name)
    
    cmd := exec.Command("choco", "install", data.Name, "-Force", "-Source", `"%cd%;http://chocolatey.org/api/v2/"`)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    if err := cmd.Start(); err != nil {
        return err
    }
    
    done := make(chan error, 1)
    
    go func(){
        done <- cmd.Wait()
    }()
    
    select {
        case <- time.After(time.Duration(data.Timeout) * time.Second):
            if err := cmd.Process.Kill(); err != nil {
                <- done
                return err
            }
            
            <- done
            return errors.New("Timeout Error")
        
        case err := <- done:
            if err != nil {
                return err
            }
    }
    
    return nil
}

func TestManual(data TestData) error {
    fmt.Println("> cd " + data.Name)
    if err := os.Chdir(data.Name); err != nil {
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
    
    if data.Install {
        if err := InstallChocoPkg(data); err != nil {
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

func TestAutomatic(data TestData) error {
    return TestKetarinAutomatic(data)
}
