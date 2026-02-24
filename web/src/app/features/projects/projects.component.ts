import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormControl, ReactiveFormsModule, Validators } from '@angular/forms';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../../core/services/api.service';
import { Project, Secret, ProjectType } from '../../core/models';
import { TranslateModule, TranslateService } from '@ngx-translate/core';

interface FileItem {
  name: string;
  path: string;
  isDir: boolean;
  excluded: boolean;
  file?: File;
  children?: FileItem[];
  expanded?: boolean;
  level?: number;
}

@Component({
  selector: 'app-projects',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, FormsModule, TranslateModule],
  template: `
    <div class="page">
      <div class="page-header">
        <div class="page-header-content">
          <h1 class="page-title">{{ 'PROJECTS.TITLE' | translate }}</h1>
          <p class="page-subtitle">{{ 'PROJECTS.SUBTITLE' | translate }}</p>
        </div>
        <button class="btn btn-primary" (click)="cancelEdit()">
          {{ showForm ? ('PROJECTS.CANCEL' | translate) : ('PROJECTS.NEW_PROJECT' | translate) }}
        </button>
      </div>

      @if (showForm) {
        <div class="form-card">
          <h3 class="form-title">{{ selectedProjectId ? ('PROJECTS.EDIT_PROJECT' | translate) : ('PROJECTS.CREATE_PROJECT' | translate) }}</h3>
          <form [formGroup]="form" (ngSubmit)="onSubmit()">
            <div class="form-grid">
              <div class="form-group">
                <label class="form-label">{{ 'PROJECTS.PROJECT_NAME' | translate }}</label>
                <input type="text" formControlName="name" class="form-input" placeholder="my-awesome-app" />
              </div>
              <div class="form-group">
                <label class="form-label">{{ 'PROJECTS.PROJECT_TYPE' | translate }}</label>
                <select formControlName="type" class="form-input" (change)="onProjectTypeChange()">
                  <option value="github">{{ 'PROJECTS.TYPE_GITHUB' | translate }}</option>
                  <option value="local">{{ 'PROJECTS.TYPE_LOCAL' | translate }}</option>
                </select>
              </div>

              @if (form.get('type')?.value === 'github') {
                <div class="form-group">
                  <label class="form-label">{{ 'PROJECTS.REPO_URL' | translate }}</label>
                  <input type="text" formControlName="repoUrl" class="form-input" placeholder="https://github.com/user/repo.git" />
                </div>
                <div class="form-group">
                  <label class="form-label">{{ 'PROJECTS.BRANCH' | translate }}</label>
                  <input type="text" formControlName="branch" class="form-input" placeholder="main" />
                </div>
              }

              @if (form.get('type')?.value === 'local') {
                <div class="form-group full-width">
                  <label class="form-label">{{ 'PROJECTS.SELECT_FOLDER' | translate }}</label>
                  <div class="folder-upload-container">
                    <button type="button" class="btn btn-secondary" (click)="folderInput.click()">
                      {{ 'PROJECTS.CHOOSE_FOLDER' | translate }}
                    </button>
                    <input type="file" #folderInput webkitdirectory directory (change)="onFolderSelected($event)" style="display:none" />
                    @if (selectedFolderName) {
                      <span class="file-name">{{ selectedFolderName }}</span>
                      <span class="file-count">({{ folderFiles.length }} {{ 'PROJECTS.RESULTS' | translate }})</span>
                    }
                  </div>

                  @if (folderFiles.length > 0) {
                    <div class="folder-preview">
                      <div class="folder-preview-header">
                        <h4>{{ 'PROJECTS.FOLDER_CONTENT' | translate }}</h4>
                        <div class="exclude-presets">
                          <span class="preset-label">{{ 'PROJECTS.QUICK_EXCLUDE' | translate }}</span>
                          <button type="button" class="preset-btn" (click)="toggleExcludePattern('node_modules')">node_modules</button>
                          <button type="button" class="preset-btn" (click)="toggleExcludePattern('venv')">venv</button>
                          <button type="button" class="preset-btn" (click)="toggleExcludePattern('__pycache__')">__pycache__</button>
                          <button type="button" class="preset-btn" (click)="toggleExcludePattern('.git')">.git</button>
                          <button type="button" class="preset-btn" (click)="toggleExcludePattern('target')">target</button>
                        </div>
                      </div>

                      <!-- Path Search Filter -->
                      <div class="path-search">
                        <input type="text" class="form-input" [ngModel]="pathSearchQuery" (ngModelChange)="onPathSearchChange($event)"
                          [ngModelOptions]="{standalone: true}"
                          [placeholder]="'PROJECTS.PATH_PLACEHOLDER' | translate" />
                        <div class="search-info">
                          <span class="search-hint">{{ 'PROJECTS.PATH_HINT' | translate }}</span>
                          <span class="search-count">@if (pathSearchQuery) { {{ getFilteredFiles().length }} {{ 'PROJECTS.RESULTS' | translate }} }</span>
                        </div>
                      </div>

                      <div class="folder-tree">
                        @for (item of getFilteredFiles(); track item.path) {
                          @if (isRootItem(item)) {
                            <ng-container>
                              @for (row of getTreeRows(item); track row.path) {
                                <div class="folder-item"
                                     [class.excluded]="row.excluded"
                                     [class.folder]="row.isDir"
                                     [style.padding-left.px]="(row.level || 0) * 20 + 12">
                                  @if (row.isDir) {
                                    <span class="expand-icon" (click)="toggleExpand(row); $event.stopPropagation()">
                                      {{ row.expanded ? '▼' : '▶' }}
                                    </span>
                                  } @else {
                                    <span class="spacer"></span>
                                  }
                                  <span class="item-icon" (click)="toggleFileExclude(row); $event.stopPropagation()">
                                    {{ row.excluded ? '☑️' : (row.isDir ? '📁' : '📄') }}
                                  </span>
                                  <span class="item-name" (click)="toggleFileExclude(row); $event.stopPropagation()"
                                        [title]="row.path">{{ row.name }}</span>
                                  <span class="item-status">{{ row.excluded ? ('PROJECTS.EXCLUDED' | translate) : ('PROJECTS.INCLUDED' | translate) }}</span>
                                </div>
                              }
                            </ng-container>
                          }
                        } @empty {
                          <div class="empty-search">
                            <span>{{ 'PROJECTS.NO_RESULTS' | translate }}</span>
                          </div>
                        }
                      </div>

                      <div class="folder-actions">
                        <div class="exclude-summary">
                          <span>{{ 'PROJECTS.INCLUDED' | translate }}: {{ getIncludedCount() }}</span>
                          <span>{{ 'PROJECTS.EXCLUDED' | translate }}: {{ getExcludedCount() }}</span>
                        </div>
                        <button type="button" class="btn btn-primary" (click)="prepareAndUploadFolder()" [disabled]="uploading || getIncludedCount() === 0">
                          @if (uploading) { {{ 'PROJECTS.CREATING_ZIP' | translate }} } @else { {{ 'PROJECTS.ZIP_AND_UPLOAD' | translate }} }
                        </button>
                      </div>
                    </div>
                  }

                  @if (form.get('localPath')?.value) {
                    <p class="form-hint success">✅ {{ form.get('localPath')?.value }}</p>
                  }
                </div>
              }
              
              <!-- Test Repository Access Button - Only for GitHub projects -->
              @if (form.get('type')?.value === 'github') {
                <div class="form-group full-width">
                  <button type="button" class="btn btn-secondary" (click)="testRepoAccess()" [disabled]="!form.get('repoUrl')?.value || testingRepo">
                  @if (testingRepo) { <span class="spinner"></span> }
                   {{ 'PROJECTS.TEST_REPO' | translate }}
                </button>
                
                @if (repoAccessResult) {
                  <div class="repo-access-result" [class.success]="repoAccessResult.accessible" [class.error]="!repoAccessResult.accessible">
                    <div class="result-header">
                      <span class="result-icon">{{ repoAccessResult.accessible ? '✅' : (repoAccessResult.isPrivate ? '🔒' : '❌') }}</span>
                      <span class="result-message">{{ repoAccessResult.message }}</span>
                    </div>
                    @if (repoAccessResult.guidance) {
                      <div class="result-guidance">
                        <h5>{{ 'PROJECTS.PRIVATE_REPO_STEPS' | translate }}</h5>
                        <pre>{{ repoAccessResult.guidance }}</pre>
                      </div>
                    }
                    @if (!repoAccessResult.branchExists && repoAccessResult.accessible) {
                      <p class="warning-text">{{ 'PROJECTS.BRANCH_NOT_FOUND' | translate }}</p>
                    }
                  </div>
                }
              </div>
              }

              @if (selectedProjectId && form.get('type')?.value === 'github') {
                @if (selectedProjectKey) {
                  <div class="form-group full-width deploy-key-card">
                    <div class="deploy-key-header">
                       <h4 class="form-title">{{ 'PROJECTS.DEPLOY_KEY_TITLE' | translate }}</h4>
                       <button type="button" class="btn btn-secondary btn-sm" (click)="copyKey()">{{ 'PROJECTS.COPY' | translate }}</button>
                    </div>
                    <p class="form-hint">{{ 'PROJECTS.DEPLOY_KEY_HINT' | translate }}</p>
                    <div class="key-box">
                      <code>{{ selectedProjectKey }}</code>
                    </div>
                  </div>
                } @else {
                  <div class="form-group full-width info-box warning">
                      <span>{{ 'PROJECTS.DEPLOY_KEY_NOT_GENERATED' | translate }}</span>
                  </div>
                }
              }

<div class="form-group full-width">
                <div class="script-label-row">
                  <label class="form-label">{{ 'PROJECTS.DEPLOY_SCRIPT_TITLE' | translate }}</label>
                  <span class="badge badge-neutral">Bash Script</span>
                </div>
                <p class="form-hint">{{ 'PROJECTS.DEPLOY_SCRIPT_HINT' | translate }}</p>
                
                <div class="template-selector">
                  @for (t of templates; track t.id) {
                    <button type="button" class="template-chip" (click)="applyTemplate(t)" title="{{t.name}} şablonunu uygula">
                      <span class="template-icon">{{ t.icon }}</span>
                      <span class="template-name">{{ t.name }}</span>
                    </button>
                  }
                </div>

                <textarea formControlName="deployScript" class="form-input form-textarea code-editor" rows="10" 
                  placeholder="# Buraya bash script kodlarınızı yazın...
# Örnek:
# npm install
# npm run build"></textarea>
                
                <div class="script-help">
                  <div class="help-item">
                    <span class="help-icon">ℹ️</span>
                    <span>{{ 'PROJECTS.SCRIPT_INFO' | translate }}</span>
                  </div>
                </div>
              </div>

              @if (selectedProjectId) {
                <div class="form-group full-width secrets-section">
                  <h4 class="form-title">{{ 'PROJECTS.ENV_VARS_TITLE' | translate }}</h4>
                  <p class="form-hint">{{ 'PROJECTS.ENV_VARS_HINT' | translate }}</p>

                  <div class="secrets-list">
                    @for (secret of secrets; track secret.id) {
                      <div class="secret-item">
                        <div class="secret-info">
                          <span class="secret-key">{{ secret.key }}</span>
                          <span class="secret-value">************</span>
                        </div>
                        <button type="button" class="btn-icon danger" (click)="deleteSecret(secret.id)" [title]="'PROJECTS.DELETE_SECRET' | translate">🗑️</button>
                      </div>
                    } @empty {
                       <p class="text-muted text-sm">{{ 'PROJECTS.NO_SECRETS' | translate }}</p>
                    }
                  </div>

                  <div class="add-secret-form">
                    <input type="text" [formControl]="secretKeyControl" placeholder="KEY (örn: DB_PASS)" class="form-input key-input">
                    <input type="password" [formControl]="secretValueControl" placeholder="VALUE" class="form-input value-input">
                    <button type="button" class="btn btn-secondary" (click)="addSecret()" [disabled]="!secretKeyControl.value || !secretValueControl.value">
                      ➕ Ekle
                    </button>
                  </div>
                </div>
              } @else {
                 <div class="form-group full-width">
                    <div class="info-box">
                        <span class="info-icon">ℹ️</span>
                        <span>{{ 'PROJECTS.ENV_VARS_INFO' | translate }}</span>
                    </div>
                 </div>
              }
            </div>
            <div class="form-actions">
              <button type="submit" class="btn btn-primary" [disabled]="form.invalid || loading">
                @if (loading) { <span class="spinner"></span> }
                {{ selectedProjectId ? ('PROJECTS.UPDATE' | translate) : ('PROJECTS.CREATE' | translate) }}
              </button>
            </div>
          </form>
        </div>
      }


      <div class="projects-grid">
        @for (project of projects; track project.id) {
          <div class="project-card">
            <div class="project-header">
              <h3 class="project-name">{{ project.name }}</h3>
              <div class="project-actions">
                <button class="btn-icon" (click)="editProject(project)" [title]="'PROJECTS.EDIT' | translate">✎</button>
                <button class="btn-icon danger" (click)="deleteProject(project.id, project.name)" [title]="'PROJECTS.DELETE' | translate">🗑</button>
              </div>
            </div>
            <p class="project-repo">
              @if (project.type === 'local') {
                <span class="badge badge-info">📁 Local</span>
                {{ project.localPath }}
              } @else {
                <span class="badge badge-neutral">🔗 GitHub</span>
                {{ project.repoUrl }}
              }
            </p>
            <div class="project-footer">
              @if (project.type === 'github') {
                <span class="badge badge-neutral">{{ project.branch }}</span>
              } @else {
                <span class="badge badge-info">Local Project</span>
              }
            </div>
          </div>
        } @empty {
          <div class="empty-state full-width">
            <span class="empty-state-icon">◫</span>
            <h3>{{ 'PROJECTS.EMPTY_TITLE' | translate }}</h3>
            <p>{{ 'PROJECTS.EMPTY_DESC' | translate }}</p>
          </div>
        }
      </div>
    </div>
  `,
  styles: [`
    .deploy-key-card {
        background: var(--bg-tertiary);
        padding: 1rem;
        border-radius: var(--radius-md);
        border: 1px solid var(--border-color);
        margin-bottom: 1.5rem;
    }
    .deploy-key-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 0.5rem;
    }
    .key-box {
        background: #1e1e1e;
        color: #a3e635;
        padding: 0.75rem;
        border-radius: var(--radius-sm);
        font-family: monospace;
        font-size: 0.75rem;
        word-break: break-all;
        overflow-y: auto;
        max-height: 100px;
    }

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

    .form-actions {
      margin-top: 1.25rem;
      display: flex;
      justify-content: flex-end;
    }

    .projects-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
      gap: 1rem;
    }

    .project-card {
      background-color: var(--bg-card);
      border: 1px solid var(--border-color);
      border-radius: var(--radius-lg);
      padding: 1.25rem;
      transition: border-color var(--transition-fast);
    }

    .project-card:hover {
      border-color: var(--accent);
    }

    .project-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 0.5rem;
    }

    .project-actions {
      display: flex;
      gap: 0.5rem;
    }

    .project-name {
      font-size: 1rem;
      font-weight: 600;
    }

    .project-repo {
      font-size: 0.75rem;
      color: var(--text-muted);
      word-break: break-all;
      margin-bottom: 1rem;
    }

    .project-footer {
      display: flex;
      gap: 0.5rem;
    }

    .empty-state {
      grid-column: 1 / -1;
    }

    .script-label-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 0.5rem;
    }

    .form-hint {
      font-size: 0.875rem;
      color: var(--text-muted);
      margin-bottom: 1rem;
    }

    .template-selector {
      display: flex;
      gap: 0.75rem;
      margin-bottom: 1rem;
      flex-wrap: wrap;
    }

    .template-chip {
      background-color: var(--bg-tertiary);
      border: 1px solid var(--border-color);
      color: var(--text-primary);
      padding: 0.5rem 0.75rem;
      border-radius: var(--radius-full);
      font-size: 0.875rem;
      cursor: pointer;
      display: flex;
      align-items: center;
      gap: 0.5rem;
      transition: all var(--transition-fast);
    }

    .template-chip:hover {
      background-color: var(--bg-card);
      border-color: var(--accent);
      transform: translateY(-1px);
    }

    .code-editor {
      font-family: 'JetBrains Mono', 'Fira Code', Consolas, monospace;
      font-size: 0.875rem;
      line-height: 1.5;
      background-color: #1e1e1e;
      color: #d4d4d4;
      border-color: #333;
    }

    .script-help {
      margin-top: 0.75rem;
      background-color: rgba(59, 130, 246, 0.1);
      border: 1px solid rgba(59, 130, 246, 0.2);
      border-radius: var(--radius-md);
      padding: 0.75rem;
    }

    .help-item {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 0.75rem;
      color: var(--text-managed);
    }

    .help-icon {
      font-size: 1rem;
    }

    .secrets-section {
      margin-top: 1.5rem;
      padding-top: 1.5rem;
      border-top: 1px dashed var(--border-color);
    }

    .secrets-list {
      margin-bottom: 1rem;
    }

    .secret-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      background-color: var(--bg-tertiary);
      padding: 0.75rem;
      border-radius: var(--radius-md);
      margin-bottom: 0.5rem;
      border: 1px solid var(--border-color);
    }

    .secret-key {
      font-weight: 600;
      color: var(--accent);
      margin-right: 1rem;
    }

    .secret-value {
      color: var(--text-muted);
      font-family: monospace;
    }

    .add-secret-form {
      display: flex;
      gap: 0.5rem;
    }
    
    .key-input { flex: 1; }
    .value-input { flex: 2; }

    .btn-icon.danger {
        color: #ef4444;
    }
    .btn-icon.danger:hover {
        background-color: rgba(239, 68, 68, 0.1);
    }

    .info-box {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem;
        background-color: var(--bg-tertiary);
        border-radius: var(--radius-md);
        font-size: 0.875rem;
        color: var(--text-muted);
    }
    .info-box.warning {
        background-color: rgba(234, 179, 8, 0.1);
        border: 1px solid rgba(234, 179, 8, 0.2);
        color: #eab308;
    }

    .file-upload-container {
        display: flex;
        gap: 0.75rem;
        align-items: center;
    }

    .file-name {
        font-size: 0.875rem;
        color: var(--text-muted);
        flex: 1;
    }

    .file-count {
        font-size: 0.75rem;
        color: var(--text-muted);
        background: var(--bg-tertiary);
        padding: 0.25rem 0.5rem;
        border-radius: var(--radius-sm);
    }

    .folder-upload-container {
        display: flex;
        gap: 0.75rem;
        align-items: center;
    }

    .folder-preview {
        margin-top: 1rem;
        border: 1px solid var(--border-color);
        border-radius: var(--radius-md);
        padding: 1rem;
        background: var(--bg-tertiary);
    }

    .path-search {
        margin-bottom: 1rem;
    }

    .path-search .form-input {
        width: 100%;
        font-family: monospace;
        font-size: 0.875rem;
    }

    .search-info {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-top: 0.5rem;
    }

    .search-hint {
        font-size: 0.75rem;
        color: var(--text-muted);
    }

    .search-count {
        font-size: 0.75rem;
        color: var(--accent);
        font-weight: 500;
    }

    .empty-search {
        padding: 2rem;
        text-align: center;
        color: var(--text-muted);
        font-size: 0.875rem;
    }

    .folder-preview-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1rem;
        flex-wrap: wrap;
        gap: 0.75rem;
    }

    .folder-preview-header h4 {
        margin: 0;
        font-size: 0.875rem;
        font-weight: 600;
    }

    .exclude-presets {
        display: flex;
        gap: 0.5rem;
        align-items: center;
        flex-wrap: wrap;
    }

    .preset-label {
        font-size: 0.75rem;
        color: var(--text-muted);
        margin-right: 0.25rem;
    }

    .preset-btn {
        background: var(--bg-card);
        border: 1px solid var(--border-color);
        color: var(--text-primary);
        padding: 0.25rem 0.5rem;
        border-radius: var(--radius-sm);
        font-size: 0.75rem;
        cursor: pointer;
        transition: all var(--transition-fast);
    }

    .preset-btn:hover {
        border-color: var(--accent);
        background: var(--bg-tertiary);
    }

    .folder-tree {
        max-height: 250px;
        overflow-y: auto;
        border: 1px solid var(--border-color);
        border-radius: var(--radius-sm);
        background: var(--bg-card);
    }

    .folder-item {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.35rem 0.75rem;
        border-bottom: 1px solid var(--border-color);
        cursor: pointer;
        transition: background var(--transition-fast);
        min-height: 32px;
    }

    .folder-item:last-child {
        border-bottom: none;
    }

    .folder-item:hover {
        background: var(--bg-tertiary);
    }

    .folder-item.excluded {
        opacity: 0.5;
        text-decoration: line-through;
    }

    .folder-item.folder {
        font-weight: 500;
    }

    .expand-icon {
        width: 16px;
        height: 16px;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        font-size: 0.625rem;
        color: var(--text-muted);
        cursor: pointer;
        user-select: none;
        transition: color var(--transition-fast);
        flex-shrink: 0;
    }

    .expand-icon:hover {
        color: var(--accent);
    }

    .spacer {
        width: 16px;
        flex-shrink: 0;
    }

    .item-icon {
        font-size: 1rem;
        flex-shrink: 0;
        cursor: pointer;
        user-select: none;
    }

    .item-name {
        flex: 1;
        font-size: 0.875rem;
        cursor: pointer;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        user-select: none;
    }

    .item-status {
        font-size: 0.7rem;
        color: var(--text-muted);
        flex-shrink: 0;
        min-width: 40px;
        text-align: right;
    }

    .folder-actions {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-top: 1rem;
        padding-top: 1rem;
        border-top: 1px dashed var(--border-color);
    }

    .exclude-summary {
        display: flex;
        gap: 1rem;
        font-size: 0.875rem;
    }

    .exclude-summary span {
        color: var(--text-muted);
    }

    .form-hint.success {
        color: #22c55e;
        background-color: rgba(34, 197, 94, 0.1);
        border: 1px solid rgba(34, 197, 94, 0.2);
        padding: 0.75rem;
        border-radius: var(--radius-md);
    }

    .badge {
        display: inline-flex;
        align-items: center;
        gap: 0.25rem;
        padding: 0.25rem 0.5rem;
        border-radius: var(--radius-sm);
        font-size: 0.7rem;
        font-weight: 500;
    }

    .badge-neutral {
        background-color: var(--bg-tertiary);
        color: var(--text-muted);
        border: 1px solid var(--border-color);
    }

    .badge-info {
        background-color: rgba(59, 130, 246, 0.1);
        color: #3b82f6;
        border: 1px solid rgba(59, 130, 246, 0.2);
    }
  `]
})
export class ProjectsComponent implements OnInit {
  private api = inject(ApiService);
  private fb = inject(FormBuilder);
  private translateService = inject(TranslateService);

