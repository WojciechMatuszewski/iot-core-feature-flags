import Amplify, { Auth, Hub } from "aws-amplify";
import React from "react";
import { useAsync } from "react-async";

Amplify.configure({
  Auth: {
    identityPoolId: "us-east-1:dfce3d0d-b161-47bc-acda-f5f3ccbcff76",
    region: "us-east-1",
    userPoolId: "us-east-1_ARxhjWdKv",
    userPoolWebClientId: "2p2r5hgunvumbc4an4ilv726jp"
  }
});

async function getCurrentUser() {
  try {
    const session = await Auth.currentSession();
    const idTokenPayload = session.getIdToken().payload;

    return {
      isAdmin: idTokenPayload["cognito:groups"].includes("admin"),
      id: idTokenPayload["sub"]
    };
  } catch (e) {
    return null;
  }
}

export function useUser() {
  const { data: user = null, isLoading, setData } = useAsync({
    promiseFn: getCurrentUser
  });

  React.useEffect(() => {
    const listener = async ({ payload: { event, data } }) => {
      switch (event) {
        case "signIn":
          setData(await getCurrentUser());
          break;
        case "signOut":
          setData(null);
          break;
        default:
          break;
      }
    };

    const cleanup = Hub.listen("auth", listener);
    return cleanup;
  }, [setData]);

  return { user, loading: isLoading };
}

export async function logout() {
  try {
    await Auth.signOut();
  } catch (e) {
    return null;
  }
}

export async function loginAdmin() {
  try {
    await Auth.signIn({
      password: "test12345",
      username: "admin@admin.com"
    });
  } catch (e) {
    return null;
  }
}
