import React from "react";
import {
  List,
  ListComponent,
  ListItem,
  OrderType,
  Wizard,
  WizardStep,
} from "@patternfly/react-core";

export const StartPage: React.FunctionComponent = () => (
  <Wizard title="Get Started">
    <WizardStep name="Overview" id="overview">
      <OverviewList />
    </WizardStep>
    <WizardStep name="User Registration" id="user-registration">
      Intentionally blank
    </WizardStep>
    <WizardStep
      name="Review"
      id="basic-review-step"
      footer={{ nextButtonText: "Finish" }}
    >
      Review step content
    </WizardStep>
  </Wizard>
);

const OverviewList: React.FunctionComponent = () => (
  <List component={ListComponent.ol} type={OrderType.number}>
    <ListItem>Review this overview.</ListItem>
    <ListItem>Register your user account.</ListItem>
    <ListItem>
      Create a new organization. Organizations have one or more projects.
    </ListItem>
    <ListItem>
      Link the organization to Github so Holos can create and manage the
      organization infrastructure code repository.
    </ListItem>
    <ListItem>
      Create a new project. Projects are the unit of isolation for multi
      tenancy.
    </ListItem>
    <ListItem>Create a new application in the project from a starter.</ListItem>
    <ListItem>Create a new cluster to deploy the application into.</ListItem>
    <ListItem>Select a Platform to deploy onto the Cluster.</ListItem>
    <ListItem>Personalize the platform with your domain name.</ListItem>
    <ListItem>Render the cluster configuration yaml into Github.</ListItem>
    <ListItem>Apply the configuration to the cluster.</ListItem>
    <ListItem>
      Link your application Github repository for Holos to deploy.
    </ListItem>
    <ListItem>Render the application starter repository.</ListItem>
    <ListItem>
      Validate Github actions run successfully for your application.
    </ListItem>
    <ListItem>
      Validate production deployment of the application into your new platform.
    </ListItem>
  </List>
);
