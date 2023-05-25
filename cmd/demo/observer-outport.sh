#!/usr/bin/env bash

# Possible values: "server"/"client"
OBSERVER_MODE="server"

CURRENT_DIR=$(pwd)
SANDBOX_PATH=$CURRENT_DIR/testnet/testnet-local/sandbox
KEY_GENERATOR_PATH=$CURRENT_DIR/testnet/mx-chain-go/cmd/keygenerator
EXTERNAL_CONFIG_DIR=$CURRENT_DIR/exporter/config/external.toml

createObserverKey(){
  pushd $CURRENT_DIR

  cd testnet/mx-chain-go/cmd/keygenerator
  go build
  ./keygenerator

  popd
}

resetWorkingDir(){
  rm -rf exporter
  mkdir "exporter/"
  cd "exporter/"
  mkdir "config/"
}

setupObserver(){
  OBSERVER_PATH=$(pwd)
  cp $SANDBOX_PATH/node/node $OBSERVER_PATH
  cp -R $SANDBOX_PATH/node/config $OBSERVER_PATH
  mv config/config_observer.toml config/config.toml
  mv $KEY_GENERATOR_PATH/validatorKey.pem config/

  sed -i "s@DestinationShardAsObserver =.*@DestinationShardAsObserver = \"$1\"@" $OBSERVER_PATH/config/prefs.toml
  sed -i '/HostDriverConfig\]/!b;n;n;c\    Enabled = true' "$EXTERNAL_CONFIG_DIR"
  sed -i "s@Mode =.*@Mode = \"$OBSERVER_MODE\"@" "$EXTERNAL_CONFIG_DIR"
  sed -i 's/MarshallerType =.*/MarshallerType = "json"/' "$EXTERNAL_CONFIG_DIR"
  sed -i 's/BlockingAckOnError =.*/BlockingAckOnError = false/' "$EXTERNAL_CONFIG_DIR"

  ./node --log-level *:INFO --rest-api-interface :10044
}


setup(){
  createObserverKey
  resetWorkingDir
  setupObserver $1
}

echoOptions(){
  echo "This script will start an observer node in the provided shard (metachain/shard).
  The only acceptable parameters, for the current test configuration, are shard (will use shard 0) and metachain."
}

main(){
    if [ $# -eq 1 ]; then
      case "$1" in
        metachain)
          setup metachain;;
        shard)
          setup 0;;
        *)
          echoOptions;;
        esac
    else
      echoOptions
    fi
}

main "$@"
