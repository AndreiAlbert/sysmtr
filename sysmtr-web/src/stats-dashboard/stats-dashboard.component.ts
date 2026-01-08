import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { SystemStats } from '../models/stats.model';
import { Subscription } from 'rxjs';
import { MonitorService } from '../app/monitor.service';
import { BaseChartDirective } from 'ng2-charts';
import { Chart, ChartConfiguration, ChartData, registerables } from 'chart.js';

Chart.register(...registerables)

@Component({
  selector: 'app-stats-dashboard',
  standalone: true,
  imports: [CommonModule, BaseChartDirective],
  templateUrl: './stats-dashboard.component.html',
  styleUrl: './stats-dashboard.component.scss'
})
export class StatsDashboardComponent implements OnInit, OnDestroy {
  @ViewChild(BaseChartDirective) chart?: BaseChartDirective;

  statsHistory: SystemStats[] = [];
  private wsSubscription?: Subscription;
  private maxDataPoints = 30;

  public lineChartData: ChartData<'line'> = {
    labels: [],
    datasets: [
      {
        data: [],
        label: 'Cpu Usage (%)',
        borderColor: '#38bdf8',
        tension: 0.1,
        fill: false
      },
      {
        data: [],
        label: "Ram Usage (%)",
        borderColor: '#fbbf24',
        tension: 0.1,
        fill: false
      }
    ]
  }

  public lineChartOptions: ChartConfiguration['options'] = {
    responsive: true,
    scales: {
      y: { min: 0, max: 100 }
    }
  }

  constructor(
    private monitorService: MonitorService,
    private cdr: ChangeDetectorRef
  ) { }

  ngOnInit(): void {
    this.monitorService.getHistory().subscribe(data => {
      this.statsHistory = data.slice(-this.maxDataPoints).reverse();
      this.updateChartData();
      this.cdr.detectChanges();
    })

    this.wsSubscription = this.monitorService.getRealTimeStats().subscribe({
      next: (stat) => {
        console.log(stat);
        this.statsHistory.push(stat);
        if (this.statsHistory.length > this.maxDataPoints) {
          this.statsHistory.shift();
        }
        this.updateChartData();
        this.cdr.detectChanges();
      },
      error: (err) => console.error('Websocket error: ', err)
    });
  }

  private updateChartData(): void {
    this.lineChartData.labels = this.statsHistory.map(s => s.created_at || '');
    this.lineChartData.datasets[0].data = this.statsHistory.map(s => s.cpu_usage);
    this.lineChartData.datasets[1].data = this.statsHistory.map(s => s.ram_usage);

    if (this.chart) {
      this.chart.update();
    }
  }

  ngOnDestroy(): void {
    this.wsSubscription?.unsubscribe();
  }

}
