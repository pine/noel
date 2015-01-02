package main

import (
    "io/ioutil"
    "encoding/json"
)

type pkgs struct {
    Manual    []string
    Automatic []string
}

func loadPkgs(path string) (*pkgs, error) {
    data, err := ioutil.ReadFile(path)
    
    if err != nil {
        return nil, err
    }
    
    var pkgs pkgs
    
    if err := json.Unmarshal(data, &pkgs); err != nil {
        return nil, err
    }
    
    return &pkgs, nil
}