import { Component, OnInit, ElementRef, ViewChild, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ApiService } from '../../core/services/api.service';
import { Stats } from '../../core/models';
import { Chart, registerables } from 'chart.js';
import { TranslateModule, TranslateService } from '@ngx-translate/core';

Chart.register(...registerables);

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  template: `
    <div class="page">
        <div class="page-header">
            <h1 class="page-title">{{ 'DASHBOARD.TITLE' | translate }}</h1>
            <p class="page-subtitle">{{ 'DASHBOARD.SUBTITLE' | translate }}</p>
        </div>

        @if (stats) {
            <div class="stats-grid">
                <div class="stat-card">
                    <h3>{{ 'DASHBOARD.TOTAL_DEPLOYS' | translate }}</h3>
                    <div class="value">{{ stats.total }}</div>
                </div>
                <div class="stat-card success">
                    <h3>{{ 'DASHBOARD.SUCCESSFUL' | translate }}</h3>
                    <div class="value">{{ stats.successful }}</div>
                </div>
                <div class="stat-card danger">
                    <h3>{{ 'DASHBOARD.FAILED' | translate }}</h3>
                    <div class="value">{{ stats.failed }}</div>
                </div>
                <div class="stat-card info">
                    <h3>{{ 'DASHBOARD.AVG_DURATION' | translate }}</h3>
                    <div class="value">{{ stats.averageDuration | number:'1.0-1' }}s</div>
                </div>
            </div>

            <div class="chart-card">
                <h3>{{ 'DASHBOARD.LAST_7_DAYS' | translate }}</h3>
                <div class="chart-container">
                    <canvas #chartCanvas></canvas>
                </div>
            </div>
        } @else {
            <div class="loading">{{ 'DASHBOARD.LOADING' | translate }}</div>
        }
    </div>
  `,
  styles: [`
    .stats-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 1.5rem;
        margin-bottom: 2rem;
    }
    .stat-card {
        background: var(--bg-card);
        padding: 1.5rem;
        border-radius: var(--radius-lg);
        border: 1px solid var(--border-color);
        box-shadow: 0 1px 3px rgba(0,0,0,0.1);
    }
    .stat-card h3 {
        font-size: 0.875rem;
        color: var(--text-muted);
        margin-bottom: 0.5rem;
        font-weight: 500;
    }
    .value { font-size: 2rem; font-weight: 700; color: var(--text-primary); }
    .stat-card.success .value { color: var(--success); }
    .stat-card.danger .value { color: var(--danger); }
    .stat-card.info .value { color: var(--info); }

    .chart-card {
        background: var(--bg-card);
        padding: 1.5rem;
        border-radius: var(--radius-lg);
        border: 1px solid var(--border-color);
        box-shadow: 0 1px 3px rgba(0,0,0,0.1);
    }
    .chart-container {
        position: relative;
        height: 300px;
        width: 100%;
        margin-top: 1rem;
    }
  `]
})
export class DashboardComponent implements OnInit {
  api = inject(ApiService);
  private translateService = inject(TranslateService);
  stats: Stats | null = null;

  @ViewChild('chartCanvas') chartCanvas!: ElementRef<HTMLCanvasElement>;
  chart: Chart | null = null;

  ngOnInit() {
    this.api.getStats().subscribe(data => {
      this.stats = data;
      setTimeout(() => this.initChart(), 0);
    });

    // Rebuild chart when language changes
    this.translateService.onLangChange.subscribe(() => {
      if (this.stats) {
        setTimeout(() => this.initChart(), 0);
      }
    });
  }

  initChart() {
    if (!this.stats || !this.chartCanvas) return;

    if (this.chart) this.chart.destroy();

    const ctx = this.chartCanvas.nativeElement.getContext('2d');
    if (!ctx) return;

    const gradient = ctx.createLinearGradient(0, 0, 0, 400);
    gradient.addColorStop(0, '#3b82f6');
    gradient.addColorStop(1, '#1e3a8a');

    this.chart = new Chart(ctx, {
      type: 'bar',
      data: {
        labels: this.stats.last7Days.dates,
        datasets: [{
          label: this.translateService.instant('DASHBOARD.DEPLOY_COUNT'),
          data: this.stats.last7Days.counts,
          backgroundColor: gradient,
          borderRadius: 6,
          borderSkipped: false,
          barThickness: 40
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: { display: false }
        },
        scales: {
          y: {
            beginAtZero: true,
            ticks: { stepSize: 1, color: '#9ca3af' },
            grid: { color: '#374151' }
          },
          x: {
            ticks: { color: '#9ca3af' },
            grid: { display: false }
          }
        }
      }
    });
  }
}
