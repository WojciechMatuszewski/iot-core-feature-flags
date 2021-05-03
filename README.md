# Learnings

- Authenticated Cognito users require you to deal with both IAM Policy that is
  associated with the identity and the IOT Policy.

- You will need to run `attach-principal-policy` command whenever you are dealing with authenticated Cognito users. I'm not sure why that is the case.

- You can listen to _IOT mqtt_ events, the simplest query `SELECT * FROM 'myTopic'`

  - You can listen to subscribe/unsubscribe events via the _topic filter_.
    Here is the list of the lifecycle events https://docs.aws.amazon.com/iot/latest/developerguide/life-cycle-events.html.
    And a simple SQL query `SELECT * FROM '$aws/events/subscriptions/subscribed/+'`

- There is a difference between _the Thing_ endpoint and the _global IOT endpoint_.
  The latter does not seem to have a policy attached to it so you cannot reach the endpoint?

- To scope down the websocket connection permissions on the frontend, use the `iot:Connect`.
  It seems that `iot:Subscribe` is totally different than `iot:Connect`. The documentation mentions that we are connecting to a _thing_ but I did not create any? (left the `*`)

- With CDK you can create custom resources using either a "framework" where you define a lambda (skipped in this case as the golang lambda library is doing this for me), or you can use sdk calls directly (through construct of course)

- If the path of the POST request is wrong, you might see an CORS error thrown from APIGW.
  While you could setup _Gateway responses_ within REST APIs that would correctly handle cors on the APIGW level, you cannot do that using HTTP APIs.
