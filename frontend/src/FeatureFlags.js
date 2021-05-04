import Amplify, { API, PubSub } from "aws-amplify";
import { AWSIoTProvider } from "@aws-amplify/pubsub";
import React, { createContext, useContext, useEffect } from "react";
import { useAsync } from "react-async";
import { Button, Checkbox, Form, Input } from "antd";

Amplify.addPluggable(
  new AWSIoTProvider({
    aws_pubsub_region: "us-east-1",
    aws_pubsub_endpoint:
      "wss://au8j1jfsb53ln-ats.iot.us-east-1.amazonaws.com/mqtt"
  })
);

API.configure({
  endpoints: [
    {
      name: "backend",
      endpoint: "https://1r2nwr46qa.execute-api.us-east-1.amazonaws.com"
    }
  ]
});

async function getFlags({ clientID }) {
  return API.get("backend", `/flags/${clientID}`);
}

async function setFlag([{ clientID, flagName, value }]) {
  return fetch(
    `https://1r2nwr46qa.execute-api.us-east-1.amazonaws.com/flag/${clientID}`,
    {
      method: "POST",
      body: JSON.stringify({
        value,
        flagName
      })
    }
  );
}

const FeatureFlagsContext = createContext({});

export function FeatureFlagsProvider({ children, clientID }) {
  const { data: flags, setData: setFlags } = useAsync({
    promiseFn: getFlags,
    clientID,
    initialValue: {}
  });

  useEffect(() => {
    const cleanup = PubSub.subscribe(`flags/${clientID}`).subscribe({
      next: ({ value: update }) => {
        setFlags({ ...flags, ...update });
      }
    });

    return () => cleanup.unsubscribe();
  }, [setFlags, clientID, flags]);

  return (
    <FeatureFlagsContext.Provider value={{ flags, clientID }}>
      {children}
    </FeatureFlagsContext.Provider>
  );
}

export function SetFlag() {
  const { run, isLoading } = useAsync({
    deferFn: setFlag,
    onReject: console.log
  });

  const clientID = useClientID();

  return (
    <Form
      layout="inline"
      onFinish={values => {
        run({ ...values, clientID });
      }}
      name="setFlag"
      initialValues={{
        flagName: undefined,
        isEnabled: false
      }}
    >
      <Form.Item label="Name of the flag" name="flagName">
        <Input type="text" placeholder="flagName" />
      </Form.Item>

      <Form.Item label="Enabled" name="value" valuePropName="checked">
        <Checkbox />
      </Form.Item>
      <Form.Item>
        <Button htmlType="submit" type="primary" loading={isLoading}>
          Save
        </Button>
      </Form.Item>
    </Form>
  );
}

export function useFeatureFlags() {
  const ctx = useContext(FeatureFlagsContext);
  if (!ctx) {
    throw new Error("boom");
  }

  return ctx.flags;
}

function useClientID() {
  const ctx = useContext(FeatureFlagsContext);
  if (!ctx) {
    throw new Error("boom");
  }

  return ctx.clientID;
}
