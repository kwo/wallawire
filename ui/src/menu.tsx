import * as React from "react";

import AccountCircle from "@material-ui/icons/AccountCircle";
import DashboardIcon from "@material-ui/icons/Dashboard";
import ExitIcon from "@material-ui/icons/ExitToApp";
import InfoIcon from "@material-ui/icons/Info";

export interface IMenuItem {
    id: string;
    name: string;
    path?: string;
    icon: React.ReactElement<any>;
}

export const menu: IMenuItem[] = [
    {
        id: "user",
        name: "Me",
        path: "/profile",
        icon: <AccountCircle/>,
    },
    {
        id: "dashboard",
        name: "Home",
        path: "/",
        icon: <DashboardIcon/>,
    },
    {
        id: "logout",
        name: "Logout",
        path: "/logout",
        icon: <ExitIcon/>,
    },
];
