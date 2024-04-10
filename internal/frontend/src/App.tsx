import { createBrowserRouter, RouterProvider } from "react-router-dom";
import "@patternfly/react-core/dist/styles/base.css";
import { nav } from "@app/routes/nav";
import { createConnectTransport } from "@connectrpc/connect-web";
import { TransportProvider } from "@connectrpc/connect-query";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Transport } from "@connectrpc/connect";

const queryClient = new QueryClient();

// AppConfig values are passed via index.html from frontend.go using template rendering
// const appConfig = window.holosAppConfig;

const location = new URL(window.location.href);
// Router basename must match the frontend.Path const in frontend.go
const router = createBrowserRouter(nav.routes, { basename: "/app/" });

function App({ transport }: { transport?: Transport }) {
  const finalTransport =
    transport ??
    createConnectTransport({
      baseUrl: `${location.protocol}//${location.host}`,
    });

  return (
    <TransportProvider transport={finalTransport}>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    </TransportProvider>
  );
}

export default App;
