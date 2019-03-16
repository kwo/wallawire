import { Divider, List, ListItem, ListItemIcon, ListItemText, SwipeableDrawer } from "@material-ui/core";
import { createStyles, withStyles, withTheme, Theme, WithStyles } from "@material-ui/core/styles";
import { inject, observer } from "mobx-react";
import * as React from "react";

import { IMenuItem } from "../menu";
import { Status } from "./Status";

const styles = (theme: Theme) => createStyles({
    avatarName: {
        flexGrow: 1,
        marginLeft: "1em",
    },
    fullList: {
        width: "auto",
    },
    list: {
        width: 250,
    },
    status: {
        bottom: "5px",
        color: theme.palette.text.primary,
        fontSize: "0.6rem",
        left: "5px",
        position: "absolute",
    },
});

export interface ISidebarProps extends WithStyles<typeof styles> {
    navigate: (path: string) => void;
    open: boolean;
    onToggle: (open: boolean) => void;
    menuItems: IMenuItem[];
    authStore?: IAuthStore;
}

interface IAuthStore {
    name: string;
    sessionExpiresIn: string;
}

@inject("authStore")
@observer
class BaseSidebar extends React.Component<ISidebarProps, {}> {

    public render() {
        const { classes, menuItems } = this.props;
        const { sessionExpiresIn } = this.props.authStore!;
        const header = this.menuItemToListItem(this.personalizeHeaderItem(menuItems[0]));
        const footer = this.menuItemToListItemWithSecondary(
            menuItems[menuItems.length - 1],
            `timeout: ${sessionExpiresIn}`,
        );
        const newMenuItems = menuItems.filter((r: IMenuItem, i: number) => i !== 0 && i !== menuItems.length - 1);
        const listItems = this.menuItemsToListItems(newMenuItems);

        return (
            <SwipeableDrawer
                open={this.props.open}
                onClose={this.toggleDrawer(false)}
                onOpen={this.toggleDrawer(true)}
            >
                <div
                    tabIndex={0}
                    role="button"
                    onClick={this.toggleDrawer(false)}
                    onKeyDown={this.toggleDrawer(false)}
                    className={classes.list}
                >
                    {header}
                    <Divider />
                    <List>
                        {listItems}
                        <Divider />
                        {footer}
                    </List>
                </div>
                <Status className={classes.status} />
            </SwipeableDrawer>
        );
    }

    private personalizeHeaderItem = (mi: IMenuItem): IMenuItem => {
        const { name } = this.props.authStore!;
        return { ...mi, name };
    }

    private menuItemToListItemWithSecondary = (mi: IMenuItem, x: any): any => {
        return (
            <ListItem id={`sidebar-${mi.id}`} key={mi.id} button={true} onClick={this.handleClick}>
                <ListItemIcon>{mi.icon}</ListItemIcon>
                <ListItemText
                    primary={mi.name}
                    secondary={x}
                    secondaryTypographyProps={{
                        variant: "caption",
                    }} />
            </ListItem>
        );
    }

    private menuItemToListItem = (mi: IMenuItem): any => {
        return mi.path
            ?
            <ListItem id={`sidebar-${mi.id}`} key={mi.id} button={true} onClick={this.handleClick}>
                <ListItemIcon>{mi.icon}</ListItemIcon>
                <ListItemText primary={mi.name} />
            </ListItem>
            :
            <ListItem id={`sidebar-${mi.id}`} key={mi.id}>
                <ListItemIcon>{mi.icon}</ListItemIcon>
                <ListItemText primary={mi.name} />
            </ListItem>
            ;
    }

    private menuItemsToListItems = (mi: IMenuItem[]): any => {
        return mi.map(this.menuItemToListItem);
    }

    private handleClick = (e: React.MouseEvent<HTMLElement>) => {
        const { id } = e.currentTarget;
        const bareId = id.split("-")[1];
        const menuItem = this.props.menuItems.find((mi: IMenuItem) => mi.id === bareId);
        const href = menuItem ? menuItem.path : undefined;
        if (href) {
            this.props.navigate(href);
        }
    }

    private toggleDrawer = (open: boolean) => () => {
        this.props.onToggle(open);
    }

}

export const Sidebar = withTheme()(withStyles(styles)(BaseSidebar));
