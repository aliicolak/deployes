import { Component, inject } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';
import { ThemeService } from '../../../core/services/theme.service';
import { TranslateModule } from '@ngx-translate/core';
import { LanguageSwitcherComponent } from '../language-switcher/language-switcher.component';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [RouterLink, TranslateModule, LanguageSwitcherComponent],
  template: `
    <header class="navbar">
      <div class="navbar-brand">
        <a routerLink="/dashboard" class="navbar-logo">
          <img 
            [src]="themeService.theme() === 'dark' ? '/logo-dark.png' : '/logo-light.png'" 
            alt="depl@yes" 
            class="logo-image"
            [class.dark-mode]="themeService.theme() === 'dark'"
          />
          <span class="logo-text">Depl&#64;yes</span>
        </a>
      </div>
      <div class="navbar-actions">
        <app-language-switcher></app-language-switcher>
        <button class="navbar-btn theme-toggle" (click)="toggleTheme()" [title]="'NAV.TOGGLE_THEME' | translate">
          <span class="icon">{{ themeService.theme() === 'light' ? '🌙' : '☀️' }}</span>
        </button>
        <button class="navbar-btn logout-btn" (click)="logout()">
          <span>{{ 'NAV.LOGOUT' | translate }}</span>
        </button>
      </div>
    </header>
  `,
  styles: [`
    .navbar {
      display: flex;
      align-items: center;
      justify-content: space-between;
      height: 60px;
      padding: 0 1.5rem;
      background-color: var(--bg-card);
      border-bottom: 1px solid var(--border-color);
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      z-index: 100;
    }

    .navbar-brand {
      display: flex;
      align-items: center;
    }

    .navbar-logo {
      display: flex;
      align-items: center;
      text-decoration: none;
    }

    .logo-image {
      height: 55px;
      width: auto;
      object-fit: contain;
      border-radius: 6px;
    }

    .logo-text {
      margin-left: 0.75rem;
      font-size: 1.5rem;
      font-weight: 700;
      color: var(--text-primary);
      letter-spacing: -0.025em;
    }

    .logo-image.dark-mode {
      mix-blend-mode: screen;
    }

    .navbar-actions {
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .navbar-btn {
      display: flex;
      align-items: center;
      gap: 0.375rem;
      padding: 0.5rem 0.75rem;
      font-size: 0.875rem;
      font-weight: 500;
      color: var(--text-secondary);
      background: transparent;
      border: none;
      border-radius: var(--radius-md);
      cursor: pointer;
      transition: all var(--transition-fast);
    }

    .navbar-btn:hover {
      background-color: var(--bg-hover);
      color: var(--text-primary);
    }

    .theme-toggle {
      font-size: 1.25rem;
    }

    .logout-btn {
      background-color: var(--danger-light);
      color: var(--danger);
    }

    .logout-btn:hover {
      background-color: var(--danger);
      color: white;
    }

    .icon {
      font-size: 1rem;
    }
  `]
})
export class NavbarComponent {
  authService = inject(AuthService);
  themeService = inject(ThemeService);
  private router = inject(Router);

  toggleTheme(): void {
    this.themeService.toggleTheme();
  }

  logout(): void {
    this.authService.logout();
  }
}
