import React from 'react';
import {
    EuiBasicTable,
    EuiButtonIcon,
    EuiFieldText,
    EuiFlexGroup,
    EuiFlexItem,
    htmlIdGenerator,
} from '@elastic/eui';
import { ICommand } from '../sideeditor';

interface Props {
    editMode: boolean
    commands: Array<ICommand>
    onChangeCommand: (s: Array<ICommand>) => void
}

interface IComponentState {
    update: boolean
    data: Array<ICommand>
}

export class CommandsTable extends React.Component<Props, IComponentState> {
    private tableKey = htmlIdGenerator('tableTest')()
    constructor(props: Props) {
        super(props);
        this.state = {
            data: props.commands,
            update: false,
        }
    }

    componentDidMount() {
    }

    componentDidUpdate() {
    }

    componentWillUnmount() {
    }

    onChangeCommand(s: Array<ICommand>) {
        this.props.onChangeCommand(s)
        this.setState({ update: !this.state.update })
    }

    update() {
        const { update } = this.state;
        this.setState({ update: !update });
    }

    getTestColumns() {
        const columns = [
            {
                field: 'id',
                name: 'Test ID',
                width: '10%',
                render: (id: string, item: ICommand) => {
                    return <EuiFieldText
                        fullWidth
                        value={id} />
                }
            },
            {
                field: 'comment',
                name: 'Comment',
                width: '15%',
                render: (comment: string, item: ICommand) => {
                    return <EuiFieldText fullWidth onChange={(e) => {
                        item.comment = e.currentTarget.value
                        this.onChangeCommand(this.state.data)
                    }
                    } value={comment} />
                }
            },
            {
                field: 'command',
                name: 'Command',
                width: '15%',
                render: (cmd: string, item: ICommand) => {
                    return <EuiFieldText fullWidth onChange={(e) => {
                        item.command = e.currentTarget.value
                        this.onChangeCommand(this.state.data)
                    }
                    } value={cmd} />
                }
            },
            {
                field: 'target',
                name: 'Target',
                width: '15%',
                render: (trg: string, item: ICommand) => {
                    return <EuiFieldText fullWidth onChange={(e) => {
                        item.target = e.currentTarget.value
                        this.onChangeCommand(this.state.data)
                    }
                    } value={trg} />
                }
            },
            {
                field: 'targets',
                name: 'Targets',
                width: '15%',
                render: (trgs: string, item: ICommand) => {
                    return <EuiFieldText fullWidth onChange={(e) => {
                        item.targets = e.currentTarget.value ? e.currentTarget.value.split(",") : []
                        this.onChangeCommand(this.state.data)
                    }
                    } value={trgs} />
                }
            },
            {
                field: 'value',
                name: 'Value',
                width: '15%',
                render: (val: string, item: ICommand) => {
                    return <EuiFieldText fullWidth onChange={(e) => {
                        item.value = e.currentTarget.value
                        this.onChangeCommand(this.state.data)
                    }
                    } value={val} />
                }
            },
            {
                name: "",
                width: '10%',
                actions: [
                    {
                        name: "",
                        render: (item: ICommand) => {
                            return (<EuiButtonIcon
                                aria-label='up'
                                onClick={() => {
                                    const idx = this.state.data.findIndex((el) => {
                                        return item.id === el.id
                                    })
                                    if (idx > 0) {
                                        const prevEl = this.state.data[idx - 1]
                                        this.state.data[idx - 1] = this.state.data[idx]
                                        this.state.data[idx] = prevEl
                                        this.onChangeCommand(this.state.data)
                                    }
                                }}
                                iconType="sortUp"
                            />)
                        },
                    },
                    {
                        name: "",
                        render: (item: ICommand) => {
                            return (<EuiButtonIcon
                                aria-label='down'
                                onClick={() => {
                                    const idx = this.state.data.findIndex((el) => {
                                        return item.id === el.id
                                    })
                                    if (idx < this.state.data.length - 1) {
                                        const nextEl = this.state.data[idx + 1]
                                        this.state.data[idx + 1] = this.state.data[idx]
                                        this.state.data[idx] = nextEl
                                        this.onChangeCommand(this.state.data)
                                    }
                                }}
                                iconType="sortDown"
                            />)
                        },
                    }]
            },
            {
                name: <EuiButtonIcon iconType="plusInCircle"
                    aria-label='add command'
                    onClick={() => {
                        const newCmd: ICommand = {
                            id: crypto.randomUUID(),
                            command: "",
                            comment: "",
                            target: "",
                            targets: [],
                            value: "",
                        }
                        this.setState({ data: [newCmd, ...this.state.data] })
                        this.onChangeCommand(this.state.data)
                    }} />,
                width: '5%',
                actions: [{
                    name: "",
                    render: (item: ICommand) => {
                        return (<EuiButtonIcon
                            aria-label='remove'
                            onClick={() => {
                                const idx = this.state.data.findIndex((el) => {
                                    return item.id === el.id
                                })
                                this.state.data.splice(idx, 1)
                                this.onChangeCommand(this.state.data)
                            }}
                            iconType="cross"
                        />)
                    }
                }
                ]
            },
        ];
        return columns;
    }

    render() {
        return (
            <EuiFlexItem grow={true}>
                <EuiFlexGroup responsive={false} gutterSize='none' direction="column">
                    <EuiBasicTable
                        key={this.tableKey}
                        items={this.state.data}
                        itemId="id"
                        columns={this.getTestColumns()}
                    />
                </EuiFlexGroup>
            </EuiFlexItem>
        );
    }
}
