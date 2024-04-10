import React from "react";
import { useRouteError } from "react-router-dom";
import { errorMessage } from "@app/utils";

const ErrorPage: React.FunctionComponent = () => {
  const error = useRouteError();
  console.error(error);

  return (
    <div id="error-page">
      <h1>Oops!</h1>
      <p>Sorry, an unexpected error has occurred.</p>
      <p>
        <i>{errorMessage(error)}</i>
      </p>
    </div>
  );
};

export default ErrorPage;
