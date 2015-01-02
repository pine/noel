Noel
--------
[![Circle CI](https://circleci.com/gh/pine613/noel.svg?style=svg)](https://circleci.com/gh/pine613/noel)

Noel is not a Blu-ray disc. Noel is a test runnner of chocolatey packages!

## Usage
### Create settings
noel.json
```json
{
    "manual": ["manual", "package", "name"],
    "automatic": ["automatic", "package", "name"]
}
```

### Run
```
$ go get github.com/pine613/noel

$ cd your_chocolatey_repo
$ noel
```

## How to test non-changed packages
Please include `[pkg-name]` to commit message. If you want to test all packages, please try `[<all>]` meta name.