package main

import (
    "fmt"
    "strings"
    "os/exec"
    "io/ioutil"
    "regexp"
    
    "code.google.com/p/go.text/encoding/japanese"
    "code.google.com/p/go.text/transform"
    
    "gopkg.in/fatih/set.v0"
)

func getChangedRootDirs(depth int) ([]string, error) {
    paths, err := getChangedFilePaths(depth)
    
    if err != nil {
        return nil, err
    }
    
    pattern, err := regexp.Compile("([^/]*)/")
    
    if err != nil {
        return nil, err
    }
    
    dirs := set.New()
    
    for _, path := range(paths) {
        matches := pattern.FindStringSubmatch(path)
        
        if len(matches) >= 2 {
            dirs.Add(matches[1])
        }
    }
    
    return set.StringSlice(dirs), nil
}

func getChangedFilePaths(depth int) ([]string, error) {
    if changes, err := getChanges(depth); err != nil {
        return nil, err
    } else {
        lines := strings.Split(string(changes), "\n")
        return lines, nil
    }
}

func getChanges(depth int) ([]byte, error) {
    rev := fmt.Sprintf("HEAD~%d", depth)
    cmd := exec.Command("git", "diff", "--name-only", rev)
    
    stdoutpipe, err := cmd.StdoutPipe()
    if err != nil {
        return nil, err
    }
    
    defer stdoutpipe.Close()
    
    if err = cmd.Start(); err != nil {
        return nil, err
    }
    
    stdout, err := ioutil.ReadAll(
        transform.NewReader(stdoutpipe, japanese.ShiftJIS.NewDecoder()))
    if err != nil {
        return nil, err
    }
    
    if err = cmd.Wait(); err != nil {
        return nil, err
    }
    
    return stdout, nil
}
