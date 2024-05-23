import React, {ReactNode} from "react";

import {
    EuiPageTemplate,
} from '@elastic/eui';

//import { defaultRest } from './libraries/rest';
//import { defaultStore } from './libraries/store';
import { defaultCookies } from './libraries/cookies';
import { bus, defaultEmitter } from './libraries/bus';
//import { bus, busMessage, defaultEmitter } from './libraries/bus';
//import { websocket } from './libraries/websocket';
//import { mainIconApp, mainIconTitle, mainIconText, mainIconInitialize } from './libraries/icons';
import { H2D } from "./components/H2D";
import { Dimon } from "./components/Dimon";
import { Toaster } from "./components/Toaster";
//import { ThemesOptions } from "./components/ThemesHandler"
//import { roundTwoDigits } from "./libraries/numbers";

interface Props {
    show?: boolean
}

interface IComponentState {
    activeId: string
    activeComponent: ReactNode | null
}

export class App extends React.Component<Props, IComponentState> {
    private bus: bus;
    //private readonly version = "2.2.0"; //TODO FROM API

    constructor (props: Props) {
        super(props);
        this.bus = new bus(defaultEmitter().getId());

        this.state =  {
            activeId: '',
            activeComponent: null
        }
    }

    componentDidMount() {
        //component
    }

    componentDidUpdate() {
        //component
    }

    componentWillUnmount() {
        this.bus.disconnect()
    }

    private renderDimon() {
        return (
            <Dimon session={defaultCookies().get("")}/>
        );
    }

    render () {
        return (
                <EuiPageTemplate>
                    <Toaster/>
                    { this.renderDimon() }
                </EuiPageTemplate>

        )
    }
}