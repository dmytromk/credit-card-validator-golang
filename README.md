# Golang Credit Card Validator API with gRPC architecture

## General Information

It's a gRPC validation API for credit cards. You can check out Python version with REST architecture [here](https://github.com/dmytromk/credit-card-validator). It checks:
1. Is card expired.
2. Does card number satisfy [Luhn check](https://en.wikipedia.org/wiki/Luhn_algorithm).
3. Card issuer and does card number length corresponds to the issuer.

## Local Deployment

### Docker startup

Inside of the repository directory run the following commands:
```shell
docker build --tag card-validator .
docker run -p 5001:5001 card-validator
```

## API Usage

API has only 1 endpoint: '/validate' POST. 
I recommend using Postman. It supports gRPC calls, you only need to import **validation.proto** file.

**Valid example:**

<img width="800" src="https://github.com/dmytromk/credit-card-validator-golang/assets/96624185/ced1021b-99c0-46e0-8b08-543530e67335">

_________________

**Invalid example:**

<img width="800" src="https://github.com/dmytromk/credit-card-validator-golang/assets/96624185/3a0bcfd1-de40-43d2-b720-fb6e797808bd">
