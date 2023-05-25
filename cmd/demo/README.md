## About

This folder contains system test demo for the ws connector, which is going to receive exported data from a local testnet

### How it works

These scripts provide a demo to help you configure and test your ws connector/driver to receive data from a local
testnet. The scripts found in this `demo` folder:

- `local-testnet.sh` creates a local testnet with one shard (0) + metachain. Each shard has 2 nodes (1 validator + 1
  observer)
- `observer-outport.sh` creates an observer in a specified shard (shard/metachain). This node acts an outport driver for
  the ws connector. When started, it will sync from genesis and export data (blocks, validators info, etc.) to your ws
  driver connector.

Inside `cmd/connector` folder, there is a `main.go` file which represents the ws connector (data receiver).

> **TIP**: Once you have experimented with this local setup, you might want to test the driver with exported data from
> an observing squad (mainnet/devnet/testnet). [Here](https://github.com/multiversx/mx-chain-observing-squad) you can
> find how to run your own observing squad. In order to enable nodes to export data, one needs to set
> **[HostDriverConfig].Enabled = true**
> [from this config file](https://github.com/multiversx/mx-chain-go/blob/master/cmd/node/config/external.toml)

## How to use

1. Create and start a local testnet:

```bash
cd cmd/demo
bash ./local-testnet.sh new
bash ./local-testnet.sh start
```

2. Start the exporter node in your desired shard (metachain/shard):

```bash
bash ./observer-outport.sh shard
```

3. Inside `cmd/connector`, start your driver to receive exported data:

```bash
go build
./connector
```

Once the setup is ready, the connected driver will start receiving data from the observer. One can see a similar log
info:

```
INFO [2023-05-24 13:53:49.234]   received payload                         topic = SaveBlock hash = de0d84d2e88f83916cfa955121938c952ceab657cd49852cac51f7f5c7c0091a 
INFO [2023-05-24 13:53:49.234]   received payload                         topic = SaveRoundsInfo 
INFO [2023-05-24 13:53:49.240]   received payload                         topic = FinalizedBlock 
 
```

After you finished testing, you can close the observer node and ws connector (can use CTRL+C) as well as the local
testnet, by executing:

```bash
bash ./local-testnet.sh stop
```