import { Component, Input, inject } from '@angular/core';
import { Observable, filter, shareReplay } from 'rxjs';
import { PlatformService } from '../../services/platform.service';
import { Platform } from '../../gen/holos/v1alpha1/platform_pb';
import { MatTab, MatTabGroup } from '@angular/material/tabs';
import { AsyncPipe, CommonModule } from '@angular/common';
import { FormGroup, ReactiveFormsModule } from '@angular/forms';
import { FormlyMaterialModule } from '@ngx-formly/material';
import { FormlyModule } from '@ngx-formly/core';
import { MatButton } from '@angular/material/button';
import { MatDivider } from '@angular/material/divider';

@Component({
  selector: 'app-platform-detail',
  standalone: true,
  imports: [
    MatTabGroup,
    MatTab,
    AsyncPipe,
    CommonModule,
    FormlyModule,
    ReactiveFormsModule,
    FormlyMaterialModule,
    MatButton,
    MatDivider,
  ],
  templateUrl: './platform-detail.component.html',
  styleUrl: './platform-detail.component.scss'
})
export class PlatformDetailComponent {
  private service = inject(PlatformService);

  platform$!: Observable<Platform>;

  form = new FormGroup({});
  model = {};

  onSubmit(model: any) {
    console.log(model);
  }

  @Input()
  set id(platformId: string) {
    this.platform$ = this.service.getPlatform(platformId).pipe(
      filter((platform): platform is Platform => platform !== undefined),
      shareReplay(1)
    )
  }
}
