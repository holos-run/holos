import { Component, Input, OnDestroy, inject } from '@angular/core';
import { Observable, Subject, shareReplay, takeUntil } from 'rxjs';
import { PlatformService } from '../../services/platform.service';
import { GetFormResponse } from '../../gen/holos/v1alpha1/platform_pb';
import { MatTab, MatTabGroup } from '@angular/material/tabs';
import { AsyncPipe, CommonModule } from '@angular/common';
import { FormGroup, ReactiveFormsModule } from '@angular/forms';
import { FormlyMaterialModule } from '@ngx-formly/material';
import { FormlyFieldConfig, FormlyFormOptions, FormlyModule } from '@ngx-formly/core';
import { MatButton } from '@angular/material/button';
import { MatDivider } from '@angular/material/divider';
import { JsonObject, JsonValue } from '@bufbuild/protobuf';

@Component({
  selector: 'app-platform-detail',
  standalone: true,
  imports: [
    AsyncPipe,
    CommonModule,
    FormlyMaterialModule,
    FormlyModule,
    MatButton,
    MatDivider,
    MatTab,
    MatTabGroup,
    ReactiveFormsModule,
  ],
  templateUrl: './platform-detail.component.html',
  styleUrl: './platform-detail.component.scss'
})
export class PlatformDetailComponent implements OnDestroy {
  private platformService = inject(PlatformService);
  private platformId: string = "";

  private destroy$: Subject<any> = new Subject<any>();
  form = new FormGroup({});
  model: JsonValue = {};
  options: FormlyFormOptions = {};
  fields: FormlyFieldConfig[] = [];

  onSubmit(model: JsonValue) {
    // if (this.form.valid) {
    console.log(model)
    this.platformService
      .putModel(this.platformId, model)
      .pipe(takeUntil(this.destroy$))
      .subscribe(resp => {
        const model = JSON.parse(JSON.stringify(resp.model))
        if (model) {
          this.model = model
        }
      })
    // }
  }

  @Input()
  set id(platformId: string) {
    this.platformId = platformId;
    this.platformService
      .getForm(platformId)
      .pipe(takeUntil(this.destroy$))
      .subscribe(resp => {
        if (resp.fields !== undefined) {
          // NOTE: We could map fields to mix in javascript functions.  Refer to
          // https://formly.dev/docs/examples/other/json-powered
          this.fields = resp.fields.map(field => field.toJson() as FormlyFieldConfig)
        }
        if (resp.model !== undefined) {
          this.model = resp.model.toJson()
        }
      })
  }

  public ngOnDestroy(): void {
    this.destroy$.next(true);
    this.destroy$.complete();
  }
}
