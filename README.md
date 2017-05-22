Bitcoin Ticker
==============

### Requirements

[Dep](https://github.com/golang/dep#usage)

### Install 

``` bash
$ mkdir -p $GOPATH/src/github.com/l-vitaly
$ git clone git@github.com:l-vitaly/btcticker.git
$ cd btcticker
$ make install
```

### The configuration 

The project contains a configuration file that allows you to use various API sources.
In addition, you can set up the exchange currency.


#### Explain toml config

``` toml
[[fetchers]]
  name = "spectrocoin"
  exchange = [
    "btc", "eur", 
    "btc", "usd", 
    "eur", "usd"
  ]
```

##### The [[fetchers]] section

This defines one fetch entry.

| Setting | Description |
| ------- | ----------- |
| `name`                | fetcher name, available btce, coindesk and spectrocoin |
| `exchange`            | currency pair from/to for exchange
