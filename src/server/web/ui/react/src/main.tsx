import React from "react";
import ReactDOM from "react-dom";
import { App } from "./App";
import { ThemesHandler } from './components/ThemesHandler';

function importAll(r: any) {
    r.keys().forEach(r);
}

importAll(require.context('@elastic/eui/es/components/icon/assets', true, /\.js$/));

ReactDOM.render(
    <React.StrictMode>
        <ThemesHandler target={ <App />}/>
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