## About

This folder contains system test demo for the firehose ws connector, which is going to receive exported data from a
local testnet

### How it works

These scripts provide a demo to help you configure and test your ws connector/driver to receive data from a local
testnet. The scripts found in this `demo` folder:

- `local-testnet.sh` creates a local testnet with one shard (0) + metachain. Each shard has 2 nodes (1 validator + 1
  observer)
- `observer-outport.sh` creates an observer in a specified shard (shard/metachain). This node acts an outport driver for
  the ws connector. When started, it will sync from genesis and export data (blocks, validators info, etc.) to your ws
  driver connector.

Inside `cmd/connector` folder, there is a `main.go` file which represents the firehose ws connector (data receiver).

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
INFO [2023-05-26 16:31:32.183]   firehose: saving block                   nonce = 16 hash = 5e394dc1e73fe9a935ce5592544038f4745770f9e406fc78b34e8c671c76321a 
FIRE BLOCK_BEGIN 16
FIRE BLOCK_END 16 bfb301080f8e60e20479acdd81eb5d99c0496de385bbabf8dd756f0d7b1e79fc 1685107882 0ab2030a83030810...
INFO [2023-05-26 16:31:32.189]   firehose: saving block                   nonce = 17 hash = 1869f099683ae84224f755afcec7564cb9e7cbf078c66118e00720d38ce2d4a6 
FIRE BLOCK_BEGIN 17
FIRE BLOCK_END 17 5e394dc1e73fe9a935ce5592544038f4745770f9e406fc78b34e8c671c76321a 1685107887 0ab2030a83030811...
```

After you finished testing, you can close the observer node and ws connector (can use CTRL+C) as well as the local
testnet, by executing:

```bash
bash ./local-testnet.sh stop
```