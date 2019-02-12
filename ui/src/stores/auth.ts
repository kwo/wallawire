import { distanceInWordsToNow } from "date-fns";
import * as Cookies from "js-cookie";
import { action, computed, observable } from "mobx";
import { AuthStatus, IAuthResponse } from "../model/auth";
import { to } from "./util";

export class AuthStore {
    @observable public name = "";
    @observable public username = "";
    @observable public roles: string[] = [];
    @observable public sessionStarted = new Date(0);
    @observable public sessionExpires = new Date(0);

    constructor() {
        this.reloadToken();
    }

    @action public setToken = (token: any) => {
        if (token) {
            const elements = token.split(".");
            if (elements.length === 3) {
                const base64UrlEncoded = elements[1];
                const base64 = base64UrlEncoded.replace("-", "+").replace("_", "/");
                const data = JSON.parse(window.atob(base64));
                if (data) {
                    this.name = data.name;
                    this.username = data.username;
                    this.roles = data.roles && data.roles.split(/\s*,\s*/) || [];
                    this.sessionStarted = data.iat && new Date(data.iat * 1000) || new Date(0);
                    this.sessionExpires = data.exp && new Date(data.exp * 1000) || new Date(0);
                }
            }
        } else {
            this.username = "";
            this.name = "";
            this.roles = [];
            this.sessionStarted = new Date(0);
            this.sessionExpires = new Date(0);
        }
    }

    @action public reloadToken = () => {
        this.setToken(Cookies.get("jwt"));
    }

    @action public login = async (username: string, password: string): Promise<IAuthResponse> => {

        const loginRequest = { username, password };

        const [rsp, errRsp] = await to(fetch("/api/login", {
            body: JSON.stringify(loginRequest),
            cache: "no-cache",
            credentials: "same-origin",
            headers: {
                "Content-Type": "application/json",
            },
            method: "POST",
            mode: "same-origin",
        }));

        this.reloadToken();

        if (errRsp) {
            // TODO: log error back to server
            console.warn("authenticate error", errRsp); // tslint:disable-line:no-console
            return { status: AuthStatus.Error, message: "An unexpected error occurred." };
        }

        const statusText: string = await rsp.text();

        // console.debug("auth.authentication", rsp, statusText); // tslint:disable-line:no-console

        switch (rsp.status) {
            case 200:
                return { status: AuthStatus.Success, message: "" };
            case 400:
                return { status: AuthStatus.BadRequest, message: statusText };
            case 401:
                return { status: AuthStatus.Unauthorized, message: statusText };
            case 403:
                return { status: AuthStatus.Forbidden, message: statusText };
            case 500:
                return { status: AuthStatus.Error, message: statusText };
        }

        // TODO: log error back to server
        console.warn("authenticate unknown status code", rsp.status); // tslint:disable-line:no-console
        return { status: AuthStatus.Unknown, message: statusText };

    }

    @action public logout = async (): Promise<any> => {
        const [rsp, errRsp] = await to(fetch("/api/logout", {
            cache: "no-cache",
            credentials: "same-origin",
            method: "POST",
            mode: "same-origin",
        }));
        this.reloadToken();
    }

    @computed get isLoggedIn(): boolean {
        const result = this.isValid() && !this.isExpired();
        // console.debug("isLoggedIn", result); // tslint:disable-line:no-console
        return result;
    }

    @computed get sessionExpiresIn(): string {
        if (this.isValid()) {
            return distanceInWordsToNow(this.sessionExpires);
        }
        return "now"; // if invalid, return 0
    }

    public isExpired = (): boolean => {
        const result = (this.sessionExpires.getTime() < Date.now());
        // console.debug("isExpired", result); // tslint:disable-line:no-console
        return result;
    }

    public isValid = (): boolean => {
        const result = (!!this.name && !!this.username && !!this.roles
            && this.sessionStarted.getTime() !== 0 && this.sessionExpires.getTime() !== 0);
        // console.debug("isValid", result); // tslint:disable-line:no-console
        return result;
    }

}
