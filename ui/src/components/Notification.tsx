import { IconButton, Snackbar } from "@material-ui/core";
import { green } from "@material-ui/core/colors";
import { createStyles, withStyles, withTheme, Theme, WithStyles } from "@material-ui/core/styles";
import { CheckCircle as CheckCircleIcon, Close as CloseIcon, Error as ErrorIcon } from "@material-ui/icons";
import * as React from "react";

// tslint:disable:object-literal-sort-keys
const styles = (theme: Theme) => createStyles({
    success: {
        color: "#FFFFFF",
        backgroundColor: green[600],
    },
    error: {
        color: "#FFFFFF",
        backgroundColor: "#FF0000",
    },
    messageIconAction: {
        padding: theme.spacing.unit / 2,
    },
    messageIconMessage: {
        marginRight: theme.spacing.unit * 2,
    },
    messageText: {
        display: "flex",
        alignItems: "center",
    },
});
// tslint:enable:object-literal-sort-keys

export interface INotificationProps extends WithStyles<typeof styles> {
    open: boolean;
    message: string;
    variant: string;
    onClose: (event: React.SyntheticEvent<any, Event>, reason?: string) => void;
}

class BaseNotification extends React.Component<INotificationProps, {}> {

    public render() {

        const { classes, message, onClose, open, variant } = this.props;
        const Icon = (variant === "success") ? CheckCircleIcon : ErrorIcon;
        const className = (variant === "success") ? classes.success : classes.error;

        return (<Snackbar
            action={[
                <IconButton
                    key="close"
                    aria-label="Close"
                    color="inherit"
                    className={classes.messageIconAction}
                    onClick={onClose}>
                    <CloseIcon />
                </IconButton>,
            ]}
            anchorOrigin={{
                horizontal: "left",
                vertical: "bottom",
            }}
            autoHideDuration={6000}
            ContentProps={{
                classes: {
                    root: className,
                },
            }}
            message={<span className={classes.messageText} id="notification-message">
                <Icon className={classes.messageIconMessage} />
                {message}
            </span>}
            onClose={onClose}
            open={open}>
        </Snackbar>);
    }

}

export const Notification = withTheme()(withStyles(styles)(BaseNotification));
