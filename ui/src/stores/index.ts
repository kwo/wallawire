import { configure } from "mobx";
import { AppStore } from "./app";
import { AuthStore } from "./auth";
import { LoginStore } from "./login";
import { MessageStore } from "./messages";
import { ProfileStore } from "./profile";
import { StatusStore } from "./status";

configure({ enforceActions: "observed" });

const authStore = new AuthStore();
const statusStore = new StatusStore();

export const stores = {
    appStore: new AppStore(),
    authStore,
    loginStore: new LoginStore(),
    messageStore: new MessageStore(authStore, statusStore),
    profileStore: new ProfileStore(authStore),
    statusStore,
};
