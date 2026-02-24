import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { ApiService } from '../../core/services/api.service';
import { Server } from '../../core/models';
import { TranslateModule, TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-servers',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, TranslateModule],
  template: `
    <div class="page">
      <div class="page-header">
        <div class="page-header-content">
          <h1 class="page-title">{{ 'SERVERS.TITLE' | translate }}</h1>
          <p class="page-subtitle">{{ 'SERVERS.SUBTITLE' | translate }}</p>
        </div>
        <button class="btn btn-primary" (click)="cancelEdit()">
          {{ showForm ? ('SERVERS.CANCEL' | translate) : ('SERVERS.NEW_SERVER' | translate) }}
        </button>
      </div>

      @if (showForm) {
        <div class="form-card">
          <h3 class="form-title">{{ selectedServerId ? ('SERVERS.EDIT_SERVER' | translate) : ('SERVERS.ADD_SERVER' | translate) }}</h3>
          <form [formGroup]="form" (ngSubmit)="onSubmit()">
            <div class="form-grid">
              <div class="form-group">
                <label class="form-label">{{ 'SERVERS.SERVER_NAME' | translate }}</label>
                <input type="text" formControlName="name" class="form-input" placeholder="Production Server" />
              </div>
              <div class="form-group">
                <label class="form-label">{{ 'SERVERS.HOST' | translate }}</label>
                <input type="text" formControlName="host" class="form-input" placeholder="192.168.1.100" />
              </div>
              <div class="form-group">
                <label class="form-label">{{ 'SERVERS.PORT' | translate }}</label>
                <input type="number" formControlName="port" class="form-input" placeholder="22" />
              </div>
              <div class="form-group">
                <label class="form-label">{{ 'SERVERS.USERNAME' | translate }}</label>
                <input type="text" formControlName="username" class="form-input" placeholder="root" />
              </div>
              <div class="form-group full-width">
                <label class="form-label">{{ 'SERVERS.SSH_KEY' | translate }}</label>
                <textarea formControlName="sshKey" class="form-input form-textarea mono" rows="5" 
                  placeholder="Private Key or Password"></textarea>
              </div>
              
              <!-- Test Connection Button in Form -->
              <div class="form-group full-width">
                <button type="button" class="btn btn-secondary" (click)="testFormConnection()" 
                  [disabled]="!form.get('host')?.value || !form.get('username')?.value || !form.get('sshKey')?.value || testingConnection">
                  @if (testingConnection) { <span class="spinner"></span> }
                   {{ 'SERVERS.TEST_CONNECTION' | translate }}
                </button>
                
                @if (connectionTestResult) {
                  <div class="connection-result" [class.success]="connectionTestResult.success" [class.error]="!connectionTestResult.success">
                    <span class="result-icon">{{ connectionTestResult.success ? '✅' : '❌' }}</span>
                    <span class="result-message">{{ connectionTestResult.message }}</span>
                    @if (connectionTestResult.success) {
                      <span class="result-latency">{{ connectionTestResult.latency }}ms</span>
                    }
                  </div>
                }
              </div>
            </div>
            <div class="form-actions">
              <button type="submit" class="btn btn-primary" [disabled]="form.invalid && !selectedServerId || loading">
                @if (loading) { <span class="spinner"></span> }
                {{ selectedServerId ? ('SERVERS.UPDATE' | translate) : ('SERVERS.ADD' | translate) }}
              </button>
            </div>
          </form>
        </div>
      }

      <div class="servers-list">
        @for (server of servers; track server.id) {
          <div class="server-card">
            <div class="server-info">
              <div class="server-name">{{ server.name }}</div>
              <div class="server-host">{{ server.username }}&#64;{{ server.host }}:{{ server.port }}</div>
            </div>
            <div class="server-status">
              @if (testingServerId === server.id) {
                <span class="spinner"></span>
              } @else if (serverTestResults[server.id]) {
                <span class="status-dot" [class.success]="serverTestResults[server.id].success" [class.error]="!serverTestResults[server.id].success"></span>
                <span class="status-text" [class.error-text]="!serverTestResults[server.id].success">
                  {{ serverTestResults[server.id].success ? ('SERVERS.CONNECTED' | translate) + ' (' + serverTestResults[server.id].latency + 'ms)' : ('SERVERS.CONNECTION_ERROR' | translate) }}
                </span>
              } @else {
                <span class="status-dot neutral"></span>
                <span class="status-text neutral-text">{{ 'SERVERS.NOT_TESTED' | translate }}</span>
              }
            </div>
            <button type="button" class="btn-icon" [title]="'SERVERS.TEST_CONNECTION' | translate" (click)="testServerConnection(server)" [disabled]="testingServerId === server.id">🔌</button>
            <button class="btn-icon" [title]="'SERVERS.SSH_TERMINAL' | translate" (click)="openTerminal(server)">💻</button>
            <button class="btn-icon" (click)="editServer(server)">✎</button>
          </div>
        } @empty {
          <div class="empty-state">
            <span class="empty-state-icon">▤</span>
            <h3>{{ 'SERVERS.EMPTY_TITLE' | translate }}</h3>
            <p>{{ 'SERVERS.EMPTY_DESC' | translate }}</p>
          </div>
        }
      </div>
    </div>
  `,
  styles: [`
    .form-card {
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

    .full-width {
      grid-column: span 2;
    }

    .form-textarea.mono {
      font-family: monospace;
      font-size: 0.8rem;
    }

    .form-actions {
      margin-top: 1.25rem;
      display: flex;
      justify-content: flex-end;
    }

    .servers-list {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .server-card {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 1rem 1.25rem;
      background-color: var(--bg-card);
      border: 1px solid var(--border-color);
      border-radius: var(--radius-lg);
      transition: border-color var(--transition-fast);
    }

    .server-card:hover {
      border-color: var(--accent);
    }

    .server-info {
      flex: 1;
    }

    .server-name {
      font-size: 0.9375rem;
      font-weight: 600;
      color: var(--text-primary);
      margin-bottom: 0.25rem;
    }

    .server-host {
      font-size: 0.75rem;
      font-family: monospace;
      color: var(--text-muted);
    }

    .server-status {
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .status-text {
      font-size: 0.75rem;
      color: var(--success);
      font-weight: 500;
    }

    .status-text.error-text {
      color: var(--error, #ef4444);
    }

    .status-text.neutral-text {
      color: var(--text-muted);
    }

    .status-dot.error {
      background-color: var(--error, #ef4444);
    }

    .status-dot.neutral {
      background-color: var(--text-muted);
    }

    .btn-secondary {
      background: var(--bg-tertiary);
      border: 1px solid var(--border-color);
      color: var(--text-primary);
      padding: 0.5rem 1rem;
      border-radius: var(--radius-md);
      font-size: 0.875rem;
      cursor: pointer;
      display: flex;
      align-items: center;
      gap: 0.5rem;
      transition: all var(--transition-fast);
    }

    .btn-secondary:hover {
      background: var(--bg-card);
      border-color: var(--accent);
    }

    .btn-secondary:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }

    .connection-result {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      margin-top: 0.75rem;
      padding: 0.75rem;
      border-radius: var(--radius-md);
      border: 1px solid var(--border-color);
    }

    .connection-result.success {
      background-color: rgba(34, 197, 94, 0.1);
      border-color: rgba(34, 197, 94, 0.3);
    }

    .connection-result.error {
      background-color: rgba(239, 68, 68, 0.1);
      border-color: rgba(239, 68, 68, 0.3);
    }

    .result-icon {
      font-size: 1rem;
    }

    .result-message {
      flex: 1;
      font-size: 0.875rem;
    }

    .result-latency {
      font-size: 0.75rem;
      color: var(--text-muted);
      font-family: monospace;
    }
  `]
})
export class ServersComponent implements OnInit {
  private api = inject(ApiService);
  private fb = inject(FormBuilder);
  private router = inject(Router);
  private translateService = inject(TranslateService);

