import { Component, OnInit, OnDestroy, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { ApiService } from '../../core/services/api.service';
import { WebSocketService } from '../../core/services/websocket.service';
import { Deployment, Project, Server, LogMessage } from '../../core/models';
import { TerminalComponent } from '../../shared/components/terminal/terminal.component';
import { Subscription } from 'rxjs';
import { TranslateModule, TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-deployments',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, TerminalComponent, TranslateModule],
  template: `
    <div class="page">
      <div class="page-header">
        <div class="page-header-content">
          <h1 class="page-title">{{ 'DEPLOYMENTS.TITLE' | translate }}</h1>
          <p class="page-subtitle">{{ 'DEPLOYMENTS.SUBTITLE' | translate }}</p>
        </div>
        <button class="btn btn-primary" (click)="showForm = !showForm">
          {{ showForm ? ('DEPLOYMENTS.CANCEL' | translate) : ('DEPLOYMENTS.NEW_DEPLOY' | translate) }}
        </button>
      </div>

      @if (showForm) {
        <div class="form-card">
          <h3 class="form-title">{{ 'DEPLOYMENTS.START_NEW' | translate }}</h3>
          <form [formGroup]="form" (ngSubmit)="onSubmit()">
            <div class="form-grid">
              <div class="form-group">
                <label class="form-label">{{ 'DEPLOYMENTS.PROJECT' | translate }}</label>
                <select formControlName="projectId" class="form-input">
                  <option value="">{{ 'DEPLOYMENTS.SELECT_PROJECT' | translate }}</option>
                  @for (project of projects; track project.id) {
                    <option [value]="project.id">{{ project.name }}</option>
                  }
                </select>
              </div>
              <div class="form-group">
                <label class="form-label">{{ 'DEPLOYMENTS.SERVER' | translate }}</label>
                <select formControlName="serverId" class="form-input">
                  <option value="">{{ 'DEPLOYMENTS.SELECT_SERVER' | translate }}</option>
                  @for (server of servers; track server.id) {
                    <option [value]="server.id">{{ server.name }} ({{ server.host }})</option>
                  }
                </select>
              </div>
            </div>
            <div class="form-actions">
              <button type="submit" class="btn btn-primary" [disabled]="form.invalid || loading">
                @if (loading) { <span class="spinner"></span> }
                {{ 'DEPLOYMENTS.START_DEPLOY' | translate }}
              </button>
            </div>
          </form>
        </div>
      }

      @if (selectedDeployment) {
        <div class="detail-card">
          <div class="detail-header">
            <div class="detail-info">
              <span class="detail-id">{{ selectedDeployment.id.slice(0, 8) }}</span>
              <span class="badge" [class]="'badge-' + getStatusClass(selectedDeployment.status)">
                {{ getStatusText(selectedDeployment.status) }}
              </span>
            </div>
            <button class="btn btn-secondary" (click)="closeDetail()">{{ 'DEPLOYMENTS.CLOSE' | translate }}</button>
          </div>
          <app-terminal [logs]="terminalLogs" [title]="'Deployment Logs'"></app-terminal>
        </div>
      }

      <div class="deployments-list">
        @for (deployment of deployments; track deployment.id) {
          <div class="deployment-card" (click)="selectDeployment(deployment)" 
               [class.active]="selectedDeployment?.id === deployment.id">
            <span class="status-dot" [class]="deployment.status"></span>
            <div class="deployment-info">
              <span class="deployment-id">{{ deployment.id.slice(0, 12) }}</span>
              <span class="deployment-time">{{ deployment.createdAt | date:'medium' }}</span>
            </div>
            <span class="badge" [class]="'badge-' + getStatusClass(deployment.status)">
              {{ getStatusText(deployment.status) }}
            </span>
            @if (deployment.status === 'success' && deployment.commitHash) {
              <button class="btn-icon rollback-btn" (click)="$event.stopPropagation(); rollback(deployment)" title="Rollback to this version">
                ↺
              </button>
            }
          </div>
        } @empty {
          <div class="empty-state">
            <span class="empty-state-icon">▷</span>
            <h3>{{ 'DEPLOYMENTS.EMPTY_TITLE' | translate }}</h3>
            <p>{{ 'DEPLOYMENTS.EMPTY_DESC' | translate }}</p>
          </div>
        }
      </div>
    </div>
  `,
  styles: [`
    .form-card, .detail-card {
      background-color: var(--bg-card);
      border: 1px solid var(--border-color);
      border-radius: var(--radius-lg);
      padding: 1.5rem;
      margin-bottom: 1.5rem;
    }

    .form-title {
      font-size: 1rem;
      font-weight: 600;
      margin-bottom: 1.25rem;
    }

    .form-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 1rem;
    }

    .form-actions {
      margin-top: 1.25rem;
      display: flex;
      justify-content: flex-end;
    }

    .detail-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 1rem;
    }

    .detail-info {
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }

    .detail-id {
      font-family: monospace;
      font-size: 0.875rem;
      font-weight: 600;
    }

    .deployments-list {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .deployment-card {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      padding: 1rem 1.25rem;
      background-color: var(--bg-card);
      border: 1px solid var(--border-color);
      border-radius: var(--radius-lg);
      cursor: pointer;
      transition: all var(--transition-fast);
    }

    .deployment-card:hover {
      border-color: var(--accent);
    }

    .deployment-card.active {
      border-color: var(--accent);
      background-color: var(--accent-light);
    }

    .deployment-info {
      flex: 1;
      display: flex;
      flex-direction: column;
      gap: 0.125rem;
    }

    .deployment-id {
      font-size: 0.875rem;
      font-weight: 500;
      font-family: monospace;
      color: var(--text-primary);
    }

    .deployment-time {
      font-size: 0.75rem;
      color: var(--text-muted);
    }

    .rollback-btn {
        margin-left: auto;
        color: var(--text-muted);
        font-size: 1.25rem;
        padding: 0.25rem 0.5rem;
    }
    .rollback-btn:hover {
        color: var(--accent);
        background: var(--bg-tertiary);
        border-radius: var(--radius-sm);
    }
  `]
})
export class DeploymentsComponent implements OnInit, OnDestroy {
  private api = inject(ApiService);
  private ws = inject(WebSocketService);
  private fb = inject(FormBuilder);
  private translateService = inject(TranslateService);

