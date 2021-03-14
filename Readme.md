# goCryptoBot

Tool to get Infos from Binance 

## Usage

```
NAME:
   goCryptoBot.exe - query data from binance

USAGE:
   goCryptoBot.exe [global options] command [command options] [arguments...]
   
COMMANDS:
account, a      show your account
listorders, lo  list your orders
listprices, lp  list the prices
help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
--apikey value, --ak value     API Key (default: "xxxx") [%API_KEY%]
--apisecret value, --as value  API Secret (default: "xxxx") [%API_SECRET%]
--help, -h                     show help (default: false)
```

API Key and API Secret can be provided vie Environment or .env File

### .env File:
Create a file with the name ".env" and the content:
```
export API_SECRET=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
export API_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### Examples
Provide API Key and API Secret via parameter:
```
.\goCryptoBot.exe  --ak xxxxxxxxx -as xxxxxxxxx lo  -sy BTCEUR,LTCEUR,ETHEUR -st filled
```

Provide API Key and API Secret via .env or Environment:
```
.\goCryptoBot.exe listorders -symbol BTCEUR,LTCEUR,ETHEUR -status filled
```

