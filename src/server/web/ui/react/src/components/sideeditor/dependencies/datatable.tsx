import React from 'react';
import {
    EuiBasicTable,
    EuiButtonIcon,
    EuiDescriptionList,
    EuiFieldText,
    EuiFlexItem,
    RIGHT_ALIGNMENT,
    htmlIdGenerator,
} from '@elastic/eui';
import { ISidefile, ITest } from '../sideeditor';
import { CommandsTable } from './commandstable';

interface Props {
    editMode: boolean
    sideFile: ISidefile
    onChangeValue: (s: any) => void
}

interface IComponentState {
    update: boolean
    data: ISidefile
}

export class DataTable extends React.Component<Props, IComponentState> {
    private mapManagersExpanded: Record<string, React.ReactNode> = {}
    private tableKey = htmlIdGenerator('tableTest')()
    constructor(props: Props) {
        super(props);
        this.state = {
            data: props.sideFile,
            update: false,
        }
    }

    componentDidMount() {
    }

    componentDidUpdate() {
    }

    componentWillUnmount() {
    }

    onChangeTest(s: ISidefile) {
        this.props.onChangeValue(s)
        this.setState({ update: !this.state.update })
    }

    update() {
        const { update } = this.state;
        this.setState({ update: !update });
    }

    private toggleDetails(el: ITest) {
        if (this.mapManagersExpanded[el.id]) {
            this.mapManagersExpanded = {}
        } else {
            this.mapManagersExpanded = {}
            const listItem = [
                {
                    title: 'Commands',
                    description: <CommandsTable editMode={this.props.editMode} commands={el.commands}
                        onChangeCommand={(s: any) => {
                            el.commands = s
                            this.onChangeTest(this.state.data)
                        }} />
                }
            ]
            this.mapManagersExpanded[el.id] = <EuiDescriptionList key={el.id} listItems={listItem} />
        }
        this.setState({ update: !this.state.update })
    }

    getTestColumns() {
        const columns = [
            {
                field: 'id',
                name: 'Test ID',
                width: "40%",
                render: (id: string, item: ITest) => {
                    return <EuiFieldText
                        fullWidth
                        value={id} />
                }
            },
            {
                field: 'name',
                name: 'Test Name',
                width: "40%",
                render: (name: string, item: ITest) => {
                    return <EuiFieldText
                        fullWidth
                        value={name}
                        onChange={(e) => {
                            item.name = e.currentTarget.value;
                            this.onChangeTest(this.state.data)
                        }} />
                }
            },
            {
                name: "",
                width: '10%',
                actions: [
                    {
                        name: '',
                        render: (item: ITest) => {
                            return (<EuiButtonIcon
                                aria-label='up'
                                onClick={() => {
                                    const idx = this.state.data.tests.findIndex((el) => {
                                        return item.id === el.id
                                    })
                                    console.log("sort up: " + idx)
                                    if (idx > 0) {
                                        const prevEl = this.state.data.tests[idx - 1]
                                        this.state.data.tests[idx - 1] = this.state.data.tests[idx]
                                        this.state.data.tests[idx] = prevEl
                                        this.onChangeTest(this.state.data)
                                    }
                                }}
                                iconType="sortUp"
                            />)
                        },
                    },
                    {
                        name: '',
                        render: (item: ITest) => {
                            return (<EuiButtonIcon
                                aria-label='down'
                                onClick={() => {
                                    const idx = this.state.data.tests.findIndex((el) => {
                                        return item.id === el.id
                                    })
                                    console.log("sort down: " + idx)
                                    if (idx < this.state.data.tests.length - 1) {
                                        const nextEl = this.state.data.tests[idx + 1]
                                        this.state.data.tests[idx + 1] = this.state.data.tests[idx]
                                        this.state.data.tests[idx] = nextEl
                                        this.onChangeTest(this.state.data)
                                    }
                                }}
                                iconType="sortDown"
                            />)
                        },
                    }]
            },
            {
                name: <EuiButtonIcon
                    aria-label='add test'
                    iconType="plusInCircle"
                    onClick={() => {
                        const newTest: ITest = {
                            id: crypto.randomUUID(),
                            name: "",
                            commands: []
                        }
                        this.state.data.tests = [newTest, ...this.state.data.tests]
                        this.onChangeTest(this.state.data)
                    }} />,
                width: '10%',
                actions: [
                    {
                        name: '',
                        render: (item: ITest) => {
                            return (<EuiButtonIcon
                                aria-label='remove'
                                onClick={() => {
                                    const idx = this.state.data.tests.findIndex((el) => {
                                        return item.id = el.id
                                    })
                                    this.state.data.tests.splice(idx, 1)
                                    this.onChangeTest(this.state.data)
                                }}
                                iconType="cross"
                            />)
                        },
                    },
                    {
                        align: RIGHT_ALIGNMENT,
                        isExpander: true,
                        name: '',
                        render: (item: ITest) => {
                            return (<EuiButtonIcon
                                onClick={() => {
                                    this.toggleDetails(item)
                                }}
                                aria-label={this.mapManagersExpanded[item.id] ? 'Collapse' : 'Expand'}
                                iconType={this.mapManagersExpanded[item.id] ? 'arrowUp' : 'arrowDown'}
                            />)
                        },
                    },
                ],
            },
        ];
        return columns;
    }

    render() {
        return (
            <EuiFlexItem grow={true} style={{ position: "absolute" }}>
                <EuiBasicTable
                    key={this.tableKey}
                    items={this.state.data.tests}
                    itemId="id"
                    itemIdToExpandedRowMap={this.mapManagersExpanded}
                    columns={this.getTestColumns()}
                />
            </EuiFlexItem>
        );
    }
}
