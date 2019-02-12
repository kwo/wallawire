import { action, observable } from "mobx";

export class LoginStore {
    @observable public username = "";
    @observable public password = "";
    @observable public loading = false;

    @action public setUsername = (username: string): void => {
        this.username = username;
    }

    @action public setPassword = (password: string): void => {
        this.password = password;
    }

    @action public setLoading = (loading: boolean): void => {
        this.loading = loading;
    }

}