  servers: Server[] = [];
  showForm = false;
  loading = false;
  selectedServerId: string | null = null;
  testingConnection = false;
  connectionTestResult: { success: boolean; message: string; latency: number } | null = null;
  testingServerId: string | null = null;
  serverTestResults: { [key: string]: { success: boolean; message: string; latency: number } } = {};

  form = this.fb.group({
    name: ['', Validators.required],
    host: ['', Validators.required],
    port: [22, [Validators.required, Validators.min(1), Validators.max(65535)]],
    username: ['', Validators.required],
    sshKey: ['', Validators.required]
  });

  ngOnInit(): void {
    this.loadServers();
  }

  loadServers(): void {
    this.api.getServers().subscribe(data => this.servers = data || []);
  }

  editServer(server: Server): void {
    this.selectedServerId = server.id;
    this.form.patchValue({
      name: server.name,
      host: server.host,
      port: server.port,
      username: server.username,
      sshKey: ''
    });
    this.form.controls.sshKey.clearValidators();
    this.form.controls.sshKey.updateValueAndValidity();
    this.showForm = true;
  }

  cancelEdit(): void {
    this.showForm = !this.showForm;
    this.resetForm();
  }

  resetForm(): void {
    this.selectedServerId = null;
    this.form.reset({ port: 22 });
    this.form.controls.sshKey.setValidators(Validators.required);
    this.form.controls.sshKey.updateValueAndValidity();
    this.connectionTestResult = null;
  }

