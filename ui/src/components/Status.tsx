import { Popover, Typography } from "@material-ui/core";
import { createStyles, withStyles, Theme, WithStyles } from "@material-ui/core/styles";
import { inject, observer } from "mobx-react";
import * as React from "react";

import { IServerStatus } from "../model/status";

const styles = (theme: Theme) => createStyles({
    paper: {
        padding: theme.spacing.unit,
    },
    popover: {
        pointerEvents: "none",
    },
});

export interface IStatusProps extends WithStyles<typeof styles> {
    className: string;
    statusStore?: IStatusStore;
}

interface IStatusStore {
    client: string;
    server: IServerStatus;
}

interface IState {
    anchorEl: HTMLElement | null;
}

@inject("statusStore")
@observer
class BaseStatus extends React.Component<IStatusProps, IState> {

    public state: IState = {
        anchorEl: null,
    };

    public render() {

        const { classes, className, statusStore } = this.props;
        const { anchorEl } = this.state;
        const open = Boolean(anchorEl);

        return (
            <React.Fragment>
                <Typography
                    className={className}
                    aria-owns={open ? "version-popover" : undefined}
                    aria-haspopup="true"
                    onMouseEnter={this.handlePopoverOpen}
                    onMouseLeave={this.handlePopoverClose}>
                    {statusStore!.client} / {statusStore!.server.version}
                </Typography>
                <Popover
                    id="version-popover"
                    className={classes.popover}
                    classes={{
                        paper: classes.paper,
                    }}
                    open={open}
                    anchorEl={anchorEl}
                    anchorOrigin={{
                        horizontal: "left",
                        vertical: "bottom",
                    }}
                    transformOrigin={{
                        horizontal: "left",
                        vertical: "top",
                    }}
                    onClose={this.handlePopoverClose}
                    disableRestoreFocus={true}>

                    <div>
                        <Typography>Uptime: {statusStore!.server.uptime}</Typography>
                        <Typography>Server: {statusStore!.server.runtime}</Typography>
                    </div>

                </Popover>
            </React.Fragment>
        );
    }

    private handlePopoverOpen = (event: React.MouseEvent<HTMLElement>) => {
        this.setState({ anchorEl: event.currentTarget });
    }

    private handlePopoverClose = () => {
        this.setState({ anchorEl: null });
    }

}

export const Status = withStyles(styles)(BaseStatus);
