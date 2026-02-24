import { HttpInterceptorFn, HttpErrorResponse } from '@angular/common/http';
import { inject } from '@angular/core';
import { ToastService } from '../services/toast.service';
import { catchError, throwError } from 'rxjs';

export const errorInterceptor: HttpInterceptorFn = (req, next) => {
    const toastService = inject(ToastService);

    // Statik varlık isteklerini (i18n JSON dosyaları vb.) hata yönetiminden hariç tut
    if (req.url.includes('/assets/') || req.url.endsWith('.json')) {
        return next(req);
    }

    return next(req).pipe(
        catchError((error: HttpErrorResponse) => {
            let errorMessage = 'Bir hata oluştu';

            if (error.error instanceof ErrorEvent) {
                // Client-side error
                errorMessage = error.error.message;
            } else {
                // Server-side error
                // Backend returns string ("failed to ...") or json
                if (typeof error.error === 'string') {
                    errorMessage = error.error;
                } else if (error.error?.message) {
                    errorMessage = error.error.message;
                } else {
                    errorMessage = `Hata Kodu: ${error.status} - ${error.statusText}`;
                }
            }

            // Don't show toasts for 401s - handled by auth interceptor
            if (error.status === 401) {
                return throwError(() => error);
            }

            toastService.error(errorMessage);
            return throwError(() => error);
        })
    );
};
