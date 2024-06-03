## Exercise #1: Defining a Custom Error

During this exercise, you will:

- Throw and handle exceptions in Temporal Workflows and Activities
- Use non-retryable errors to fail an Activity

Make your changes to the code in the `practice` subdirectory (look for
`TODO` comments that will guide you to where you should make changes to
the code). If you need a hint or want to verify your changes, look at
the complete version in the `solution` subdirectory.

## Part A: Throw a non-retryable Application Error to fail an Activity

In this part of the exercise, you will throw a non-retryable Application Failure 
that will fail your Activities.

Application Failures are used to communicate application-specific failures in
Workflows and Activities. In Activities, returning an `ApplicationError` will
cause the Activity to fail. However, this unless this Activity is specified as
non-retryable, it will retry according to the Retry Policy. To have an Activity
fail when an `ApplicationError` is returned, set it as non-retryable. Any other
error that is returned in Go is automatically converted to an `ActivityError`.

1. Start by opening and reviewing the `activities.go` file, familiarizing yourself with
   the Activity functions.
2. In the `SendBill` Activity, notice how an error is thrown if the
   `chargeAmount` is less than 0. If the calculated amount to charge the
   customer is negative, the Activity returns with a non-`nil` error.
   Specifically, it returns with `fmt.Errorf("invalid charge amount: %d (< 1)",
   chargeAmount)`. In Go, returning `fmt.Errorf()` lets you return an error with
   customized output without having to define a new custom error class — it will
   be converted to an `ActivityError` by Temporal. By default, depending on
   this Workflow's retry policy, this Activity may be retried if it fails, but
   since an invalid charge amount would probably reflect bad data being
   provided, you may not want it to retry. To do that, you can replace the
   `fmt.Errorf()` return value with
   `temporal.NewNonRetryableApplicationError(fmt.Sprintf("invalid charge amount: %d (< 1)", chargeAmount), "invalidChargeError", nil, nil)`. This will use the same formatting, but will
   mark the error as a non-retryable `invalidChargeError`.
3. Go to `ProcessCreditCard` Activity, where you will return an
   `ApplicationError` if the credit card fails its validation step. In this
   Activity, you will throw an error if the entered credit card number does not
   have 16 digits. Use the `NewNonRetryableApplicationError()` code from the
   previous step as a reference. Designate this as a `CreditCardError`.
4. Save your file.

## Part B: Catch the Activity Failure

In this part of the exercise, you will catch the `ApplicationError` that was
thrown from the `ProcessCreditCard` Activity and handle it.

1. Open `workflow.go` in your text editor.
2. `ProcessCreditCard` is run like so:
   `err = workflow.ExecuteActivity(ctx, ProcessCreditCard, confirmation).Get(ctx, &chargestatus)`.
   Note that the results are returned to `err`, and the next code block, `if err != nil {}`,
   handles any errors that have been returned from the Activity. By default, if a non-retryable
   `ApplicationFailure` is thrown, the Workflow will fail. However, it is possible to catch this
   failure and either handle it or continue to propagate it up.
3. Within the `if err != nil {}` block, check for an `ApplicationError` like so:
   ```
   var applicationErr *temporal.ApplicationError
		if errors.As(err, &applicationErr) {}
   ```			
   Within that block, you can handle the error, add additional logging, and so on. For example, you could add:
   ```
   println("Billing timestamp of failed order:", confirmation.BillingTimestamp)
   logger.Error("Unable to charge credit card", "Error", err)
   ```
   Finally, make sure that you add `return OrderConfirmation{}, err` at the end of the
   `if err != nil {}` block, so the error is returned.
4. Save your file.

## Part C: Run the Workflow

In this part of the exercise, you will run your Workflow and see both your
Workflow and Activity succeed and fail.

The starter used to run this Workflow, `start/main.go`, contains a "valid"
16-digit credit card number which will be used by default, causing the Workflow
to complete successfully.

**First, run the Workflow successfully:**

1. In one terminal, start the Worker by running:
   ```bash
   `go run worker/main.go`
   ```
2. In another terminal, start the Workflow by executing `start/main.go`:
   ```bash
   `go run start/main.go`
   ```
3. In the Web UI, verify that the Workflow ran successfully to completion.

**Next, run the Workflow with a bad credit card number to trigger a failure:**

1. You do not need to restart the Worker. In the terminal where you ran
   `start/main.go`, run it again with an additional parameter, `--creditcard`.
   Supply an obviously wrong credit card number like "1234":
   ```bash
   `go run start/main.go --creditcard 1234`
   ```
2. You should see the Workflow fail in the terminal where you executed `go run start/main.go`.
   You can also check the Web UI and view the failure there.

### This is the end of the exercise.







Old / unused:

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

Exercise #2 would obviously have to cover these -- because the primary reason to use custom Error types (in Go at least) is, as far as I can tell, so you can specify `nonRetryableErrorTypes` in your Retry Policy. But I wonder if this is enough error handling to scaffold in Exercise #1 before we go there.
