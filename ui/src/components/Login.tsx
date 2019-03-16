import { Button, TextField, Typography } from "@material-ui/core";
import { createStyles, withStyles, withTheme, Theme, WithStyles } from "@material-ui/core/styles";
import { get as _get } from "lodash";
import { inject, observer } from "mobx-react";
import * as React from "react";
import { RouteComponentProps } from "react-router";

import { AuthStatus, IAuthResponse } from "../model/auth";
import { Status } from "./Status";

const styles = (theme: Theme) => createStyles({
    button: {
    },
    buttonContainer: {
        marginTop: theme.spacing.unit * 2,
        position: "relative",
        width: "100%",

    },
    container: {
        alignItems: "center",
        display: "flex",
        flexFlow: "row wrap",
        width: "400px",

    },
    field: {
        width: "100%",
    },
    fieldsContainer: {
        width: "100%",
    },
    messageContainer: {
        marginTop: theme.spacing.unit * 2,
    },
    messageContent: {
    },
    root: {
        alignItems: "center",
        display: "flex",
        justifyContent: "center",
    },
    titleContainer: {
        width: "100%",
    },
    status: {
        bottom: "5px",
        color: theme.palette.primary.main,
        fontSize: "0.6rem",
        position: "absolute",
        right: "5px",
    },

});

export interface ILoginProps extends RouteComponentProps, WithStyles<typeof styles> {
    navigate: (path: string, replace?: boolean) => void;
    title: string;
    appStore?: IAppStore;
    authStore?: IAuthStore;
    loginStore?: ILoginStore;
}

interface IAppStore {
    showNotification(message: string, variant: string): void;
}

interface IAuthStore {
    login(username: string, password: string): Promise<IAuthResponse>;
}

interface ILoginStore {
    username: string;
    password: string;
    loading: boolean;
    setUsername(username: string): void;
    setPassword(password: string): void;
    setLoading(loading: boolean): void;
}

@inject("appStore", "authStore", "loginStore")
@observer
class BaseLogin extends React.Component<ILoginProps, {}> {

    public render() {
        const { classes } = this.props;
        const { username, password, loading } = this.props.loginStore!;
        const allowLogin: boolean = (!username || !password);
        const buttonText = loading ? "Waiting..." : "Login";

        return (
            <React.Fragment>
                <Status className={classes.status} />
                <div className={classes.root}>
                    <form className={classes.container} noValidate={false} autoComplete="off">
                        <div className={classes.titleContainer}>
                            <Typography align="center" color="primary" variant="h5">
                                {this.props.title}
                            </Typography>
                        </div>
                        <TextField
                            autoComplete="username"
                            className={classes.field}
                            disabled={loading}
                            id="login-username"
                            label="Username"
                            margin="normal"
                            required={true}
                            value={username}
                            onChange={this.handleChange}
                        />
                        <TextField
                            autoComplete="current-password"
                            className={classes.field}
                            disabled={loading}
                            id="login-password"
                            label="Password"
                            margin="normal"
                            required={true}
                            type="password"
                            value={password}
                            onChange={this.handleChange}
                        />
                        <div className={classes.buttonContainer}>
                            <Button
                                className={classes.button}
                                color="primary"
                                disabled={allowLogin || loading}
                                fullWidth={true}
                                size="large"
                                type="submit"
                                variant="outlined"
                                onClick={this.doLogin}>
                                {buttonText}
                            </Button>
                        </div>
                    </form>
                </div>
            </React.Fragment>
        );
    }

    private handleChange = (e: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement | HTMLSelectElement>) => {
        e.preventDefault();
        const { setUsername, setPassword } = this.props.loginStore!;

        const { id, value } = e.currentTarget;
        if (id === "login-username") {
            setUsername(value);
        } else if (id === "login-password") {
            setPassword(value);
        } else {
            console.warn("handleChange: no handler for id: ", id); // tslint:disable-line:no-console
        }
    }

    private doLogin = (e: React.MouseEvent<Element>): Promise<boolean> => {
        e.preventDefault();

        const { showNotification } = this.props.appStore!;
        const { username, password, setLoading } = this.props.loginStore!;
        setLoading(true);

        return this.props.authStore!.login(username, password).then((rsp: IAuthResponse) => {
            setLoading(false);
            if (rsp.status === AuthStatus.Success) {
                // set in Root.AllowableRoute
                this.props.navigate(_get(this.props, "location.state.referrer.pathname", "/"), true);
            } else {
                showNotification(rsp.message, "error");
            }
            return true;
        });

    }

}

export const Login = withTheme()(withStyles(styles)(BaseLogin));
