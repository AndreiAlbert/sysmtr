import { CommonModule } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { SystemStats } from '../models/stats.model';
import { Subscription } from 'rxjs';
import { MonitorService } from '../app/monitor.service';

@Component({
  selector: 'app-stats-dashboard',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './stats-dashboard.component.html',
  styleUrl: './stats-dashboard.component.scss'
})
export class StatsDashboardComponent implements OnInit, OnDestroy {
  statsHistory: SystemStats[] = [];
  private wsSubscription?: Subscription;

  constructor(private monitorService: MonitorService) { }

  ngOnInit(): void {
    this.monitorService.getHistory().subscribe(data => {
      this.statsHistory = data;
    })

    this.wsSubscription = this.monitorService.getRealTimeStats().subscribe({
      next: (stat) => {
        this.statsHistory.unshift(stat);
      },
      error: (err) => console.error('Websocket error: ', err)
    });
  }

  ngOnDestroy(): void {
    this.wsSubscription?.unsubscribe();
  }

}
