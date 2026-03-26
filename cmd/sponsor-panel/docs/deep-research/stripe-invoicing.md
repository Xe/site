# Stripe invoicing for EU sponsors in Go

**Stripe's Invoicing API combined with the stripe-go SDK provides a complete solution for sending invoices to EU sponsors with SEPA and bank transfer payment options, automatic VAT handling, and webhook-driven status tracking.** The approach centers on creating invoices via the API (not Checkout), configuring `payment_settings` for EU payment methods, enabling Stripe Tax for automatic VAT calculation, and using the Hosted Invoice Page as the payment UI. This guide covers every component with production-ready Go code using the current stripe-go v84 SDK and its `stripe.NewClient()` pattern.

The key architectural decision: use the **Invoice API directly** rather than Stripe Checkout for invoice-based billing. Checkout has no "invoice mode" — it generates post-payment receipts, not payable invoices. The Invoice API creates proper invoices with payment terms, sends them via email, and provides a Hosted Invoice Page where sponsors choose their payment method.

---

## Creating customers and invoices with the stripe-go SDK

The current stripe-go SDK is **v84**, using the `stripe.NewClient()` pattern introduced in v82.1. This replaces the old global `stripe.Key` approach with a thread-safe client. All service methods now take `context.Context` as the first argument and live under `sc.V1*` namespaces.

```go
import "github.com/stripe/stripe-go/v84"

sc := stripe.NewClient(os.Getenv("STRIPE_SECRET_KEY"))
```

**Creating a sponsor as a Stripe Customer** requires a name, email (mandatory for sending invoices), and for EU tax purposes, a full billing address:

```go
custParams := &stripe.CustomerCreateParams{
    Name:  stripe.String("Acme GmbH"),
    Email: stripe.String("billing@acme.de"),
    Address: &stripe.AddressParams{
        Line1:      stripe.String("Friedrichstraße 123"),
        City:       stripe.String("Berlin"),
        PostalCode: stripe.String("10117"),
        Country:    stripe.String("DE"),
    },
    Metadata: map[string]string{
        "sponsor_id":   "sp_12345",
        "sponsor_tier": "gold",
    },
}
custParams.SetIdempotencyKey("create-customer-sp_12345")
cust, err := sc.V1Customers.Create(ctx, custParams)
```

**Creating a draft invoice** with payment terms uses `DaysUntilDue` (the API field `days_until_due`). This only works when `CollectionMethod` is set to `"send_invoice"` — the mode designed for invoice-based payment flows where customers pay after receiving the invoice:

```go
invParams := &stripe.InvoiceCreateParams{
    Customer:         stripe.String(cust.ID),
    CollectionMethod: stripe.String(string(stripe.InvoiceCollectionMethodSendInvoice)),
    DaysUntilDue:     stripe.Int64(30),  // net-30; use 60 for net-60
    Currency:         stripe.String(string(stripe.CurrencyEUR)),
    AutoAdvance:      stripe.Bool(false), // manual control over finalization
    Description:      stripe.String("Gold Sponsorship - Q1 2025"),
    Metadata: map[string]string{"sponsor_id": "sp_12345"},
}
invParams.SetIdempotencyKey("invoice-sp_12345-Q1-2025")
inv, err := sc.V1Invoices.Create(ctx, invParams)
```

**Adding line items** to the draft invoice can use either a pre-existing Price ID or inline amounts. For sponsor panels with variable amounts, inline is simpler:

```go
itemParams := &stripe.InvoiceItemCreateParams{
    Customer:    stripe.String(cust.ID),
    Invoice:     stripe.String(inv.ID),
    Amount:      stripe.Int64(500000), // €5,000.00 in cents
    Currency:    stripe.String(string(stripe.CurrencyEUR)),
    Description: stripe.String("Gold Sponsorship Package - Q1 2025"),
}
_, err = sc.V1InvoiceItems.Create(ctx, itemParams)
```

**Finalizing and sending** transitions the invoice from `draft` → `open`, generates a PDF, creates a PaymentIntent, and emails it to the customer. The `SendInvoice` method finalizes automatically if the invoice is still in draft:

```go
sendParams := &stripe.InvoiceSendInvoiceParams{}
inv, err = sc.V1Invoices.SendInvoice(ctx, inv.ID, sendParams)
// inv.HostedInvoiceURL → payment page for the sponsor
// inv.InvoicePDF → downloadable PDF link
```

