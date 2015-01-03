package main

import (
    "io/ioutil"
    "encoding/json"
)

type Conf struct {
    Manual    []string
    Automatic []string
}

func LoadConf(path string) (*Conf, error) {
    data, err := ioutil.ReadFile(path)
    
    if err != nil {
        return nil, err
    }
    
    var pkgs Conf
    
    if err := json.Unmarshal(data, &pkgs); err != nil {
        return nil, err
    }
    
    return &pkgs, nil
}