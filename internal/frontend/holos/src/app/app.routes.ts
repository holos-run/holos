import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { ErrorNotFoundComponent } from './error-not-found/error-not-found.component';
import { PlatformsComponent } from './views/platforms/platforms.component'
import { PlatformDetailComponent } from './views/platform-detail/platform-detail.component';

export const routes: Routes = [
  { path: 'platform/:id', component: PlatformDetailComponent },
  { path: 'platforms', component: PlatformsComponent },
  { path: 'home', component: HomeComponent },
  { path: '', redirectTo: '/platforms', pathMatch: 'full' },
  { path: '**', component: ErrorNotFoundComponent },
];
