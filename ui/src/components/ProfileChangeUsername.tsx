import {
    Button, Divider, ExpansionPanel, ExpansionPanelActions, ExpansionPanelDetails, ExpansionPanelSummary,
    TextField, Typography,
} from "@material-ui/core";
import { createStyles, withStyles, withTheme, Theme, WithStyles } from "@material-ui/core/styles";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import { inject, observer } from "mobx-react";
import * as React from "react";

import { ChangeProfileStatus, IChangeProfileResponse } from "../model/profile";

// tslint:disable:object-literal-sort-keys
const styles = (theme: Theme) => createStyles({
    heading: {
        fontSize: theme.typography.pxToRem(16),
        fontWeight: theme.typography.fontWeightRegular,
    },
    field: {
        marginBottom: theme.spacing.unit * .5,
        marginLeft: theme.spacing.unit,
        marginRight: theme.spacing.unit,
        marginTop: theme.spacing.unit * .5,
        width: "100%",
    },
});
// tslint:enable:object-literal-sort-keys

export interface IProfileChangeUsernameProps extends WithStyles<typeof styles> {
    appStore?: IAppStore;
    profileStore?: IProfileStore;
}

interface IAppStore {
    messageOpen: boolean;
    messageText: string;
    messageVariant: string;
    sidebarOpen: boolean;
    showNotification(message: string, variant: string): void;
    closeNotification(): void;
    setSidebarOpen(value: boolean): void;
}

interface IProfileStore {
    loading: boolean;
    panelExpanded: string;
    setPanelExpanded: (panelName: string) => void;
    changeusername: (password: string, newusername: string) => Promise<IChangeProfileResponse>;

    username: string;
    password: string;
    setUsername: (value: string) => void;
    setPassword: (value: string) => void;
    resetUsername: () => void;

}

@inject("appStore", "profileStore")
@observer
class BaseProfileChangeUsername extends React.Component<IProfileChangeUsernameProps, {}> {

    private static panelName = "username";

    public render() {
        const { classes } = this.props;
        const { username, password, loading } = this.props.profileStore!;
        const allowUpdate: boolean = (!!username && !!password);

        const { panelExpanded } = this.props.profileStore!;
        const expanded = (panelExpanded === BaseProfileChangeUsername.panelName);

        return (
            <ExpansionPanel expanded={expanded} onChange={this.onPanelChange}>
                <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
                    <Typography className={classes.heading}>Change Username</Typography>
                </ExpansionPanelSummary>
                <ExpansionPanelDetails>
                    <TextField
                        autoComplete="new-username"
                        className={classes.field}
                        disabled={loading}
                        id="profile-username"
                        label="New Username"
                        margin="normal"
                        value={username}
                        onChange={this.handleChange}
                    />
                    <TextField
                        autoComplete="password"
                        className={classes.field}
                        disabled={loading}
                        id="profile-password"
                        label="Password"
                        margin="normal"
                        type="password"
                        value={password}
                        onChange={this.handleChange}
                    />
                </ExpansionPanelDetails>
                <Divider />
                <ExpansionPanelActions>
                    <Button
                        size="small"
                        disabled={loading}
                        onClick={this.doReset}>
                        Cancel
                    </Button>
                    <Button
                        color="primary"
                        disabled={!allowUpdate || loading}
                        size="small"
                        type="submit"
                        onClick={this.doChangeUsername}>
                        Update
                    </Button>
                </ExpansionPanelActions>
            </ExpansionPanel>
        );
    }

    private onPanelChange = (e: any, expanded: boolean): void => {
        this.props.profileStore!.setPanelExpanded(expanded ? BaseProfileChangeUsername.panelName : "");
    }

    private handleChange = (e: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement | HTMLSelectElement>) => {
        e.preventDefault();
        const { id, value } = e.currentTarget;
        const { setUsername, setPassword } = this.props.profileStore!;

        if (id === "profile-username") {
            setUsername(value);
        } else if (id === "profile-password") {
            setPassword(value);
        } else {
            console.warn("handleChange: no handler for id: ", id); // tslint:disable-line:no-console
        }
    }

    private doReset = (e: React.MouseEvent<Element>) => {
        e.preventDefault();
        this.props.profileStore!.resetUsername();
    }

    private doChangeUsername = (e: React.MouseEvent<Element>) => {
        e.preventDefault();

        const { showNotification } = this.props.appStore!;
        const { username, password } = this.props.profileStore!;

        this.props.profileStore!.changeusername(username, password).then((rsp: IChangeProfileResponse) => {
            switch (rsp.status) {
                default:
                    showNotification(rsp.message, "error");
                    break;
                case ChangeProfileStatus.Success:
                    showNotification("Username updated succcessfully", "success");
                    break;
            }
        });

    }

}

export const ProfileChangeUsername = withTheme()(withStyles(styles)(BaseProfileChangeUsername));