  projects: Project[] = [];
  showForm = false;
  loading = false;
  selectedProjectId: string | null = null;
  selectedProjectKey: string | undefined = undefined;
  // showExamples removed
  testingRepo = false;
  repoAccessResult: { accessible: boolean; isPrivate: boolean; message: string; guidance?: string; repoType: string; branchExists: boolean } | null = null;

  secrets: Secret[] = [];
  secretKeyControl = new FormControl('');
  secretValueControl = new FormControl('');

  // Local project upload - folder based
  selectedFolderName: string | null = null;
  folderFiles: FileItem[] = [];
  excludedPatterns: string[] = ['node_modules', 'venv', '__pycache__', '.git', 'target', '.env', 'dist', 'build'];
  uploading = false;

  // Path search filtering
  pathSearchQuery: string = '';

  // Legacy single file support
  uploadedFile: File | null = null;
  uploadProgress = 0;

  // Folder upload with File API (client-side zipping)
  folderFileList: FileList | null = null;

  readonly templates = [
    {
      id: 'node-pm2',
      name: 'Node.js + PM2',
      icon: '🟢',
      script: `# Node.js PM2 Deployment
echo "🚀 Starting deployment..."
npm install
npm run build
# PM2 varsa restart eder, yoksa başlatır
pm2 reload ecosystem.config.js || pm2 restart all || echo "PM2 not configured"
echo "✅ Deployment finished"`
    },
    {
      id: 'python',
      name: 'Python / Django',
      icon: '🐍',
      script: `# Python Deployment
echo "🚀 Starting deployment..."
# Venv varsa aktif et
if [ -d "venv" ]; then source venv/bin/activate; fi
pip install -r requirements.txt
python manage.py migrate
systemctl restart gunicorn || echo "Systemctl failed, skipping restart"
echo "✅ Deployment finished"`
    },
    {
      id: 'docker',
      name: 'Docker Compose',
      icon: '🐳',
      script: `# Docker Compose Deployment
echo "🚀 Starting deployment..."
docker-compose pull
docker-compose up -d --build
docker system prune -f
echo "✅ Deployment finished"`
    },
    {
      id: 'static',
      name: 'Static / React / Vue',
      icon: '⚛️',
      script: `# Static Site Deployment
echo "🚀 Starting deployment..."
npm install
npm run build
# Deploy to Nginx folder (adjust path!)
# rsync -av --delete ./dist/ /var/www/html/
echo "✅ Build finished. Uncomment rsync line to deploy."`
    },
    {
      id: 'go',
      name: 'Go Lang',
      icon: '🐹',
      script: `# Go Deployment
echo "🚀 Starting deployment..."
go build -o app
systemctl restart myapp.service
echo "✅ Deployment finished"`
    },
    {
      id: 'dotnet',
      name: '.NET / ASP.NET Core',
      icon: '🔷',
      script: `# .NET / ASP.NET Core Deployment
echo "🚀 Starting deployment..."
# Restore dependencies
dotnet restore
# Build the project in Release mode
dotnet build --configuration Release
# Publish the application
dotnet publish --configuration Release --output ./publish
# Restart the service (adjust service name!)
systemctl restart myapp.service || echo "Systemctl failed, skipping restart"
echo "✅ Deployment finished"`
    }
  ];

