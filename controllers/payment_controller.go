package controllers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"griddy-payment-service/models/payment"
	"griddy-payment-service/services"
	httpResponse "griddy-payment-service/utilities/http"
	"net/http"
)

// PaymentController defines the interface for working with payments
type PaymentController interface {
	List(writer http.ResponseWriter, request *http.Request)
	Post(writer http.ResponseWriter, request *http.Request)
}

type PaymentControllerImpl struct {
	paymentService services.PaymentService
}

var _ PaymentController = (*PaymentControllerImpl)(nil)

// RenderJSON is an alias of method to render JSON for easy mocking in tests
var RenderJSON = render.JSON

func GetPaymentController() PaymentController {
	return &PaymentControllerImpl{
		paymentService: services.GetPaymentService(),
	}
}

func (controller *PaymentControllerImpl) List(w http.ResponseWriter, req *http.Request) {
	accountId := getParam(req, "account_id")
	if accountId == "" {
		w.WriteHeader(http.StatusBadRequest)
		RenderJSON(w, req, httpResponse.BadRequest("Invalid URL parameter \"account_id\" received"))
		return
	}

	payments, err := controller.paymentService.GetPayments(req.Context(), accountId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		RenderJSON(w, req, httpResponse.BadRequest(err.Error()))
		return
	}

	if payments == nil {
		httpResponse.RenderJSON(w, req, []payment.DTO{})
		return
	}

	httpResponse.RenderJSON(w, req, payments)
}
// Post creates a PaymentIntent within Stripe and adds a simplified result to the DB
// This endpoint assumes that amount is a float within a string
func (controller *PaymentControllerImpl) Post(w http.ResponseWriter, req *http.Request) {
	paymentReq, err := getPostBody(req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		RenderJSON(w, req, httpResponse.BadRequest(err.Error()))
		return
	}

	if paymentReq.CustomerId == "" || paymentReq.Amount == "" {
		w.WriteHeader(http.StatusBadRequest)
		RenderJSON(w, req, httpResponse.BadRequest("\"account_id\" and \"amount\" are required fields in the POST body"))
		return
	}

	createdPayment, err := controller.paymentService.CreatePayment(req.Context(), *paymentReq)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		RenderJSON(w, req, httpResponse.BadRequest(err.Error()))
		return
	}

	httpResponse.RenderJSON(w, req, createdPayment)
}

// getPostBody wraps operations for extracting request body for easy mocking in tests
var getPostBody = func(r *http.Request) (*payment.PostBody, error) {
	var data payment.PostBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	return &data, err
}

// getParam wraps chi.URLParam for easy mocking in test
var getParam = func(r *http.Request, p string) string {
	return chi.URLParam(r, p)
}
