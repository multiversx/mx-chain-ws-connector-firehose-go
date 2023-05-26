# mx-chain-ws-connector-firehose-go

## Introduction

In today's rapidly evolving blockchain ecosystem, the demand for real-time access to exported data from blockchain nodes
has significantly increased. To meet this growing need within our expanding system, we have developed the Outport
Driver, a powerful websocket connection system specifically designed for the MultiversX blockchain.

The Outport Driver acts as a crucial bridge between the MultiversX blockchain node and various services that rely on
receiving up-to-date data from the node. By establishing a websocket connection, the Outport Driver enables seamless
transmission of essential information, including blocks, validator details, account changes, and other processing
results.

This repository utilizes the [ws connector template](https://github.com/multiversx/mx-chain-ws-connector-template-go)
to serve as a robust data provider for the firehose ingestion processor. It receives incoming data and seamlessly
streams it to the standard output. The streamed data is appropriately prefixed with markers (such as `FIRE BLOCK BEGIN`)
, enabling easy identification and integration with downstream systems.

## How to use

### Use cases

To export data from a specific shard, you need to enable an observer node within that shard to export data. One can do
that
by enabling the **[HostDriverConfig].Enabled** flag to true
[from this config file](https://github.com/MultiversX/mx-chain-go/blob/master/cmd/node/config/external.toml)

In our case of a network consisting of three shards and a metachain, you would require four separate websocket receiver
binaries, each responsible for receiving data from its respective shard. This distributed approach ensures data
integrity and allows for shard-specific data processing.

However, the websocket connection in the Outport Driver is highly flexible and parameterizable. By customizing the
configuration, you can configure the observer nodes to function as clients and the ws firehose connector as a
server. This configuration reduces the number of driver receiver binaries to just one, acting as a server that receives
data from each observer, from each shard. To set your receiver as a server, you need set **mode**
from `cmd/connector/config/config.toml` to **"server"**.

Alternatively, you have the option to set up the observer node as a server capable of handling multiple clients. This
setup proves advantageous when multiple services within your ecosystem need to receive exported data.

### Illustration

1. In this scenario, each observer node acts as a server, sending data to one or more receivers. The receivers act as a
   client and receive data from each observer node and processes it accordingly.

```
        +----------------+   +-------------------+
        |   Observer 1   |-->|     Receiver 1    |
        |                |   +-------------------+
        |                |
        |                |   +-------------------+
        |                |-->|     Receiver 2    |
        +----------------+   +-------------------+
        
        +----------------+   +-------------------+
        |   Observer 2   |-->|     Receiver 3    |
        +----------------+   +-------------------+
        
        +----------------+   +-------------------+
        |   Observer 3   |-->|     Receiver 4    |
        +----------------+   +-------------------+                                                 
```

2. In this scenario, the observer nodes act as clients, while a single receiver collects data from all of them. The
   receiver acts as the server, receiving data from each observer node.

```
        +----------------+   +-------------------+
        |   Observer 1   |-->|                   |
        +----------------+   |                   |
                             |                   |
        +----------------+   |                   |
        |   Observer 2   |-->|      Receiver     |
        +----------------+   |                   |
                             |                   |
        +----------------+   |                   |
        |   Observer 3   |-->|                   |
        +----------------+   +-------------------+
```

### Example

If you want to test the firehose ws receiver within a local setup testnet environment, you can use the provided demo
within [demo directory](cmd/demo).
