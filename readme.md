## MRR Calculator

Monthly Recurring Revenue (MRR) is a key metric for subscription-based businesses, measuring predictable monthly revenue from active subscriptions. Imagine running a subscription business; MRR reflects your financial health and growth. For example, an MRR of 678.12 EUR means your business is currently generating 678.12 EUR in revenue each month. This amount will remain constant if no customers change or cancel their subscriptions.
## Task Description

Your task is to create a CLI application that calculates the  current MRR from a given set of subscription data. The application should be able to parse the input data and compute the following:

- The net MRR value
- The MRR breakdown, including:
    - New business
    - Upgrades
    - Downgrades
    - Churn
    - Reactivations

For an additional challenge, the application should handle subscriptions purchased in different currencies and use an exchange rate API to compute the values in a unified currency.

For an extra challenge, the application should also plot a chart or table displaying the daily MRR movement over a specified period.

## Input Data

The input data for the MRR calculator will be a list of subscription records. Each record will include the following fields:

- `subscription_id`: A unique identifier for the subscription
- `customer_id`: A unique identifier for the customer
- `start_at`: The date and time when the subscription started (in RFC3339 format, e.g., YYYY-MM-DDTHH:MM:SSZ)
- `end_at`: The date and time when the subscription ended, if applicable (in RFC3339 format, e.g., YYYY-MM-DDTHH:MM:SSZ)
- `amount`: The subscription amount in decimal form (e.g., "100.00")
- `currency`: The currency code (e.g., USD, EUR)
- `interval`: The billing interval (e.g., month, year)
- `status`: The current status of the subscription (active, cancelled, amended)
- `cancelled_at`: The date and time when the subscription was cancelled, if applicable (in RFC3339 format, e.g., YYYY-MM-DDTHH:MM:SSZ)

### Statuses

- `active`: The subscription is currently active and generating revenue.
- `cancelled`: The subscription has been cancelled but might still be active until the end of the current billing period.
- `amended`: The subscription has been modified (e.g., upgraded or downgraded).

Here is an example of how the input data might look in JSON format:

```json
[
  {
    "subscription_id": "sub_001",
    "customer_id": "cust_001",
    "start_at": "2024-01-01T00:00:00Z",
    "end_at": null,
    "amount": "100.00",
    "currency": "USD",
    "interval": "month",
    "status": "active",
    "cancelled_at": null
  },
  {
    "subscription_id": "sub_002",
    "customer_id": "cust_002",
    "start_at": "2024-02-01T00:00:00Z",
    "end_at": "2024-06-01T00:00:00Z",
    "amount": "50.00",
    "currency": "EUR",
    "interval": "month",
    "status": "cancelled",
    "cancelled_at": "2024-05-01T00:00:00Z"
  },
  {
    "subscription_id": "sub_003",
    "customer_id": "cust_003",
    "start_at": "2024-03-01T00:00:00Z",
    "end_at": "2024-04-01T00:00:00Z",
    "amount": "75.00",
    "currency": "USD",
    "interval": "month",
    "status": "amended",
    "cancelled_at": null
  },
  {
    "subscription_id": "sub_004",
    "customer_id": "cust_002",
    "start_at": "2024-07-15T09:45:43Z",
    "end_at": null,
    "amount": "55.00",
    "currency": "EUR",
    "interval": "month",
    "status": "active",
    "cancelled_at": null
  },
  {
    "subscription_id": "sub_005",
    "customer_id": "cust_003",
    "start_at": "2024-04-01T00:00:00Z",
    "end_at": null,
    "amount": "60.00",
    "currency": "GBP",
    "interval": "year",
    "status": "active",
    "cancelled_at": null
  }
]
```

## Output Data

The output of the MRR calculator must include the following:

1. **Present MRR Net Value**: The current net MRR value.
2. **MRR Breakdown**: A detailed breakdown of the MRR, including:
- New business
- Upgrades
- Downgrades
- Churn
- Reactivations
3. **Daily MRR Values Table**: A table displaying the daily MRR values for the past n months.

The CLI application must accept the following arguments:

- `currency`: The currency in which to display the MRR values.
- `period`: The number of months for which to plot the daily MRR values in the table.

### Example Usage

```sh
./mrr-calc --currency USD --period 3 --input subscriptions.json
```

### Example Output

```
Present MRR Net Value: 2230.00 USD

Present MRR Breakdown:
- New Business: 2000.00 USD
- Upgrades: 500.00 USD
- Downgrades: -200.00 USD
- Churn: -170.00 USD
- Reactivations: 100.00 USD

Daily MRR:
|------------|------------------|
| Date       | MRR Value (USD)  |
|------------|------------------|
| 2024-03-01 | 4500.00          |
| 2024-03-02 | 4550.00          |
| 2024-03-03 | 4550.00          |
| ...        | ...              |
| 2024-06-01 | 5000.00          |
|------------|------------------|
```
