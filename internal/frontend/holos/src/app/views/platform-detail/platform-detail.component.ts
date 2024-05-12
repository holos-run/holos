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
import { Platform } from '../../gen/holos/platform/v1alpha1/platform_pb';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';

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
  private _snackBar: MatSnackBar = inject(MatSnackBar);
  private platformId: string = "";
  private platform: Platform = new Platform();
  isLoading = false;

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
      this.isLoading = true;
      this.platformService
        .updateModel(this.platform.id, model)
        .pipe(takeUntil(this.destroy$))
        .subscribe(platform => {
          this._snackBar.open('Saved ' + platform?.displayName, 'OK', {
            duration: 5000,
          }).afterDismissed().subscribe(() => {
            this.isLoading = false;
          })
        })
    }
  }

  @Input()
  set id(platformId: string) {
    this.platformId = platformId;
    this.platformService
      .getPlatform(platformId)
      .pipe(takeUntil(this.destroy$))
      .subscribe(platform => {
        if (platform !== undefined) {
          this.platform = platform
        }
        if (platform?.spec?.model !== undefined) {
          this.setModel(platform.spec.model.toJson())
        }
        if (platform?.spec?.form !== undefined) {
          // NOTE: We could mix functions into the json data via mapped fields,
          // but consider carefully before doing so.  Refer to
          // https://formly.dev/docs/examples/other/json-powered
          this.fields = platform.spec.form.fieldConfigs.map(fieldConfig => fieldConfig.toJson() as FormlyFieldConfig)
        }
      })
  }

  public ngOnDestroy(): void {
    this.destroy$.next(true);
    this.destroy$.complete();
  }
}
