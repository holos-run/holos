import { Component } from '@angular/core';
import { MatGridListModule } from '@angular/material/grid-list';

@Component({
  selector: 'app-error-not-found',
  standalone: true,
  imports: [MatGridListModule],
  templateUrl: './error-not-found.component.html',
  styleUrl: './error-not-found.component.scss'
})
export class ErrorNotFoundComponent {

}
