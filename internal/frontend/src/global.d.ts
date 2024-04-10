export {};

type AppConfig = {
  oidcIssuer: string;
};

declare global {
  interface Window {
    holosAppConfig: AppConfig;
  }
}
