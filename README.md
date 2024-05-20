# Shopify-Integrator  • [![ci](https://github.com/Keenan-Faure/Shopify-Integrator/actions/workflows/ci.yml/badge.svg)](https://github.com/Keenan-Faure/Shopify-Integrator/actions/workflows/ci.yml) ![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/Keenan-Faure/Shopify-Integrator/ci.yml) ![GitHub Release](https://img.shields.io/github/v/release/Keenan-Faure/Shopify-Integrator) ![Docker Image Version](https://img.shields.io/docker/v/keenansame/shopify-integrator) ![GitHub License](https://img.shields.io/github/license/Keenan-Faure/Shopify-Integrator)

Integration application for Shopify that makes use of it's powerful [API](https://shopify.dev/docs/api/admin-rest) to perform CRUD operations.

## 🚀 Quick start

After completing the [prerequisites](#prerequisites) below simply run in the (cloned) project directory:

```bash
./scripts/install.sh
```

This will download all necessary docker images for you and start on the `APP_PORT` on your localhost which by default is `3000`.

## Prerequisites

### Configuring .env file

The environment variable file contains variables that are used by the application. Of course, if removed, the application will not function correctly.

The default `.example.env` file contains the default values that needs to be changed. Ensure that a copy is made and renamed to `.env`. Thereafter, the values may be updated for personal use

#### PostgresSQL

```txt
DB_USER
DB_PSW
```

#### Shopify

```txt
SHOPIFY_STORE_NAME
SHOPIFY_API_KEY
SHOPIFY_API_PASSWORD
SHOPIFY_API_VERSION
```

### Installing Docker

It is required to have a valid docker installation to run the containerized application. Please see a guide on how to install [Docker](https://www.docker.com/) on your respective operating system.

### Setting up application passwords for emails

Please read this guide [here](https://support.google.com/mail/answer/185833?hl=en) on how to setup application passwords. These would need to be saved in the `.env` file

```txt
EMAIL=
EMAIL_PSW=
```

### Creating Ngrok account and authtoken

To properly use the ngrok container, an account needs to be created on their [website](https://dashboard.ngrok.com). You can either link your Github account or create a new account.

After successfully creating an account, an `authtoken` needs to be retrieved and saved into the ngrok config file located in

```bash
${pwd}/ngrok/ngrok.yml
```

Note that the `ngrok.example.yml` file needs to be **renamed** to `ngrok.yml` in order to be recongized by the application.

## Configuring custom ngrok URL for orders

_This step assumes that you have a valid ngrok account and an `authtoken` populated in the ngrok config file._

This application, which runs on your local machine `localhost`, cannot be accessed over the internet. Hence, we use an
application called `ngrok` to do this for us.

The `install.sh` script installs ngrok for you, however, to create the webhook url that is used when creating orders
you can just retrieve it in the settings page of the application on the front-end or over the API.

Please see the small guide below on how to setup the webhook URL on Shopify.

_note that this assumes that you have a shopify store with a valid ngrok authToken_

- [Guide on how to link Ngrok with your Shopify webhook](https://ngrok.com/docs/integrations/shopify/webhooks/)

**The ngrok domain changes each time the container is reloaded when using a free ngrok account plan**

## Creating OAuth Access Tokens in Google

The application (v2) supports OAuth2.0 with google which allows users, if they have a google account, to login using their existing account with google.

To use this feature `OAUTH_CLIENT_ID` AND `OAUTH_SECRET` needs to be configured. To obtain these access tokens, please follow the guide outlined [here](https://developers.google.com/identity/protocols/oauth2)

## ⚙️ Usage

### To install the application (uses docker)

```bash
./scripts/install.sh
```

### To reset the application

```bash
./scripts/reset.sh
```

Note that there exists an additional `rmi` argument that can be added to the `reset.sh` script.

### To stop, remove containers and remove any images downloaded

```bash
./scripts/reset.sh rmi
```

### To update the current database to the latest migration

Available flags:

- `production/development` - The database you want to update
- `up/reset/down` - The Goose command that you wish to do on the database

```bash
./scripts/update.sh production up
```

#### Examples

To update the production database inside docker to the latest version:

```bash
./scripts/update.sh production up
```

To update the local development database to the latest version:

```bash
./scripts/update.sh development up
```

To reset all migrations on the production database:

```bash
./scripts/update.sh production reset
```

Note that the above will remove all data on the current database

To migrate one down on the development database:

```bash
./scripts/update.sh development down
```

### To install new node modules on the docker app container

```bash
./scripts/app.update.sh ./app/package.json
./scripts/app.update.sh ./app/package-lock.json
```

This copies the respective files and places them inside the app directory on the docker container.

Essentially anything can be copied, but note that it will only place them inside the app directory and attempt to run `npm install`

Lastly, note that the volumes created will not be deleted upon running the `reset.sh` script. You may manually delete the volume should you feel the need to.

## 🤝 Contributing

### Clone the repo

```bash
git clone https://github.com/Keenan-Faure/Shopify-Integrator
cd Shopify-Integrator
```

Then complete prerequisites the project

### Run the project

```bash
./scripts/install.sh
```

### Run the tests

```bash
go test ./...
```

Ensure that the docker container is running.

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch. Feel free to let me know what you think can be done to make it awesome.

## Description

Web based application that pulls product data from Shopify at a predefined intervals using [goroutines](https://go.dev/tour/concurrency/1). Furthermore, it allows products to be added, both gradually and using a bulk feature.

The application has it's own RESTful API that allows communication between the application and one or more software components. An example is that it allows the syncing of orders from Shopify into the application through the use of a webhook and a predefined endpoint.

For a entire list of available features, please see [List of Features](https://google.com).

## Goals of this project

Shopify is one of the biggest E-Commerce platforms world wide. Hence, it only makes sense that many companies will be using their platform for online sales, and furthermore there is an existing need to make this easier for them. From Products and Orders to customer data, Shopify handles all of that in a single web application.

Now, this project uses the wonderful [API](https://shopify.dev/docs/api/admin-rest) of Shopify, more specifically the REST Admin API, to perform CRUD operations on the respective objects to automate the syncing of data to the web application. This seeks to remove any manualy work on the shopify interface, like changing a product's pricing etc.

## What is already done

With the use of [Docker](https://www.docker.com/) containers, many of the prerequistes are installed once the shell install script runs. This includes, but not limited to:

- Golang (programming language the server it built on)
- Goose (used to manage database migrations)
- Node Package Manager (NPM)
- React (Popular HTML framework with which the front-end is built upon)
- Ngrok

## What is next

1. Updating code base of application to make use of the [Gin Web Framework](https://gin-gonic.com/docs/introduction/)
2. Adding additional features that support the use of a 3rd Party ERP system with features that
   1. Pull Product, Order & customer information
   2. Push Product, Order information back into the ERP
   3. Also adding OAuth2.0 to login to the respective 3rd Party ERP for simplified authentication.
3. OAuth2.0 Login with Github Accounts

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
