import React from "react";
import ReactDOM from "react-dom";
import { App } from "./App";
import { EuiThemeProvider } from "@elastic/eui";

function importAll(r: any) {
    r.keys().forEach(r);
}

importAll(require.context('@elastic/eui/es/components/icon/assets', true, /\.js$/));

ReactDOM.render(
    <React.StrictMode>
        <EuiThemeProvider>
            <App />
        </EuiThemeProvider>
    </React.StrictMode>,

    document.getElementById("root")
);

/*
ReactDOM.render(
    <React.StrictMode>
        <App />
    </React.StrictMode>,
    document.getElementById("root")
);

*/