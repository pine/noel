package main

import (
    "strings"
    "io/ioutil"
    "encoding/json"
    
    "gopkg.in/yaml.v2"
)

type Conf struct {
    Manual         []string
    Automatic      []string
    WithoutInstall []string `yaml:"withoutInstall"`
}

func LoadConf(path string) (*Conf, error) {
    data, err := ioutil.ReadFile(path)
    
    if err != nil {
        return nil, err
    }
    
    var pkgs Conf
    
    unmarshal := json.Unmarshal
    
    // 拡張子で判定
    if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
        unmarshal = yaml.Unmarshal
    }
    
    if err := unmarshal(data, &pkgs); err != nil {
        return nil, err
    }
    
    return &pkgs, nil
}