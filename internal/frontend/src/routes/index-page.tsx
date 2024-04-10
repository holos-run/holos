import React from "react";
import {
  EmptyState,
  EmptyStateActions,
  EmptyStateBody,
  EmptyStateFooter,
  EmptyStateHeader,
  EmptyStateIcon,
} from "@patternfly/react-core";
import CubesIcon from "@patternfly/react-icons/dist/esm/icons/cubes-icon";

const IndexPage: React.FunctionComponent = () => (
  <EmptyState>
    <EmptyStateHeader
      titleText="Holos Platform"
      headingLevel="h4"
      icon={<EmptyStateIcon icon={CubesIcon} />}
    />
    <EmptyStateBody>Welcome to Holos!</EmptyStateBody>
    <EmptyStateFooter>
      <EmptyStateActions></EmptyStateActions>
    </EmptyStateFooter>
  </EmptyState>
);

export default IndexPage;
