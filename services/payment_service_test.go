package services

import (
	"github.com/stretchr/testify/assert"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	"griddy-payment-service/configurations"
	"testing"
)

func TestThatStripeAPIKeysWorkProperly(t *testing.T) {
	expect := assert.New(t)

	stripe.Key = configurations.GetConfiguration().StripeKey

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(1000),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		ReceiptEmail: stripe.String("michalnmeyers@gmail.com"),
	}
	actual, err := paymentintent.New(params)
	expect.Nil(err)
	expect.NotNil(actual)
}
