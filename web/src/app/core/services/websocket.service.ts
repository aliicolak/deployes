import { Injectable, inject, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { Subject, Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { AuthService } from './auth.service';
import { LogMessage } from '../models';

@Injectable({
    providedIn: 'root'
})
export class WebSocketService {
    private authService = inject(AuthService);
    private platformId = inject(PLATFORM_ID);
    private socket: WebSocket | null = null;
    private messagesSubject = new Subject<LogMessage>();

    connect(deploymentId: string): Observable<LogMessage> {
        if (!isPlatformBrowser(this.platformId)) {
            return this.messagesSubject.asObservable();
        }

        const token = this.authService.getToken();
        if (!token) {
            console.error('No auth token available for WebSocket connection');
            return this.messagesSubject.asObservable();
        }

        const wsUrl = `${environment.wsUrl}/deployments/${deploymentId}/logs/stream?token=${token}`;

        this.disconnect(); // Close any existing connection

        this.socket = new WebSocket(wsUrl);

        this.socket.onopen = () => {
            console.log('WebSocket connected for deployment:', deploymentId);
        };

        this.socket.onmessage = (event) => {
            try {
                const message: LogMessage = JSON.parse(event.data);
                this.messagesSubject.next(message);
            } catch (e) {
                console.error('Failed to parse WebSocket message:', e);
            }
        };

        this.socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        this.socket.onclose = (event) => {
            console.log('WebSocket closed:', event.code, event.reason);
        };

        return this.messagesSubject.asObservable();
    }

    disconnect(): void {
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }
    }

    isConnected(): boolean {
        return this.socket?.readyState === WebSocket.OPEN;
    }
}
