import { Component, inject, PLATFORM_ID, OnInit } from '@angular/core';
import { CommonModule, isPlatformBrowser } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink, ActivatedRoute } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';
import { ThemeService } from '../../../core/services/theme.service';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { LanguageSwitcherComponent } from '../../../shared/components/language-switcher/language-switcher.component';

/**
 * RegisterComponent: /register rotası bu bileşeni yükler.
 * Login ve Register aynı "Sliding Overlay Glassmorphism" tasarımını kullanır.
 * Tek farkı sayfanın ilk açılışta .right-panel-active sınıfını default aktif tutmasıdır.
 */
@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterLink, TranslateModule, LanguageSwitcherComponent],
  template: `
  <div class="auth-page">
    <div class="container" [class.right-panel-active]="isRegisterMode">

    <!-- SIGN UP FORM -->
    <div class="form-container sign-up-container">
      <form [formGroup]="registerForm" (ngSubmit)="onRegisterSubmit()">
        <div class="top-actions">
          <a routerLink="/" class="icon-btn" [title]="'AUTH.HOME' | translate"><i class='bx bx-home-alt'></i></a>
          <app-language-switcher></app-language-switcher>
          <button type="button" class="icon-btn theme-btn" (click)="toggleTheme()" [title]="isDarkMode ? ('AUTH.LIGHT_MODE' | translate) : ('AUTH.DARK_MODE' | translate)">
            <i class='bx' [class.bx-sun]="isDarkMode" [class.bx-moon]="!isDarkMode"></i>
          </button>
        </div>

        <img [src]="isDarkMode ? '/logo-dark.png' : '/logo-light.png'" alt="depl@yes" class="auth-logo" />

        <h1 class="font-bold">{{ 'AUTH.CREATE_ACCOUNT' | translate }}</h1>
        <span class="subtext">{{ 'AUTH.REGISTER_SUBTITLE' | translate }}</span>

        @if (error && isRegisterMode) {
          <div class="error-box"><i class='bx bx-error-circle'></i> {{ error }}</div>
        }

        <div class="input-wrapper">
          <i class='bx bx-envelope'></i>
          <input type="email" formControlName="email" [placeholder]="'AUTH.EMAIL' | translate" />
        </div>

        <div class="input-wrapper">
          <i class='bx bx-lock-alt'></i>
          <input type="password" formControlName="password" [placeholder]="'AUTH.PASSWORD' | translate" />
        </div>

        <div class="input-wrapper">
          <i class='bx bx-lock-alt'></i>
          <input type="password" formControlName="confirmPassword" [placeholder]="'AUTH.CONFIRM_PASSWORD' | translate" />
        </div>

        <button class="primary-btn" type="submit" [disabled]="loading || registerForm.invalid">
          @if (loading && isRegisterMode) { <span class="spinner"></span> }
          {{ 'AUTH.REGISTER' | translate }}
        </button>
      </form>
    </div>

    <!-- SIGN IN FORM -->
    <div class="form-container sign-in-container">
      <form [formGroup]="loginForm" (ngSubmit)="onLoginSubmit()">
        <div class="top-actions">
          <a routerLink="/" class="icon-btn" [title]="'AUTH.HOME' | translate"><i class='bx bx-home-alt'></i></a>
          <app-language-switcher></app-language-switcher>
          <button type="button" class="icon-btn theme-btn" (click)="toggleTheme()" [title]="isDarkMode ? ('AUTH.LIGHT_MODE' | translate) : ('AUTH.DARK_MODE' | translate)">
            <i class='bx' [class.bx-sun]="isDarkMode" [class.bx-moon]="!isDarkMode"></i>
          </button>
        </div>

        <img [src]="isDarkMode ? '/logo-dark.png' : '/logo-light.png'" alt="depl@yes" class="auth-logo" />

        <h1 class="font-bold">{{ 'AUTH.SIGN_IN' | translate }}</h1>
        <span class="subtext">{{ 'AUTH.LOGIN_SUBTITLE' | translate }}</span>

        @if (error && !isRegisterMode) {
          <div class="error-box"><i class='bx bx-error-circle'></i> {{ error }}</div>
        }

        <div class="input-wrapper">
          <i class='bx bx-envelope'></i>
          <input type="email" formControlName="email" [placeholder]="'AUTH.EMAIL' | translate" />
        </div>

        <div class="input-wrapper">
          <i class='bx bx-lock-alt'></i>
          <input type="password" formControlName="password" [placeholder]="'AUTH.PASSWORD' | translate" />
        </div>

        <a routerLink="/forgot-password" class="forgot-password">{{ 'AUTH.FORGOT_PASSWORD' | translate }}</a>

        <button class="primary-btn" type="submit" [disabled]="loading || loginForm.invalid">
          @if (loading && !isRegisterMode) { <span class="spinner"></span> }
          {{ 'AUTH.LOGIN' | translate }}
        </button>
      </form>
    </div>

        <!-- OVERLAY PANEL -->
        <div class="overlay-container">
            <div class="overlay">

                <div class="overlay-panel overlay-left">
                    <h1>{{ 'AUTH.WELCOME_BACK' | translate }}</h1>
                    <p>{{ 'AUTH.WELCOME_BACK_DESC' | translate }}</p>
                    <button class="ghost-btn" (click)="switchMode()">{{ 'AUTH.LOGIN' | translate }}</button>
                </div>

                <div class="overlay-panel overlay-right">
                    <h1>{{ 'AUTH.HELLO_FRIEND' | translate }}</h1>
                    <p>{{ 'AUTH.HELLO_FRIEND_DESC' | translate }}</p>
                    <button class="ghost-btn" (click)="switchMode()">{{ 'AUTH.REGISTER' | translate }}</button>
                </div>

            </div>
        </div>
    </div>
</div>
`,
  styleUrl: './register.component.css'
})
export class RegisterComponent {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);
  private themeService = inject(ThemeService);
  private router = inject(Router);
  private platformId = inject(PLATFORM_ID);
  private translateService = inject(TranslateService);

  isRegisterMode = true; // /register pathinde doğrudan bu panel aktif
  loading = false;
  error = '';

  loginForm = this.fb.group({
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(6)]]
  });

  registerForm = this.fb.group({
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(6)]],
    confirmPassword: ['', [Validators.required]]
  });

  get isDarkMode(): boolean {
    return this.themeService.theme() === 'dark';
  }

  toggleTheme(): void {
    this.themeService.toggleTheme();
  }

  switchMode(): void {
    this.isRegisterMode = !this.isRegisterMode;
    this.error = '';
    const newPath = this.isRegisterMode ? '/register' : '/login';
    window.history.replaceState({}, '', newPath);
  }

  onLoginSubmit(): void {
    if (this.loginForm.invalid) return;
    this.loading = true;
    this.error = '';

    const { email, password } = this.loginForm.value;
    this.authService.login({ email: email!, password: password! }).subscribe({
      next: () => this.router.navigate(['/dashboard']),
      error: (err) => {
        this.loading = false;
        this.error = err.error || this.translateService.instant('AUTH.LOGIN_FAILED');
      }
    });
  }

  onRegisterSubmit(): void {
    if (this.registerForm.invalid) return;

    const { email, password, confirmPassword } = this.registerForm.value;
    if (password !== confirmPassword) {
      this.error = this.translateService.instant('AUTH.PASSWORDS_MISMATCH');
      return;
    }

    this.loading = true;
    this.error = '';

    this.authService.register({ email: email!, password: password! }).subscribe({
      next: () => this.router.navigate(['/dashboard']),
      error: (err) => {
        this.loading = false;
        this.error = err.error || this.translateService.instant('AUTH.REGISTER_FAILED');
      }
    });
  }
}
