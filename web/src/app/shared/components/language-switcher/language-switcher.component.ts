import { Component, inject, OnInit, PLATFORM_ID, HostListener, ElementRef } from '@angular/core';
import { CommonModule, isPlatformBrowser } from '@angular/common';
import { TranslateService } from '@ngx-translate/core';

interface Language {
    code: string;
    flag: string;
    label: string;
}

@Component({
    selector: 'app-language-switcher',
    standalone: true,
    imports: [CommonModule],
    template: `
    <div class="lang-switcher" (click)="toggleDropdown($event)">
      <span class="current-lang">{{ currentFlag }} {{ currentCode }}</span>
      <span class="arrow">▾</span>
      @if (open) {
        <div class="lang-dropdown">
          @for (lang of languages; track lang.code) {
            <button class="lang-option" [class.active]="lang.code === currentCode"
              (click)="switchLang(lang, $event)">
              <span class="lang-flag">{{ lang.flag }}</span>
              <span class="lang-label">{{ lang.label }}</span>
            </button>
          }
        </div>
      }
    </div>
  `,
    styles: [`
    .lang-switcher {
      position: relative;
      display: flex;
      align-items: center;
      gap: 0.25rem;
      padding: 0.375rem 0.625rem;
      font-size: 0.8125rem;
      font-weight: 500;
      color: var(--text-secondary);
      background: transparent;
      border: 1px solid var(--border-color);
      border-radius: var(--radius-md, 8px);
      cursor: pointer;
      transition: all 0.15s ease;
      user-select: none;
    }

    .lang-switcher:hover {
      background-color: var(--bg-hover);
      border-color: var(--accent);
      color: var(--text-primary);
    }

    .current-lang {
      display: flex;
      align-items: center;
      gap: 0.375rem;
    }

    .arrow {
      font-size: 0.625rem;
      opacity: 0.6;
    }

    .lang-dropdown {
      position: absolute;
      top: calc(100% + 4px);
      right: 0;
      min-width: 140px;
      background: var(--bg-card, #fff);
      border: 1px solid var(--border-color);
      border-radius: var(--radius-md, 8px);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
      z-index: 1000;
      overflow: hidden;
    }

    .lang-option {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      width: 100%;
      padding: 0.5rem 0.75rem;
      font-size: 0.8125rem;
      color: var(--text-secondary);
      background: transparent;
      border: none;
      cursor: pointer;
      transition: all 0.1s ease;
    }

    .lang-option:hover {
      background: var(--bg-hover);
      color: var(--text-primary);
    }

    .lang-option.active {
      background: var(--accent-light);
      color: var(--accent);
      font-weight: 600;
    }

    .lang-flag {
      font-size: 1rem;
    }

    .lang-label {
      flex: 1;
    }
  `]
})
export class LanguageSwitcherComponent implements OnInit {
    private translate = inject(TranslateService);
    private platformId = inject(PLATFORM_ID);
    private elementRef = inject(ElementRef);
    open = false;

    languages: Language[] = [
        { code: 'en', flag: '🇺🇸', label: 'English' },
        { code: 'tr', flag: '🇹🇷', label: 'Türkçe' },
        { code: 'zh', flag: '🇨🇳', label: '中文' },
        { code: 'es', flag: '🇪🇸', label: 'Español' }
    ];

    currentCode = 'en';
    currentFlag = '🇺🇸';

    // Use HostListener instead of document.addEventListener for proper Angular change detection
    @HostListener('document:click', ['$event'])
    onDocumentClick(event: Event): void {
        if (!this.elementRef.nativeElement.contains(event.target)) {
            this.open = false;
        }
    }

    ngOnInit(): void {
        let savedLang = 'en';

        if (isPlatformBrowser(this.platformId)) {
            savedLang = localStorage.getItem('app_lang') || this.detectBrowserLanguage();
        }

        this.setLanguage(savedLang);
    }

    private detectBrowserLanguage(): string {
        if (!isPlatformBrowser(this.platformId)) return 'en';
        const browserLang = navigator.language?.slice(0, 2)?.toLowerCase() || 'en';
        const supported = this.languages.map(l => l.code);
        return supported.includes(browserLang) ? browserLang : 'en';
    }

    toggleDropdown(event: Event): void {
        event.stopPropagation();
        this.open = !this.open;
    }

    switchLang(lang: Language, event: Event): void {
        event.stopPropagation();
        this.setLanguage(lang.code);
        this.open = false;
        if (isPlatformBrowser(this.platformId)) {
            localStorage.setItem('app_lang', lang.code);
        }
    }

    private setLanguage(code: string): void {
        const lang = this.languages.find(l => l.code === code) || this.languages[0];
        this.currentCode = lang.code.toUpperCase();
        this.currentFlag = lang.flag;
        this.translate.use(lang.code);
    }
}
