import React from 'react';
import {
    EuiFlexItem,
} from '@elastic/eui';
import { IBorderRadius, MonacoContainer } from '../../monacocontainer/monacocontainer'

interface Props {
    editMode: boolean
    data: string
    onChangeScript: (s: string) => void
}

interface IComponentState {
    update: boolean
    data: string
}

export class Data extends React.Component<Props, IComponentState> {
    private borderRadius: IBorderRadius = {
        topLeft: 4,
        topRight: 4,
        bottomRight: 4,
        bottomLeft: 4
    }

    constructor(props: Props) {
        super(props);
        this.state = {
            data: props.data,
            update: false,
        }
    }

    componentDidMount() {
    }

    componentDidUpdate() {
    }

    componentWillUnmount() {
    }

    onChangeScript(s: string) {
        this.props.onChangeScript(s)
        this.setState({ data: s });
    }

    update() {
        const { update } = this.state;
        this.setState({ update: !update });
    }

    render() {
        return (
            <EuiFlexItem grow={true}>
                <MonacoContainer hasShadow={true} border={this.borderRadius} editMode={this.props.editMode} show={true} value={this.state.data} onChangeScript={(x: string) => this.onChangeScript(x)} />
            </EuiFlexItem>
        );
    }
}
