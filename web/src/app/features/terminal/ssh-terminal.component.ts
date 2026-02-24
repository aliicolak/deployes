import { Component, ElementRef, OnDestroy, OnInit, ViewChild, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { environment } from '../../../environments/environment';

@Component({
    selector: 'app-ssh-terminal',
    standalone: true,
    imports: [CommonModule],
    template: `
    <div class="terminal-wrapper">
        <div class="terminal-container" #terminalDiv></div>
        <div class="status-bar" [class.connected]="connected">
            <span class="status-indicator"></span>
            {{ connected ? 'Connected' : 'Disconnected' }}
        </div>
    </div>
  `,
    styles: [`
    :host {
        display: block;
        height: 100%;
        background: #1e1e1e;
        overflow: hidden;
    }
    .terminal-wrapper {
        display: flex;
        flex-direction: column;
        height: 100%;
    }
    .terminal-container {
        flex: 1;
        overflow: hidden;
        padding: 8px;
        background: #000;
    }
    .status-bar {
        height: 28px;
        background: #2d2d2d;
        color: #aaa;
        font-family: monospace;
        font-size: 12px;
        display: flex;
        align-items: center;
        padding: 0 12px;
        border-top: 1px solid #444;
    }
    .status-bar.connected { color: #f0f0f0; }
    .status-indicator {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: #dc2626;
        margin-right: 8px;
    }
    .status-bar.connected .status-indicator {
        background: #22c55e;
        box-shadow: 0 0 8px rgba(34, 197, 94, 0.4);
    }
  `]
})
export class SshTerminalComponent implements OnInit, OnDestroy {
    @ViewChild('terminalDiv', { static: true }) terminalDiv!: ElementRef;

    private route = inject(ActivatedRoute);
    private term!: Terminal;
    private fitAddon!: FitAddon;
    private ws!: WebSocket;
    connected = false;

    ngOnInit() {
        const serverId = this.route.snapshot.paramMap.get('id');
        if (!serverId) return;

        this.initTerminal();
        this.connect(serverId);
    }

    initTerminal() {
        this.term = new Terminal({
            cursorBlink: true,
            theme: {
                background: '#000000',
                foreground: '#f0f0f0',
                cursor: '#ffffff'
            },
            fontFamily: '"Cascadia Code", "JetBrains Mono", monospace',
            fontSize: 14,
            allowProposedApi: true
        });

        this.fitAddon = new FitAddon();
        this.term.loadAddon(this.fitAddon);

        this.term.open(this.terminalDiv.nativeElement);
        // Initial fit might fail if container has no dimensions yet
        setTimeout(() => {
            this.fitAddon.fit();
            this.term.focus();
        }, 100);

        this.term.onData(data => {
            if (this.ws && this.ws.readyState === WebSocket.OPEN) {
                this.ws.send(data);
            }
        });

        window.addEventListener('resize', this.onWindowResize);
    }

    onWindowResize = () => {
        if (this.fitAddon) {
            this.fitAddon.fit();
            if (this.ws && this.ws.readyState === WebSocket.OPEN && this.term) {
                const msg = { type: 'resize', rows: this.term.rows, cols: this.term.cols };
                this.ws.send(JSON.stringify(msg));
            }
        }
    };

    connect(serverId: string) {
        const token = localStorage.getItem('token');
        const url = `${environment.wsUrl}/ws/terminal?serverId=${serverId}&token=${token}`;

        this.ws = new WebSocket(url);
        this.ws.binaryType = 'arraybuffer';

        this.ws.onopen = () => {
            this.connected = true;
            this.term.writeln('\x1b[32mTarget System Connection Established.\r\n\x1b[0m');
            // Send initial resize
            this.onWindowResize();
        };

        this.ws.onmessage = (event) => {
            if (typeof event.data === 'string') {
                this.term.write(event.data.replace(/\n/g, '\r\n'));
            } else {
                const u8 = new Uint8Array(event.data);
                this.term.write(u8);
            }
        };

        this.ws.onclose = () => {
            this.connected = false;
            this.term.writeln('\r\n\x1b[31mConnection closed.\x1b[0m');
        };

        this.ws.onerror = (err) => {
            console.error('WS Error', err);
        }
    }

    ngOnDestroy() {
        window.removeEventListener('resize', this.onWindowResize);
        if (this.ws) this.ws.close();
        if (this.term) this.term.dispose();
    }
}
