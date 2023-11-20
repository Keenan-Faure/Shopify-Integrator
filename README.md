# Shopify-Integrator  • ![code coverage badge][def]

Pushes and pulls data from Shopify with some additional features.

1. Description of project
2. Setup Project
   1. What is already done
   2. Configuring .env file
   3. Setting up application passwords for emails
   4. Configuring app and shopify setting
   5. Installing Docker
   6. Creating Ngrok account and authtoken
3. Configuring custom [ngrok](https://ngrok.com/) URL for orders
4. How to run the project
5. List of features
6. What is next (inc)

## Description of project

Web based application that pulls product data from Shopify at a predefined intervals using [goroutines](https://go.dev/tour/concurrency/1). Furthermore, it allows products to be added, both gradually and using a bulk feature.

The application has it's own RESTful API that allows communication between the application and one or more software components. An example is that it allows the syncing of orders from Shopify into the application through the use of a webhook and a predefined endpoint.

For a entire list of available features, please see [List of Features](https://google.com).

## Setup Project

### What is already done

With the use of [Docker](https://www.docker.com/) containers, many of the prerequistes are installed once the shell install script runs. This includes, but not limited to:

- Golang (programming language the server it built on)
- Goose (used to manage database migrations)
- Node Package Manager (NPM)
- React (Popular HTML framework with which the front-end is built upon)
- Ngrok

### Configuring .env file

The environment variable file contains variables that are used by the application. Of course, if one is removed, please do not expect the application to work correctly.

The default `.env` file contains the default values that needs to be changed and can be found below:

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

### Installing Golang

Installing golang is required to compile and build the server script. Follow the setup guide below to install Golang on your machine.

[How to install Golang](https://go.dev/doc/install)

### Installing Docker

It is required to have a valid docker installation. Please see a guide on how to install [Docker](https://www.docker.com/).

### Creating Ngrok account and authtoken

To use ngrok, an account needs to be created on their [website](https://dashboard.ngrok.com). You can either link you github account or create a new one.

After successfully creating an account, an `authtoken` needs to be retrieved and saved into the ngrok config file located in

```bash
${pwd}/ngrok/ngrok.yml
```

**Please dont alter any of the data in the `ngrok.yml` file**

## Configuring custom ngrok URL for orders

_This step assumes that you have a valid ngrok account and an `authtoken` populated in the ngrok config file._

This application, which runs on your local machine `localhost`, cannot be accessed over the internet. Hence, we use an
application called `ngrok` to do this for us.

The `install.sh` script installs ngrok for you, however, to create the webhook url that is used when creating orders
you can just retrieve it in the dashboard of the application.

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
- Automatic customer creation from orders.
- RESTful API endpoint and documentation.
- Workers that automatically fetch products from Shopify at a set interval to keep constant syncronization of product data.
- Adjustable settings via the API or on the app.

[def]: https://github.com/keenan-faure/learn-cicd-starter/actions/workflows/ci.yml/badge.svg
