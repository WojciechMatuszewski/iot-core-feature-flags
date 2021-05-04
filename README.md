# IOT MQTT Feature Flags

Not finished and most likely never be.

## Deployment

1. Make sure you have all your AWS stuff set up

2. ```sh
    cd backend
    make bootstrap
    make deploy
   ```

3. Copy outputs to the right places in the frontend app

## Learnings

- Authenticated Cognito users require you to deal with both IAM Policy that is
  associated with the identity and the IOT Policy.

- You will need to run `attach-principal-policy` command whenever you are dealing with authenticated Cognito users. I'm not sure why that is the case.

- You can listen to _IOT MQTT_ events, the simplest query `SELECT * FROM 'myTopic'`

- You can listen to subscribe/unsubscribe events via the _topic filter_.
  Here is the list [of the lifecycle events](https://docs.aws.amazon.com/iot/latest/developerguide/life-cycle-events.html).
  And a simple SQL query `SELECT * FROM '$aws/events/subscriptions/subscribed/+'`

- There is a difference between _the Thing_ endpoint and the _global IOT endpoint_.

- To scope down the websocket connection permissions on the frontend, use the `iot:Connect`.
  It seems that `iot:Subscribe` is totally different than `iot:Connect`. The documentation mentions that we are connecting to a _thing_ but I did not create any?

- With CDK you can create custom resources using either a "provider framework", or you can use sdk calls directly (through a construct of course)

- If the path of the POST request is wrong, you might see an CORS error thrown from APIGW.
  For the REST APIs you can use _Gateway responses_ to specify headers returned when request fails on the APIGW REST API level.
  For the HTTP APIs one might try [parameter mappings](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-parameter-mapping.html). I'm yet to give it a try.
