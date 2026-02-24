import { Injectable, signal, effect, PLATFORM_ID, inject } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';

export type Theme = 'light' | 'dark';

@Injectable({
    providedIn: 'root'
})
export class ThemeService {
    private platformId = inject(PLATFORM_ID);
    theme = signal<Theme>('light');

    constructor() {
        if (isPlatformBrowser(this.platformId)) {
            // Load saved theme or default to system preference
            const savedTheme = localStorage.getItem('theme') as Theme;
            if (savedTheme) {
                this.theme.set(savedTheme);
            } else if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
                this.theme.set('dark');
            }
        }

        // Effect to update DOM and storage when theme changes
        effect(() => {
            const currentTheme = this.theme();
            if (isPlatformBrowser(this.platformId)) {
                document.documentElement.setAttribute('data-theme', currentTheme);
                localStorage.setItem('theme', currentTheme);
            }
        });
    }

    toggleTheme() {
        this.theme.update(t => t === 'light' ? 'dark' : 'light');
    }
}
