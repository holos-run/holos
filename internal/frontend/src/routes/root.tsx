import React from "react";
import { Outlet } from "react-router-dom";
import {
  Avatar,
  Brand,
  Button,
  ButtonVariant,
  Divider,
  Dropdown,
  DropdownGroup,
  DropdownItem,
  DropdownList,
  Masthead,
  MastheadBrand,
  MastheadContent,
  MastheadMain,
  MastheadToggle,
  MenuToggle,
  MenuToggleElement,
  Page,
  PageSection,
  PageSectionVariants,
  PageSidebar,
  PageSidebarBody,
  PageToggleButton,
  SkipToContent,
  Text,
  TextContent,
  Toolbar,
  ToolbarContent,
  ToolbarGroup,
  ToolbarItem,
} from "@patternfly/react-core";
import BarsIcon from "@patternfly/react-icons/dist/esm/icons/bars-icon";
import BellIcon from "@patternfly/react-icons/dist/esm/icons/bell-icon";
import CogIcon from "@patternfly/react-icons/dist/esm/icons/cog-icon";
import HelpIcon from "@patternfly/react-icons/dist/esm/icons/help-icon";
import QuestionCircleIcon from "@patternfly/react-icons/dist/esm/icons/question-circle-icon";
import EllipsisVIcon from "@patternfly/react-icons/dist/esm/icons/ellipsis-v-icon";
import imgAvatar from "@patternfly/react-core/dist/styles/assets/images/img_avatar.svg";
import appBranding from "@public/logos/logo-light.svg";
import { nav } from "@app/routes/nav";

