import { Component } from '@angular/core';
import { FormGroup, ReactiveFormsModule } from '@angular/forms';
import { MatTabsModule } from '@angular/material/tabs';
import { MatButton } from '@angular/material/button';
import { MatCard, MatCardActions, MatCardContent, MatCardHeader, MatCardTitle } from '@angular/material/card';
import { FormlyModule, FormlyFieldConfig } from '@ngx-formly/core';
import { FormlyMaterialModule } from '@ngx-formly/material';

@Component({
  selector: 'app-platform-config',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    FormlyMaterialModule,
    FormlyModule,
    MatTabsModule,
    MatCard,
    MatCardHeader,
    MatCardTitle,
    MatCardContent,
    MatCardActions,
    MatButton,
  ],
  templateUrl: './platform-config.component.html',
  styleUrl: './platform-config.component.scss'
})
export class PlatformConfigComponent {
  form = new FormGroup({});
  model: any = {};
  fields: FormlyFieldConfig[] = [
    {
      key: 'name',
      type: 'input',
      props: {
        label: 'Name',
        placeholder: 'example',
        required: true,
        description: "DNS label, e.g. 'example'"
      }
    },
    {
      key: 'domain',
      type: 'input',
      props: {
        label: 'Domain',
        placeholder: 'example.com',
        required: true,
        description: "DNS domain, e.g. 'example.com'"
      }
    },
    {
      key: 'displayName',
      type: 'input',
      props: {
        label: 'Display Name',
        placeholder: 'Example Organization',
        required: true,
        description: "Display name, e.g. 'My Organization'"
      }
    },
    {
      key: 'contactEmail',
      type: 'input',
      props: {
        label: 'Contact Email',
        placeholder: '',
        required: true,
        description: "Organization technical contact."
      }
    },
  ];

  integrationFields: FormlyFieldConfig[] = [
    {
      key: 'cloudflareEmail',
      type: 'input',
      props: {
        label: 'Cloudflare Account',
        placeholder: 'example@example.com',
        required: true,
        description: "Cloudflare account email address."
      }
    },
    {
      key: 'githubPrimaryOrg',
      type: 'input',
      props: {
        label: 'Github Organization',
        placeholder: 'ExampleOrg',
        required: true,
        description: "Github organization, e.g. 'ExampleOrg'"
      }
    }
  ];

  provisionerFields: FormlyFieldConfig[] = [
    {
      key: 'provisionerCABundle',
      type: 'input',
      props: {
        label: 'Provisioner API CA Bundle',
        placeholder: 'LS0tLS1CRUdJTiBDRVJUSUZJQXXXXXXXXXXXXXXXXXXXXXXX',
        required: true,
        description: "kubectl config view --minify --flatten -ojsonpath='{.clusters[0].cluster.certificate-authority-data}'"
      }
    },
    {
      key: 'provisionerURL',
      type: 'input',
      props: {
        label: 'Provisioner API URL',
        placeholder: 'https://1.2.3.4',
        required: true,
        description: "kubectl config view --minify --flatten -ojsonpath='{.clusters[0].cluster.server}'"
      }
    }
  ]

  onSubmit(model: any) {
    console.log(model);
  }
}