  form = this.fb.group({
    name: ['', Validators.required],
    type: ['github' as ProjectType, Validators.required],
    repoUrl: [''],
    branch: ['main'],
    localPath: [''],
    deployScript: ['', Validators.required]
  });

  ngOnInit(): void {
    this.loadProjects();
  }

  loadProjects(): void {
    this.api.getProjects().subscribe(data => this.projects = data || []);
  }

  editProject(project: Project): void {
    this.selectedProjectId = project.id;
    this.form.patchValue({
      name: project.name,
      type: project.type || 'github',
      repoUrl: project.repoUrl || '',
      branch: project.branch || 'main',
      localPath: project.localPath || '',
      deployScript: project.deployScript
    });
    this.selectedProjectKey = project.scmPublicKey;
    this.showForm = true;
    this.loadSecrets(project.id);
  }

  loadSecrets(projectId: string): void {
    this.api.getSecrets(projectId).subscribe({
      next: (data) => this.secrets = data || [],
      error: () => this.secrets = []
    });
  }

  addSecret(): void {
    if (!this.selectedProjectId || !this.secretKeyControl.value || !this.secretValueControl.value) return;

    this.api.createSecret({
      projectId: this.selectedProjectId,
      key: this.secretKeyControl.value || '',
      value: this.secretValueControl.value || ''
    }).subscribe(secret => {
      this.secrets.push(secret);
      this.secretKeyControl.reset();
      this.secretValueControl.reset();
    });
  }

