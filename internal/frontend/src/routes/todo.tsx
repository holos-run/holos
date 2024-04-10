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
} from "@patternfly/react-core";
import CubesIcon from "@patternfly/react-icons/dist/esm/icons/cubes-icon";

const Todo: React.FunctionComponent = () => (
  <EmptyState>
    <EmptyStateHeader
      titleText="Holos Platform"
      headingLevel="h4"
      icon={<EmptyStateIcon icon={CubesIcon} />}
    />
    <EmptyStateBody>Under Construction</EmptyStateBody>
    <EmptyStateFooter>
      <EmptyStateActions>
        <Link to={`/`}>
          <Button variant="primary">Go Home</Button>
        </Link>
      </EmptyStateActions>
    </EmptyStateFooter>
  </EmptyState>
);

export default Todo;
