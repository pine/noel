package main

import "testing"

func TestLoadConf(t *testing.T) {
    conf, err := LoadConf("test_files/noel.json")
    
    if err != nil {
        t.Errorf("Can't load noel.json")
        t.Errorf("%s", err)
        return
    }
    
    if len(conf.Manual) < 3 || conf.Manual[0] != "manual" {
        t.Errorf("Can't load manual property")
    }
    
    if len(conf.Automatic) < 3 || conf.Automatic[0] != "automatic" {
        t.Errorf("Can't load automatic property")
    }
    
    if len(conf.WithoutInstall) < 3 || conf.WithoutInstall[0] != "without" {
        t.Errorf("Can't load withoutInstall property")
    }
}