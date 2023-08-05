# Description

Test assignment https://trustwallet.notion.site/Backend-Homework-Tx-Parser-abd431fca950427db75d73d90a0244a8

## How to install
```shell

# copy env file from default
cp .env-default .env

# build service
make build
```

## How to run
```
make run
```

## How to test

```shell
# health check
curl --location 'http://localhost:8818/health'

# get latest processed block
curl --location 'http://localhost:8818/blocks/current'

# subscribe address
curl --location --request POST 'http://localhost:8818/subscriptions?address=0xae2fc483527b8ef99eb5d9b44875f005ba1fae13'

# get transactions by the subscribed address
curl --location 'http://localhost:8818/transactions?address=0x1a2a5cbd389f7a80ea85f304101954d98c927e8a'
```

There are env variables `START_BLOCK` and `DEFAULT_SUBSCRIBED` 
which allow to set known subscribed addresses and start scanning from the block in the past 
