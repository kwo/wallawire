import green from "@material-ui/core/colors/green";
import purple from "@material-ui/core/colors/purple";
import { createMuiTheme } from "@material-ui/core/styles";

export const theme = createMuiTheme({
    palette: {
        error: {
            contrastText: "#ffffff",
            dark: "#aa2e25",
            light: "#f44336",
            main: "#f44336",
        },
        primary: {
            contrastText: "#ffffff",
            dark: "#2c387e",
            light: "#6573c3",
            main: "#3f51b5",
        },
        secondary: {
            contrastText: "#ffffff",
            dark: "#9500ae",
            light: "#dd33fa",
            main: "#d500f9",
        },
    },
    typography: {
        useNextVariants: true,
    },
});
