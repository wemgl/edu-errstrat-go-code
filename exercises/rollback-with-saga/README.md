## Exercise #3: Rollback with the Saga Pattern

During this exercise, you will:

- Orchestrate microservices using a Saga pattern to implement compensating transactions
- Handle failures with rollback logic

Make your changes to the code in the `practice` subdirectory (look for `TODO`
comments that will guide you to where you should make changes to the code). If
you need a hint or want to verify your changes, look at the complete version in
the `solution` subdirectory.

## Part A: Create a Test Error

In this part of the exercise, you will define a Test Error that you will use to
test the rolling back of compensations with the Saga Pattern.

1. Edit the `activities.go` file.
2. At the very top of the `SendBill` Activity, return an Application Failure
   that contains the message: `Test Error`. You will return this error in the
   `SendBill` Activity. Set the error's `nonRetryable` key to `true`. This way,
   when you trigger this intentional error, you can demonstrate rollback of an
   Actiivty Failure without your Activity first attempting to retry.
3. Save the file.

## Part B: Create your Compensation Activities

A new `UpdateInventory` Activity has been added to the code from the previous
exercise. This Activity reduces the stock from the pizza inventory once the
pizza order comes through. This step is done before the `SendBill` Activity is
called.

Imagine that there is an error in the `SendBill` Activity. You would then need
to cancel the billing step by invoking a `RefundCustomer` Activity. You would
also need to roll back the `UpdateInventory` Activity by invoking a
`RevertInventory` Activity, which would add the ingredients back into the pizza
inventory.

In this part of the exercise, you will create your compensation Activities. When
one of the Activities fails, the Workflow will "compensate" by calling
Activities which reverse the successful calls up to that point.

1. Edit the `activities.go` file.
2. Uncomment the `RevertInventory` Activity.
3. Uncomment the `RefundCustomer` Activity.
4. Pass in `bill` with the type `Bill` (imported from `shared.go`) into the `RefundCustomer` Activity.
5. Edit the `workflow.go` file.
6. Add your `RefundCustomer` and `RevertInventory` Activities in the `ProxyActivities`.

## Part C: Create Your Compensation List

In this part of the exercise, you will create an array of compensation objects.
Each compensation object will include an Activity that would cause the rolling
back of the Activity that would fail.

We want the array of compensation objects to take on a shape like this:
```
[{message: 'unable to call Activity A successfully',
  fn: revertActivityA()},
{message: 'unable to call Activity B successfully',
  fn: revertActivityB()}]
```

1. Edit `shared.go` file.
2. Create a new `struct` called `Compensation`.
3. It should take a `string` key called `message`.
4. Save the file.

## Part D: Create Your Compensation Function

In this part of the exercise, you will create a function which will loop through
an array of compensation objects. In the case of an error, you will invoke this
function to roll back on Activities you need to undo.

1. Edit the `activties.go` file.
2. Import your `Compensation` interface from `shared.go` that you defined in Part C.
3. Note that there is already an `ErrorMessage` function which takes in the
   error message from a failing Activity and displays it in a more readable
   fashion.
4. Now, look at the next function: `Compensate`. This function will take in a
   list of the `Compensation` objects that you defined in part C. It will then
   iterate through a list of `Compensation` objects, log the error message
   provided in the `Compensation` object, and call the function provided in the
   `Compensation` object.
5. Add a `for` loop that iterates over the `compensations` array.
6. Save the file.

## Part E: Fill in Your Compensation Array

In this part of the exercise, you will fill in the `compensations` array that
you will call the `Compensate` function on.

Before we call an Activity, we want to add the correlating compensation Activity
into the `compensations` list. For example, before we call `SendBill`, we want
to add `RefundCustomer` into the list of compensations.

Then, if `SendBill` throws an error, we call the `Compensate` function which
rolls back on the `SendBill` Activity by calling `RefundCustomer`.

1. Edit the `workflow.go` file.
2. Import your `Compensation` interface from `shared.go` that you defined in
   Part C. Notice that on line 20, we define a variable called `compensations`,
   which is a list of `Compensation` objects (and defaulted as an empty array).
3. Look at the first compensation (compensation for `UpdateInventory`) which is
   provided for you on line 67. Before we call `UpdateInventory`, we add its
   compensating counterpart - `RevertInventory` into the array of
   `compensations`. We use the `unshift` method, which adds an item in the
   beginning of an array. This ensures that the compensations are executed in
   the reverse order of their addition, which is important for correctly
   reversing the steps of the Workflow.
4. Following the pattern in step 3, add a compensation for an error in the
   `SendBill` Activity. On line 93, add in a compensation object for `SendBill`
   by calling `RefundCustomer` which takes in a `bill` argument.
5. At this point, as you go through the pizza Workflow, your `compensations`
  array should look like this:
  ```
  [{message: 'reversing send bill: ',
    fn: refundCustomer
  },
  {message: 'reversing update inventory: ',
    fn: revertInventory
  }]
  ```

## Part F: Call the `compensate` Function

In this part of the exercise, you will call the `Compensate` function that you defined in Part D.

1. Edit the `workflow.go` file.
2. Import your `Compensate` function and `ErrorMessage` function from the
   `activities.go` file you looked at in part D.
3. In the `defer` block of your `sendBill` Activity, add an `if err != nil {}`
   block. Now if `SendBill` fails, you can roll back `SendBill` by calling
   `RefundCustomer`, then roll back on `UpdateInventory` by calling
   `RevertInventory`.
4. Save the file.

## Part G: Test the Rollback of Your Activities

To run the Workflow:

Next, let's run the Workflow.

1. In one terminal, start the Worker by running:
   ```bash
   `go run worker/main.go`
   ```
2. In another terminal, start the Workflow by executing `start/main.go`:
   ```bash
   `go run start/main.go`
   ```
3. In your Web UI, you should see a `WorkflowExecutionFailed` Event to indicate
   that the Workflow failed. After the `SendBill` Activity, we then
   called the Activities: `RefundCustomer` and `RevertInventory`.

### This is the end of the exercise.
