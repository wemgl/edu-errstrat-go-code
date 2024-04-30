## Exercise #1: Defining a Custom Error

During this exercise, you will:

Make your changes to the code in the `practice` subdirectory (look for
`TODO` comments that will guide you to where you should make changes to
the code). If you need a hint or want to verify your changes, look at
the complete version in the `solution` subdirectory.

Some thoughts on what I'm doing here (not final enough to develop Practice steps yet):

- I've added another activity, `ChargeCreditCard`, that will always fail and return `errors.New("Credit Card Charge Error")`. When an Activity does this in Go, you're supposed to add a bunch of handling to the Activity call in the Workflow definition where you can check for different errors and do things with them, so I have this:

```go
	err = workflow.ExecuteActivity(ctx, ChargeCreditCard, confirmation).Get(ctx, &chargestatus)
	if err != nil {
		var applicationErr *temporal.ApplicationError
		if errors.As(err, &applicationErr) {
			// You could be pushing individual values to a logging system here
			println("Billing timestamp of failed order:", confirmation.BillingTimestamp)
			logger.Error("Unable to charge credit card", "Error", err)
		}

		var canceledErr *temporal.CanceledError
		if errors.As(err, &canceledErr) {
			// handle cancellation
		}

		return OrderConfirmation{}, err
	}
```

This also shows a bit of advanced error handling and some other ways to act on error data rather than just returning it and failing. I've changed the replay policy from the default for the pizza workflow to 3x, so that a user can see this retry a certain number of times and then fail.

- I've *also* added the ability to override the customer's street address in the starter via a flag (`var streetAddress = flag.String("address", "1 Main St", "Provide a street address"`). This means that, if you supply a sufficiently long street address (because of the goofy way we calculate delivery distance), the `GetDistance` Activity can fail. This is a nice contrast from the `ChargeCreditCard` failure because it calls `errors.New()` from the *Workflow* rather than the Activity code. In Go, doing this automatically fails the Workflow Execution. The idea would be that a user could first run this with the default, working address to get the Activity Errors, which retry, and then run it with a longer address, to throw the Workflow Error, which doesn't.

Is this a good start? Do we still want to show a Custom Error type on top of this? Temporal turns `errors.New()` into `ApplicationError` pretty elegantly (and lets you override their params without having to actually define a new type), so I wanted a head check -- especially since anything else to do with Error types would mostly be a Go lesson.