import React, { ReactNode } from "react";

import {
    EuiPageTemplate,
} from '@elastic/eui';

//import { defaultCookies } from './libraries/cookies';
import { bus, defaultEmitter } from './libraries/bus';
import { SideEditor } from "./components/sideeditor/sideeditor";


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

    constructor(props: Props) {
        super(props);
        this.bus = new bus(defaultEmitter().getId());

        this.state = {
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
            <div></div>
        );
    }

    render() {
        return (
            <EuiPageTemplate>
                {/* <Toaster/> */}
                <SideEditor editMode={true}/>
            </EuiPageTemplate>

        )
    }
}