The invoice lifecycle follows five statuses: **`draft` → `open` → `paid`** is the happy path. From `open`, invoices can also move to `void` (cancelled) or `uncollectible` (written off). Both `void` and `paid` are terminal. Draft invoices are fully editable; once finalized, monetary values become immutable.

---

## Enabling SEPA Direct Debit and bank transfers for EU invoices

EU sponsors typically pay via bank transfer or SEPA Direct Debit. Both are configured through the invoice's `PaymentSettings` field. A critical distinction: **bank transfers use the `customer_balance` payment method type**, not `bank_transfer`. Stripe's bank transfer system works through a Customer Cash Balance — Stripe assigns each customer a unique virtual bank account number (VBAN), and incoming transfers are auto-reconciled against open invoices.

Here is the complete configuration for an invoice accepting card, SEPA, and EU bank transfer:

```go
invParams := &stripe.InvoiceCreateParams{
    Customer:         stripe.String(cust.ID),
    CollectionMethod: stripe.String(string(stripe.InvoiceCollectionMethodSendInvoice)),
    DaysUntilDue:     stripe.Int64(30),
    Currency:         stripe.String(string(stripe.CurrencyEUR)),
    PaymentSettings: &stripe.InvoicePaymentSettingsParams{
        PaymentMethodTypes: stripe.StringSlice([]string{
            "card",
            "sepa_debit",
            "customer_balance",
        }),
        PaymentMethodOptions: &stripe.InvoicePaymentSettingsPaymentMethodOptionsParams{
            CustomerBalance: &stripe.InvoicePaymentSettingsPaymentMethodOptionsCustomerBalanceParams{
                FundingType: stripe.String("bank_transfer"),
                BankTransfer: &stripe.InvoicePaymentSettingsPaymentMethodOptionsCustomerBalanceBankTransferParams{
                    Type: stripe.String("eu_bank_transfer"),
                    EUBankTransfer: &stripe.InvoicePaymentSettingsPaymentMethodOptionsCustomerBalanceBankTransferEUBankTransferParams{
                        Country: stripe.String("NL"), // BE, DE, ES, FR, IE, or NL
                    },
                },
            },
        },
    },
}
inv, err := sc.V1Invoices.Create(ctx, invParams)
```

When the sponsor opens the Hosted Invoice Page, they see all three payment options. **SEPA Direct Debit** requires EUR currency and collects a mandate automatically through the hosted page. Funds arrive in approximately **5 business days**, and customers can dispute within 13 months. **Bank transfer** shows the VBAN and transfer instructions on the hosted page and PDF — the sponsor initiates a standard bank transfer, and Stripe reconciles it automatically by matching the reference code, amount, or oldest open invoice.

The `EUBankTransfer.Country` parameter determines where the virtual bank account is domiciled. Choose the country closest to your primary customer base to minimize transfer fees and settlement time. Supported EU countries: **BE, DE, ES, FR, IE, NL**.

---

## Why Checkout Sessions don't replace the Invoice API

Stripe Checkout supports three modes — `payment`, `subscription`, and `setup` — but **there is no `invoice` mode**. For invoice-based B2B billing, Checkout is the wrong tool. Here's why:

Checkout's `invoice_creation` feature (setting `InvoiceCreation.Enabled = true` on a payment-mode session) generates an invoice **after** payment succeeds. This produces a receipt, not a payable invoice — the opposite of what EU sponsors paying on net-30 terms need. The correct approach is to create the invoice first via the API, then direct sponsors to `inv.HostedInvoiceURL`, which functions as a payment page with all configured payment methods.

That said, Checkout does support SEPA and bank transfer if you need it for other flows. Add `"sepa_debit"` or `"customer_balance"` to `PaymentMethodTypes` on the session. Bank transfers in Checkout require an existing Customer ID. But for the sponsor invoice use case, **the Invoice API + Hosted Invoice Page is the recommended architecture** — it provides full control over payment terms, line items, tax handling, and the invoice lifecycle.

---

## Handling invoice webhooks in a Go HTTP handler

Webhook handling follows a strict pattern: read the raw body, verify the signature, parse the event, handle it, and return **200 OK immediately**. Stripe expects a response within 10 seconds and retries failed deliveries for up to 3 days.

