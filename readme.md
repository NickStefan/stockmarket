# Stock Market Clone

_now in golang!_

### TODO
- [x] ledger service  
  - [x] listen for http (from orderbook)  
  - [x] track quantity and asset for each user (cash is an asset)  

- [x] orderbook service  
  - [x] listen for http (for submitting of orders)  
  - [ ] use redis ordered sets for buy and sell priority queues
  - [x] priotiy queues for buy and sell orders  
  - [x] dequeue priority queues into trades  
  - [x] message to ledger service  
  - [x] message to ticker service  


- [x] ticker service (chart and quote stream)  
  - [x] use redis to for accumulating trades into tick data
  - [x] on price data http, update minute high, low, open, close, vol information  
  - [x] every 60 seconds, persist 1minute period tick data to DB  
  - [x] every 1 second, publish 1second period tick data to websockets
  - [?] deal with 0d out values in ticker values
  - [x] API for charts

- [ ] web service  
  - [x] serve front end javascript  
  - [x] socket messages to interested clients
  - [ ] handle accounts and authentication  

- [ ] web client (front end)  
  - [x] graph "CHART" data into a stock chart  
  - [x] listen to "TICKER" channels  
  - [x] append new data to chart  
  - [ ] display account info  
  - [ ] place orders  

- [ ] bot service (automate trades)  
  - [ ] trading bots   
  - [ ] front end UI for creating bots  
  - [ ] rule based trading using JSON {$when: ..., $buy: ... }  

- [ ] discover service
  - [ ] proxy each service
  - [ ] load balance between multiple instances of each service
  - [ ] health check API for each service
  - [ ] api for adding and removing the proxied urls

- [ ] nginx reverse proxy all services to one domain  
- [ ] message queue between services  
  
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

