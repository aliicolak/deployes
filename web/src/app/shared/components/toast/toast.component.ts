import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ToastService } from '../../../core/services/toast.service';

@Component({
    selector: 'app-toast',
    standalone: true,
    imports: [CommonModule],
    template: `
    <div class="toast-container">
      @for (toast of toastService.toasts(); track toast.id) {
        <div class="toast" [ngClass]="toast.type" (click)="toastService.remove(toast.id)">
          <span class="message">{{ toast.message }}</span>
          <span class="close">&times;</span>
        </div>
      }
    </div>
  `,
    styles: [`
    .toast-container {
      position: fixed;
      top: 1rem;
      right: 1rem;
      z-index: 9999;
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
      pointer-events: none;
    }

    .toast {
      padding: 1rem 1.5rem;
      border-radius: 8px;
      background: var(--bg-secondary, #1f2937);
      color: white;
      box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 1rem;
      min-width: 300px;
      max-width: 400px;
      pointer-events: auto;
      cursor: pointer;
      animation: slideIn 0.3s ease-out;
      border-left: 4px solid transparent;
    }

    .toast.success { border-left-color: #10b981; background: #064e3b; }
    .toast.error { border-left-color: #ef4444; background: #450a0a; } 
    .toast.warning { border-left-color: #f59e0b; background: #451a03; }
    .toast.info { border-left-color: #3b82f6; background: #172554; }

    .close {
      opacity: 0.7;
      font-size: 1.25rem;
    }

    @keyframes slideIn {
      from { transform: translateX(100%); opacity: 0; }
      to { transform: translateX(0); opacity: 1; }
    }
  `]
})
export class ToastComponent {
    toastService = inject(ToastService);
}
