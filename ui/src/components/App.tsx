import { createStyles, withStyles, withTheme, Theme, WithStyles } from "@material-ui/core/styles";
import { inject, observer } from "mobx-react";
import * as React from "react";

import { IMenuItem } from "../menu";
import { IAuthResponse } from "../model/auth";
import { Appbar } from "./Appbar";
import { Sidebar } from "./Sidebar";

const styles = (theme: Theme) => createStyles({
    main: {
        margin: "1em",
    },
    root: {
    },
});

export interface IAppProps extends WithStyles<typeof styles> {
    name: string;
    navigate: (path: string) => void;
    menu: IMenuItem[];
    appStore?: IAppStore;
    authStore?: IAuthStore;
}

interface IAppStore {
    sidebarOpen: boolean;
    setSidebarOpen(value: boolean): void;
}

interface IAuthStore {
    isLoggedIn: boolean;
}

@inject("appStore", "authStore")
@observer
class BaseApp extends React.Component<IAppProps, {}> {

    public render() {

        const { name, classes, navigate, menu } = this.props;
        const { isLoggedIn } = this.props.authStore!;
        const { sidebarOpen } = this.props.appStore!;

        const authenticatedControls = (
            <div>
                <Appbar title={name} onMenuClick={this.appbarOnMenuClick} />
                <Sidebar
                    navigate={navigate}
                    open={sidebarOpen}
                    onToggle={this.sidebarOnToggle}
                    menuItems={menu} />
            </div>
        );
        const nonAuthenticatedControls = (<div />);
        const controls = isLoggedIn ? authenticatedControls : nonAuthenticatedControls;

        return (
            <div className={classes.root}>
                {controls}
                <main className={classes.main}>
                    {this.props.children}
                </main>
            </div>
        );

    }

    private appbarOnMenuClick = () => {
        const { sidebarOpen } = this.props.appStore!;
        this.props.appStore!.setSidebarOpen(!sidebarOpen);
    }

    private sidebarOnToggle = (open: boolean) => {
        this.props.appStore!.setSidebarOpen(open);
    }

}

export const App = withTheme()(withStyles(styles)(BaseApp));
