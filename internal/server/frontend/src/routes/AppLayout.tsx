import React from "react";
import { Outlet } from "react-router-dom";

interface AppLayoutProps {}

export default function AppLayout(
  props: React.PropsWithChildren<AppLayoutProps>,
) {
  return props.children ?? <Outlet />;
}