```go
func handleStripeWebhook(w http.ResponseWriter, req *http.Request) {
    const MaxBodyBytes = int64(65536)
    req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)

    body, err := io.ReadAll(req.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
    event, err := webhook.ConstructEvent(
        body,
        req.Header.Get("Stripe-Signature"),
        endpointSecret,
    )
    if err != nil {
        log.Printf("Webhook signature verification failed: %v", err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    switch event.Type {
    case stripe.EventTypeInvoicePaid:
        var invoice stripe.Invoice
        if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        // Activate sponsor benefits, update database
        log.Printf("Invoice %s paid: %d %s", invoice.ID,
            invoice.AmountPaid, invoice.Currency)

    case stripe.EventTypeInvoicePaymentFailed:
        var invoice stripe.Invoice
        json.Unmarshal(event.Data.Raw, &invoice)
        // Notify sponsor, log failure, trigger follow-up

    case stripe.EventTypeInvoiceFinalized:
        var invoice stripe.Invoice
        json.Unmarshal(event.Data.Raw, &invoice)
        // Record hosted URL, update sponsor panel status

    case stripe.EventTypeInvoicePaymentActionRequired:
        var invoice stripe.Invoice
        json.Unmarshal(event.Data.Raw, &invoice)
        // Notify sponsor to complete SCA/3DS authentication

    default:
        log.Printf("Unhandled event type: %s", event.Type)
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"received": true}`))
}
```

The `webhook` package (`github.com/stripe/stripe-go/v84/webhook`) provides `ConstructEvent` with a **default 300-second tolerance** for timestamp validation. For testing, use `ConstructEventIgnoringTolerance` or `GenerateTestSignedPayload` to create signed test payloads. Recent SDK versions also check API version compatibility — use `ConstructEventWithOptions` with `IgnoreAPIVersionMismatch: true` if your endpoint receives events from a different API version.

The essential events for an invoice-based sponsor panel are:

- **`invoice.paid`** — the critical event for activating sponsor benefits
- **`invoice.payment_failed`** — triggers dunning/notification logic
- **`invoice.finalized`** — confirms the invoice is open and payable
- **`invoice.payment_action_required`** — handles SCA/3DS for card payments
- **`invoice.sent`** — confirms email delivery to the sponsor
- **`invoice.finalization_failed`** — catches tax calculation errors (missing address data)

Two best practices matter most: **idempotency** (track processed event IDs since Stripe may deliver duplicates) and **async processing** (acknowledge immediately, process in a background goroutine or queue for anything that takes more than a few seconds).

---

## Letting sponsors self-manage via the Customer Portal

The Stripe Customer Portal is a hosted UI where sponsors can view invoice history, download PDFs, update payment methods, and edit billing information including tax IDs. Set it up by creating a portal configuration and then generating short-lived session URLs.

**Creating a portal configuration** (do this once, store the configuration ID):

```go
configParams := &stripe.BillingPortalConfigurationCreateParams{
    Features: &stripe.BillingPortalConfigurationCreateFeaturesParams{
        InvoiceHistory: &stripe.BillingPortalConfigurationCreateFeaturesInvoiceHistoryParams{
            Enabled: stripe.Bool(true),
        },
        CustomerUpdate: &stripe.BillingPortalConfigurationCreateFeaturesCustomerUpdateParams{
            Enabled: stripe.Bool(true),
            AllowedUpdates: stripe.StringSlice([]string{
                "email", "address", "phone", "tax_id",
            }),
        },
        PaymentMethodUpdate: &stripe.BillingPortalConfigurationCreateFeaturesPaymentMethodUpdateParams{
            Enabled: stripe.Bool(true),
        },
    },
    BusinessProfile: &stripe.BillingPortalConfigurationCreateBusinessProfileParams{
        Headline: stripe.String("Manage your sponsorship billing"),
    },
    DefaultReturnURL: stripe.String("https://yoursite.com/sponsors/dashboard"),
}
config, err := sc.V1BillingPortalConfigurations.Create(ctx, configParams)
```

**Generating a session** when a sponsor clicks "Manage Billing":

```go
func handlePortalRedirect(w http.ResponseWriter, r *http.Request) {
    customerID := getAuthenticatedSponsorStripeID(r) // your auth logic

    params := &stripe.BillingPortalSessionCreateParams{
        Customer:      stripe.String(customerID),
        ReturnURL:     stripe.String("https://yoursite.com/sponsors/dashboard"),
        Configuration: stripe.String("bpc_1abc..."), // your config ID
    }
    session, err := sc.V1BillingPortalSessions.Create(r.Context(), params)
    if err != nil {
        http.Error(w, "Failed to create portal session", 500)
        return
    }
    http.Redirect(w, r, session.URL, http.StatusSeeOther)
}
```

The portal URL is short-lived and single-use. Sponsors can view all past invoices, download PDFs, add or change payment methods (including SEPA mandates), and update their VAT number — all without you building custom UI. Listen for `customer.updated`, `payment_method.attached`, and `customer.tax_id.created` webhooks to sync changes back to your database. One limitation: the portal **cannot be embedded in an iframe** — it must open as a redirect.

---

## Automatic VAT with Stripe Tax for EU invoices

Stripe Tax automatically calculates and applies the correct VAT rate at invoice finalization. Enable it with a single parameter on the invoice and by adding tax registrations in the Stripe Dashboard.

**Setup requirements:**

1. Configure your **head office address** in Dashboard → Settings → Tax
2. Add **tax registrations** for each EU country where you collect VAT (Dashboard → Tax → Registrations)
3. Set a **preset product tax code** (e.g., `txcd_10201000` for SaaS)
4. Set **tax behavior** on your Prices to `exclusive` or `inclusive`

**Creating a tax-enabled invoice:**

```go
// Add your own VAT number to invoices
acctTaxID, _ := sc.V1TaxIDs.Create(ctx, &stripe.TaxIDCreateParams{
    Type:  stripe.String(string(stripe.TaxIDTypeEUVAT)),
    Value: stripe.String("DE123456789"), // your VAT number
})

