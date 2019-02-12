import * as jwt from "jsonwebtoken";
import { AuthStore } from "./auth";

describe("AuthStore", () => {

    test("success", () => {

        const maxAge = (3600 * 24 * 1000); // one day in milliseconds
        const expires = Math.floor((Date.now() + maxAge) / 1000);

        const data = {
            username: "barkerb",
            name: "Barker, Bob",
            roles: "user, friends",
            exp: expires,
        };
        const value = jwt.sign(data, "secret");

        const authStore = new AuthStore();
        authStore.setToken(value);

        expect(authStore.isLoggedIn).toEqual(true);
        expect(authStore.isValid()).toEqual(true);
        expect(authStore.isExpired()).toEqual(false);
        expect(authStore.name).toEqual("Barker, Bob");
        expect(authStore.username).toEqual("barkerb");
        expect(authStore.roles).toEqual(["user", "friends"]);
        expect(authStore.sessionStarted.getTime()).toEqual((expires * 1000) - maxAge);
        expect(authStore.sessionExpires.getTime()).toEqual(expires * 1000);
        expect(authStore.sessionExpires.getTime() - authStore.sessionStarted.getTime()).toEqual(maxAge);

    });

    test("bogus - no name", () => {

        const maxAge = (3600 * 24 * 1000); // one day in milliseconds
        const expires = Math.floor((Date.now() + maxAge) / 1000);

        const data = {
            username: "barkerb",
            roles: "user, friends",
            exp: expires,
        };
        const value = jwt.sign(data, "secret");

        const authStore = new AuthStore();
        authStore.setToken(value);

        expect(authStore.isLoggedIn).toEqual(false);
        expect(authStore.isValid()).toEqual(false);
        expect(authStore.isExpired()).toEqual(false);

    });

    test("bogus - no expires", () => {

        const data = {
            username: "barkerb",
            name: "Barker, Bob",
            roles: "user, friends",
        };
        const value = jwt.sign(data, "secret");

        const authStore = new AuthStore();
        authStore.setToken(value);

        expect(authStore.isLoggedIn).toEqual(false);
        expect(authStore.isValid()).toEqual(false);
        expect(authStore.isExpired()).toEqual(true);
        expect(authStore.sessionExpires.getTime()).toEqual(0);

    });

    test("empty string", () => {

        const authStore = new AuthStore();
        authStore.setToken("");

        expect(authStore.isLoggedIn).toEqual(false);
        expect(authStore.isValid()).toEqual(false);
        expect(authStore.isExpired()).toEqual(true);
        expect(authStore.name).toEqual("");
        expect(authStore.username).toEqual("");
        expect(authStore.roles).toEqual([]);
        expect(authStore.sessionStarted.getTime()).toEqual(0);
        expect(authStore.sessionExpires.getTime()).toEqual(0);

    });

    test("null", () => {

        const authStore = new AuthStore();
        authStore.setToken(null);

        expect(authStore.isLoggedIn).toEqual(false);
        expect(authStore.isValid()).toEqual(false);
        expect(authStore.isExpired()).toEqual(true);
        expect(authStore.name).toEqual("");
        expect(authStore.username).toEqual("");
        expect(authStore.roles).toEqual([]);
        expect(authStore.sessionStarted.getTime()).toEqual(0);
        expect(authStore.sessionExpires.getTime()).toEqual(0);

    });

    test("undefined", () => {

        const authStore = new AuthStore();
        authStore.setToken(undefined);

        expect(authStore.isLoggedIn).toEqual(false);
        expect(authStore.isValid()).toEqual(false);
        expect(authStore.isExpired()).toEqual(true);
        expect(authStore.name).toEqual("");
        expect(authStore.username).toEqual("");
        expect(authStore.roles).toEqual([]);
        expect(authStore.sessionStarted.getTime()).toEqual(0);
        expect(authStore.sessionExpires.getTime()).toEqual(0);

    });

    test("empty roles", () => {

        const maxAge = (3600 * 24 * 1000); // one day in milliseconds
        const expires = Math.floor((Date.now() + maxAge) / 1000);

        const data = {
            username: "barkerb",
            name: "Barker, Bob",
            roles: "",
            exp: expires,
        };
        const value = jwt.sign(data, "secret");

        const authStore = new AuthStore();
        authStore.setToken(value);

        expect(authStore.isLoggedIn).toEqual(true);
        expect(authStore.roles).toEqual([]);
        expect(!!authStore.roles).toEqual(true);

    });

});
