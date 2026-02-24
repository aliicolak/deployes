export interface User {
    id: string;
    email: string;
    createdAt?: string;
}

export interface LoginRequest {
    email: string;
    password: string;
}

export interface RegisterRequest {
    email: string;
    password: string;
}

export interface AuthResponse {
    accessToken: string;
    refreshToken: string;
    user?: User;
}

export type ProjectType = 'github' | 'local';

export interface Project {
    id: string;
    name: string;
    type: ProjectType;
    repoUrl?: string;
    branch?: string;
    localPath?: string;
    deployScript: string;
    includePatterns?: string;
    excludePatterns?: string;
    preservePatterns?: string;
    scmPublicKey?: string;
    createdAt?: string;
}

export interface CreateProjectRequest {
    name: string;
    type: ProjectType;
    repoUrl?: string;
    branch?: string;
    localPath?: string;
    deployScript: string;
    includePatterns?: string;
    excludePatterns?: string;
    preservePatterns?: string;
}

export interface Server {
    id: string;
    name: string;
    host: string;
    port: number;
    username: string;
    createdAt?: string;
}

export interface CreateServerRequest {
    name: string;
    host: string;
    port: number;
    username: string;
    sshKey: string;
}

export interface Deployment {
    id: string;
    projectId: string;
    serverId: string;
    serverIds?: string[]; // For multi-server deployments
    status: 'queued' | 'running' | 'success' | 'failed';
    logs: string;
    createdAt?: string;
    startedAt?: string;
    finishedAt?: string;
    commitHash?: string;
    rollbackFromId?: string;
}

export interface CreateDeploymentRequest {
    projectId: string;
    serverId: string;
}

export interface Webhook {
    id: string;
    projectId: string;
    serverIds: string[]; // Changed from serverId to support multiple servers
    secret?: string;
    webhookUrl: string;
    isActive: boolean;
    createdAt?: string;
}

export interface CreateWebhookRequest {
    projectId: string;
    serverIds: string[]; // Changed from serverId to support multiple servers
}

export interface LogMessage {
    deploymentId: string;
    message: string;
    timestamp: number;
}

export interface Secret {
    id: string;
    key: string;
    createdAt: string;
}

export interface CreateSecretRequest {
    projectId: string;
    key: string;
    value: string;
}

export interface Stats {
    total: number;
    successful: number;
    failed: number;
    averageDuration: number;
    last7Days: {
        dates: string[];
        counts: number[];
    };
}
