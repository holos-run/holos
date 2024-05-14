import { BreakpointObserver, Breakpoints } from '@angular/cdk/layout';
import { AsyncPipe, NgIf } from '@angular/common';
import { Component, OnDestroy, OnInit, inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatListModule } from '@angular/material/list';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatToolbarModule } from '@angular/material/toolbar';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { Observable, Subject } from 'rxjs';
import { map, shareReplay, takeUntil } from 'rxjs/operators';
import { Organization } from '../gen/holos/organization/v1alpha1/organization_pb';
import { User } from '../gen/holos/user/v1alpha1/user_pb';
import { ProfileButtonComponent } from '../profile-button/profile-button.component';
import { OrganizationService } from '../services/organization.service';
import { UserService } from '../services/user.service';
import { VersionButtonComponent } from '../version-button/version-button.component';

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
    VersionButtonComponent,
  ]
})
export class NavComponent implements OnInit, OnDestroy {
  private breakpointObserver = inject(BreakpointObserver);
  private userService = inject(UserService);
  private orgService = inject(OrganizationService);
  private destroy$: Subject<boolean> = new Subject<boolean>();

  user$!: Observable<User | null>;
  org$!: Observable<Organization | undefined>;

  isHandset$: Observable<boolean> = this.breakpointObserver.observe(Breakpoints.Handset)
    .pipe(
      map(result => result.matches),
      shareReplay()
    );

  refreshOrg(): void {
    this.orgService.refreshOrganizations()
  }

  ngOnInit(): void {
    this.user$ = this.userService.getUser().pipe(takeUntil(this.destroy$));
    this.org$ = this.orgService.activeOrg().pipe(takeUntil(this.destroy$));
  }

  public ngOnDestroy(): void {
    this.destroy$.next(true);
    this.destroy$.complete();
  }
}
