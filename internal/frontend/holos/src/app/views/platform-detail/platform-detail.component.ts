import { Component, Input, inject } from '@angular/core';
import { Observable, map, shareReplay } from 'rxjs';
import { Model, PlatformService } from '../../services/platform.service';
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
  private platformId: string = "";

  platform$!: Observable<Platform>;

  form = new FormGroup({});
  model: Model = {};

  onSubmit(model: Model) {
    console.log(model)
    // if (this.form.valid) {
    this.service.putConfig(this.platformId, model).pipe(shareReplay(1)).subscribe()
    // }
  }

  @Input()
  set id(platformId: string) {
    this.platformId = platformId;
    this.platform$ = this.service.getPlatform(platformId).pipe(
      map(project => {
        // Initialize the model container for each section of the form config
        project.config?.form?.spec?.sections.forEach(section => {
          this.model[section.name] = {}
        })
        // Load existing values into the form
        const sections = project.config?.values?.sections
        if (sections !== undefined) {
          Object.keys(sections).forEach(sectionName => {
            Object.keys(sections[sectionName].fields).forEach(fieldName => {
              this.model[sectionName][fieldName] = sections[sectionName].fields[fieldName].toJson()
            })
          })
        }
        return project
      }),
      shareReplay(1)
    )
  }
}
