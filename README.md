# Shopify-Integrator

Pushes and pulls data from Shopify with some additional features

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

- logs
  - error logging in db
  - logs contain errors/warnings/info messages on when products/orders were processed.
  - logs section on front-end to be presented in the dashboard.