  deleteSecret(id: string): void {
    if (!confirm('Bu değişkeni silmek istediğinize emin misiniz?')) return;
    this.api.deleteSecret(id).subscribe(() => {
      this.secrets = this.secrets.filter(s => s.id !== id);
    });
  }

  cancelEdit(): void {
    this.showForm = !this.showForm;
    this.resetForm();
  }

  resetForm(): void {
    this.selectedProjectId = null;
    this.selectedProjectKey = undefined;
    this.form.reset({ type: 'github', branch: 'main' });
    this.repoAccessResult = null;
    this.secrets = [];
    this.secretKeyControl.reset();
    this.secretValueControl.reset();
    this.uploadedFile = null;
    this.uploadProgress = 0;
    // Reset folder upload state
    this.selectedFolderName = null;
    this.folderFiles = [];
    this.excludedPatterns = ['node_modules', 'venv', '__pycache__', '.git', 'target', '.env', 'dist', 'build'];
    this.pathSearchQuery = '';
  }

  applyTemplate(template: any): void {
    const currentVal = this.form.get('deployScript')?.value;
    if (currentVal && !confirm('Mevcut scriptin üzerine yazılacak. Devam etmek istiyor musunuz?')) {
      return;
    }
    this.form.patchValue({ deployScript: template.script });
  }

