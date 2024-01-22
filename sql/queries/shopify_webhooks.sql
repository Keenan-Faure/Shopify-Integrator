-- name: CreateShopifyWebhook :exec
INSERT INTO shopify_webhooks(
    id,
    shopify_webhook_id,
    webhook_url,
    topic,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: UpdateShopifyWebhook :exec
UPDATE shopify_webhooks
SET
    webhook_url = $1
    topic = $2
WHERE shopify_webhook_id = $3;

-- name: RemoveShopifyWebhook :exec
DELETE FROM shopify_webhooks
WHERE shopify_webhook_id = $1;
