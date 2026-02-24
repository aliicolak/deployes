import { Component, Input, OnChanges, SimpleChanges, ElementRef, ViewChild, AfterViewChecked } from '@angular/core';

@Component({
    selector: 'app-terminal',
    standalone: true,
    template: `
    <div class="terminal" #terminalContainer>
      <div class="terminal-header">
        <div class="terminal-dots">
          <span class="dot red"></span>
          <span class="dot yellow"></span>
          <span class="dot green"></span>
        </div>
        <span class="terminal-title">{{ title }}</span>
      </div>
      <div class="terminal-body" #terminalBody>
        <pre class="terminal-output">{{ logs }}</pre>
      </div>
    </div>
  `,
    styles: [`
    .terminal {
      background: #0d1117;
      border-radius: 0.75rem;
      overflow: hidden;
      border: 1px solid rgba(255, 255, 255, 0.1);
      font-family: 'JetBrains Mono', 'Fira Code', 'Monaco', monospace;
    }

    .terminal-header {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 0.75rem 1rem;
      background: rgba(255, 255, 255, 0.05);
      border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    }

    .terminal-dots {
      display: flex;
      gap: 0.5rem;
    }

    .dot {
      width: 12px;
      height: 12px;
      border-radius: 50%;
    }

    .dot.red { background: #ff5f56; }
    .dot.yellow { background: #ffbd2e; }
    .dot.green { background: #27ca40; }

    .terminal-title {
      color: rgba(255, 255, 255, 0.6);
      font-size: 0.8rem;
    }

    .terminal-body {
      padding: 1rem;
      max-height: 400px;
      overflow-y: auto;
    }

    .terminal-output {
      margin: 0;
      color: #c9d1d9;
      font-size: 0.875rem;
      line-height: 1.6;
      white-space: pre-wrap;
      word-break: break-word;
    }

    .terminal-body::-webkit-scrollbar {
      width: 8px;
    }

    .terminal-body::-webkit-scrollbar-track {
      background: transparent;
    }

    .terminal-body::-webkit-scrollbar-thumb {
      background: rgba(255, 255, 255, 0.2);
      border-radius: 4px;
    }
  `]
})
export class TerminalComponent implements AfterViewChecked {
    @Input() logs: string = '';
    @Input() title: string = 'Deployment Logs';
    @ViewChild('terminalBody') terminalBody!: ElementRef;

    private shouldScroll = true;

    ngAfterViewChecked(): void {
        if (this.shouldScroll && this.terminalBody) {
            this.scrollToBottom();
        }
    }

    private scrollToBottom(): void {
        try {
            this.terminalBody.nativeElement.scrollTop = this.terminalBody.nativeElement.scrollHeight;
        } catch (err) { }
    }
}
