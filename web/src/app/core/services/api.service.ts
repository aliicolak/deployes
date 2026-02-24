import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import {
    Project, CreateProjectRequest,
    Server, CreateServerRequest,
    Deployment, CreateDeploymentRequest,
    Webhook, CreateWebhookRequest,
    Secret, CreateSecretRequest,
    Stats
} from '../models';

export interface UploadResponse {
    localPath: string;
    message: string;
}

@Injectable({
    providedIn: 'root'
})
export class ApiService {
    private http = inject(HttpClient);
    private apiUrl = environment.apiUrl;

    // Projects
    getProjects(): Observable<Project[]> {
        return this.http.get<Project[]>(`${this.apiUrl}/projects`);
    }

    createProject(project: CreateProjectRequest): Observable<Project> {
        return this.http.post<Project>(`${this.apiUrl}/projects`, project);
    }

    updateProject(id: string, project: CreateProjectRequest): Observable<Project> {
        return this.http.put<Project>(`${this.apiUrl}/projects?id=${id}`, project);
    }

    // Servers
    getServers(): Observable<Server[]> {
        return this.http.get<Server[]>(`${this.apiUrl}/servers`);
    }

    createServer(server: CreateServerRequest): Observable<Server> {
        return this.http.post<Server>(`${this.apiUrl}/servers`, server);
    }

    updateServer(id: string, server: CreateServerRequest): Observable<Server> {
        return this.http.put<Server>(`${this.apiUrl}/servers?id=${id}`, server);
    }

    // Deployments
    getDeployments(): Observable<Deployment[]> {
        return this.http.get<Deployment[]>(`${this.apiUrl}/deployments/list`);
    }

    getDeployment(id: string): Observable<Deployment> {
        return this.http.get<Deployment>(`${this.apiUrl}/deployments?id=${id}`);
    }

    createDeployment(deployment: CreateDeploymentRequest): Observable<Deployment> {
        return this.http.post<Deployment>(`${this.apiUrl}/deployments`, deployment);
    }

    rollbackDeployment(id: string): Observable<Deployment> {
        return this.http.post<Deployment>(`${this.apiUrl}/deployments/rollback?id=${id}`, {});
    }

    getStats(): Observable<Stats> {
        return this.http.get<Stats>(`${this.apiUrl}/dashboard/stats`);
    }

    // Webhooks
    getWebhooks(): Observable<Webhook[]> {
        return this.http.get<Webhook[]>(`${this.apiUrl}/webhooks`);
    }

    createWebhook(webhook: CreateWebhookRequest): Observable<Webhook> {
        return this.http.post<Webhook>(`${this.apiUrl}/webhooks`, webhook);
    }

    updateWebhook(id: string, data: { isActive: boolean }): Observable<Webhook> {
        return this.http.put<Webhook>(`${this.apiUrl}/webhooks?id=${id}`, data);
    }

    deleteWebhook(id: string): Observable<void> {
        return this.http.delete<void>(`${this.apiUrl}/webhooks?id=${id}`);
    }

    // Server Connection Test
    testServerConnection(data: { host?: string; port?: number; username?: string; sshKey?: string; serverId?: string }): Observable<{ success: boolean; message: string; latency: number }> {
        return this.http.post<{ success: boolean; message: string; latency: number }>(`${this.apiUrl}/servers/test-connection`, data);
    }

    // Encryption Status
    getEncryptionStatus(): Observable<{ active: boolean; algorithm: string; keyLength: number }> {
        return this.http.get<{ active: boolean; algorithm: string; keyLength: number }>(`${this.apiUrl}/encryption/status`);
    }

    // Test Repository Access
    testRepoAccess(data: { repoUrl: string; branch?: string }): Observable<{
        accessible: boolean;
        isPrivate: boolean;
        message: string;
        guidance?: string;
        repoType: string;
        branchExists: boolean
    }> {
        return this.http.post<{ accessible: boolean; isPrivate: boolean; message: string; guidance?: string; repoType: string; branchExists: boolean }>(`${this.apiUrl}/projects/test-access`, data);
    }

    // Secrets
    getSecrets(projectId: string): Observable<Secret[]> {
        return this.http.get<Secret[]>(`${this.apiUrl}/secrets?projectId=${projectId}`);
    }

    createSecret(secret: CreateSecretRequest): Observable<Secret> {
        return this.http.post<Secret>(`${this.apiUrl}/secrets`, secret);
    }

    deleteSecret(id: string): Observable<void> {
        return this.http.delete<void>(`${this.apiUrl}/secrets?id=${id}`);
    }

    // Upload Local Project
    uploadLocalProject(file: File, projectName?: string): Observable<UploadResponse> {
        const formData = new FormData();
        formData.append('file', file);
        if (projectName) {
            formData.append('projectName', projectName);
        }
        return this.http.post<UploadResponse>(`${this.apiUrl}/projects/upload`, formData);
    }

    // Upload Local Project Files (multiple files for folder upload)
    uploadLocalProjectFiles(formData: FormData): Observable<UploadResponse> {
        return this.http.post<UploadResponse>(`${this.apiUrl}/projects/upload`, formData);
    }

    // Delete Project
    deleteProject(id: string): Observable<void> {
        return this.http.delete<void>(`${this.apiUrl}/projects/delete?id=${id}`);
    }
}

