import { Button, Spin } from "antd";
import { Fragment } from "react";
import { FeatureFlagsProvider, SetFlag, useFeatureFlags } from "./FeatureFlags";
import { Nav } from "./Nav";
import { useUser } from "./user";

function App() {
  const { user, loading } = useUser();
  if (loading) {
    return <Spin />;
  }

  const clientID = user ? user.id : "unknown";
  return (
    <Fragment>
      <Nav />
      <FeatureFlagsProvider clientID={clientID}>
        <SomeFeature />
        <SetFlag />
      </FeatureFlagsProvider>
    </Fragment>
  );
}

function SomeFeature() {
  const { showButton } = useFeatureFlags();

  if (!showButton) {
    return null;
  }

  return <Button>I'm a button</Button>;
}

export default App;
