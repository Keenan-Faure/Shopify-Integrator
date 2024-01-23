-- name: UpdateShopifyWebhook :exec
UPDATE shopify_webhooks
SET
    shopify_webhook_id = $1,
    webhook_url = $2,
    topic = $3
WHERE id = $4;

-- name: RemoveShopifyWebhook :exec
DELETE FROM shopify_webhooks
WHERE shopify_webhook_id = $1;

-- name: GetShopifyWebhooks :one
SELECT * FROM shopify_webhooks
LIMIT 1;
