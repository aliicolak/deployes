import { Injectable, signal, computed, PLATFORM_ID, inject } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable, tap, catchError, throwError, BehaviorSubject, filter, take, switchMap } from 'rxjs';
import { environment } from '../../../environments/environment';
import { LoginRequest, RegisterRequest, AuthResponse, User } from '../models';

@Injectable({
    providedIn: 'root'
})
export class AuthService {
    private http = inject(HttpClient);
    private router = inject(Router);
    private platformId = inject(PLATFORM_ID);

    private tokenSignal = signal<string | null>(null);
    private refreshTokenSignal = signal<string | null>(null);
    private userSignal = signal<User | null>(null);

    private isRefreshing = false;
    private refreshTokenSubject: BehaviorSubject<string | null> = new BehaviorSubject<string | null>(null);

    isAuthenticated = computed(() => !!this.tokenSignal());
    currentUser = computed(() => this.userSignal());
    token = computed(() => this.tokenSignal());

    constructor() {
        if (isPlatformBrowser(this.platformId)) {
            const storedToken = localStorage.getItem('token');
            const storedRefreshToken = localStorage.getItem('refreshToken');
            if (storedToken) {
                this.tokenSignal.set(storedToken);
            }
            if (storedRefreshToken) {
                this.refreshTokenSignal.set(storedRefreshToken);
            }
        }
    }

    register(request: RegisterRequest): Observable<AuthResponse> {
        return this.http.post<AuthResponse>(`${environment.apiUrl}/auth/register`, request).pipe(
            tap(response => this.handleAuthResponse(response))
        );
    }

    login(request: LoginRequest): Observable<AuthResponse> {
        return this.http.post<AuthResponse>(`${environment.apiUrl}/auth/login`, request).pipe(
            tap(response => this.handleAuthResponse(response))
        );
    }

    logout(): void {
        if (isPlatformBrowser(this.platformId)) {
            localStorage.removeItem('token');
            localStorage.removeItem('refreshToken');
        }
        this.tokenSignal.set(null);
        this.refreshTokenSignal.set(null);
        this.userSignal.set(null);
        this.router.navigate(['/login']);
    }

    refreshAccessToken(): Observable<AuthResponse> {
        const refreshToken = this.refreshTokenSignal();
        if (!refreshToken) {
            return throwError(() => new Error('No refresh token available'));
        }

        if (!this.isRefreshing) {
            this.isRefreshing = true;
            this.refreshTokenSubject.next(null);

            return this.http.post<AuthResponse>(`${environment.apiUrl}/auth/refresh`, { refreshToken }).pipe(
                tap(response => {
                    this.isRefreshing = false;
                    this.handleAuthResponse(response);
                    this.refreshTokenSubject.next(response.accessToken);
                }),
                catchError(error => {
                    this.isRefreshing = false;
                    this.logout();
                    return throwError(() => error);
                })
            );
        } else {
            return this.refreshTokenSubject.pipe(
                filter(token => token !== null),
                take(1),
                switchMap(() => {
                    return this.http.post<AuthResponse>(`${environment.apiUrl}/auth/refresh`, { refreshToken });
                })
            );
        }
    }

    getIsRefreshing(): boolean {
        return this.isRefreshing;
    }

    getRefreshTokenSubject(): BehaviorSubject<string | null> {
        return this.refreshTokenSubject;
    }

    private handleAuthResponse(response: AuthResponse): void {
        if (response.accessToken) {
            if (isPlatformBrowser(this.platformId)) {
                localStorage.setItem('token', response.accessToken);
                if (response.refreshToken) {
                    localStorage.setItem('refreshToken', response.refreshToken);
                }
            }
            this.tokenSignal.set(response.accessToken);
            if (response.refreshToken) {
                this.refreshTokenSignal.set(response.refreshToken);
            }
            if (response.user) {
                this.userSignal.set(response.user);
            }
        }
    }

    getToken(): string | null {
        return this.tokenSignal();
    }

    getRefreshToken(): string | null {
        return this.refreshTokenSignal();
    }
}

