import { Component, Input, inject } from '@angular/core';
import { Observable, filter, map, shareReplay } from 'rxjs';
import { PlatformService } from '../../services/platform.service';
import { ConfigValues, ConfigValuesAny, ConfigValuesSection, ConfigValuesSectionAny, Platform } from '../../gen/holos/v1alpha1/platform_pb';
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
  private platformId: string = "";

  platform$!: Observable<Platform>;

  form = new FormGroup({});
  model = new ConfigValuesAny({});

  onSubmit(model: ConfigValuesAny) {
    if (this.form.valid) {
      const jsonEncodedValues = new ConfigValues({})
      const sections = Object.entries(this.model.sections)
      for (const [sectionName, section] of sections) {
        jsonEncodedValues.sections[sectionName] = new ConfigValuesSection({})
        const fields = Object.entries(section.values)
        for (const [fieldName, value] of fields) {
          jsonEncodedValues.sections[sectionName].values[fieldName] = JSON.stringify(value)
        }
      }
      // Make the RPC call to put the values.
      this.service.putConfig(this.platformId, jsonEncodedValues).pipe(shareReplay(1)).subscribe()
    }
  }

  @Input()
  set id(platformId: string) {
    this.platformId = platformId;
    this.platform$ = this.service.getPlatform(platformId).pipe(
      map(project => {
        // Initialize the model container for each section of the form config
        project.config?.form?.spec?.sections.forEach(section => {
          this.model.sections[section.name] = new ConfigValuesSectionAny({})
        })
        // Copy current values into the model
        if (project.config?.values?.sections !== undefined) {
          const sections = Object.entries(project.config.values.sections)
          for (const [sectionName, section] of sections) {
            const fields = Object.entries(section.values)
            for (const [fieldName, value] of fields) {
              this.model.sections[sectionName].values[fieldName] = JSON.parse(value)
            }
          }
        }

        return project
      }),
      shareReplay(1)
    )
  }

}
