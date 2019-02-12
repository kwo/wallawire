import { createStyles, withStyles, withTheme, Theme, WithStyles } from "@material-ui/core/styles";
import { inject, observer } from "mobx-react";
import * as React from "react";
import { withRouter, Redirect, Route, RouteComponentProps } from "react-router";
import { Switch } from "react-router-dom";

import { menu } from "../menu";
import { About } from "./About";
import { App } from "./App";
import { Hello } from "./Hello";
import { Login } from "./Login";
import { Logout } from "./Logout";
import { Notification } from "./Notification";
import { Profile } from "./Profile";

const appName = "Wallawire";

const styles = (theme: Theme) => createStyles({
    main: {
        margin: "1em",
    },
});

export interface IRootProps extends RouteComponentProps, WithStyles<typeof styles> {
    appStore?: IAppStore;
    authStore?: IAuthStore;
    statusStore?: IStatusStore;
}

interface IAppStore {
    messageOpen: boolean;
    messageText: string;
    messageVariant: string;
    closeNotification(): void;
}

interface IAuthStore {
    isLoggedIn: boolean;
    logout(): Promise<any>;
}

interface IStatusStore {
    start: () => void;
}

@inject("appStore", "authStore", "statusStore")
@observer
class BaseRoot extends React.Component<IRootProps, {}> {

    constructor(props: Readonly<IRootProps>) {
        super(props);
    }

    public componentDidMount() {
        this.props.statusStore!.start();
    }

    public render() {
        const { isLoggedIn } = this.props.authStore!;
        const { messageOpen, messageText, messageVariant } = this.props.appStore!;

        const InternalRoute:
            React.SFC</*RouteProps*/ any> = ({ component: Component, render, path, ...routeProps }) =>
                <Route {...routeProps} path={path} render={(props: RouteComponentProps) => {
                    if (path !== "/login" ? isLoggedIn : !isLoggedIn) {
                        return render ? render(props) : <Component {...props} />;
                    }
                    return <Redirect to={{
                        pathname: (path !== "/login" ? "/login" : "/"),
                        state: {
                            referrer: props.location,
                        },
                    }} />;
                }} />;

        return (
            <App name={appName} navigate={this.navigate} menu={menu}>
                <Switch>
                    <InternalRoute path="/" exact={true} strict={true}
                        render={(props: RouteComponentProps) => <Hello caption="Welcome" />} />
                    <InternalRoute path="/about" exact={false} strict={true} component={About} />
                    <InternalRoute path="/login" exact={true} strict={true}
                        render={(props: RouteComponentProps) =>
                            <Login
                                navigate={this.navigate}
                                title={`${appName} Login`}
                                {...props}
                            />}
                    />
                    <InternalRoute path="/logout" exact={true} strict={true}
                        render={(props: RouteComponentProps) => <Logout logout={this.logout} {...props} />} />
                    />
                    <InternalRoute path="/profile" exact={true} strict={true}
                        render={(props: RouteComponentProps) =>
                            <Profile {...props} />}
                    />
                    <InternalRoute
                        render={(props: RouteComponentProps) => <Hello caption="Not Found!" />} />
                </Switch>
                <Notification
                    message={messageText}
                    open={messageOpen}
                    onClose={this.onMessageClose}
                    variant={messageVariant} />
            </App>
        );
    }

    private onMessageClose = (event: React.SyntheticEvent<any, Event>, reason?: string) => {
        const { closeNotification } = this.props.appStore!;
        if (reason === "clickaway") {
            return;
        }
        closeNotification();
    }

    private logout = (): Promise<any> => {
        return this.props.authStore!.logout().then(() => {
            window.location.pathname = "/";
        });
    }

    private navigate = (path: string, replace?: boolean) => {
        if (replace) {
            this.props.history.replace(path);
        } else {
            this.props.history.push(path);
        }
    }

}

export const Root = withTheme()(withStyles(styles)(withRouter(BaseRoot)));