// Add the sponsor's VAT number
custTaxID, _ := sc.V1TaxIDs.Create(ctx, &stripe.TaxIDCreateParams{
    Type:  stripe.String(string(stripe.TaxIDTypeEUVAT)),
    Value: stripe.String("FR12345678901"), // sponsor's VAT number
    Owner: &stripe.TaxIDOwnerParams{
        Type:     stripe.String("customer"),
        Customer: stripe.String(cust.ID),
    },
})

// Create the invoice with automatic tax
invParams := &stripe.InvoiceCreateParams{
    Customer:         stripe.String(cust.ID),
    CollectionMethod: stripe.String(string(stripe.InvoiceCollectionMethodSendInvoice)),
    DaysUntilDue:     stripe.Int64(30),
    Currency:         stripe.String(string(stripe.CurrencyEUR)),
    AutomaticTax: &stripe.InvoiceAutomaticTaxParams{
        Enabled: stripe.Bool(true),
    },
    AccountTaxIDs: stripe.StringSlice([]string{acctTaxID.ID}),
}
```

**Stripe Tax handles intra-EU B2B reverse charge automatically.** When a sponsor has an `eu_vat` tax ID stored and is in a different EU country than your business, Stripe applies reverse charge — zero tax calculated, with "Reverse charge" notation on the invoice PDF. The sponsor's VAT number is validated asynchronously against the European Commission's **VIES** system. Listen for the `customer.tax_id.updated` webhook to track verification status (`verified`, `unverified`, `unavailable`, `pending`).

Key detail: Stripe Tax applies reverse charge based on the tax ID's **format validity**, not VIES verification status. Even a pending or unverified ID triggers reverse charge if it matches the expected pattern. Tax behavior on Prices is **immutable** after creation — choose `exclusive` (tax added on top, standard for B2B) or `inclusive` (tax included, common for B2C in EU) carefully.

Stripe Tax pricing follows a pay-as-you-go model at **0.5% per transaction** for no-code integrations (invoicing, billing), with a Tax Complete plan starting at **$90/month** that bundles registrations, calculations, and automated filing through Taxually.

---

## Conclusion

The architecture for an EU sponsor invoice panel in Go comes down to five connected pieces. The **Invoice API** with `collection_method: "send_invoice"` and `days_until_due` creates proper payable invoices with net-30/60 terms. **Payment settings** with `sepa_debit` and `customer_balance` (configured for `eu_bank_transfer`) enable the payment methods EU companies expect. **Stripe Tax** with `automatic_tax: {enabled: true}` handles VAT calculation and reverse charge without manual rate management. **Webhooks** centered on `invoice.paid` and `invoice.payment_failed` drive your sponsor activation logic. The **Customer Portal** offloads billing self-service entirely to Stripe's hosted UI.

The most important implementation detail is using `customer_balance` — not `bank_transfer` — as the payment method type for bank transfers, combined with the nested `BankTransfer.Type: "eu_bank_transfer"` configuration. The second is that Stripe Checkout has no invoice mode; the Hosted Invoice Page at `inv.HostedInvoiceURL` is the correct payment UI for this workflow. Build your sponsor panel around these two facts and the rest follows naturally from the stripe-go v84 SDK's `sc.V1Invoices` service methods.
