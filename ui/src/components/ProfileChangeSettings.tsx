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

export interface IProfileChangeSettingsProps extends WithStyles<typeof styles> {
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
    changesettings: (displayname: string) => Promise<IChangeProfileResponse>;

    displayname: string;
    setDisplayname: (value: string) => void;
    resetSettings: () => void;
}

@inject("appStore", "profileStore")
@observer
class BaseProfileChangeProfile extends React.Component<IProfileChangeSettingsProps, {}> {

    private static panelName = "settings";

    public render() {
        const { classes } = this.props;
        const { displayname, loading } = this.props.profileStore!;
        const allowUpdate: boolean = (!!displayname);

        const { panelExpanded } = this.props.profileStore!;
        const expanded = (panelExpanded === BaseProfileChangeProfile.panelName);

        return (
            <ExpansionPanel expanded={expanded} onChange={this.onPanelChange}>
                <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
                    <Typography className={classes.heading}>Change Settings</Typography>
                </ExpansionPanelSummary>
                <ExpansionPanelDetails>
                    <TextField
                        className={classes.field}
                        disabled={loading}
                        id="profile-displayname"
                        label="display name"
                        margin="normal"
                        value={displayname}
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
                        onClick={this.doChangeSettings}>
                        Update
                    </Button>
                </ExpansionPanelActions>
            </ExpansionPanel>
        );
    }

    private onPanelChange = (e: any, expanded: boolean): void => {
        this.props.profileStore!.setPanelExpanded(expanded ? BaseProfileChangeProfile.panelName : "");
    }

    private handleChange = (e: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement | HTMLSelectElement>) => {
        e.preventDefault();
        const { id, value } = e.currentTarget;

        if (id === "profile-displayname") {
            this.props.profileStore!.setDisplayname(value);
        } else {
            console.warn("handleChange: no handler for id: ", id); // tslint:disable-line:no-console
        }
    }

    private doReset = (e: React.MouseEvent<Element>) => {
        e.preventDefault();
        this.props.profileStore!.resetSettings();
    }

    private doChangeSettings = (e: React.MouseEvent<Element>) => {
        e.preventDefault();

        const { showNotification } = this.props.appStore!;
        const { displayname } = this.props.profileStore!;

        this.props.profileStore!.changesettings(displayname).then((rsp: IChangeProfileResponse) => {
            switch (rsp.status) {
                default:
                    showNotification(rsp.message, "error");
                    break;
                case ChangeProfileStatus.Success:
                    showNotification("Settings updated succcessfully", "success");
                    break;
            }
        });

    }

}

export const ProfileChangeProfile = withTheme()(withStyles(styles)(BaseProfileChangeProfile));
