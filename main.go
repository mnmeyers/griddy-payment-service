package main

import (
	"griddy-payment-service/configurations"
	"griddy-payment-service/routes"
	"griddy-payment-service/utilities"
	"net/http"
)

func main() {
	configuration := configurations.GetConfiguration()

	port := configuration.Port

	utilities.PrintStartMessage(port)

	router := routes.GetRouter()

	panic(http.ListenAndServe(":"+port, router))
}
