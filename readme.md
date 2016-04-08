# Stock Market Clone

_now in golang!_

### TODO
- [x] ledger service  
  - [x] listen for http (from orderbook)
  - [x] track quantity and asset for each user (cash is an asset)  

- [x] orderbook service  
  - [ ] listen for http (for submitting of orders)  
  - [x] priotiy queues for buy and sell orders
  - [x] dequeue priority queues into trades
  - [x] message to ledger service  
  - [x] message to ticker service  


- [ ] ticker service (chart and quote stream)  
  - [ ] on price data, publish quote to "QUOTE" channel; rate limit 1 / second  
  - [ ] on price data, add to cache of last 60 seconds of trades  
  - [ ] every 60 seconds, calc minute data, persist to DB, publish to "CHART" channel  

- [ ] web service (front end data)  
  - [ ] serve front end javascript  
  - [ ] handle accounts and oath  

- [ ] web client (front end ui)  
  - [ ] enable listening to "QUOTE" and "CHART" channels  
  - [ ] graph "CHART" data into a stock chart  
  - [ ] display account info  
  - [ ] place orders  

- [ ] bot service (automate trades)  
  - [ ] write trading bots   
  - [ ] make a front end UI for creating bots

- [ ] kafka (message queue between services)  

