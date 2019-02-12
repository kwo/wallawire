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

export interface IProfileChangePasswordProps extends WithStyles<typeof styles> {
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
    changepassword: (newpassword: string, oldpassword: string) => Promise<IChangeProfileResponse>;

    passwordOld: string;
    passwordNew: string;
    passwordConfirm: string;
    warning: string;
    warning1: boolean;
    warning2: boolean;
    setPasswordOld: (value: string) => void;
    setPasswordNew: (value: string) => void;
    setPasswordConfirm: (value: string) => void;
    setWarning: (value: string) => void;
    setWarning1: (value: boolean) => void;
    setWarning2: (value: boolean) => void;
    resetPassword: () => void;
}

@inject("appStore", "profileStore")
@observer
class BaseProfileChangePassword extends React.Component<IProfileChangePasswordProps, {}> {

    private static panelName = "password";

    public render() {
        const { classes } = this.props;
        const { loading, passwordOld, passwordNew, passwordConfirm, warning, warning1, warning2 }
            = this.props.profileStore!;
        const allowChangePassword: boolean = (!!passwordOld && !!passwordNew && !!passwordConfirm
            && passwordNew === passwordConfirm);

        const { panelExpanded } = this.props.profileStore!;
        const expanded = (panelExpanded === BaseProfileChangePassword.panelName);

        return (
            <ExpansionPanel expanded={expanded} onChange={this.onPanelChange}>
                <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
                    <Typography className={classes.heading}>Change Password</Typography>
                </ExpansionPanelSummary>
                <ExpansionPanelDetails>
                    <TextField
                        autoComplete="old-password"
                        className={classes.field}
                        disabled={loading}
                        id="cp-oldpassword"
                        label="Old Password"
                        margin="normal"
                        type="password"
                        value={passwordOld}
                        onChange={this.handleChange}
                    />
                </ExpansionPanelDetails>
                <ExpansionPanelDetails>
                    <TextField
                        autoComplete="new-password"
                        className={classes.field}
                        disabled={loading}
                        error={warning1}
                        helperText={warning}
                        id="cp-newpassword"
                        label="New Password"
                        margin="normal"
                        type="password"
                        value={passwordNew}
                        onChange={this.handleChange}
                    />
                </ExpansionPanelDetails>
                <ExpansionPanelDetails>
                    <TextField
                        autoComplete="confirm-password"
                        className={classes.field}
                        disabled={loading}
                        error={warning2}
                        helperText={warning}
                        id="cp-confirmpassword"
                        label="Confirm Password"
                        margin="normal"
                        type="password"
                        value={passwordConfirm}
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
                        disabled={!allowChangePassword || loading}
                        size="small"
                        type="submit"
                        onClick={this.doChangePassword}>
                        Update
                    </Button>
                </ExpansionPanelActions>
            </ExpansionPanel>
        );
    }

    private onPanelChange = (e: any, expanded: boolean): void => {
        this.props.profileStore!.setPanelExpanded(expanded ? BaseProfileChangePassword.panelName : "");
    }

    private handleChange = (e: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement | HTMLSelectElement>) => {
        e.preventDefault();
        const { id, value } = e.currentTarget;
        const { passwordNew, passwordConfirm, setPasswordNew, setPasswordOld, setPasswordConfirm,
            setWarning, setWarning1, setWarning2 } = this.props.profileStore!;
        // console.debug("handleChange", id, value); // tslint:disable-line:no-console

        if (id === "cp-oldpassword") {
            setPasswordOld(value);
            setWarning("");
        } else if (id === "cp-newpassword") {
            if (value !== passwordConfirm && passwordConfirm) {
                setWarning("passwords do not match");
                setPasswordNew(value);
                setWarning1(true);
                setWarning2(false);
            } else {
                setWarning("");
                setPasswordNew(value);
                setWarning1(false);
                setWarning2(false);
            }
        } else if (id === "cp-confirmpassword") {
            if (value !== passwordNew) {
                setPasswordConfirm(value);
                setWarning("passwords do not match");
                setWarning1(false);
                setWarning2(true);
            } else {
                setPasswordConfirm(value);
                setWarning("");
                setWarning1(false);
                setWarning2(false);
            }
        } else {
            console.warn("handleChange: no handler for id: ", id); // tslint:disable-line:no-console
        }
    }

    private doReset = (e: React.MouseEvent<Element>) => {
        e.preventDefault();
        this.props.profileStore!.resetPassword();
    }

    private doChangePassword = (e: React.MouseEvent<Element>) => {
        e.preventDefault();

        const { showNotification } = this.props.appStore!;
        const { passwordOld, passwordNew, passwordConfirm, setWarning }
            = this.props.profileStore!;

        // clear message, set loading
        setWarning("");

        // confirm password
        if (passwordNew !== passwordConfirm) {
            setWarning("passwords do not match");
            return;
        }

        this.props.profileStore!.changepassword(passwordNew, passwordOld).then((rsp: IChangeProfileResponse) => {
            switch (rsp.status) {
                default:
                    showNotification(rsp.message, "error");
                    break;
                case ChangeProfileStatus.Success:
                    showNotification("Password updated succcessfully", "success");
                    break;
            }
        });

    }

}

export const ProfileChangePassword = withTheme()(withStyles(styles)(BaseProfileChangePassword));