  deployments: Deployment[] = [];
  projects: Project[] = [];
  servers: Server[] = [];
  selectedDeployment: Deployment | null = null;
  terminalLogs = '';
  showForm = false;
  loading = false;
  private wsSubscription: Subscription | null = null;
  private pollInterval: ReturnType<typeof setInterval> | null = null;
  private listPollInterval: ReturnType<typeof setInterval> | null = null;

  form = this.fb.group({
    projectId: ['', Validators.required],
    serverId: ['', Validators.required]
  });

  ngOnInit(): void {
    this.loadData();
    this.startListPolling();
  }

  ngOnDestroy(): void {
    this.wsSubscription?.unsubscribe();
    this.ws.disconnect();
    if (this.pollInterval) clearInterval(this.pollInterval);
    if (this.listPollInterval) clearInterval(this.listPollInterval);
  }

  loadData(): void {
    this.api.getDeployments().subscribe(data => this.deployments = data || []);
    this.api.getProjects().subscribe(data => this.projects = data || []);
    this.api.getServers().subscribe(data => this.servers = data || []);
  }

  selectDeployment(deployment: Deployment): void {
    if (this.selectedDeployment?.id === deployment.id) return;

    this.wsSubscription?.unsubscribe();
    this.ws.disconnect();

    this.selectedDeployment = deployment;
    this.terminalLogs = deployment.logs || '';

    if (deployment.status === 'running' || deployment.status === 'queued') {
      this.wsSubscription = this.ws.connect(deployment.id).subscribe({
        next: (msg: LogMessage) => {
          this.terminalLogs += msg.message;
        }
      });
      this.pollDeployment(deployment.id);
    }
  }

  private pollDeployment(id: string): void {
    if (this.pollInterval) {
      clearInterval(this.pollInterval);
    }

    this.pollInterval = setInterval(() => {
      this.api.getDeployment(id).subscribe(d => {
        if (this.selectedDeployment?.id === id) {
          // Update status but preserve WebSocket-appended logs
          const previousStatus = this.selectedDeployment.status;
          this.selectedDeployment = { ...d, logs: this.terminalLogs || d.logs };

          // Only sync logs from API if WebSocket hasn't been updating
          if (!this.terminalLogs || this.terminalLogs.length < (d.logs?.length || 0)) {
            this.terminalLogs = d.logs || '';
          }

          if (d.status === 'success' || d.status === 'failed') {
            if (this.pollInterval) clearInterval(this.pollInterval);
            this.loadData();
            // Final sync of logs when deployment completes
            this.terminalLogs = d.logs || '';
          }
        } else {
          if (this.pollInterval) clearInterval(this.pollInterval);
        }
      });
    }, 2000);
  }

  private startListPolling(): void {
    if (this.listPollInterval) {
      clearInterval(this.listPollInterval);
    }

    this.listPollInterval = setInterval(() => {
      this.api.getDeployments().subscribe(data => {
        this.deployments = data || [];
      });
    }, 3000);
  }

  closeDetail(): void {
    this.wsSubscription?.unsubscribe();
    this.ws.disconnect();
    this.selectedDeployment = null;
    this.terminalLogs = '';
  }

  onSubmit(): void {
    if (this.form.invalid) return;
    this.loading = true;

    this.api.createDeployment(this.form.value as any).subscribe({
      next: (deployment) => {
        this.deployments.unshift(deployment);
        this.form.reset();
        this.showForm = false;
        this.loading = false;
        this.selectDeployment(deployment);
      },
      error: () => this.loading = false
    });
  }

  rollback(d: Deployment): void {
    if (!confirm(this.translateService.instant('DEPLOYMENTS.ROLLBACK_CONFIRM', { hash: d.commitHash?.slice(0, 7) }))) return;

    this.api.rollbackDeployment(d.id).subscribe(() => {
      // Toast will be shown by interceptor if success message returned, or we can just reload
      this.loadData();
    });
  }

  getStatusText(status: string): string {
    const map: Record<string, string> = {
      'queued': this.translateService.instant('DEPLOYMENTS.STATUS_QUEUED'),
      'running': this.translateService.instant('DEPLOYMENTS.STATUS_RUNNING'),
      'success': this.translateService.instant('DEPLOYMENTS.STATUS_SUCCESS'),
      'failed': this.translateService.instant('DEPLOYMENTS.STATUS_FAILED')
    };
    return map[status] || status;
  }

  getStatusClass(status: string): string {
    const map: Record<string, string> = {
      'queued': 'warning',
      'running': 'info',
      'success': 'success',
      'failed': 'danger'
    };
    return map[status] || 'neutral';
  }
}
