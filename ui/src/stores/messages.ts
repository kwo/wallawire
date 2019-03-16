import { action, autorun, observable } from "mobx";
import { IServerStatus } from "../model/status";

interface IAuthStore {
    isLoggedIn: boolean;
}

interface IStatusStore {
    setServerStatus(status: IServerStatus): void;
}

const blankServerStatus = {
    service: "",
    version: "",
    runtime: "",
    buildTime: new Date(0),
    start: new Date(0),
    time: new Date(0),
    uptime: "",
};

export class MessageStore {

    @observable public connected: boolean = false;
    @observable public status: IServerStatus = blankServerStatus;

    private authStore: IAuthStore | null = null;
    private statusStore: IStatusStore | null = null;
    private client: EventSource | null = null;

    private checkInterval: any = null;
    private backoff: number = 1000;
    private lastCheck: number = 0;

    constructor(authStore: IAuthStore, statusStore: IStatusStore) {
        this.authStore = authStore;
        this.statusStore = statusStore;
        autorun(() => {
            if (this.authStore!.isLoggedIn) {
                this.start();
            } else {
                this.stop();
            }
        });
    }

    public start = () => {
        console.log("notifications starting"); // tslint:disable-line:no-console
        this.connect();
        this.checkInterval = setInterval(this.check, 1000);
    }

    public stop = () => {
        console.debug("notifications stopping"); // tslint:disable-line:no-console
        clearInterval(this.checkInterval);
        this.checkInterval = null;
        this.disconnect();
    }

    @action private setConnected = (connected: boolean) => {
        this.connected = connected;
        if (connected) {
            this.backoff = 1000;
        }
    }

    private connect = () => {
        this.client = new EventSource("/api/inbox");
        this.client.onopen = (msg: MessageEvent) => {
            this.setConnected(true);
            console.debug("notifications opened"); // tslint:disable-line:no-console
        };
        this.client.onerror = (msg: MessageEvent) => {
            this.setConnected(false);
            console.debug("notifications error"); // tslint:disable-line:no-console
        };
        this.client.onmessage = (msg: MessageEvent) => {
            const data = JSON.parse(msg.data);
            console.debug("notifications untyped", data); // tslint:disable-line:no-console
        };
        this.client.addEventListener("heartbeat", (evt: any) => {
            const status = JSON.parse(evt.data);
            this.statusStore!.setServerStatus(status);
            console.debug("notifications heartbeat"); // tslint:disable-line:no-console
        });
        console.debug("notifications connecting..."); // tslint:disable-line:no-console
    }

    private disconnect = () => {
        if (this.client) {
            this.client.close();
            this.client = null;
        }
        this.setConnected(false);
        console.debug("notifications disconnect"); // tslint:disable-line:no-console
    }

    private check = () => {

        if (this.connected) {
            return;
        }
        if (Date.now() - this.lastCheck < this.backoff) {
            return;
        }

        console.log("notifications checking", this.backoff); // tslint:disable-line:no-console

        this.lastCheck = Date.now();
        this.backoff = this.backoff * 2;
        if (this.backoff > 60000) {
            this.backoff = 60000; // max one minute
        }

        this.disconnect();
        this.connect();

    }

}
