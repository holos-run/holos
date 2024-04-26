import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { ClusterListComponent } from './cluster-list/cluster-list.component';
import { ErrorNotFoundComponent } from './error-not-found/error-not-found.component';
import { PlatformConfigComponent } from './views/platform-config/platform-config.component';
import { AddressFormComponent } from './examples/address-form/address-form.component';

export const routes: Routes = [
  { path: 'platform-config', component: PlatformConfigComponent },
  { path: 'address-form', component: AddressFormComponent },
  { path: 'home', component: HomeComponent },
  { path: 'clusters', component: ClusterListComponent },
  { path: '', redirectTo: '/home', pathMatch: 'full' },
  { path: '**', component: ErrorNotFoundComponent },
];
