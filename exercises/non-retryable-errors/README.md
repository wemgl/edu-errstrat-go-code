## Exercise #2: Defining Non-Retryable Errors Using Custom Retry Policies

During this exercise, you will:

- Configure non-retryable errors for Activities
- Implement customize retry policies for Activities
- Develop Workflow logic for fallback strategies in the case of Activity failure
- Add Heartbeats and Heartbeat timeouts to help users monitor the health of Activities

Make your changes to the code in the `practice` subdirectory (look for
`TODO` comments that will guide you to where you should make changes to
the code). If you need a hint or want to verify your changes, look at
the complete version in the `solution` subdirectory.

## Part A: Configure Retry Policies of an Error

In this part of the exercise, we will configure the retry policies of an error.

- Initial Interval: Amount of time that must elapse before the first retry occurs
- Backoff Coefficient: How much the retry interval increases (default is 2.0)
- Maximum Interval: The maximum interval between retries
- Maximum Attempts: The maximum number of execution attempts that can be made in the presence of failures

1. In `activities.go`, notice that we added a new Activity called
   `NotifyInternalDeliveryDriver`. This Activity simulates a case where an
   internal driver is not available and forces a hard coded error. We will now
   configure this error.
2. Edit `workflow.go` and set the retry policy to retry once per second for
   five seconds. In the `retry` object of your `proxyActivities`, add in the
   values for `initialInterval`, `backoffCoefficient`, `maximumInterval`,
   `maximumAttempts` that would allow for this.
3. If, after retrying `NotifyInternalDeliveryDriver` once per second for five
   seconds, the Activity is still unsuccessful, you can invoke
   `PollExternalDeliveryDriver`. This Activity will poll a microservice looking
   for external drivers (imagine polling UberEats, Grubhub, DoorDash and so on).
4. Save and close the file.

## Part B: Add Heartbeats

In this part of the exercise, we will add heartbeating to our `PollExternalDeliveryDriver` Activity.

1. Edit `activities.go`. In the `PollExternalDeliveryDriver` Activity, notice
   the `StartingPoint` variable. This variable is set to last incremented
   heartbeat, or 1, if the heartbeating has not began.
2. TK TK like this: https://github.com/temporalio/samples-go/blob/main/retryactivity/retry_activity_workflow.go
3. Save and close the file.

## Part C: Add a Heartbeat Timeout

In this part of the exercise, we will add a Heartbeat Timeout to your Activities.

1. Edit `workflow.go`.
2. Below the `StartToCloseTimeout`, add a `HeartbeatTimeout` and set it to ten
   seconds like so: `HeartbeatTimeout:    10 * time.Second,`. This sets the
   maximum time between Activity Heartbeats. If an Activity times out (e.g., due
   to a missed Heartbeat), the next attempt can use this payload to continue
   from where it left off.
3. Save and close the file.

## Part D: Run the Workflow

Next, let's run the Workflow.

1. In one terminal, start the Worker by running:
   ```bash
   `go run worker/main.go`
   ```
2. In another terminal, start the Workflow by executing `start/main.go`:
   ```bash
   `go run start/main.go`
   ```
3. In your Web UI, you will see that there is an `ActivityTaskFailed` Event for
   `NotifyInternalDeliveryDriver`. In the terminal where your Worker is running,
   you can see that there were five attempts to run this Activity before moving
   onto the `PollExternalDeliveryDriver` Activity.

### This is the end of the exercise.
