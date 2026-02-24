import { HttpInterceptorFn, HttpErrorResponse, HttpRequest, HttpHandlerFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, throwError, switchMap, filter, take } from 'rxjs';
import { AuthService } from '../services/auth.service';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
    const authService = inject(AuthService);
    const router = inject(Router);

    // Statik varlık isteklerini (i18n JSON dosyaları vb.) auth sürecinden hariç tut
    if (req.url.includes('/assets/') || req.url.endsWith('.json')) {
        return next(req);
    }

    const token = authService.getToken();

    if (token) {
        req = req.clone({
            setHeaders: {
                Authorization: `Bearer ${token}`
            }
        });
    }

    return next(req).pipe(
        catchError((error: HttpErrorResponse) => {
            // Don't try to refresh if this is a refresh request or auth request
            if (error.status === 401 && !req.url.includes('/auth/')) {
                const refreshToken = authService.getRefreshToken();

                if (refreshToken && !authService.getIsRefreshing()) {
                    return authService.refreshAccessToken().pipe(
                        switchMap(() => {
                            // Retry the original request with new token
                            const newToken = authService.getToken();
                            const clonedReq = req.clone({
                                setHeaders: {
                                    Authorization: `Bearer ${newToken}`
                                }
                            });
                            return next(clonedReq);
                        }),
                        catchError((refreshError) => {
                            authService.logout();
                            router.navigate(['/login']);
                            return throwError(() => refreshError);
                        })
                    );
                } else if (authService.getIsRefreshing()) {
                    // Wait for the refresh to complete
                    return authService.getRefreshTokenSubject().pipe(
                        filter(token => token !== null),
                        take(1),
                        switchMap(() => {
                            const newToken = authService.getToken();
                            const clonedReq = req.clone({
                                setHeaders: {
                                    Authorization: `Bearer ${newToken}`
                                }
                            });
                            return next(clonedReq);
                        })
                    );
                } else {
                    authService.logout();
                    router.navigate(['/login']);
                }
            }
            return throwError(() => error);
        })
    );
};

