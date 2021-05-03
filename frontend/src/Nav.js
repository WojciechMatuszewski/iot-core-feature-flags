import { Alert, Menu, Spin } from "antd";
import { Fragment } from "react";
import { useAsync } from "react-async";
import { loginAdmin, logout, useUser } from "./user";

export function Nav() {
  const { user, loading } = useUser();

  const { run: logoutUser, isLoading: isLoggingOut } = useAsync({
    deferFn: logout
  });

  const { run: loginAsAdmin, isLoading: isLoggingAsAdmin } = useAsync({
    deferFn: loginAdmin
  });

  if (loading) {
    return null;
  }

  return (
    <Fragment>
      <Alert
        banner={true}
        message={`User ID: ${user ? user.id : "unknown"}`}
      ></Alert>
      <Menu mode="horizontal">
        <Menu.Item>Home</Menu.Item>
        {!user && (
          <Menu.Item onClick={loginAsAdmin}>
            Login as admin
            {isLoggingAsAdmin && (
              <Spin size="small" style={{ marginLeft: 12 }} />
            )}
          </Menu.Item>
        )}

        {user?.isAdmin && <Menu.Item>Admin panel</Menu.Item>}
        {user && (
          <Menu.Item onClick={logoutUser}>
            Logout
            {isLoggingOut && <Spin size="small" />}
          </Menu.Item>
        )}
      </Menu>
    </Fragment>
  );
}
