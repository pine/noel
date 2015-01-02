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

func getPackageNamesOfCommitMessage() ([]string, error) {
    msg, err := getLastCommitMessage()
    
    if err != nil {
        return nil, err
    }
    
    if pattern, err := regexp.Compile(`\[([A-Za-z0-9<>_\.-]+)\]`); err != nil {
        return nil, err
    } else {
        matches := pattern.FindAllStringSubmatch(msg, -1)
        names := []string{}
        
        for _, pair := range(matches) {
            names = append(names, pair[1])
        }
        
        return names, nil
    }
}

func getLastCommitMessage() (string, error) {
    cmd := exec.Command("git", "log", "-1", "--pretty=format:%s")
    msg, err := getCommandStdout(cmd)
    
    if err != nil {
        return "", err
    }
    
    return msg, nil
}

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
        lines := strings.Split(changes, "\n")
        return lines, nil
    }
}

func getChanges(depth int) (string, error) {
    rev := fmt.Sprintf("HEAD~%d", depth)
    cmd := exec.Command("git", "diff", "--name-only", rev)
    
    return getCommandStdout(cmd)
}

func getCommandStdout(cmd *exec.Cmd) (string, error) {
    stdoutpipe, err := cmd.StdoutPipe()
    if err != nil {
        return "", err
    }
    
    defer stdoutpipe.Close()
    
    if err = cmd.Start(); err != nil {
        return "", err
    }
    
    stdout, err := ioutil.ReadAll(
        transform.NewReader(stdoutpipe, japanese.ShiftJIS.NewDecoder()))
    if err != nil {
        return "", err
    }
    
    if err = cmd.Wait(); err != nil {
        return "", err
    }
    
    return string(stdout), nil
}
