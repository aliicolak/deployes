import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators, FormArray, FormControl } from '@angular/forms';
import { ApiService } from '../../core/services/api.service';
import { Webhook, Project, Server } from '../../core/models';
import { TranslateModule, TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-webhooks',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, TranslateModule],
  template: `
    <div class="page">
      <div class="page-header">
        <div class="page-header-content">
          <h1 class="page-title">{{ 'WEBHOOKS.TITLE' | translate }}</h1>
          <p class="page-subtitle">{{ 'WEBHOOKS.SUBTITLE' | translate }}</p>
        </div>
        <button class="btn btn-primary" (click)="showForm = !showForm">
          {{ showForm ? ('WEBHOOKS.CANCEL' | translate) : ('WEBHOOKS.NEW_WEBHOOK' | translate) }}
        </button>
      </div>

      @if (showForm) {
        <div class="form-card">
          <h3 class="form-title">{{ 'WEBHOOKS.CREATE_WEBHOOK' | translate }}</h3>
          <form [formGroup]="form" (ngSubmit)="onSubmit()">
            <div class="form-grid">
              <div class="form-group">
                <label class="form-label">{{ 'WEBHOOKS.PROJECT' | translate }}</label>
                <select formControlName="projectId" class="form-input">
                  <option value="">{{ 'WEBHOOKS.SELECT_PROJECT' | translate }}</option>
                  @for (project of projects; track project.id) {
                    <option [value]="project.id">{{ project.name }}</option>
                  }
                </select>
              </div>
              <div class="form-group">
                <label class="form-label">{{ 'WEBHOOKS.SERVERS_MULTI' | translate }}</label>
                <div class="server-checkboxes">
                  @for (server of servers; track server.id; let i = $index) {
                    <label class="checkbox-label">
                      <input type="checkbox" 
                             [checked]="isServerSelected(server.id)"
                             (change)="onServerCheckChange(server.id, $event)">
                      <span class="checkbox-text">{{ server.name }}</span>
                      <span class="checkbox-host">({{ server.host }})</span>
                    </label>
                  }
                </div>
                @if (servers.length === 0) {
                  <p class="no-servers">{{ 'WEBHOOKS.NO_SERVERS' | translate }}</p>
                }
              </div>
            </div>
            <div class="form-actions">
              <span class="selected-count" *ngIf="selectedServerIds.length > 0">
                {{ selectedServerIds.length }} {{ 'WEBHOOKS.SELECTED_COUNT' | translate }}</span>
              <button type="submit" class="btn btn-primary" [disabled]="form.invalid || selectedServerIds.length === 0 || loading">
                @if (loading) { <span class="spinner"></span> }
                {{ 'WEBHOOKS.CREATE' | translate }}
              </button>
            </div>
          </form>
        </div>
      }

      @if (newWebhook) {
        <div class="success-card">
          <div class="success-header">
            <span class="success-icon">✓</span>
            <h3>{{ 'WEBHOOKS.CREATED_TITLE' | translate }}</h3>
          </div>
          <div class="webhook-details">
            <div class="detail-row">
              <label>{{ 'WEBHOOKS.WEBHOOK_URL' | translate }}</label>
              <div class="copy-field">
                <code>{{ newWebhook.webhookUrl }}</code>
                <button class="btn-icon" (click)="copyToClipboard(newWebhook.webhookUrl)">⧉</button>
              </div>
            </div>
            <div class="detail-row">
              <label>{{ 'WEBHOOKS.SECRET' | translate }}</label>
              <div class="copy-field">
                <code>{{ newWebhook.secret }}</code>
                <button class="btn-icon" (click)="copyToClipboard(newWebhook.secret || '')">⧉</button>
              </div>
            </div>
            <div class="detail-row">
              <label>{{ 'WEBHOOKS.LINKED_SERVERS' | translate }}</label>
              <div class="server-tags">
                @for (serverId of newWebhook.serverIds; track serverId) {
                  <span class="server-tag">{{ getServerName(serverId) }}</span>
                }
              </div>
            </div>
          </div>
          <div class="warning-box">
            {{ 'WEBHOOKS.SECRET_WARNING' | translate }}
          </div>
          <button class="btn btn-secondary" (click)="newWebhook = null">{{ 'WEBHOOKS.CLOSE' | translate }}</button>
        </div>
      }

      <div class="webhooks-list">
        @for (webhook of webhooks; track webhook.id) {
          <div class="webhook-card">
            <div class="webhook-info">
              <span class="webhook-url">{{ webhook.webhookUrl }}</span>
              <span class="webhook-meta">{{ getProjectName(webhook.projectId) }} → {{ getServerNames(webhook.serverIds) }}</span>
            </div>
            <div class="webhook-actions">
              <button class="toggle-btn" [class.active]="webhook.isActive" (click)="toggleWebhook(webhook)">
                {{ webhook.isActive ? ('WEBHOOKS.ACTIVE' | translate) : ('WEBHOOKS.INACTIVE' | translate) }}
              </button>
              <button class="btn-icon danger" (click)="deleteWebhook(webhook)">✕</button>
            </div>
          </div>
        } @empty {
          <div class="empty-state">
            <span class="empty-state-icon">◇</span>
            <h3>{{ 'WEBHOOKS.EMPTY_TITLE' | translate }}</h3>
            <p>{{ 'WEBHOOKS.EMPTY_DESC' | translate }}</p>
          </div>
        }
      </div>
    </div>
  `,
  styles: [`
    .form-card, .success-card {
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

    .server-checkboxes {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
      max-height: 200px;
      overflow-y: auto;
      padding: 0.75rem;
      background-color: var(--bg-secondary);
      border: 1px solid var(--border-color);
      border-radius: var(--radius-md);
    }

    .checkbox-label {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      cursor: pointer;
      padding: 0.375rem 0.5rem;
      border-radius: var(--radius-sm);
      transition: background-color var(--transition-fast);
    }

    .checkbox-label:hover {
      background-color: var(--bg-tertiary);
    }

    .checkbox-label input[type="checkbox"] {
      width: 1rem;
      height: 1rem;
      accent-color: var(--primary);
    }

    .checkbox-text {
      font-size: 0.875rem;
      font-weight: 500;
      color: var(--text-primary);
    }

    .checkbox-host {
      font-size: 0.75rem;
      color: var(--text-muted);
    }

    .no-servers {
      font-size: 0.75rem;
      color: var(--text-muted);
      font-style: italic;
    }

    .selected-count {
      font-size: 0.75rem;
      color: var(--text-secondary);
      margin-right: 1rem;
    }

    .form-actions {
      margin-top: 1.25rem;
      display: flex;
      justify-content: flex-end;
      align-items: center;
    }

    .success-card {
      border-color: var(--success-border);
      background-color: var(--success-light);
    }

    .success-header {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      margin-bottom: 1rem;
    }

    .success-icon {
      color: var(--success);
      font-weight: bold;
    }

    .success-header h3 {
      font-size: 1rem;
      color: var(--success);
    }

    .webhook-details {
      margin-bottom: 1rem;
    }

    .detail-row {
      margin-bottom: 0.75rem;
    }

    .detail-row label {
      display: block;
      font-size: 0.75rem;
      color: var(--text-secondary);
      margin-bottom: 0.25rem;
    }

    .copy-field {
      display: flex;
      gap: 0.5rem;
      align-items: center;
    }

    .copy-field code {
      flex: 1;
      padding: 0.5rem 0.75rem;
      background-color: var(--bg-primary);
      border: 1px solid var(--border-color);
      border-radius: var(--radius-md);
      font-size: 0.75rem;
      word-break: break-all;
    }

    .server-tags {
      display: flex;
      flex-wrap: wrap;
      gap: 0.375rem;
    }

    .server-tag {
      padding: 0.25rem 0.5rem;
      background-color: var(--bg-primary);
      border: 1px solid var(--border-color);
      border-radius: var(--radius-sm);
      font-size: 0.75rem;
      color: var(--text-primary);
    }

    .warning-box {
      padding: 0.75rem;
      background-color: var(--warning-light);
      border: 1px solid var(--warning-border);
      border-radius: var(--radius-md);
      color: var(--warning);
      font-size: 0.75rem;
      margin-bottom: 1rem;
    }

    .webhooks-list {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .webhook-card {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 1rem 1.25rem;
      background-color: var(--bg-card);
      border: 1px solid var(--border-color);
      border-radius: var(--radius-lg);
    }

    .webhook-info {
      flex: 1;
      display: flex;
      flex-direction: column;
      gap: 0.25rem;
    }

    .webhook-url {
      font-size: 0.75rem;
      font-family: monospace;
      color: var(--text-primary);
      word-break: break-all;
    }

    .webhook-meta {
      font-size: 0.75rem;
      color: var(--text-muted);
    }

    .webhook-actions {
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .toggle-btn {
      padding: 0.375rem 0.75rem;
      font-size: 0.75rem;
      font-weight: 500;
      background-color: var(--bg-secondary);
      color: var(--text-secondary);
      border: 1px solid var(--border-color);
      border-radius: var(--radius-md);
      cursor: pointer;
      transition: all var(--transition-fast);
    }

    .toggle-btn.active {
      background-color: var(--success-light);
      color: var(--success);
      border-color: var(--success-border);
    }

    .btn-icon.danger {
      color: var(--danger);
    }

    .btn-icon.danger:hover {
      background-color: var(--danger-light);
    }
  `]
})
export class WebhooksComponent implements OnInit {
  private api = inject(ApiService);
  private fb = inject(FormBuilder);
  private translateService = inject(TranslateService);

