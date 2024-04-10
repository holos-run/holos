import React from "react";
import { Link } from "react-router-dom";
import {
  Button,
  EmptyState,
  EmptyStateActions,
  EmptyStateBody,
  EmptyStateFooter,
  EmptyStateHeader,
  EmptyStateIcon,
  EmptyStateVariant,
} from "@patternfly/react-core";
import PathMissingIcon from "@patternfly/react-icons/dist/esm/icons/path-missing-icon";

const PageNotFound: React.FunctionComponent = () => (
  <EmptyState variant={EmptyStateVariant.xl}>
    <EmptyStateHeader
      titleText="404: Page no longer exists"
      headingLevel="h4"
      icon={<EmptyStateIcon icon={PathMissingIcon} />}
    />
    <EmptyStateBody>
      Another page might have the information you need.
    </EmptyStateBody>
    <EmptyStateFooter>
      <EmptyStateActions>
        <Link to={`/`}>
          <Button variant="primary">Return to home page</Button>
        </Link>
      </EmptyStateActions>
    </EmptyStateFooter>
  </EmptyState>
);

export default PageNotFound;
