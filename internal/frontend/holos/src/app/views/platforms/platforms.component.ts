import { Platform } from '../../gen/holos/platform/v1alpha1/platform_pb';
import { Component, OnInit, inject } from '@angular/core';
import { MatListItem, MatNavList } from '@angular/material/list';
import { Observable, filter } from 'rxjs';
import { Organization } from '../../gen/holos/organization/v1alpha1/organization_pb';
import { OrganizationService } from '../../services/organization.service';
import { PlatformService } from '../../services/platform.service';
import { AsyncPipe, CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-platforms',
  standalone: true,
  imports: [
    MatNavList,
    MatListItem,
    AsyncPipe,
    CommonModule,
    RouterLink,
  ],
  templateUrl: './platforms.component.html',
  styleUrl: './platforms.component.scss'
})
export class PlatformsComponent implements OnInit {
  private orgSvc = inject(OrganizationService);
  private platformSvc = inject(PlatformService);

  org$!: Observable<Organization | undefined>;
  platforms$!: Observable<Platform[]>;

  ngOnInit() {
    this.org$ = this.orgSvc.activeOrg();
    this.platforms$ = this.platformSvc.listPlatforms(this.org$.pipe(
      filter((org): org is Organization => org !== undefined)
    ))
  }
}
