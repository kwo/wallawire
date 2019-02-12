import * as React from "react";
import { RouteComponentProps } from "react-router";

export interface ILogoutProps extends RouteComponentProps {
    logout: () => Promise<any>;
}

export class Logout extends React.Component<ILogoutProps, {}> {

    public componentDidMount() {
        this.props.logout();
    }

    public render() {
        return <div/>;
    }

}