  webhooks: Webhook[] = [];
  projects: Project[] = [];
  servers: Server[] = [];
  newWebhook: Webhook | null = null;
  showForm = false;
  loading = false;
  selectedServerIds: string[] = [];

  form = this.fb.group({
    projectId: ['', Validators.required]
  });

  ngOnInit(): void {
    this.loadData();
  }

  loadData(): void {
    this.api.getWebhooks().subscribe(data => this.webhooks = data || []);
    this.api.getProjects().subscribe(data => this.projects = data || []);
    this.api.getServers().subscribe(data => this.servers = data || []);
  }

  isServerSelected(serverId: string): boolean {
    return this.selectedServerIds.includes(serverId);
  }

  onServerCheckChange(serverId: string, event: Event): void {
    const checkbox = event.target as HTMLInputElement;
    if (checkbox.checked) {
      if (!this.selectedServerIds.includes(serverId)) {
        this.selectedServerIds.push(serverId);
      }
    } else {
      this.selectedServerIds = this.selectedServerIds.filter(id => id !== serverId);
    }
  }

  onSubmit(): void {
    if (this.form.invalid || this.selectedServerIds.length === 0) return;
    this.loading = true;

    const request = {
      projectId: this.form.value.projectId as string,
      serverIds: this.selectedServerIds
    };

    this.api.createWebhook(request).subscribe({
      next: (webhook) => {
        this.newWebhook = webhook;
        this.webhooks.unshift(webhook);
        this.form.reset();
        this.selectedServerIds = [];
        this.showForm = false;
        this.loading = false;
      },
      error: () => this.loading = false
    });
  }

