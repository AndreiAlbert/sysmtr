import { Component } from '@angular/core';
import { StatsDashboardComponent } from '../stats-dashboard/stats-dashboard.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [StatsDashboardComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  title = 'sysmtr-web';
}
