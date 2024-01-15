# Shopify-Integrator  • [![ci](https://github.com/Keenan-Faure/Shopify-Integrator/actions/workflows/ci.yml/badge.svg)](https://github.com/Keenan-Faure/Shopify-Integrator/actions/workflows/ci.yml)

Integration application for Shopify.

1. [Description](#description)
2. [Why integrate with Shopify?](#why-integrate-with-shopify)
3. [What is already done](#what-is-already-done)
4. [Quick start ⚙️](#quick-start-⚙️)
   1. [Configuring .env file](#configuring-env-file)
   2. [Installing Docker](#installing-docker)
   3. [Setting up application passwords for emails](#setting-up-application-passwords-for-emails)
   4. [Configuring app and shopify setting](#configuring-app-and-shopify-setting)
   5. [Creating Ngrok account and authtoken](#creating-ngrok-account-and-authtoken)
5. [Configuring custom ngrok URL for orders](#configuring-custom-ngrok-url-for-orders)
6. [Usage](#usage)
7. [Contributing](#contributing)
7. [List of available features](#list-of-available-features)
8. [What is next](#what-is-next)

## Description

Web based application that pulls product data from Shopify at a predefined intervals using [goroutines](https://go.dev/tour/concurrency/1). Furthermore, it allows products to be added, both gradually and using a bulk feature.

The application has it's own RESTful API that allows communication between the application and one or more software components. An example is that it allows the syncing of orders from Shopify into the application through the use of a webhook and a predefined endpoint.

For a entire list of available features, please see [List of Features](https://google.com).

## Why Integrate with Shopify?

Shopify is one of the biggest E-Commerce platforms world wide. Hence, it only makes sense that many companies will be using their platform for online sales, and furthermore there is an existing need to make this easier for them. From Products and Orders to customer data, Shopify handles all of that in a single web application.

Now, this application simply uses the wonderful [API](https://shopify.dev/docs/api/admin-rest) of Shopify, more specifically the REST Admin API, to perform CRUD operations on the respective objects.

## What is already done

With the use of [Docker](https://www.docker.com/) containers, many of the prerequistes are installed once the shell install script runs. This includes, but not limited to:

- Golang (programming language the server it built on)
- Goose (used to manage database migrations)
- Node Package Manager (NPM)
- React (Popular HTML framework with which the front-end is built upon)
- Ngrok

## Quick start ⚙️

After completing the steps below simply run in the (cloned) project directory:

```bash
./scripts/install.sh
```

### Configuring .env file

The environment variable file contains variables that are used by the application. Of course, if removed, the application will not function correctly.

The default `.example.env` file contains the default values that needs to be changed. Ensure that a copy is made and renamed to `.env`. Thereafter, the values may be updated for personal use

#### PostgresSQL

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

### Installing Docker

It is required to have a valid docker installation to run the containerized application. Please see a guide on how to install [Docker](https://www.docker.com/) on your respective operating system.

### Setting up application passwords for emails

This application consists of a `preregistrater` and `registrater` endpoints. This means that the user that wishes to register as a user to the application will first need enter his email address and after a token has been sent to the respective email, and received, it can then be used in the register endpoint.

The email sending these tokens would need to be set up and the `application password` saved in the `.env` file. Small guide found below:

Please read this guide [here](https://support.google.com/mail/answer/185833?hl=en)

### Creating Ngrok account and authtoken

To use ngrok, an account needs to be created on their [website](https://dashboard.ngrok.com). You can either link you github account or create a new account.

After successfully creating an account, an `authtoken` needs to be retrieved and saved into the ngrok config file located in

```bash
${pwd}/ngrok/ngrok.yml
```

**Please dont alter any of the other data in the `ngrok.yml` file when replacing the `authToken`**

### Configuring app and shopify setting

The application consists of settings that needs to be configured, which are, of course, important in the functions of each feature. These can be done either over the API using a client like [Postman](https://www.postman.com) or using the front-end of the application.

## Configuring custom ngrok URL for orders

_This step assumes that you have a valid ngrok account and an `authtoken` populated in the ngrok config file._

This application, which runs on your local machine `localhost`, cannot be accessed over the internet. Hence, we use an
application called `ngrok` to do this for us.

The `install.sh` script installs ngrok for you, however, to create the webhook url that is used when creating orders
you can just retrieve it in the settings page of the application on the front-end or over the API.

Please see the small guide below on how to setup the webhook URL on Shopify.

_note that this assumes that you have a shopify store with a valid ngrok authToken_

- [Guide on how to link Ngrok with your Shopify webhook](https://ngrok.com/docs/integrations/shopify/webhooks/)

**Note that your ngrok domain name can be found on the logs of the docker container. Also the ngrok domain changes each time when using a free ngrok account plan**

## Usage

To install the application (uses docker):

```bash
./scripts/install.sh
```

To reset the application:

```bash
./scripts/reset.sh
```

Note that there exists an additional `rmi` argument that can be added to the `reset.sh` script.

To stop, remove containers and remove any images downloaded:

```bash
./scripts/reset.sh rmi
```

Lastly, note that the volumes created will not be deleted upon running the `reset.sh` script. You may manually delete the volume should you feel the need to.

## Contributing

Shopify-Integrator is currently an entry-level project, hence, it is not any accepting contributions at the moment. Many thanks in advance for your consideration 😄

## What is next

1. Adding additional features that support the use of a 3rd Party ERP system with features that
   1. Pull Product, Order & customer information
   2. Push Product, Order information back into the ERP
2. Creating a CD pipeline that uses Docker images and Github Actions to streamline the work, instead of manually creating an update script

## List of available features

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
