import { action, observable } from "mobx";
import { IServerStatus } from "../model/status";
import { to } from "./util";

const blankServerStatus = {
    service: "",
    version: "",
    runtime: "",
    buildTime: new Date(0),
    start: new Date(0),
    time: new Date(0),
    uptime: "",
};

export class StatusStore {
    @observable public client = "";
    @observable public server = blankServerStatus;

    @action public refresh = () => {
        this.updateClientVersion();
        this.updateServerStatus();
    }

    @action public setServerStatus(serverStatus: IServerStatus) {
        this.server = serverStatus;
    }

    @action public updateClientVersion = () => {

        const versionTag =
            Array.from(document.getElementsByTagName("meta"))
                .find((tag: HTMLMetaElement) => tag.getAttribute("name") === "version");

        this.client = versionTag ? versionTag.getAttribute("content") || "UNKNOWN" : "MISSING";

    }

    @action public updateServerStatus = async () => {

        const [rsp, errRsp] = await to(fetch("/api/status", {
            cache: "no-cache",
            credentials: "same-origin",
            method: "GET",
            mode: "same-origin",
        }));

        if (errRsp) {
            // TODO: log error back to server
            console.warn("get server status error", errRsp); // tslint:disable-line:no-console
            return blankServerStatus;
        }

        const serverStatus = await rsp.json();
        this.setServerStatus(serverStatus);

    }

}
