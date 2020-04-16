package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	"griddy-payment-service/configurations"
	"griddy-payment-service/models/payment"
	"griddy-payment-service/repositories"
	"griddy-payment-service/utilities"
	"strconv"
	"sync"
	"time"
)

// PaymentService defines the interface for the service that deals with payments
type PaymentService interface {
	GetPayments(ctx context.Context, customerId string) ([]payment.DTO, error)
	CreatePayment(ctx context.Context, postBody payment.PostBody) (*payment.DTO, error)
}

// PaymentServiceImpl defines an instance of the payment service
type PaymentServiceImpl struct {
	paymentRepository repositories.PaymentRepository
}

var _ PaymentService = (*PaymentServiceImpl)(nil)
var paymentService PaymentService
var oncePaymentService sync.Once

// GetPaymentService returns a thread-safe singleton of the payment service.
func GetPaymentService() PaymentService {
	oncePaymentService.Do(func() {
		paymentService = &PaymentServiceImpl{
			paymentRepository: repositories.GetDatabase()}
	})

	return paymentService
}

func (service *PaymentServiceImpl) GetPayments(ctx context.Context, customerId string) ([]payment.DTO, error) {
	ch := make(chan []payment.DAO, 1)
	defer close(ch)

	var payments []payment.DAO
	var err error
	go func() {
		payments = service.getPaymentsFromStripe(customerId)
		ch <- payments
	}()

	select {
	case paymentsFromStripe := <-ch:
		fmt.Println("Successfully connected to Stripe API")
		if len(paymentsFromStripe) > 0 {
			go service.paymentRepository.UpdatePayments(ctx, customerId, paymentsFromStripe)
		}
	case <-time.After(2 * time.Second):
		fmt.Println("Stripe API exceeded time limit of 1 second")
		payments, err = service.paymentRepository.GetPaymentsByCustomerId(ctx, customerId)
	}
	if err != nil {
		return nil, err
	}

	paymentDTOs, err := service.convertPaymentDAOsToDTOs(payments)

	if err != nil {
		return nil, err
	}
	return paymentDTOs, nil
}

func (service *PaymentServiceImpl) CreatePayment(ctx context.Context, postBody payment.PostBody) (*payment.DTO, error) {
	result, err := service.postPaymentIntentToStripe(postBody)
	if err != nil {
		return nil, err
	}
	paymentDao, err := service.convertPostBodyToDAO(postBody, result)
	if err != nil {
		return nil, err
	}
	go service.paymentRepository.CreatePayment(ctx, *paymentDao)

	paymentDto, err := service.convertPaymentDAOToDTO(*paymentDao)
	if err != nil {
		return nil, err
	}
	return paymentDto, nil
}

func (service *PaymentServiceImpl) getPaymentsFromStripe(customerId string) []payment.DAO {
	var payments []payment.DAO

	if stripe.Key == "" {
		stripe.Key = configurations.GetConfiguration().StripeKey
	}

	params := &stripe.PaymentIntentListParams{
		Customer: &customerId,
	}
	i := paymentintent.List(params)
	for i.Next() {
		pi := i.PaymentIntent()
		payments = append(payments, payment.DAO{
			Amount:     pi.Amount,
			CustomerId: pi.Customer.ID,
			PaymentId:  pi.ID,
			Status:     pi.Status,
		})
	}
	return payments
}

func (service *PaymentServiceImpl) convertPostBodyToDAO(body payment.PostBody, paymentIntent *stripe.PaymentIntent) (*payment.DAO, error) {
	amount, err := utilities.ConvertStringToInt64InPennies(body.Amount)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid amount received: %s", body.Amount))
	}
	return &payment.DAO{
		Amount:     amount,
		CustomerId: body.CustomerId,
		PaymentId:  paymentIntent.ID,
		Status:     paymentIntent.Status,
	}, nil
}

func (service *PaymentServiceImpl) convertPaymentDAOsToDTOs(daos []payment.DAO) ([]payment.DTO, error) {
	var dtos []payment.DTO

	for _, dao := range daos {
		dto, err := service.convertPaymentDAOToDTO(dao)
		if err != nil {
			return nil, err
		}
		dtos = append(dtos, *dto)
	}
	return dtos, nil
}

func (service *PaymentServiceImpl) convertPaymentDAOToDTO(dao payment.DAO) (*payment.DTO, error) {
	formattedAmount, err := utilities.FormatStringIntoDollarsAndCents(strconv.FormatInt(dao.Amount, 10))
	if err != nil {
		return nil, err
	}
	return &payment.DTO{
		PaymentId: dao.PaymentId,
		Amount:    formattedAmount,
		Status:    fmt.Sprint(dao.Status),
	}, nil
}

func (service *PaymentServiceImpl) postPaymentIntentToStripe(body payment.PostBody) (*stripe.PaymentIntent, error) {
	amount, err := utilities.ConvertStringToInt64InPennies(body.Amount)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid amount received: %s", body.Amount))
	}

	if stripe.Key == "" {
		stripe.Key = configurations.GetConfiguration().StripeKey
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Customer:     &body.CustomerId,
	}
	paymentIntent, err := paymentintent.New(params)

	if err != nil {
		return nil, err
	}

	return paymentIntent, nil
}
