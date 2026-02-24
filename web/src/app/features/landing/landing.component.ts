import { Component, OnInit, inject, PLATFORM_ID, computed } from '@angular/core';
import { CommonModule, isPlatformBrowser } from '@angular/common';
import { RouterModule } from '@angular/router';
import { ThemeService } from '../../core/services/theme.service';
import { TranslateModule } from '@ngx-translate/core';
import { LanguageSwitcherComponent } from '../../shared/components/language-switcher/language-switcher.component';

@Component({
    selector: 'app-landing',
    standalone: true,
    imports: [CommonModule, RouterModule, TranslateModule, LanguageSwitcherComponent],
    templateUrl: './landing.component.html',
    styleUrl: './landing.component.css'
})
export class LandingComponent implements OnInit {
    private platformId = inject(PLATFORM_ID);
    themeService = inject(ThemeService);

    // Sinyalden türetilmiş hesaplanmış değer: dark mı yoksa light mı?
    isDark = computed(() => this.themeService.theme() === 'dark');

    ngOnInit(): void {
        if (isPlatformBrowser(this.platformId)) {
            setTimeout(() => this.initScrollBehavior(), 100);
            setTimeout(() => this.initRevealAnimations(), 150);
            setTimeout(() => this.initTypingEffect(), 600);
        }
    }

    // Tema geçiş butonu
    toggleTheme(): void {
        this.themeService.toggleTheme();
    }

    // Mobil menü toggle
    toggleMobileMenu(): void {
        if (isPlatformBrowser(this.platformId)) {
            const menu = document.getElementById('mobile-menu');
            menu?.classList.toggle('open');
        }
    }

    // Navbar şeffaflaşması: kaydırıldığında blur arka plan
    private initScrollBehavior(): void {
        const navbar = document.getElementById('landing-navbar');
        if (!navbar) return;
        window.addEventListener('scroll', () => {
            navbar.classList.toggle('scrolled', window.scrollY > 40);
        }, { passive: true });
    }

    // Intersection Observer: staggered reveal animasyonu
    private initRevealAnimations(): void {
        const targets = document.querySelectorAll('.reveal, .stagger-group');
        const observer = new IntersectionObserver((entries, obs) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    entry.target.classList.add('visible');
                    obs.unobserve(entry.target);
                }
            });
        }, { threshold: 0.12, rootMargin: '0px 0px -60px 0px' });
        targets.forEach(el => observer.observe(el));
    }

    // Hero terminal: yazı yazma animasyonu
    private initTypingEffect(): void {
        const output = document.getElementById('terminal-output');
        if (!output) return;

        const lines = [
            { text: '$ depl@yes deploy --project myapp --server prod-01', color: 'term-cmd' },
            { text: 'Connecting to prod-01 (203.0.113.42) via SSH...', color: 'term-info' },
            { text: 'Cloning aliicolak/myapp@main...', color: 'term-info' },
            { text: 'Running deploy script: docker-compose up -d --build', color: 'term-info' },
            { text: 'Container myapp_app  Started', color: 'term-ok' },
            { text: 'Container myapp_db   Started', color: 'term-ok' },
            { text: '✔ Deployed successfully in 14.2s', color: 'term-success' },
        ];

        let lineIdx = 0;
        let charIdx = 0;
        let currentDiv: HTMLDivElement | null = null;

        const type = () => {
            if (lineIdx >= lines.length) return;

            const line = lines[lineIdx];
            if (charIdx === 0) {
                currentDiv = document.createElement('div');
                currentDiv.className = `term-line ${line.color}`;
                output.appendChild(currentDiv);
            }

            if (currentDiv && charIdx < line.text.length) {
                currentDiv.textContent += line.text[charIdx];
                charIdx++;
                setTimeout(type, lineIdx === 0 ? 35 : 18);
            } else {
                // Satır bitti: yeni satıra geç
                charIdx = 0;
                lineIdx++;
                setTimeout(type, lineIdx === 1 ? 300 : 150);
            }
        };

        type();
    }
}
