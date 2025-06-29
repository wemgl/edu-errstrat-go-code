package pizza

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/temporal"
	"time"

	"go.temporal.io/sdk/activity"
)

func GetDistance(ctx context.Context, address Address) (Distance, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("GetDistance invoked; determining distance to customer address")

	// this is a simulation, which calculates a fake (but consistent)
	// distance for a customer address based on its length. The value
	// will therefore be different when called with different addresses,
	// but will be the same across all invocations with the same address.
	kilometers := len(address.Line1) + len(address.Line2) - 10
	if kilometers < 1 {
		kilometers = 5
	}

	distance := Distance{
		Kilometers: kilometers,
	}

	logger.Debug("GetDistance complete", "Distance", distance.Kilometers)
	return distance, nil
}

func SendBill(ctx context.Context, bill Bill) (OrderConfirmation, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("SendBill invoked", "Customer", bill.CustomerID, "Amount", bill.Amount)

	chargeAmount := bill.Amount

	// This month's special offer: Get $5 off all orders over $30
	if bill.Amount > 3000 {
		logger.Info("Applying discount")

		chargeAmount -= 500 // reduce amount charged by 500 cents
	}

	// reject invalid amounts before calling the payment processor
	if chargeAmount < 0 {
		return OrderConfirmation{},
			temporal.NewNonRetryableApplicationError(
				fmt.Sprintf("invalid charge amount: %d (< 1)", chargeAmount),
				"InvalidChargeError", nil, nil)
	}

	// pretend we called a payment processing service here :-)

	confirmation := OrderConfirmation{
		OrderNumber:        bill.OrderNumber,
		ConfirmationNumber: "AB9923",
		Status:             "SUCCESS",
		BillingTimestamp:   time.Now().Unix(),
		Amount:             chargeAmount,
	}

	logger.Debug("SendBill complete", "ConfirmationNumber", confirmation.ConfirmationNumber)

	return confirmation, nil
}

func ProcessCreditCard(ctx context.Context, address Address) (ChargeStatus, error) {
	// pretend to charge card here
	chargestatus := ChargeStatus{
		Success: true,
	}

	if len(address.CardNumber) != 16 {
		return chargestatus, temporal.NewNonRetryableApplicationError("Credit card does not contain 16 digits", "CreditCardError", nil, nil)
	} else {
		return chargestatus, nil
	}
}
