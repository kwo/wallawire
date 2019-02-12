import { action, observable } from "mobx";

export class AppStore {
    @observable public sidebarOpen = false;
    @observable public messageOpen = false;
    @observable public messageText = "";
    @observable public messageVariant = "";

    @action public setSidebarOpen = (value: boolean): void => {
        this.sidebarOpen = value;
    }

    @action public showNotification = (message: string, variant: string) => {
        this.messageOpen = true;
        this.messageText = message;
        this.messageVariant = variant;
    }

    @action public closeNotification = (): void => {
        this.messageOpen = false;
    }

}
