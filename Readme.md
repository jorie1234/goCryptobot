# goCryptoBot

Tool to get deal with Binance 

## Usage

```
NAME:
   goCryptoBot.exe - query data from binance

USAGE:
   goCryptoBot.exe [global options] command [command options] [arguments...]

COMMANDS:
   account, a            show your account
   depotinfo, di         list your new or filled and not selled orders
   listorders, lo        list your orders
   listprices, lp        list the prices
   sellbot, sb           let the sellbot sell all your open orders
   sellorder, so         sell an existing order
   showsingleorder, sso  show a single order in raw API data
   watchbot, wb          the watchbot waits for new filled buy orders and send them to telegram
   help, h               Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --apikey value, --ak value     API Key (default: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX") [%API_KEY%]
   --apisecret value, --as value  API Secret (default: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX") [%API_SECRET%]
   --help, -h                     show help (default: false)
```

API Key, API Secret and Telegram ChatID can be provided via Environment or .env File

### .env File:
Create a file with the name ".env" and the content:
```
export API_SECRET=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
export API_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### Examples for API Key/Secret
Provide API Key and API Secret via parameter:
```
.\goCryptoBot.exe  --ak xxxxxxxxx -as xxxxxxxxx lo  -sy BTCEUR,LTCEUR,ETHEUR -st filled
```

Provide API Key and API Secret via .env or Environment:
```
.\goCryptoBot.exe listorders -symbol BTCEUR,LTCEUR,ETHEUR -status filled
```

### Use Cases

#### Show general Depot Infos (command depotinfo or di)

```
.\goCryptoBot.exe di -symbol ADAEUR,BTCEUR -last 2h
```
This lists:
* all open buy orders (status: `NEW`), not matter how old they are
* all filled buy orders (status: `FILLED`), that are not sold, not matter how old they are
* all buy orders that are already sold (sell order in status `FILLED`, where the buy or sell timestamp is no older than 2h (`--last,-l` parameter)

```
┏━━━━━━━━━━┳━━━━━━━━━━━━━━┳━━━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━┳━━━━━━┳━━━━━━━━┳━━━━━━━━━━━━━━━━━┳━━━━━━━━━━┳━━━━━━━━━━━┳━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━┓
┃  ORDERID ┃ PRICE        ┃ EXQNT      ┃ CUMQUOTEQNT ┃ PROFIT ┃ SIDE ┃ STATUS ┃ TIME            ┃ SELLSTAT ┃ SELLORDER ┃ SELLPRICE    ┃ SELLTIME        ┃
┣━━━━━━━━━━╋━━━━━━━━━━━━━━╋━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━╋━━━━━━╋━━━━━━━━╋━━━━━━━━━━━━━━━━━╋━━━━━━━━━━╋━━━━━━━━━━━╋━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━┫
┃ 9XXXXXX6 ┃ 170.00000000 ┃ 0.27500000 ┃ 46.75000000 ┃ 0.25   ┃ BUY  ┃ FILLED ┃ Mar 15 10:24:12 ┃ NEW      ┃  90925115 ┃ 178.50000000 ┃ Mar 15 10:28:33 ┃
┃ 9XXXXXX5 ┃ 160.00000000 ┃ 0.00000000 ┃ 0.00000000  ┃  -     ┃ BUY  ┃ NEW    ┃ Mar 17 07:38:04 ┃          ┃         0 ┃              ┃                 ┃
┃ 9XXXXXX5 ┃ 155.00000000 ┃ 0.00000000 ┃ 0.00000000  ┃  -     ┃ BUY  ┃ NEW    ┃ Mar 19 19:38:09 ┃          ┃         0 ┃              ┃                 ┃
┗━━━━━━━━━━┻━━━━━━━━━━━━━━┻━━━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━┻━━━━━━┻━━━━━━━━┻━━━━━━━━━━━━━━━━━┻━━━━━━━━━━┻━━━━━━━━━━━┻━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━┛
```

#### List all Orders

Lists all of your orders of the specified Symbol

```
 .\goCryptoBot.exe listorders --symbol DOTEUR
