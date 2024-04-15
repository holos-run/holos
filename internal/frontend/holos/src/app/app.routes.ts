import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { ClusterListComponent } from './cluster-list/cluster-list.component';
import { ErrorNotFoundComponent } from './error-not-found/error-not-found.component';

export const routes: Routes = [
  { path: 'home', component: HomeComponent },
  { path: 'clusters', component: ClusterListComponent },
  { path: '', redirectTo: '/home', pathMatch: 'full' },
  { path: '**', component: ErrorNotFoundComponent },
];
