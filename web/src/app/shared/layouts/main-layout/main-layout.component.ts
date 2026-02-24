import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { NavbarComponent } from '../../components/navbar/navbar.component';
import { SidebarComponent } from '../../components/sidebar/sidebar.component';

@Component({
  selector: 'app-main-layout',
  standalone: true,
  imports: [RouterOutlet, NavbarComponent, SidebarComponent],
  template: `
    <div class="layout">
      <app-navbar></app-navbar>
      <div class="layout-body">
        <app-sidebar></app-sidebar>
        <main class="main-content">
          <router-outlet></router-outlet>
        </main>
      </div>
    </div>
  `,
  styles: [`
    .layout {
      min-height: 100vh;
      background-color: var(--bg-primary);
    }

    .layout-body {
      display: flex;
      padding-top: 60px;
    }

    .main-content {
      flex: 1;
      margin-left: 220px;
      padding: 1.5rem 2rem;
      min-height: calc(100vh - 60px);
    }
  `]
})
export class MainLayoutComponent { }
