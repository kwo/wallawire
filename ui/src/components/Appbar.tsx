import { AppBar, IconButton, Toolbar, Typography } from "@material-ui/core";
import { createStyles, withStyles, withTheme, Theme, WithStyles } from "@material-ui/core/styles";
import MenuIcon from "@material-ui/icons/Menu";
import * as React from "react";

const styles = (theme: Theme) => createStyles({
    grow: {
        flexGrow: 1,
    },
    menuButton: {
        marginLeft: -12,
        marginRight: 20,
    },
});

export interface IAppbarProps extends WithStyles<typeof styles> {
    title: string;
    onMenuClick: () => void;
}

class BaseAppbar extends React.Component<IAppbarProps, {}> {

    public render() {
        const { classes } = this.props;

        return (
            <AppBar position="static">
                <Toolbar>
                    <IconButton
                        className={classes.menuButton} color="inherit" aria-label="Menu"
                        onClick={(/*e*/) => this.props.onMenuClick()} >
                        <MenuIcon />
                    </IconButton>
                    <Typography variant="h6" color="inherit" className={classes.grow}>
                        {this.props.title}
                    </Typography>
                </Toolbar>
            </AppBar>
        );
    }
}

export const Appbar = withTheme()(withStyles(styles)(BaseAppbar));
