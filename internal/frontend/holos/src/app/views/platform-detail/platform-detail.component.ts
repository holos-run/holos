import { AsyncPipe, CommonModule } from '@angular/common';
import { Component, Input, OnDestroy, inject } from '@angular/core';
import { FormGroup, ReactiveFormsModule } from '@angular/forms';
import { MatButton } from '@angular/material/button';
import { MatDivider } from '@angular/material/divider';
import { MatTab, MatTabGroup } from '@angular/material/tabs';
import { JsonValue } from '@bufbuild/protobuf';
import { FormlyFieldConfig, FormlyFormOptions, FormlyModule } from '@ngx-formly/core';
import { FormlyMaterialModule } from '@ngx-formly/material';
import { Subject, takeUntil } from 'rxjs';
import { PlatformService } from '../../services/platform.service';

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

  private destroy$: Subject<boolean> = new Subject<boolean>();
  form = new FormGroup({});
  fields: FormlyFieldConfig[] = [];
  model: JsonValue = {};
  // Use form state to store the model for nested forms
  // Refer to https://formly.dev/docs/examples/form-options/form-state/
  options: FormlyFormOptions = {
    formState: {
      model: this.model,
    },
  };

  private setModel(model: JsonValue) {
    if (model) {
      this.model = model
      this.options.formState.model = model
    }
  }

  onSubmit(model: JsonValue) {
    if (this.form.valid) {
      console.log(model)
      this.platformService
        .putModel(this.platformId, model)
        .pipe(takeUntil(this.destroy$))
        .subscribe(resp => {
          if (resp.model !== undefined) {
            this.setModel(resp.model.toJson())
          }
        })
    }
  }

  @Input()
  set id(platformId: string) {
    this.platformId = platformId;
    this.platformService
      .getForm(platformId)
      .pipe(takeUntil(this.destroy$))
      .subscribe(resp => {
        if (resp.model !== undefined) {
          this.setModel(resp.model.toJson())
        }
        if (resp.fields !== undefined) {
          // NOTE: We could mix functions into the json data via mapped fields,
          // but consider carefully before doing so.  Refer to
          // https://formly.dev/docs/examples/other/json-powered
          this.fields = resp.fields.map(field => field.toJson() as FormlyFieldConfig)
        }
      })
  }

  public ngOnDestroy(): void {
    this.destroy$.next(true);
    this.destroy$.complete();
  }
}
