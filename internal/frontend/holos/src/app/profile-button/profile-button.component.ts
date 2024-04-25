import { Component, Input } from '@angular/core';
import { Claims } from '../gen/holos/v1alpha1/user_pb';
import { MatButtonModule } from '@angular/material/button';
import { MatMenuModule } from '@angular/material/menu';
import { Observable, map, takeLast } from 'rxjs';
import { AsyncPipe, NgIf, NgStyle } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { Organization } from '../gen/holos/v1alpha1/organization_pb';

@Component({
  selector: 'app-profile-button',
  standalone: true,
  imports: [
    MatButtonModule,
    MatMenuModule,
    MatIconModule,
    MatCardModule,
    AsyncPipe,
    NgIf,
    NgStyle,
  ],
  templateUrl: './profile-button.component.html',
  styleUrl: './profile-button.component.scss'
})
export class ProfileButtonComponent {
  @Input({ required: true }) claims$!: Observable<Claims | null>;
  @Input({ required: true }) organizations$!: Observable<Organization[] | null>;
}
