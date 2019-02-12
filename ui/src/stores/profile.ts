import { action, autorun, observable } from "mobx";
import { ChangeProfileStatus, IChangeProfileResponse } from "../model/profile";
import { to } from "./util";

interface IAuthStore {
    username: string;
    name: string;
    reloadToken(): void;
}

export class ProfileStore {

    // common
    @observable public loading = false;
    @observable public panelExpanded = "";

    // username
    @observable public username = ""; // must be in sync with authStore
    @observable public password = "";

    // password
    @observable public passwordOld = "";
    @observable public passwordNew = "";
    @observable public passwordConfirm = "";
    @observable public warning = "";
    @observable public warning1 = false;
    @observable public warning2 = false;

    // settings
    @observable public displayname = "";

    private authStore: IAuthStore | null = null;

    constructor(authStore: IAuthStore) {
        this.authStore = authStore;
        autorun(() => {
            this.setDisplayname(this.authStore!.name);
            this.setUsername(this.authStore!.username);
        });
    }

    // common
    @action public setLoading = (value: boolean): void => {
        this.loading = value;
    }
    @action public setPanelExpanded = (panelName: string): void => {
        this.panelExpanded = panelName;
    }

    // username
    @action public setUsername = (value: string): void => {
        this.username = value;
    }
    @action public setPassword = (value: string): void => {
        this.password = value;
    }
    @action public resetUsername = (): void => {
        this.username = this.authStore!.username;
        this.password = "";
    }

    // password
    @action public setPasswordOld = (value: string): void => {
        this.passwordOld = value;
    }
    @action public setPasswordNew = (value: string): void => {
        this.passwordNew = value;
    }
    @action public setPasswordConfirm = (value: string): void => {
        this.passwordConfirm = value;
    }
    @action public setWarning = (value: string): void => {
        this.warning = value;
    }
    @action public setWarning1 = (value: boolean): void => {
        this.warning1 = value;
    }
    @action public setWarning2 = (value: boolean): void => {
        this.warning2 = value;
    }
    @action public resetPassword = (): void => {
        this.passwordOld = "";
        this.passwordNew = "";
        this.passwordConfirm = "";
        this.warning = "";
        this.warning1 = false;
        this.warning2 = false;
    }

    // settings
    @action public setDisplayname = (value: string): void => {
        this.displayname = value;
    }
    @action public resetSettings = (): void => {
        this.displayname = this.authStore!.name;
    }

    @action public changeusername =
        async (newusername: string, password: string): Promise<IChangeProfileResponse> => {

            this.setLoading(true);

            const changeUsernameRequest = { newusername, password };

            const [rsp, errRsp] = await to(fetch("/api/changeusername", {
                body: JSON.stringify(changeUsernameRequest),
                cache: "no-cache",
                credentials: "same-origin",
                headers: {
                    "Content-Type": "application/json",
                },
                method: "POST",
                mode: "same-origin",
            }));

            if (errRsp) {
                // TODO: log error back to server
                console.warn("authenticate error", errRsp); // tslint:disable-line:no-console
                return { status: ChangeProfileStatus.Error, message: "An unexpected error occurred." };
            }

            let data: any = {};
            if (rsp.headers.get("Content-Type") === "application/json") {
                data = await rsp.json();
            }

            this.setLoading(false);

            switch (rsp.status) {
                case 200:
                    this.authStore!.reloadToken();
                    this.resetUsername();
                    return { status: ChangeProfileStatus.Success, message: "" };
                case 400:
                    return { status: ChangeProfileStatus.BadRequest, message: data.message };
                case 401:
                    return { status: ChangeProfileStatus.Unauthorized, message: "Error: not authenticated." };
                case 403:
                    return { status: ChangeProfileStatus.Forbidden, message: "Error: forbidden." };
            }

            // TODO: log error back to server
            console.warn("authenticate unknown status code", rsp.status); // tslint:disable-line:no-console
            return { status: ChangeProfileStatus.Unknown, message: `An unknown status code occurred: ${rsp.status}` };

        }

    @action public changepassword =
        async (newpassword: string, oldpassword: string): Promise<IChangeProfileResponse> => {

            this.setLoading(true);
            const changePasswordRequest = { newpassword, oldpassword };

            const [rsp, errRsp] = await to(fetch("/api/changepassword", {
                body: JSON.stringify(changePasswordRequest),
                cache: "no-cache",
                credentials: "same-origin",
                headers: {
                    "Content-Type": "application/json",
                },
                method: "POST",
                mode: "same-origin",
            }));

            if (errRsp) {
                // TODO: log error back to server
                console.warn("authenticate error", errRsp); // tslint:disable-line:no-console
                return { status: ChangeProfileStatus.Error, message: "An unexpected error occurred." };
            }

            let data: any = {};
            if (rsp.headers.get("Content-Type") === "application/json") {
                data = await rsp.json();
            }

            this.setLoading(false);

            switch (rsp.status) {
                case 200:
                    this.authStore!.reloadToken();
                    this.resetPassword();
                    return { status: ChangeProfileStatus.Success, message: "" };
                case 400:
                    return { status: ChangeProfileStatus.BadRequest, message: data.message };
                case 401:
                    return { status: ChangeProfileStatus.Unauthorized, message: "Error: not authenticated." };
                case 403:
                    return { status: ChangeProfileStatus.Forbidden, message: "Error: forbidden." };
            }

            // TODO: log error back to server
            console.warn("authenticate unknown status code", rsp.status); // tslint:disable-line:no-console
            return { status: ChangeProfileStatus.Unknown, message: `An unknown status code occurred: ${rsp.status}` };

        }

    @action public changesettings = async (displayname: string): Promise<IChangeProfileResponse> => {

        this.setLoading(true);
        const changeProfileRequest = { displayname };

        const [rsp, errRsp] = await to(fetch("/api/changeprofile", {
            body: JSON.stringify(changeProfileRequest),
            cache: "no-cache",
            credentials: "same-origin",
            headers: {
                "Content-Type": "application/json",
            },
            method: "POST",
            mode: "same-origin",
        }));

        if (errRsp) {
            // TODO: log error back to server
            console.warn("authenticate error", errRsp); // tslint:disable-line:no-console
            return { status: ChangeProfileStatus.Error, message: "An unexpected error occurred." };
        }

        let data: any = {};
        if (rsp.headers.get("Content-Type") === "application/json") {
            data = await rsp.json();
        }

        this.setLoading(false);

        switch (rsp.status) {
            case 200:
                this.authStore!.reloadToken();
                this.resetSettings();
                return { status: ChangeProfileStatus.Success, message: "" };
            case 400:
                return { status: ChangeProfileStatus.BadRequest, message: data.message };
            case 401:
                return { status: ChangeProfileStatus.Unauthorized, message: "Error: not authenticated." };
            case 403:
                return { status: ChangeProfileStatus.Forbidden, message: "Error: forbidden." };
        }

        // TODO: log error back to server
        console.warn("authenticate unknown status code", rsp.status); // tslint:disable-line:no-console
        return { status: ChangeProfileStatus.Unknown, message: `An unknown status code occurred: ${rsp.status}` };

    }

}
