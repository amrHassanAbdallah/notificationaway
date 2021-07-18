# Notificationaway
> Sending notifications to your users supporting different messages languages and different sending interfaces such as push notifications, sms, webhook. 


## Getting started
### Manually
1. install [golang](https://golang.org/)
1. start a mongo instance locally
1. install kafka
1. run the app
   ```
    $ make generate
    $ make build
    $ ./bin/app 
   ```
### Using docker-compose
1. make sure that you have docker, docker-compose installed
1. run
   ```
   $ docker-compose up
   ```

### Check the API
Use this [file](https://github.com/amrHassanAbdallah/notificationaway/blob/master/api/api.yml) content and paste it inside this [viewer](https://editor.swagger.io/)


## Features
* [x] Add a message
* [x] Consume notification
* [x] Send notifications
* [x] Add Light size integration test
* [ ] Update a message
* [ ] Query messages
* [ ] Add unit tests
* [ ] Handle message with templates keys
* [ ] Add HTTP logger
* [ ] Add CI/CD to run the test/test file after every commit (almost done)
* [ ] Add a circuit breaker over the third party notifications sender providers


Maybe will add more depending on this [design document]()
  
