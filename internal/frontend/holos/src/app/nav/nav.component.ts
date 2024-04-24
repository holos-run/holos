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
import { GetCallerClaimsResponse, User } from '../gen/holos/v1alpha1/user_pb';
import { HolosUserService } from '../holos-user.service';
import { ConnectModule } from "../../connect/connect.module";

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
  ]
})
export class NavComponent implements OnInit {
  private breakpointObserver = inject(BreakpointObserver);
  private userService = inject(HolosUserService);

  user$!: Observable<GetCallerClaimsResponse>;

  isHandset$: Observable<boolean> = this.breakpointObserver.observe(Breakpoints.Handset)
    .pipe(
      map(result => result.matches),
      shareReplay()
    );

  ngOnInit(): void {
    this.user$ = this.userService.getUser();
  }
}
