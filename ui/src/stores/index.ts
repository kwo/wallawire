import { configure } from "mobx";
import { AppStore } from "./app";
import { AuthStore } from "./auth";
import { LoginStore } from "./login";
import { ProfileStore } from "./profile";
import { StatusStore } from "./status";

configure({ enforceActions: "observed" });

const authStore = new AuthStore();

export const stores = {
    appStore: new AppStore(),
    authStore,
    loginStore: new LoginStore(),
    profileStore: new ProfileStore(authStore),
    statusStore: new StatusStore(),
};