Loaded 1823 orders
2021/03/20 12:28:02 found 3 orders for symbol DOTEUR
┏━━━┳━━━━━━━━┳━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━━━━━━┳━━━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━━━━━━┳━━━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━━┳━━━━━━━┳━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃ # ┃ SYMBOL ┃  ORDERID ┃ PRICE       ┃ ORIGQUANTITY ┃ EXQNT      ┃ CUMQUOTEQNT ┃ AVGPRICE NOW ┃ VALUE NOW  ┃ PROFIT      ┃ STATUS   ┃ TYPE  ┃ SIDE ┃ TIME                          ┃
┣━━━╋━━━━━━━━╋━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━━━╋━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━━━╋━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━╋━━━━━━━╋━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ 1 ┃ DOTEUR ┃ 3XXXXXX3 ┃ 30.20000000 ┃ 1.50300000   ┃ 0.00000000 ┃ 0.00000000  ┃ 33.17021428  ┃ 0.00000000 ┃  0.00000000 ┃ NEW      ┃ LIMIT ┃ BUY  ┃ 2021-03-20 11:54:31 +0100 CET ┃
┃ 2 ┃ DOTEUR ┃ 3XXXXXX1 ┃ 30.10000000 ┃ 1.50300000   ┃ 0.00000000 ┃ 0.00000000  ┃ 33.17021428  ┃ 0.00000000 ┃  0.00000000 ┃ CANCELED ┃ LIMIT ┃ BUY  ┃ 2021-03-20 11:54:32 +0100 CET ┃
┃ 3 ┃ DOTEUR ┃ 3XXXXXX5 ┃ 30.00000000 ┃ 1.51100000   ┃ 0.00000000 ┃ 0.00000000  ┃ 33.17021428  ┃ 0.00000000 ┃  0.00000000 ┃ NEW      ┃ LIMIT ┃ BUY  ┃ 2021-03-20 12:13:26 +0100 CET ┃
┣━━━╋━━━━━━━━╋━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━━━╋━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━━━╋━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━╋━━━━━━━╋━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃   ┃        ┃          ┃             ┃              ┃            ┃             ┃              ┃            ┃  0.00000000 ┃          ┃       ┃      ┃                               ┃
┗━━━┻━━━━━━━━┻━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━━━━━━━┻━━━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━━━━━━━┻━━━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━━━┻━━━━━━━┻━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
```


#### Show single order

Shows all details of a single order

```
 .\goCryptoBot.exe sso --help
NAME:
   goCryptoBot.exe showsingleorder - show a single order in raw API data

USAGE:
   goCryptoBot.exe showsingleorder [command options] [arguments...]

OPTIONS:
   --order value, -o value     order ID to fetch (default: "all")
   --symbol value, --sy value  order symbol like  (default: "BTCEUR")
   --help, -h                  show help (default: false)
```

###### Example

```   
 .\goCryptoBot.exe sso --symbol DOTEUR -o 3XXXXXX3
Loaded 1824 orders
(*binance.Order)(0xc0006a0000)({
 Symbol: (string) (len=6) "DOTEUR",
 OrderID: (int64) 33872113,
 ClientOrderID: (string) (len=32) "x-J6MCRYME5F4EXXXXXXXXXXXXF436",
 Price: (string) (len=11) "30.20000000",
 OrigQuantity: (string) (len=10) "1.50300000",
 ExecutedQuantity: (string) (len=10) "0.00000000",
 CummulativeQuoteQuantity: (string) (len=10) "0.00000000",
 Status: (binance.OrderStatusType) (len=3) "NEW",
 TimeInForce: (binance.TimeInForceType) (len=3) "GTC",
 Type: (binance.OrderType) (len=5) "LIMIT",
 Side: (binance.SideType) (len=3) "BUY",
 StopPrice: (string) (len=10) "0.00000000",
 IcebergQuantity: (string) (len=10) "0.00000000",
 Time: (int64) 1616237671915,
 UpdateTime: (int64) 1616237671915,
 IsWorking: (bool) true,
 IsIsolated: (bool) false
})
```
#### Start Watchbot

The watchbot repeatedly scans for new FILLED Buy Order and informs you via telegram. You have to provide the Telegram ChatID (via parameter or .env File)

```
PS C:\Users\live\Sync\code\goCryptoBot> .\goCryptoBot.exe watchbot --help
NAME:
   goCryptoBot.exe watchbot - the watchbot waits for new filled buy orders and send them to telegram

USAGE:
   goCryptoBot.exe watchbot [command options] [arguments...]

OPTIONS:
   --symbol value, --sy value            order symbol like  (default: "BTCEUR")
   --last value, -l value                time period e.g. 24h for the last 24 hours
   --telegrambotkey value, --tbk value   the bot key for telegram (default: "1601384507:AAEn4V3bmL06skYXUWrnSqdCLPiTBslJyoc") [%TELEGRAMBOTKEY%]
   --repeat value, -r value              run forever and check orders every <repeat> duration, eg. 30s or 1m or 1h, should somehow match the last duration
   --telegramchatid value, --tcid value  the chat id for telegram (default: "XXXXXX") [%TELEGRAMCHATID%]
   --help, -h                            show help (default: false)
```



###### Example

```
 .\goCryptoBot.exe watchbot -sy ADAEUR,DOGEEUR,ETHEUR,BTCEUR,LTCEUR,DATAUSDT,FIOUSDT,COTIUSDT,BTCSTUSDT,DOTEUR -l 1h  --repeat 1m
Loaded 1XXX orders
2021/03/20 12:37:46 found X orders for symbol ADAEUR
2021/03/20 12:37:46 found X orders for symbol DOGEEUR
2021/03/20 12:37:46 found X orders for symbol ETHEUR
2021/03/20 12:37:47 found X orders for symbol BTCEUR
2021/03/20 12:37:47 found X orders for symbol LTCEUR
2021/03/20 12:37:47 found X orders for symbol DATAUSDT
2021/03/20 12:37:48 found X orders for symbol FIOUSDT
2021/03/20 12:37:48 found X orders for symbol COTIUSDT
2021/03/20 12:37:48 found X orders for symbol BTCSTUSDT
2021/03/20 12:37:48 found X orders for symbol DOTEUR
```

This Command runs until it is stopped with CTRL-C

