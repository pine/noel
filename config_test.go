package main

import "testing"

func TestLoadJsonConf(t *testing.T) {
    testLoadConfByPath(t, "test_files/noel.json")
}

func TestLoadYamlConf(t *testing.T) {
    testLoadConfByPath(t, "test_files/noel.yml")
    testLoadConfByPath(t, "test_files/noel.yaml")
}

func testLoadConfByPath(t *testing.T, path string) {
    conf, err := LoadConf(path)
    
    if err != nil {
        t.Errorf("Can't load noel.json")
        t.Error(err)
        t.Error(conf)
        return
    }
    
    t.Log(conf)
    
    if len(conf.Manual) < 3 || conf.Manual[0] != "manual" {
        t.Errorf("Can't load manual property")
        t.Error(conf.Manual)
    }
    
    if len(conf.Automatic) < 3 || conf.Automatic[0] != "automatic" {
        t.Errorf("Can't load automatic property")
        t.Error(conf.Automatic)
    }
    
    if len(conf.WithoutInstall) < 3 || conf.WithoutInstall[0] != "without" {
        t.Errorf("Can't load withoutInstall property")
        t.Error(conf.WithoutInstall)
    }
}