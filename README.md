# Shopify-Integrator  • [![ci](https://github.com/Keenan-Faure/Shopify-Integrator/actions/workflows/ci.yml/badge.svg)](https://github.com/Keenan-Faure/Shopify-Integrator/actions/workflows/ci.yml)

Pushes and pulls data from Shopify with some additional features.

1. [Description of project](#description-of-project)
2. [Setup Project](#setup-project)
   1. [What is already done](#what-is-already-done)
   2. [Configuring .env file](#configuring-env-file)
   3. [Setting up application passwords for emails](#setting-up-application-passwords-for-emails)
   4. [Configuring app and shopify setting](#configuring-app-and-shopify-setting)
   5. [Installing Golang](#installing-golang)
   6. [Installing Docker](#installing-docker)
3. [Configuring custom ngrok URL for orders](#configuring-custom-ngrok-url-for-orders)
4. [How to run the application](#how-to-run-the-project)
5. [How to run a fresh install of the application](#how-to-run-a-fresh-install-of-the-application)
6. [List of features](#list-of-features)
7. [What is next](#what-is-next)

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

The environment variable file contains variables that are used by the application. Of course, if one is not set correctly or removed, please do not expect the application to work correctly. Once the repository has been cloned, there exists a file `.example.env` that needs to be renamed to `.env` and thereafter  needs to be replaced with the respective values

The default `.example.env` file contains the default values that needs to be changed and can be found below:

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

**Please do not rename the variables or alter any of the ones not mentioned above**

### Setting up application passwords for emails

This application consists of a `preregistrater` and `registrater` endpoints. This means that the user that wishes to register as a user to the application will first need enter his email address and after a token has been sent to the respective email, and received, it can then be used in the register endpoint.

The email sending these tokens would need to be set up and the `application password` saved in the `.env` file. Small guide found below:

Please read this guide [here](https://support.google.com/mail/answer/185833?hl=en)

### Configuring app and shopify setting

The application consists of settings that needs to be configured, which are, of course, important in the functions of each feature. These can be done either over the API using a client like [Postman](https://www.postman.com) or using the front-end of the application.

Configuring on the front-end can be done by simply heading to the settings page after logging in and then adjusting the setting values.

Unsaved setting values will not be updated.

### Installing Golang

Installing golang is required to compile and build the server script. Follow the setup guide below to install Golang on your machine.

[How to install Golang](https://go.dev/doc/install)

### Installing Docker

It is required to have a valid docker installation. Please see a guide on how to install [Docker](https://www.docker.com/).

### Creating Ngrok account and authtoken

To use ngrok, an account needs to be created on their [website](https://dashboard.ngrok.com). You can either link you github account or create a new account.

After successfully creating an account, an `authtoken` needs to be retrieved and saved into the ngrok config file located in

```bash
${pwd}/ngrok/ngrok.yml
```

This ngrok `authToken` is used in the step below

**Please dont alter any of the other data in the `ngrok.yml` file when replacing the `authToken`**

## Configuring custom ngrok URL for orders

_This step assumes that you have a valid ngrok account and an `authtoken` populated in the ngrok config file._

This application, which runs on your local machine `localhost`, cannot be accessed over the internet. Hence, we use an
application called `ngrok` to do this for us.

The `install.sh` script installs ngrok for you, however, to create the webhook url that is used when creating orders
you can just retrieve it in the settings page of the application on the front-end or over the API.

Please see the small guide below on how to setup the webhook URL on Shopify.

_Note that this assumes that you have a shopify store with a valid ngrok `authToken`_

- [Guide on how to link Ngrok with your Shopify webhook](https://ngrok.com/docs/integrations/shopify/webhooks/)

**Note that your ngrok domain name can be found on the web interface of the ngrok docker container. Also the ngrok domain changes each time when using a free ngrok account plan**

## How to run the project

To run the application simply open the (cloned) local version of the application in your favorite command line interface, then run:

```bash
./scripts/install.sh
```

## How to run a fresh install of the application

```bash
./scripts/reset.sh
```

There exists an optional `rmi` parameter that is used to also remove local images. It is only recommended to use the `rmi` parameter if you wish to remove the application entirely.

```bash
./scripts/reset.sh rmi
```

## What is next

1. Adding additional features that support the use of a 3rd Party ERP system with features that
   1. Pull Product, Order & customer information
   2. Push Product, Order information back into the ERP
2. Creating a CD pipeline that uses Docker images and Github Actions to streamline the work, instead of manually creating an update script

## List of features

A friendly list of features currently supported and available:

- Adding, disabling, updating, and viewing products from the application.
- Adding, updating, and viewing order data. Note that this is only via the API.
- Viewing of customer data (generated from order data).
- Import/export function available for products.
- Filter & search features for products, orders and customers.
- Automatic customer creation from orders.
- RESTful API (JSON) endpoint and documentation.
- Workers that automatically fetch products from Shopify at a set interval to keep constant syncronization of product data.
- Adjustable settings via the API or on the app front-end.
- Neat and easy to follow front-end design using Popular HTML Framework React.
- Dockerized for portibility - _If it works on your PC, then we'll ship your PC_.
- Additional features, like restrictions & global warehousing, to make pushing and fetching shopify data easier and more controllable
- Automatic Ngrok integration with your Shopify Webhook.
- Queue feature