const Root: React.FunctionComponent = () => {
  const [isDropdownOpen, setIsDropdownOpen] = React.useState(false);
  const [isKebabDropdownOpen, setIsKebabDropdownOpen] = React.useState(false);
  const [isFullKebabDropdownOpen, setIsFullKebabDropdownOpen] =
    React.useState(false);

  const onDropdownToggle = () => {
    setIsDropdownOpen(!isDropdownOpen);
  };

  const onDropdownSelect = () => {
    setIsDropdownOpen(false);
  };

  const onKebabDropdownToggle = () => {
    setIsKebabDropdownOpen(!isKebabDropdownOpen);
  };

  const onKebabDropdownSelect = () => {
    setIsKebabDropdownOpen(false);
  };

  const onFullKebabDropdownToggle = () => {
    setIsFullKebabDropdownOpen(!isFullKebabDropdownOpen);
  };

  const onFullKebabDropdownSelect = () => {
    setIsFullKebabDropdownOpen(false);
  };

  const kebabDropdownItems = (
    <>
      <DropdownItem>
        <CogIcon /> Settings
      </DropdownItem>
      <DropdownItem>
        <HelpIcon /> Help
      </DropdownItem>
    </>
  );
  const userDropdownItems = (
    <>
      <DropdownItem key="group 2 profile">My profile</DropdownItem>
      <DropdownItem key="group 2 user">User management</DropdownItem>
      <DropdownItem key="group 2 logout">Logout</DropdownItem>
    </>
  );

  const headerToolbar = (
    <Toolbar id="toolbar" isFullHeight isStatic>
      <ToolbarContent>
        <ToolbarGroup
          variant="icon-button-group"
          align={{ default: "alignRight" }}
          spacer={{ default: "spacerNone", md: "spacerMd" }}
        >
          <ToolbarItem>
            <Button
              aria-label="Notifications"
              variant={ButtonVariant.plain}
              icon={<BellIcon />}
            />
          </ToolbarItem>
          <ToolbarGroup
            variant="icon-button-group"
            visibility={{ default: "hidden", lg: "visible" }}
          >
            <ToolbarItem>
              <Button
                aria-label="Settings"
                variant={ButtonVariant.plain}
                icon={<CogIcon />}
              />
            </ToolbarItem>
            <ToolbarItem>
              <Button
                aria-label="Help"
                variant={ButtonVariant.plain}
                icon={<QuestionCircleIcon />}
              />
            </ToolbarItem>
          </ToolbarGroup>
          <ToolbarItem
            visibility={{ default: "hidden", md: "visible", lg: "hidden" }}
          >
            <Dropdown
              isOpen={isKebabDropdownOpen}
              onSelect={onKebabDropdownSelect}
              onOpenChange={(isOpen: boolean) => {
                setIsKebabDropdownOpen(isOpen);
              }}
              popperProps={{ position: "right" }}
              toggle={(toggleRef: React.Ref<MenuToggleElement>) => (
                <MenuToggle
                  ref={toggleRef}
                  onClick={onKebabDropdownToggle}
                  isExpanded={isKebabDropdownOpen}
                  variant="plain"
                  aria-label="Settings and help"
                >
                  <EllipsisVIcon aria-hidden="true" />
                </MenuToggle>
              )}
            >
              <DropdownList>{kebabDropdownItems}</DropdownList>
            </Dropdown>
          </ToolbarItem>
          <ToolbarItem visibility={{ md: "hidden" }}>
            <Dropdown
              isOpen={isFullKebabDropdownOpen}
              onSelect={onFullKebabDropdownSelect}
              onOpenChange={(isOpen: boolean) => {
                setIsFullKebabDropdownOpen(isOpen);
              }}
              popperProps={{ position: "right" }}
              toggle={(toggleRef: React.Ref<MenuToggleElement>) => (
                <MenuToggle
                  ref={toggleRef}
                  onClick={onFullKebabDropdownToggle}
                  isExpanded={isFullKebabDropdownOpen}
                  variant="plain"
                  aria-label="Toolbar menu"
                >
                  <EllipsisVIcon aria-hidden="true" />
                </MenuToggle>
              )}
            >
              <DropdownGroup key="group 2" aria-label="User actions">
                <DropdownList>{userDropdownItems}</DropdownList>
              </DropdownGroup>
              <Divider />
              <DropdownList>{kebabDropdownItems}</DropdownList>
            </Dropdown>
          </ToolbarItem>
        </ToolbarGroup>
        <ToolbarItem visibility={{ default: "hidden", md: "visible" }}>
          <Dropdown
            isOpen={isDropdownOpen}
            onSelect={onDropdownSelect}
            onOpenChange={(isOpen: boolean) => {
              setIsDropdownOpen(isOpen);
            }}
            popperProps={{ position: "right" }}
            toggle={(toggleRef: React.Ref<MenuToggleElement>) => (
              <MenuToggle
                ref={toggleRef}
                onClick={onDropdownToggle}
                isFullHeight
                isExpanded={isDropdownOpen}
                icon={<Avatar src={imgAvatar} alt="" />}
              >
                FIXME
              </MenuToggle>
            )}
          >
            <DropdownList>{userDropdownItems}</DropdownList>
          </Dropdown>
        </ToolbarItem>
      </ToolbarContent>
    </Toolbar>
  );

  const masthead = (
    <Masthead>
      <MastheadToggle>
        <PageToggleButton variant="plain" aria-label="Global navigation">
          <BarsIcon />
        </PageToggleButton>
      </MastheadToggle>
      <MastheadMain>
        <MastheadBrand>
          <Brand src={appBranding} alt="Holos" heights={{ default: "36px" }} />
        </MastheadBrand>
      </MastheadMain>
      <MastheadContent>{headerToolbar}</MastheadContent>
    </Masthead>
  );

  const sidebar = (
    <PageSidebar>
      <PageSidebarBody>{nav.Navigation()}</PageSidebarBody>
    </PageSidebar>
  );

  const mainContainerId = "main-content";

  const pageSkipToContent = (
    <SkipToContent href={`#${mainContainerId}`}>Skip to content</SkipToContent>
  );

  return (
    <Page
      header={masthead}
      sidebar={sidebar}
      isManagedSidebar
      skipToContent={pageSkipToContent}
      mainContainerId={mainContainerId}
      isBreadcrumbWidthLimited
      isBreadcrumbGrouped
      additionalGroupedContent={
        <PageSection variant={PageSectionVariants.light} isWidthLimited>
          <TextContent>
            <Text component="h1">Holos Platform</Text>
            <Text component="p">Get started with holos</Text>
          </TextContent>
        </PageSection>
      }
      groupProps={{
        stickyOnBreakpoint: { default: "top" },
      }}
    >
      <PageSection variant={PageSectionVariants.light} isWidthLimited>
        <Outlet />
      </PageSection>
    </Page>
  );
};

export default Root;