  testRepoAccess(): void {
    const repoUrl = this.form.get('repoUrl')?.value;
    const branch = this.form.get('branch')?.value || 'main';

    if (!repoUrl) return;

    this.testingRepo = true;
    this.repoAccessResult = null;

    this.api.testRepoAccess({ repoUrl, branch }).subscribe({
      next: (result) => {
        this.repoAccessResult = result;
        this.testingRepo = false;
      },
      error: () => {
        this.repoAccessResult = {
          accessible: false,
          isPrivate: false,
          message: 'Repository erişim testi başarısız oldu.',
          repoType: 'unknown',
          branchExists: false
        };
        this.testingRepo = false;
      }
    });
  }

  copyKey(): void {
    if (this.selectedProjectKey) {
      navigator.clipboard.writeText(this.selectedProjectKey);
      alert(this.translateService.instant('WEBHOOKS.COPIED'));
    }
  }

  // Handle project type change
  onProjectTypeChange(): void {
    const type = this.form.get('type')?.value;
    if (type === 'local') {
      this.form.get('repoUrl')?.clearValidators();
      this.form.get('branch')?.clearValidators();
      this.form.get('localPath')?.setValidators(Validators.required);
    } else {
      this.form.get('repoUrl')?.setValidators(Validators.required);
      this.form.get('branch')?.setValidators(Validators.required);
      this.form.get('localPath')?.clearValidators();
    }
    this.form.get('repoUrl')?.updateValueAndValidity();
    this.form.get('branch')?.updateValueAndValidity();
    this.form.get('localPath')?.updateValueAndValidity();
  }

