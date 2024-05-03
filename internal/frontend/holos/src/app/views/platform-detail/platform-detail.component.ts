import { Component, Input, inject } from '@angular/core';
import { Observable, map, shareReplay } from 'rxjs';
import { Model, PlatformService } from '../../services/platform.service';
import { Platform } from '../../gen/holos/v1alpha1/platform_pb';
import { MatTab, MatTabGroup } from '@angular/material/tabs';
import { AsyncPipe, CommonModule } from '@angular/common';
import { FormGroup, ReactiveFormsModule } from '@angular/forms';
import { FormlyMaterialModule } from '@ngx-formly/material';
import { FormlyFieldConfig, FormlyModule } from '@ngx-formly/core';
import { MatButton } from '@angular/material/button';
import { MatDivider } from '@angular/material/divider';
import { Value } from '@bufbuild/protobuf';

type Section = {
  name: string
  displayName: string
  description: string
  fields: FormlyFieldConfig[]
}

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
export class PlatformDetailComponent {
  private service = inject(PlatformService);
  private platformId: string = "";

  platform$!: Observable<Platform>;

  form = new FormGroup({});
  model: Model = {};
  sections: Section[] = [];

  onSubmit(model: Model) {
    if (this.form.valid) {
      this.service.putConfig(this.platformId, model).pipe(shareReplay(1)).subscribe()
    }
  }

  // rpcFieldsToFormlyFields adapts a protobuf Value[] to a typescript
  // FormlyFieldConfig[].  Necessary to ensure select fields get scalar values,
  // not google.protobuf.Value objects.  If objects are passed instead of
  // scalars you may receive `ERROR Error: `compareWith` must be a function.`
  // Refer to
  // https://github.com/ngx-formly/ngx-formly/issues/2408#issuecomment-666686211
  rpcFieldsToFormlyFields(fields: Value[]): FormlyFieldConfig[] {
    const ffcs: FormlyFieldConfig[] = []
    fields.forEach(field => { ffcs.push(JSON.parse(JSON.stringify(field))) })
    return ffcs
  }

  @Input()
  set id(platformId: string) {
    this.platformId = platformId;
    this.platform$ = this.service.getPlatform(platformId).pipe(
      map(project => {
        // Initialize the model container for each section of the form config
        const formSections = project.config?.form?.spec?.sections
        if (formSections !== undefined) {
          formSections.forEach(section => {
            // Map the pb FieldConfig[] to es FormlyFieldConfig[] so we get
            // scalar values, not objects as are provided by
            // google.protobuf.Value.
            const jsSection: Section = {
              name: section.name,
              displayName: section.displayName,
              description: section.description,
              fields: this.rpcFieldsToFormlyFields(section.fieldConfigs),
            }
            this.sections.push(jsSection)
            // Initialize the model for the section
            this.model[section.name] = {}
          })
        }

        // Upstream docs recommend loading values into the model instead of
        // using the defaultValue FormlyFormConfig property.
        const modelSections = project.config?.values?.sections
        if (modelSections !== undefined) {
          Object.keys(modelSections).forEach(sectionName => {
            Object.keys(modelSections[sectionName].fields).forEach(fieldName => {
              this.model[sectionName][fieldName] = modelSections[sectionName].fields[fieldName].toJson()
            })
          })
        }
        return project
      }),
      shareReplay(1)
    )
  }
}
