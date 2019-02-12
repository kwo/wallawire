import { Typography } from "@material-ui/core";
import * as React from "react";

export interface IHelloProps {
    caption: string;
}

export const Hello: React.SFC<IHelloProps> = (props) => {
    return <Typography variant="h6" color="inherit">{props.caption}</Typography>;
};