  toggleWebhook(webhook: Webhook): void {
    this.api.updateWebhook(webhook.id, { isActive: !webhook.isActive }).subscribe({
      next: (updated) => {
        const idx = this.webhooks.findIndex(w => w.id === webhook.id);
        if (idx > -1) this.webhooks[idx] = updated;
      }
    });
  }

  deleteWebhook(webhook: Webhook): void {
    if (!confirm(this.translateService.instant('WEBHOOKS.DELETE_CONFIRM'))) return;

    this.api.deleteWebhook(webhook.id).subscribe({
      next: () => {
        this.webhooks = this.webhooks.filter(w => w.id !== webhook.id);
      }
    });
  }

  getProjectName(id: string): string {
    return this.projects.find(p => p.id === id)?.name || this.translateService.instant('WEBHOOKS.UNKNOWN');
  }

  getServerName(id: string): string {
    return this.servers.find(s => s.id === id)?.name || this.translateService.instant('WEBHOOKS.UNKNOWN');
  }

  getServerNames(ids: string[]): string {
    if (!ids || ids.length === 0) return this.translateService.instant('WEBHOOKS.NO_SERVER');
    return ids.map(id => this.getServerName(id)).join(', ');
  }

  copyToClipboard(text: string): void {
    navigator.clipboard.writeText(text).then(() => {
      alert(this.translateService.instant('WEBHOOKS.COPIED'));
    });
  }
}

