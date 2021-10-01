# LogTweets

![License](https://img.shields.io/github/license/Pac23/LogTweets)

Logs Replies,Mentions and all other Twitter [Account Activity](https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/overview) in Realtime using Webhooks,to a json file viewable and downloadable in the browser.

---

* [Install](#install)
* [Setup](#setup)
* [Usage](#usage)
* [Full Example](#full-example)

---

## Install 


With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
go get -u github.com/Pac23/LogTweets
```

## Setup 

Obtain Your Api auth Keys from Twitter by Creating a [Developer account](https://www.extly.com/docs/autotweetng_joocial/tutorials/how-to-auto-post-from-joomla-to-twitter/apply-for-a-twitter-developer-account/#apply-for-a-developer-account) and [Create a Application](https://docs.inboundnow.com/guide/create-twitter-application/) and set the enviroment name aswell.Set the Callback in the application feild on the twitter app dashboard to `yourhost:port/webhook/twitter`


Once that's done Set your Credentials in the Env File accordingly

```
CONSUMER_KEY=your consumer public key
CONSUMER_SECRET=your consumer secret
ACCESS_TOKEN_KEY=your acess token public
ACCESS_TOKEN_SECRET=your acess token secret
WEBHOOK_ENV=Your enviroment name
APP_URL=url of this bot
SERVER_PORT=Port you wish to run this on 
```

Once the Above Setup is Done Hit the Registerwebhook url as below to register and subscribe to the webhook.

```
host:port/registerWebhook
```

That's it,Once you have done the above steps reregistration ever again unless you change your twitter app enviromentname or the host of this bot will not be required.

## Usage 

To run the bot simply run the following command in your temrinal 

```
LogTweets
```

#### Check if the Bot is up 
To check if the server is running
```
host:port
```

#### Veiw/Download Twitter Data
The bot dumps all the recived json in a json file which can be viewd in the browser itself and can also be saved(refer your specific browser tools for that) using 

```
host:port/data
```

#### View System Logs
Logger Logs everything the server does.Can be accesed using

```
host:port/systemlogs
```

#### Check Crc 
Crc can be checked using 

```
host:port/webhook/twitter?crc_token=test 
```

