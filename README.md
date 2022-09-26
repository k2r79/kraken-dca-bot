# Kraken DCA Bot

> A simple DCA investment bot to make sure you never miss the opportunity to lose money with cryptos 😁

...Table of content...

## How it works

The process is pretty simple : the bot invests in multiple crypto tokens at a given frequency.
This bot is made to be used with the Kraken exchange.

At the moment, tokens are bought at their market price using Kraken's ticker.
Tokens are bought in the order they are declared in the configuration.

This bot is based on [beldur/kraken-go-api-client](https://github.com/beldur/kraken-go-api-client), a big thanks to the 
developers of these libs.

### Error management

Investment rounds If ever the bot fails to invest in a token, the other tokens won't be invested in.
This is a choice that has been made to make sure token priority is respected in the case of insufficient funding for one
of the declared pairs.

## Configuration

...
Sample config in test/data

## Running the bot