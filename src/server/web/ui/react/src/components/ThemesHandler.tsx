//import '@elastic/eui/dist/eui_theme_light.css';
//import '@elastic/eui/dist/eui_theme_dark.css';

import '../css/app.css';
import '../css/common.css';

import React from "react";
import {EuiComboBox, EuiProvider, EuiThemeProvider} from "@elastic/eui";
import {EuiThemeColorMode, EuiThemeModifications} from "@elastic/eui/src/services/theme/types";
import {EuiComboBoxOptionOption} from "@elastic/eui/src/components/combo_box/types";

import themeLight from '../scss/light.scss';
//import themeDark from '../scss/dark.scss';

//eui/src-docs/src/
//eui/src-docs/src/services/theme/

const themes: any = {};

export function registerTheme(theme: string, cssFiles: any) {
    themes[theme] = cssFiles;
    console.log('registering theme')
    console.log(cssFiles)
}

export function applyTheme(newTheme: string) {
    Object.keys(themes).forEach((theme) =>
        themes[theme].forEach((cssFile: any) => {
            cssFile.unuse()
        })
    );
    themes[newTheme]?.forEach((cssFile: any) => {
        cssFile.use()
    });
}

//TODO in webpack.config.js
//scss -> scss_disabled

registerTheme('light', [themeLight]);
//registerTheme('dark', [themeDark]);

//applyTheme('light')

//emptyShade => default Panel Color
//shadow = default Shadow
//body: lightseagreen


interface theme extends EuiThemeModifications {
    name: string
    mode: EuiThemeColorMode
}

//"ink": "black",
//"ghost": "blue",

const romaTheme: theme = {
    name: "roma",
    mode: "dark",
    colors: {
        LIGHT: {},
        DARK: {
            "body": "lightseagreen",
            "emptyShade": "lightblue",
            "shadow": "green",

            "primary": "red",
            "accent": "pink",
            "success": "green",
            "warning": "yellow",
            "danger": "red",

            "lightestShade": "gray",
            "lightShade": "blue",

            "mediumShade": "darkgray",

            "darkShade": "darkblue",
            "darkestShade": "blue",

            "fullShade": "brown",

            "highlight": "#fff9e8",
            "disabled": "#ABB4C4",

            "disabledText": "#a2abba",
            "primaryText": "#006bb8",
            "accentText": "#ba3d76",
            "successText": "#007871",
            "warningText": "#83650a",
            "dangerText": "#bd271e",
            "text": "#343741",

            "subduedText": "#646a77",

            "title": "#1a1c21",

            "link": "#006bb8",
        }
    },
};

const defaultTheme: theme = {
    name: "default",
    mode: "light",
    colors: {
        LIGHT: {
        },
        DARK: {
        }
    }
}

interface IProps {
    target: React.ReactNode
}

interface IComponentState {
    current: theme
    currentId: string
}

export class ThemesHandler extends React.Component<IProps, IComponentState> {
    private static instance: ThemesHandler | null = null
    private avail = new Map<string, theme>;

    constructor (props: IProps) {
        super(props);

        this.state = { currentId: "default", current: defaultTheme };
        this.avail.set(this.state.currentId, this.state.current);
        this.avail.set("roma", romaTheme);
        ThemesHandler.instance = this;
    }

    public static cId() : string {
        return "ThemesHandler"
    }

    public static get() : Array<string> {
        const themes : Array<string> = [];
        const self = ThemesHandler.instance;
        if (!self) {
            console.log("MISSING THEMES INSTANCE!");
            return themes;
        }
        self.avail.forEach((v, k) => {
            themes.push(k)
        })
        return themes;
    }

    public static getCurrent() : string {
        const self = ThemesHandler.instance;
        if (!self) {
            console.log("MISSING THEMES INSTANCE!");
            return "";
        }
        return self.state.currentId
    }

    public static set(id: string) : boolean {
        const self = ThemesHandler.instance;
        if (!self) {
            console.log("MISSING THEMES INSTANCE!");
            return false;
        }
        const theme = self.avail.get(id)
        if (!theme) {
            return false;
        }
        self.setState({current: theme, currentId: id})
        return true
    }

    render() {
        //const { euiTheme } = useEuiTheme();

        const { current } = this.state;
        return(
            <EuiProvider colorMode={current.mode}>
                <EuiThemeProvider  modify={current}>
                    { this.props.target }
                </EuiThemeProvider>
            </EuiProvider>
        )
    }
}



interface IHelperProps {
}

interface IHelperComponentState {
    update: boolean
}

export class ThemesOptions extends React.Component<IHelperProps, IHelperComponentState> {
    constructor (props: IHelperProps) {
        super(props);
        this.state = {update: false};
    }

    private getOptions() : Array<EuiComboBoxOptionOption>{
        const themes : Array<EuiComboBoxOptionOption> = [];
        ThemesHandler.get().forEach(k => { themes.push({label: k})})
        return themes
    }

    private setOption(theme: string) {
        if (ThemesHandler.set(theme)) {
            this.setState({update: !this.state.update})
        }
    }

    private static getOptionCurrent() : Array<EuiComboBoxOptionOption> {
        return [{label: ThemesHandler.getCurrent()}]
    }

    render() {
        return(
            <EuiComboBox
                aria-label="Accessible screen reader label"
                placeholder="Select a Theme"
                singleSelection={{ asPlainText: true }}
                options={this.getOptions()}
                selectedOptions={ ThemesOptions.getOptionCurrent() }
                onChange={(options)=> {
                    if (options && options.length > 0) {
                        this.setOption(options[0].label)
                    }
                }}
            />
        )
    }
}