import { AsyncPipe, NgIf, NgStyle } from '@angular/common';
import { Component, OnDestroy, OnInit, inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { Observable, Subject, of, startWith, switchMap, takeUntil } from 'rxjs';
import { Version } from '../gen/holos/system/v1alpha1/system_pb';
import { SystemService } from '../services/system.service';
import { TruncatePipe } from '../truncate.pipe';
import { MatDivider } from '@angular/material/divider';

@Component({
  selector: 'app-version-button',
  standalone: true,
  imports: [
    AsyncPipe,
    MatButtonModule,
    MatCardModule,
    MatDivider,
    MatIconModule,
    MatMenuModule,
    NgIf,
    NgStyle,
    TruncatePipe,
  ],
  templateUrl: './version-button.component.html',
  styleUrl: './version-button.component.scss'
})
export class VersionButtonComponent implements OnInit, OnDestroy {
  private destroy$: Subject<boolean> = new Subject<boolean>();
  private refreshVersion$ = new Subject<boolean>();
  private systemService = inject(SystemService);
  version$!: Observable<Version | undefined>;
  isLoading = false;

  refreshVersion(): void {
    this.refreshVersion$.next(true);
  }

  ngOnInit(): void {
    this.version$ = this.refreshVersion$.pipe(
      takeUntil(this.destroy$),
      startWith(true),
      switchMap(() => {
        this.isLoading = true;
        return this.systemService.getVersion().pipe(
          switchMap((version) => { this.isLoading = false; return of(version); })
        );
      }),
    )
  }

  public ngOnDestroy(): void {
    this.destroy$.next(true);
    this.destroy$.complete();
  }
}
