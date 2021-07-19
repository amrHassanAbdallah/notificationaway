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


### How to use
1. Create a message through the API
2. Publish a message into the `notifications` topic with the expected payload (type, to, language), make sure that the type and language match a message you already created in order to complete the whole cycle
3. if the type is webhook, you will get a message over that webhook, otherwise are dummy adapter the consumer will consume the message and ignore it.


**Or**

You can run the test/test.go that would do all of that to you.


## Features
* [x] Add a message
* [x] Consume notification
* [x] Send notifications (webhook type)
* [x] Add Light size integration test
* [ ] Update a message
* [ ] Query messages
* [ ] Add unit tests
* [ ] Handle message with templates keys
* [ ] Add HTTP logger
* [ ] Add CI/CD to run the test/test file after every commit (almost done)
* [ ] Add a circuit breaker over the third party notifications sender providers


## Architecture
![Blank diagram](https://user-images.githubusercontent.com/15635708/126164351-e290f676-6886-4ffc-9a9c-3a347b0c62f7.png)
More details in [design document](https://drive.google.com/file/d/11mnoDKF4rNicQAYUzmcNaiIbpKJ5doZo/view?usp=sharing)
  
