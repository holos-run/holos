import { ApplicationConfig, importProvidersFrom } from '@angular/core';
import { provideRouter, withComponentInputBinding } from '@angular/router';
import { FormlyModule } from '@ngx-formly/core';
// import { provideHttpClient, withFetch } from '@angular/common/http';
import { routes } from './app.routes';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { ConnectModule } from '../connect/connect.module';
import { provideClient } from "../connect/client.provider";
import { UserService } from './gen/holos/v1alpha1/user_connect';
import { OrganizationService } from './gen/holos/v1alpha1/organization_connect';
import { PlatformService } from './gen/holos/v1alpha1/platform_connect';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes, withComponentInputBinding()),
    provideAnimationsAsync(),
    // provideHttpClient(withFetch()),
    provideClient(UserService),
    provideClient(OrganizationService),
    provideClient(PlatformService),
    importProvidersFrom(
      ConnectModule.forRoot({
        baseUrl: window.location.origin
      }),
      FormlyModule.forRoot(),
    ),
  ]
};