  testFormConnection(): void {
    const host = this.form.get('host')?.value;
    const port = this.form.get('port')?.value || 22;
    const username = this.form.get('username')?.value;
    const sshKey = this.form.get('sshKey')?.value;

    if (!host || !username || !sshKey) return;

    this.testingConnection = true;
    this.connectionTestResult = null;

    this.api.testServerConnection({ host, port, username, sshKey }).subscribe({
      next: (result) => {
        this.connectionTestResult = result;
        this.testingConnection = false;
      },
      error: () => {
        this.connectionTestResult = {
          success: false,
          message: this.translateService.instant('SERVERS.TEST_FAILED'),
          latency: 0
        };
        this.testingConnection = false;
      }
    });
  }

  testServerConnection(server: Server): void {
    this.testingServerId = server.id;

    this.api.testServerConnection({ serverId: server.id }).subscribe({
      next: (result) => {
        this.serverTestResults[server.id] = result;
        this.testingServerId = null;
      },
      error: () => {
        this.serverTestResults[server.id] = {
          success: false,
          message: this.translateService.instant('SERVERS.TEST_FAILED'),
          latency: 0
        };
        this.testingServerId = null;
      }
    });
  }

  openTerminal(server: Server): void {
    this.router.navigate(['/terminal', server.id]);
  }

  onSubmit(): void {
    if (this.selectedServerId && !this.form.value.sshKey) {
      this.form.controls.sshKey.clearValidators();
      this.form.controls.sshKey.updateValueAndValidity();
    }

    if (this.form.invalid) return;
    this.loading = true;

    if (this.selectedServerId) {
      this.api.updateServer(this.selectedServerId, this.form.value as any).subscribe({
        next: (updated) => {
          const idx = this.servers.findIndex(s => s.id === updated.id);
          if (idx !== -1) this.servers[idx] = updated;
          this.showForm = false;
          this.loading = false;
          this.resetForm();
        },
        error: () => this.loading = false
      });
    } else {
      this.api.createServer(this.form.value as any).subscribe({
        next: (server) => {
          this.servers.unshift(server);
          this.showForm = false;
          this.loading = false;
          this.resetForm();
        },
        error: () => this.loading = false
      });
    }
  }
}
