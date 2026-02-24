import { Routes } from '@angular/router';
import { authGuard, guestGuard } from './core/guards/auth.guard';

export const routes: Routes = [
    {
        path: '',
        loadComponent: () => import('./features/landing/landing.component').then(m => m.LandingComponent),
        pathMatch: 'full'
    },
    {
        path: 'login',
        loadComponent: () => import('./features/auth/login/login.component').then(m => m.LoginComponent),
        canActivate: [guestGuard]
    },
    {
        path: 'register',
        loadComponent: () => import('./features/auth/register/register.component').then(m => m.RegisterComponent),
        canActivate: [guestGuard]
    },
    {
        path: '',
        loadComponent: () => import('./shared/layouts/main-layout/main-layout.component').then(m => m.MainLayoutComponent),
        canActivate: [authGuard],
        children: [
            {
                path: 'dashboard',
                loadComponent: () => import('./features/dashboard/dashboard.component').then(m => m.DashboardComponent)
            },
            {
                path: 'projects',
                loadComponent: () => import('./features/projects/projects.component').then(m => m.ProjectsComponent)
            },
            {
                path: 'servers',
                loadComponent: () => import('./features/servers/servers.component').then(m => m.ServersComponent)
            },
            {
                path: 'deployments',
                loadComponent: () => import('./features/deployments/deployments.component').then(m => m.DeploymentsComponent)
            },
            {
                path: 'webhooks',
                loadComponent: () => import('./features/webhooks/webhooks.component').then(m => m.WebhooksComponent)
            },
            {
                path: 'terminal/:id',
                loadComponent: () => import('./features/terminal/ssh-terminal.component').then(m => m.SshTerminalComponent)
            }
        ]
    },
    {
        path: '**',
        redirectTo: 'dashboard'
    }
];
