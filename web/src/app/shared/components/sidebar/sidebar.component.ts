import { Component } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [RouterLink, RouterLinkActive, TranslateModule],
  template: `
    <aside class="sidebar">
      <nav class="sidebar-nav">
        <a routerLink="/dashboard" routerLinkActive="active" [routerLinkActiveOptions]="{exact: true}" class="nav-item">
          <span class="nav-icon">▣</span>
          <span class="nav-text">{{ 'NAV.DASHBOARD' | translate }}</span>
        </a>
        <a routerLink="/projects" routerLinkActive="active" class="nav-item">
          <span class="nav-icon">◫</span>
          <span class="nav-text">{{ 'NAV.PROJECTS' | translate }}</span>
        </a>
        <a routerLink="/servers" routerLinkActive="active" class="nav-item">
          <span class="nav-icon">▤</span>
          <span class="nav-text">{{ 'NAV.SERVERS' | translate }}</span>
        </a>
        <a routerLink="/deployments" routerLinkActive="active" class="nav-item">
          <span class="nav-icon">▷</span>
          <span class="nav-text">{{ 'NAV.DEPLOYMENTS' | translate }}</span>
        </a>
        <a routerLink="/webhooks" routerLinkActive="active" class="nav-item">
          <span class="nav-icon">◇</span>
          <span class="nav-text">{{ 'NAV.WEBHOOKS' | translate }}</span>
        </a>
      </nav>
    </aside>
  `,
  styles: [`
    .sidebar {
      width: 220px;
      height: calc(100vh - 60px);
      position: fixed;
      top: 60px;
      left: 0;
      background-color: var(--bg-card);
      border-right: 1px solid var(--border-color);
      padding: 1rem 0.75rem;
    }

    .sidebar-nav {
      display: flex;
      flex-direction: column;
      gap: 2px;
    }

    .nav-item {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      padding: 0.625rem 0.75rem;
      border-radius: var(--radius-md);
      text-decoration: none;
      color: var(--text-secondary);
      font-size: 0.875rem;
      font-weight: 500;
      transition: all var(--transition-fast);
    }

    .nav-item:hover {
      background-color: var(--bg-hover);
      color: var(--text-primary);
    }

    .nav-item.active {
      background-color: var(--accent-light);
      color: var(--accent);
    }

    .nav-icon {
      font-size: 1rem;
      width: 1.25rem;
      text-align: center;
      opacity: 0.8;
    }

    .nav-item.active .nav-icon {
      opacity: 1;
    }

    .nav-text {
      flex: 1;
    }
  `]
})
export class SidebarComponent { }
