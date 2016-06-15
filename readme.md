# Stock Market Clone

_now in golang!_

TOP TODO
- ticker HTTP GET API for bid / ask spread
- make a market maker
- make a trading bot

### TODO
- [ ] ledger service  
  - [ ] rename to account service
  - [x] listen for http (from orderbook)  
  - [ ] use django and mysql for auth, accounts and assets
  - [ ] handle accounts and authentication  
  - [x] track quantity and asset for each user (cash is an asset)  

- [x] orderbook service  
  - [x] listen for http (for submitting of orders)  
  - [x] redlock on orderbook ticker keys
  - [x] containerize, load balance
  - [x] use redis ordered sets for buy and sell priority queues
  - [x] priotiy queues for buy and sell orders  
  - [x] dequeue priority queues into trades  
  - [x] message to web service anonymized trades and orderbook bid ask
  - [x] message to ledger service trades
  - [x] message to ticker service anonymized trades 


- [ ] ticker service (chart and quote stream)  
  - [ ] rename to tickdata service
  - [x] redlock on second and minute ticker keys
  - [x] containerize and load balance
  - [?] need leader elect for ticker
  - [x] use redis to for accumulating trades into tick data
  - [x] on price data http, update minute high, low, open, close, vol information  
  - [x] every 60 seconds, persist 1minute period tick data to DB  
  - [x] API for charts

- [ ] web service  
  - [x] serve front end javascript  
  - [x] socket messages to interested clients
  - [ ] scale messages across multiple instances with web sockets?  

- [ ] web client (front end)  
  - [x] graph "CHART" data into a stock chart  
  - [x] listen to "TICKER" channels  
  - [x] append new data to chart  
  - [ ] display account info  
  - [x] place orders  

- [ ] bot service (automate trades)  
  - [ ] trading bots   
  - [ ] front end UI for creating bots  
  - [ ] rule based trading using JSON {$when: ..., $buy: ... }  

- [x] web load balancer 
  - [x] nginx
  - [x] proxy each service to single domain

Ideas:  
 - Ticker names: Kryptonite, Adamantium, Puppies  
 - rule based trading:  
  
    {  
      $when: {  
        $comparison: {  
          itemA: { $ticker: "STOCKA"},  
          itemB: 20.55,  
          type: "gte"  
        }  
      },  
      $do: {  
        $order: { kind: "LIMIT", intent: "BUY", bid: 20.60, shares: 100 }  
      },  
      $then: {  
        $when: {  
          $comparison: {  
            itemA: { $ticker: "STOCKA"},  
            itemB: 21.00,  
            type: "gte"  
          }  
        },   
        $do: {  
          $order: { kind: "LIMIT", intent: "SELL", ask: 20.90, shares 100 }  
        }  
      }  
    }  

