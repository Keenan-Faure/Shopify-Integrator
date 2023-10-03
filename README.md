# Shopify-Integrator

Pushes and pulls data from Shopify with some additional features

-   Create product fetch (automatic fetch)
-   Create product push
-   Create product endpoints (adding new products locally)
    -   update products locally (including variants)
-   Additional Features

    -   last fetch time (endpoint)
    -   product search
    -   product listing (pagination)
    -   product fetch configuration (filter/limiter)
    -   product push configuration (filter/limiter)

-   pull orders from Shopify

    -   use ngrok to expose an endpoint
    -   webhook token per user & validation
    -   automatic quantity deduction when orders are processed
    -   automatic product push when update occurs in database (or per fetch_limit)
    -   order pushed back to shopify when updated (option to exist - push to Shopify)

-   settings:

    -   Fetch products from Shopify [true/false]
    -   Fetch products timer (minutes)
    -   Enable Webhook [true/false]
    -   Enable product variations [true/false]
    -   Enable push to Shopufy
    -   Enable creation of customers from orders

-   Test endpoints
-   Create Product export/import features

### Setting up Application password for emails

Please read this guide [here](https://support.google.com/mail/answer/185833?hl=en)
