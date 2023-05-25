# mx-chain-ws-connector-template-go

## Introduction

In today's rapidly evolving blockchain ecosystem, the demand for real-time access to exported data from blockchain nodes
has significantly increased. To meet this growing need within our expanding system, we have developed the Outport
Driver, a powerful websocket connection system specifically designed for the MultiversX blockchain.

The Outport Driver acts as a crucial bridge between the MultiversX blockchain node and various services that rely on
receiving up-to-date data from the node. By establishing a websocket connection, the Outport Driver enables seamless
transmission of essential information, including blocks, validator details, account changes, and other processing
results.

One of the key features of the MultiversX blockchain node is its flexible configuration file. By simply setting a flag
to "true" within this file, the node can initiate the delivery of data via a websocket connection to any connected
websocket client. This capability provides a reliable and efficient method for receiving live updates from the
MultiversX blockchain.

To simplify the process of setting up your own data receiver from the MultiversX blockchain node, we have created this
repository as a comprehensive template. It serves as a valuable resource, guiding you through the necessary steps to
establish a successful connection with the MultiversX blockchain node and receive the exported data effortlessly.

Whether you are building a decentralized application, monitoring the blockchain network, or conducting in-depth data
analysis, this template repository equips you with the foundation to swiftly integrate and leverage the exported data
from the MultiversX blockchain node.

Begin your journey towards seamless data integration and take advantage of the Outport Driver's robust websocket
connector by exploring the contents of this repository. We look forward to witnessing the innovative applications and
services that will emerge as a result of this collaborative ecosystem.

_Note: The Outport Driver is an open-source project, and we encourage active contributions and feedback from the
community to further enhance its capabilities and compatibility with different blockchain networks._

## How to use

### Use cases

The Outport Driver operates on a robust websocket connection architecture, supporting both server and client roles. The
system is designed to seamlessly export data from the MultiversX blockchain node, which adopts a sharded architecture
where each shard represents a separate running chain, interconnected by the meta chain.

To export data from a specific shard, you need to enable an observer node within that shard to export data. One can do
that
by enabling the **[HostDriverConfig].Enabled** flag to true
[from this config file](https://github.com/MultiversX/mx-chain-go/blob/master/cmd/node/config/external.toml)

In our case of a network consisting of three shards and a metachain, you would require four separate websocket receiver
binaries, each responsible for receiving data from its respective shard. This distributed approach ensures data
integrity and allows for shard-specific data processing.

However, the websocket connection in the Outport Driver is highly flexible and parameterizable. By customizing the
configuration, you can configure the observer nodes to function as clients and the driver (this template binary) as a
server. This configuration reduces the number of driver receiver binaries to just one, acting as a server that receives
data from each observer, from each shard. To set your receiver as a server, you need set **mode**
from `cmd/connector/config/config.toml` to **"server"**.

Alternatively, you have the option to set up the observer node as a server capable of handling multiple clients. This
setup proves advantageous when multiple services within your ecosystem need to receive exported data. For instance, a
single node could export data to various services, such as an elastic indexer, a monitoring tool, a notifier, or other
implementations.

By tailoring the configuration of the observer nodes and the driver (this template binary) to suit your specific
requirements, you can effectively streamline and centralize the reception of data from the MultiversX blockchain node.
This flexibility empowers you to design a data distribution strategy that aligns with your ecosystem's needs while
ensuring efficient and reliable data transmission.

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
        +----------------+   |                   |
                             |                   |
        +----------------+   |                   |
        |    Receiver    |-->|                   |
        +----------------+   +-------------------+
```

### Start building

Once you have an observer node running and configured to export data (either server or client) you can test your
receiver.

1. Start by navigating to the `cmd/connector` directory in your terminal or file explorer. This is where you'll find
   the `main.go` file.

2. Build the binary by running the `go build` command. This will compile the code and create the `connector` binary
   executable.

3. After building, you can start the `connector` binary by executing the generated executable file. This will launch the
   websocket data receiver.

4. In the `cmd/connector/config/config.toml` file, you can configure the settings for your receiver. Adjust the
   necessary parameters such as the url, port, or any other relevant options.

5. The provided code is a dummy implementation that logs events upon receiving data from the blockchain node. To
   customize the actions taken when receiving exported data, you need to define your own logic in the `process` package.

6. In the `process` package, there is a `logDataProcessor.go`. This processor is capable of handling data from the
   websocket outport driver and logging events.

7. One can see that `operationHandlers` defines a map of actions and their handler functions for the received payload.
   All you need to do is to replace the dummy code (which only logs events) with your specific use case on how to handle
   received data  (`saveBlock`, `revertIndexedBlock`, etc.). Modify the logic according to your requirements,
   such as saving to a database, triggering other processes, or performing any necessary data transformations.

8. Once you have customized the code to handle the received data as desired, you can build and run the `connector`
   binary again to see your changes in action.

_If you want to test the ws receiver within a local setup testnet environment, you can use the provided demo
within `cmd/demo` directory._