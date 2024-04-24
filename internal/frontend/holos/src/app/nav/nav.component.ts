import { Component, OnInit, inject } from '@angular/core';
import { BreakpointObserver, Breakpoints } from '@angular/cdk/layout';
import { AsyncPipe, NgIf } from '@angular/common';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { Observable } from 'rxjs';
import { map, shareReplay } from 'rxjs/operators';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { Claims } from '../gen/holos/v1alpha1/user_pb';
import { HolosUserService } from '../holos-user.service';
import { ProfileButtonComponent } from '../profile-button/profile-button.component';

@Component({
  selector: 'app-nav',
  templateUrl: './nav.component.html',
  styleUrl: './nav.component.scss',
  standalone: true,
  imports: [
    MatToolbarModule,
    MatButtonModule,
    MatSidenavModule,
    MatListModule,
    MatIconModule,
    NgIf,
    AsyncPipe,
    RouterLink,
    RouterLinkActive,
    RouterOutlet,
    MatCardModule,
    ProfileButtonComponent,
  ]
})
export class NavComponent implements OnInit {
  private breakpointObserver = inject(BreakpointObserver);
  private userService = inject(HolosUserService);

  claims$!: Observable<Claims | null>;

  isHandset$: Observable<boolean> = this.breakpointObserver.observe(Breakpoints.Handset)
    .pipe(
      map(result => result.matches),
      shareReplay()
    );

  ngOnInit(): void {
    this.claims$ = this.userService.getClaims();
  }
}
