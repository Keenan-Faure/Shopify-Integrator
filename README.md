# Shopify-Integrator

Pushes and pulls data from Shopify with some additional features.

1. Description of project
2. Setup Project
    1. What is already done
    2. Configuring .env file
    3. Setting up application passwords for emails
    4. Configuring app and shopify setting
    5. Installing Docker
3. Configuring custom [ngrok](https://ngrok.com/) URL for orders
4. How to run the project
5. List of features
6. What is next (inc)

## Description of project

Web based application that pulls product data from Shopify at a predefined intervals using [goroutines](https://go.dev/tour/concurrency/1). Furthermore, it allows products to be added, both gradually and using a bulk feature.

The application has it's own RESTful API that allows communication between the application and one or more software components. An example is that it allows the syncing of orders from Shopify into the application through the use of a webhook and a predefined endpoint.

For a entire list of available features, please see [List of Features]().

## Setup Project

### What is already done

With the use of [Docker](https://www.docker.com/) containers, many of the prerequistes are installed once the shell install script runs. This includes, but not limited to:

- Golang (programming language the server it built on)
- Goose (used to manage database migrations)
- Node Package Manager (NPM)
- React (Popular HTML framework with which the front-end is built upon)

### Configuring .env file

The environment variable file contains variables that are used by the application. Of course, if one is removed, please do not expect the application to work correctly.

The default `.env` file contains the default values that needs to be changed.

#### Postgresql

```txt
DB_USER - Database username
DB_PSW - Database password
```

#### Shopify

```txt
SHOPIFY_STORE_NAME - Shopify Store name
SHOPIFY_API_KEY - API key generated on Shopify
SHOPIFY_API_PASSWORD - API Password generated on Shopify
SHOPIFY_API_VERSION - Shopify API version to use
```

### Setting up application passwords for emails

This application consists of a `preregistrater` and `registrater` endpoints. This means that the user that wishes to register as a user to the application will first need enter his email address and after a token has been sent to the respective email, and received, it can then be used in the register endpoint.

The email sending these tokens would need to be set up and the `application password` saved in the `.env` file. Small guide found below:

Please read this guide [here](https://support.google.com/mail/answer/185833?hl=en)

### Configuring app and shopify setting

The application consists of settings that needs to be configured, which are, of course, import in the functions of each feature. These are outlined below:

### Instlling Golang

Installing golang is required to compile and build the server script. Follow the setup guide below to install Golang on your machine.

[How to install Golang](https://go.dev/doc/install)

### Installing Docker

It is required to have a valid docker installation. Please see a guide on how to install [Docker](https://www.docker.com/).

## Configuring custom [ngrok](https://ngrok.com/) URL for orders

Coming soon.

## How to run the project

To run the application simply open the (cloned) local version of the application in your favorite command line interface, then run:

```bash
./scripts/run.sh
```

## List of features

A friendly list of features currently supported and available:

- Adding, removing, updating, and viewing products from the application.
- Adding, updating, and viewing order and customer data.
- Adding, updating, and viewing of customer data.
- Import/export function available for products.
- Filter search for products, orders and customers.
- Automatic customer created from orders.

## What is next (inc)

- Create product fetch (automatic fetch)
- Create product push
- Create product endpoints (adding new products locally)
  - update products locally (including variants)
- Additional Features

  - last fetch time (endpoint)
  - product search
  - product listing (pagination)
  - product fetch configuration (filter/limiter)
  - product push configuration (filter/limiter)

- pull orders from Shopify

  - use ngrok to expose an endpoint
  - webhook token per user & validation
  - automatic quantity deduction when orders are processed
  - automatic product push when update occurs in database (or per fetch_limit)
  - order pushed back to shopify when updated (option to exist - push to Shopify)

- settings:

  - Fetch products from Shopify [true/false]
  - Fetch products timer (minutes)
  - Enable Webhook [true/false]
  - Enable product variations [true/false]
  - Enable push to Shopufy
  - Enable creation of customers from orders

- Test endpoints
- Create Product export/import features

### Setting up Application password for emails

Please read this guide [here](https://support.google.com/mail/answer/185833?hl=en)
