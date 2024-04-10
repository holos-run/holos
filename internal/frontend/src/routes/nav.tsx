import React from "react";
import Root from "./root";
import ErrorPage from "./error-page";
import IndexPage from "./index-page";
import Todo from "./todo";
import PageNotFound from "./page-not-found";
import { Nav, NavList } from "@patternfly/react-core";
import { NavLink } from "react-router-dom";
import { StartPage } from "./start";

const linkClassName = ({ isActive }: { isActive: boolean }): string => {
  return isActive ? "pf-v5-c-nav__link pf-m-current" : "pf-v5-c-nav__link";
};

interface AppRouteNav {
  label: string;
  to: string;
}

interface AppRouteBase {
  path?: string;
  element?: React.ReactNode;
  errorElement?: React.ReactNode;
  nav?: AppRouteNav;
}

interface IndexAppRoute extends AppRouteBase {
  index?: true;
  children?: undefined;
}

interface NonIndexAppRoute extends AppRouteBase {
  index?: false;
  children?: AppRoute[];
}

type AppRoute = IndexAppRoute | NonIndexAppRoute;

interface INavRoutes {
  routes: AppRoute[];
  Navigation(): React.ReactNode;
}

const getNavItems = (route: AppRoute): AppRouteNav[] => {
  const navs: AppRouteNav[] = [];
  if (route.nav) {
    navs.push(route.nav);
  }
  if (route.children) {
    route.children.forEach((child) => {
      getNavItems(child).forEach((nav) => {
        navs.push(nav);
      });
    });
  }
  return navs;
};

const renderNavItem = (route: AppRoute): React.ReactElement[] => {
  return getNavItems(route).map((nav, idx) => (
    <li key={`nav-${nav.label}-${idx}`}>
      <NavLink to={nav.to} className={linkClassName}>
        {nav.label}
      </NavLink>
    </li>
  ));
};

class NavRoutes implements INavRoutes {
  routes: AppRoute[];

  constructor(routes: AppRoute[]) {
    this.routes = routes;
  }

  // Sidebar Nav component
  Navigation(): React.ReactNode {
    return (
      <Nav id="nav-primary" theme="dark">
        <NavList id="nav-list">
          {this.routes.map((route) => renderNavItem(route))}
        </NavList>
      </Nav>
    );
  }
}

const routes: AppRoute[] = [
  {
    path: "/",
    element: <Root />,
    errorElement: <ErrorPage />,
    children: [
      {
        errorElement: <ErrorPage />,
        children: [
          {
            index: true,
            element: <IndexPage />,
            nav: {
              label: "Home",
              to: "/",
            },
          },
          {
            path: "start",
            element: <StartPage />,
            nav: {
              label: "Start",
              to: "start",
            },
          },
          {
            path: "connections",
            element: <Todo />,
            nav: {
              label: "Connections",
              to: "connections",
            },
          },
          {
            path: "profile",
            element: <Todo />,
            nav: {
              label: "Profile",
              to: "profile",
            },
          },
          {
            path: "*",
            element: <PageNotFound />,
          },
        ],
      },
    ],
  },
];

export const nav = new NavRoutes(routes);