  // Handle file selection
  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      this.uploadedFile = input.files[0];
    }
  }

  // Handle folder selection
  onFolderSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      this.folderFileList = input.files;
      this.selectedFolderName = this.extractFolderName(input.files[0].webkitRelativePath);

      // Build hierarchical file tree
      this.folderFiles = this.buildFileTree(Array.from(input.files));
    }
  }

  // Build hierarchical file tree from FileList
  private buildFileTree(files: File[]): FileItem[] {
    const rootItems: FileItem[] = [];

    files.forEach(file => {
      const relativePath = file.webkitRelativePath;
      const pathParts = relativePath.split('/');

      // Start from root and build tree
      let currentLevel = rootItems;
      let currentPath = '';

      for (let i = 0; i < pathParts.length; i++) {
        const part = pathParts[i];
        currentPath = currentPath ? `${currentPath}/${part}` : part;
        const isDir = i < pathParts.length - 1;

        // Find existing item at this level
        let item = currentLevel.find(x => x.name === part);

        if (!item) {
          item = {
            name: part,
            path: currentPath,
            isDir,
            excluded: this.isExcluded(currentPath, isDir),
            file: isDir ? undefined : file,
            children: isDir ? [] : undefined,
            expanded: false,
            level: i
          };
          currentLevel.push(item);
        }

        // Move to next level if directory
        if (isDir && item.children) {
          currentLevel = item.children;
        }
      }
    });

    // Sort each level: directories first, then files alphabetically
    const sortLevel = (items: FileItem[]) => {
      items.sort((a, b) => {
        if (a.isDir && !b.isDir) return -1;
        if (!a.isDir && b.isDir) return 1;
        return a.name.localeCompare(b.name);
      });
      items.forEach(item => {
        if (item.children) {
          sortLevel(item.children);
        }
      });
    };
    sortLevel(rootItems);

    return rootItems;
  }

  // Extract folder name from webkitRelativePath
  private extractFolderName(path: string): string {
    const parts = path.split('/');
    return parts[0] || 'Seçilen Klasör';
  }

  // Check if path should be excluded
  private isExcluded(path: string, isDir: boolean): boolean {
    return this.excludedPatterns.some(pattern => {
      // Check if path starts with pattern (for directory exclusion)
      if (path.includes(pattern + '/')) return true;
      if (path.endsWith(pattern)) return true;
      if (isDir && path === pattern) return true;
      return false;
    });
  }

  // Toggle file/folder exclusion (recursive)
  toggleFileExclude(item: FileItem): void {
    item.excluded = !item.excluded;

    // If it's a directory, also toggle all children recursively
    if (item.isDir && item.children) {
      this.toggleChildrenExcluded(item.children, item.excluded);
    }
  }

  // Recursively toggle children exclusion
  private toggleChildrenExcluded(children: FileItem[], excluded: boolean): void {
    children.forEach(child => {
      child.excluded = excluded;
      if (child.isDir && child.children) {
        this.toggleChildrenExcluded(child.children, excluded);
      }
    });
  }

  // Toggle exclude pattern (quick buttons)
  toggleExcludePattern(pattern: string): void {
    if (this.excludedPatterns.includes(pattern)) {
      this.excludedPatterns = this.excludedPatterns.filter(p => p !== pattern);
    } else {
      this.excludedPatterns.push(pattern);
    }

    // Update file list recursively
    this.updateExcludedPatterns(this.folderFiles);
  }

  // Recursively update excluded patterns
  private updateExcludedPatterns(items: FileItem[]): void {
    items.forEach(item => {
      item.excluded = this.isExcluded(item.path, item.isDir);
      if (item.children) {
        this.updateExcludedPatterns(item.children);
      }
    });
  }

  // Get count of included files (recursive)
  getIncludedCount(): number {
    return this.countIncluded(this.folderFiles);
  }

  // Recursively count included items
  private countIncluded(items: FileItem[]): number {
    let count = 0;
    items.forEach(item => {
      if (!item.excluded) {
        count++;
      }
      if (item.children) {
        count += this.countIncluded(item.children);
      }
    });
    return count;
  }

  // Get count of excluded files (recursive)
  getExcludedCount(): number {
    return this.countExcluded(this.folderFiles);
  }

  // Recursively count excluded items
  private countExcluded(items: FileItem[]): number {
    let count = 0;
    items.forEach(item => {
      if (item.excluded) {
        count++;
      }
      if (item.children) {
        count += this.countExcluded(item.children);
      }
    });
    return count;
  }

  // Get filtered files based on search query
  getFilteredFiles(): FileItem[] {
    if (!this.pathSearchQuery.trim()) {
      return this.folderFiles;
    }

    const query = this.pathSearchQuery.toLowerCase();
    // Return items that match search (with their parents expanded)
    return this.getMatchingItems(this.folderFiles, query);
  }

  // Recursively get matching items with their parent paths
  private getMatchingItems(items: FileItem[], query: string): FileItem[] {
    const result: FileItem[] = [];

    for (const item of items) {
      const matches = item.path.toLowerCase().includes(query) ||
        item.name.toLowerCase().includes(query);

      if (matches) {
        // Auto-expand matched directories
        if (item.isDir) {
          item.expanded = true;
        }
        result.push(item);

        // Include all children of matched items
        if (item.children) {
          result.push(...this.getAllDescendants(item.children));
        }
      } else if (item.children) {
        // Check children for matches
        const childMatches = this.getMatchingItems(item.children, query);
        if (childMatches.length > 0) {
          item.expanded = true;
          result.push(item);
          result.push(...childMatches);
        }
      }
    }

    return result;
  }

  // Get all descendants of items (flattened)
  private getAllDescendants(items: FileItem[]): FileItem[] {
    const result: FileItem[] = [];
    for (const item of items) {
      result.push(item);
      if (item.children) {
        result.push(...this.getAllDescendants(item.children));
      }
    }
    return result;
  }

  // Check if item is a root level item (no parent in the list)
  isRootItem(item: FileItem): boolean {
    return this.folderFiles.includes(item);
  }

  // Get tree rows (flattened view with only visible items)
  getTreeRows(rootItem: FileItem): FileItem[] {
    const rows: FileItem[] = [rootItem];

    if (rootItem.expanded && rootItem.children) {
      for (const child of rootItem.children) {
        rows.push(child);
        if (child.expanded && child.children) {
          rows.push(...this.getDescendantRows(child));
        }
      }
    }

    return rows;
  }

  // Get descendant rows recursively
  private getDescendantRows(item: FileItem): FileItem[] {
    const rows: FileItem[] = [];
    if (item.expanded && item.children) {
      for (const child of item.children) {
        rows.push(child);
        if (child.expanded && child.children) {
          rows.push(...this.getDescendantRows(child));
        }
      }
    }
    return rows;
  }

  // Toggle expand/collapse for directory
  toggleExpand(item: FileItem): void {
    if (item.isDir) {
      item.expanded = !item.expanded;
    }
  }

  // Handle path search change
  onPathSearchChange(query: string): void {
    this.pathSearchQuery = query;
  }

  // Recursively auto-exclude matching paths
  private autoExcludeMatchingPaths(items: FileItem[], query: string): void {
    items.forEach(item => {
      if (item.path.toLowerCase().includes(query) && !item.excluded) {
        item.excluded = true;
        // Also exclude all children
        if (item.children) {
          this.excludeAllChildren(item.children);
        }
      }
      if (item.children) {
        this.autoExcludeMatchingPaths(item.children, query);
      }
    });
  }

  // Recursively exclude all children
  private excludeAllChildren(items: FileItem[]): void {
    items.forEach(item => {
      item.excluded = true;
      if (item.children) {
        this.excludeAllChildren(item.children);
      }
    });
  }

  // Prepare and upload folder (zip on server)
  prepareAndUploadFolder(): void {
    if (!this.folderFileList || this.getIncludedCount() === 0) return;

    this.uploading = true;

    // Collect only included files from hierarchical tree
    const includedFiles = this.collectIncludedFiles(this.folderFiles);

    // Upload files (server will zip them)
    this.uploadMultipleFiles(includedFiles);
  }

  // Recursively collect included files from tree
  private collectIncludedFiles(items: FileItem[]): File[] {
    const files: File[] = [];
    items.forEach(item => {
      if (!item.excluded && !item.isDir && item.file) {
        files.push(item.file);
      }
      if (item.children) {
        files.push(...this.collectIncludedFiles(item.children));
      }
    });
    return files;
  }

  // Upload multiple files
  private uploadMultipleFiles(files: File[]): void {
    const projectName = this.form.get('name')?.value || undefined;

    // Create FormData with multiple files and their paths
    const formData = new FormData();

    // First, add all files with their relative paths
    files.forEach((file, index) => {
      formData.append('files', file);
      // Store the relative path (webkitRelativePath contains full path)
      formData.append(`file_path_${index}`, (file as any).webkitRelativePath || file.name);
    });

    formData.append('fileCount', files.length.toString());
    if (projectName) {
      formData.append('projectName', projectName);
    }

    this.uploading = true;

    this.api.uploadLocalProjectFiles(formData).subscribe({
      next: (response) => {
        this.form.patchValue({ localPath: response.localPath });
        this.uploading = false;
        alert('Klasör başarıyla yüklendi: ' + response.localPath);
        // Clear selection
        this.selectedFolderName = null;
        this.folderFiles = [];
        this.pathSearchQuery = '';
      },
      error: () => {
        this.uploading = false;
        alert('Dosya yüklenirken hata oluştu!');
      }
    });
  }

  // Legacy: Upload single file (kept for backward compatibility)
  uploadFile(): void {
    if (!this.uploadedFile) return;

    this.uploading = true;
    this.uploadProgress = 0;

    const projectName = this.form.get('name')?.value || undefined;

    this.api.uploadLocalProject(this.uploadedFile, projectName).subscribe({
      next: (response) => {
        this.form.patchValue({ localPath: response.localPath });
        this.uploading = false;
        this.uploadProgress = 100;
        alert('Dosya başarıyla yüklendi: ' + response.localPath);
      },
      error: () => {
        this.uploading = false;
        this.uploadProgress = 0;
        alert('Dosya yüklenirken hata oluştu!');
      }
    });
  }

  onSubmit(): void {
    if (this.form.invalid) return;
    this.loading = true;

    if (this.selectedProjectId) {
      this.api.updateProject(this.selectedProjectId, this.form.value as any).subscribe({
        next: (updated) => {
          const idx = this.projects.findIndex(p => p.id === updated.id);
          if (idx !== -1) this.projects[idx] = updated;
          this.showForm = false;
          this.loading = false;
          this.resetForm();
        },
        error: () => this.loading = false
      });
    } else {
      this.api.createProject(this.form.value as any).subscribe({
        next: (project) => {
          this.projects.unshift(project);
          this.showForm = false;
          this.loading = false;
          this.resetForm();
        },
        error: () => this.loading = false
      });
    }
  }

  // Delete project
  deleteProject(projectId: string, projectName: string): void {
    if (!confirm(this.translateService.instant('PROJECTS.DELETE_CONFIRM', { name: projectName }))) {
      return;
    }

    this.api.deleteProject(projectId).subscribe({
      next: () => {
        this.projects = this.projects.filter(p => p.id !== projectId);
      },
      error: () => {
        alert('Proje silinirken hata oluştu!');
      }
    });
  }
}
