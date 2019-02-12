import { createStyles, withStyles, withTheme, Theme, WithStyles } from "@material-ui/core/styles";
import * as React from "react";
import { RouteComponentProps } from "react-router";

import { ProfileChangePassword } from "./ProfileChangePassword";
import { ProfileChangeProfile } from "./ProfileChangeSettings";
import { ProfileChangeUsername } from "./ProfileChangeUsername";

// tslint:disable:object-literal-sort-keys
const styles = (theme: Theme) => createStyles({
    root: {
        maxWidth: 600,
    },
});
// tslint:enable:object-literal-sort-keys

export interface IProfileProps extends RouteComponentProps, WithStyles<typeof styles> { }

class BaseProfile extends React.Component<IProfileProps, {}> {

    public render() {
        const { classes } = this.props;

        return (
            <div className={classes.root}>
                <ProfileChangePassword />
                <ProfileChangeUsername />
                <ProfileChangeProfile />
            </div>
        );
    }

}

export const Profile = withTheme()(withStyles(styles)(BaseProfile));
