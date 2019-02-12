import "typeface-roboto";

import * as React from "react";
import * as ReactDOM from "react-dom";

import { MuiThemeProvider } from "@material-ui/core/styles";
import CssBaseline from "@material-ui/core/CssBaseline";
import { BrowserRouter as Router } from "react-router-dom";

import { Root } from "./components/Root";
import { theme } from "./themes";

import { Provider } from "mobx-react";
import { stores } from "./stores";

ReactDOM.render(
    <React.Fragment>
        <CssBaseline />
        <MuiThemeProvider theme={theme}>
            <Router>
                <Provider {...stores}>
                    <Root />
                </Provider>
            </Router>
        </MuiThemeProvider>
    </React.Fragment>,
    document.getElementById("app"),
);
