import React from 'react';
import {
    EuiButton,
    EuiFilePicker,
    EuiFlexGroup,
    EuiFlexItem,
    EuiPanel,
    EuiTab,
    EuiTabbedContentTab,
    EuiTabs,
    htmlIdGenerator,
} from '@elastic/eui';
import { Data } from './dependencies/data';
import { DataTable } from './dependencies/datatable';

export interface ICommand {
    id: string
    comment: string
    command: string
    target: string
    targets: Array<any>
    value: string
}

export interface ITest {
    id: string
    name: string
    commands: Array<ICommand>
}

export interface ISuite {
    id: string
    name: string
    persistSession: boolean
    parallel: boolean
    timeout: number
    tests: Array<string>
}

export interface ISidefile {
    id: string
    name: string
    url: string
    tests: Array<ITest>
    suites: Array<ISuite>
    urls: Array<string>
    plugins: Array<any>
    version: string
}

interface Props {
    editMode: boolean
}

interface IComponentState {
    editMode: boolean
    update: boolean
}

interface IStyle {
    minHeight: number
    minWidth: number
}

const personalStyle: IStyle = {
    minHeight: 450,
    minWidth: 660
}

export class SideEditor extends React.Component<Props, IComponentState> {
    private idxTabClick: string = 'view';
    private tabs: Array<EuiTabbedContentTab> = []
    private sidefile: ISidefile = {
        id: "",
        name: "",
        plugins: [],
        suites: [],
        tests: [],
        url: "",
        urls: [],
        version: ""
    }
    private stringFile: string = ""
    private sideFileEmpty: ISidefile = {
        id: "",
        name: "",
        plugins: [],
        suites: [],
        tests: [],
        url: "",
        urls: [],
        version: ""
    }

    constructor(props: Props) {
        super(props);
        this.stringFile = JSON.stringify(this.sideFileEmpty, null, 2)
        this.state = {
            editMode: true,
            update: false,
        }
    }

    componentDidMount() {
    }

    componentDidUpdate() {
    }

    componentWillUnmount() {
    }

    onClickApply() {
        console.log("click apply")
    }

    onClickCancel() {
    }

    setEditMode(mode: boolean) {
        this.setState({ editMode: mode })
    }

    onClickGEE() {
        this.setEditMode(true)
    }

    onClickGVE() {
        this.setEditMode(false)
    }

    onChangeName(e: React.ChangeEvent<HTMLInputElement>) {
        console.log("change name")
    }

    buildTabScript() {
        return {
            id: 'script',
            name: 'Editor JSON',
            content: (
                <EuiFlexGroup responsive={false} style={{ paddingTop: '10px', paddingBottom: '30px' }}
                    gutterSize='s' direction="row" className="eui-yScrollWithShadows">
                    <Data key={htmlIdGenerator("datascript")()} editMode={this.state.editMode} data={this.stringFile}
                        onChangeScript={(s: string) => {
                            this.sidefile = JSON.parse(s)
                            this.stringFile = s
                            this.update()
                        }} />;
                </EuiFlexGroup>
            ),
        };
    }

    buildTabView() {
        return {
            id: 'view',
            name: 'View',
            content: (
                <EuiFlexGroup responsive={false} style={{ position: 'relative', height: "100%", width: "100%", paddingTop: '10px', paddingBottom: '30px' }}
                    gutterSize='s' direction="column" className="eui-yScrollWithShadows">
                    <DataTable key={htmlIdGenerator("datatable")()} editMode={this.state.editMode} sideFile={this.sidefile}
                        onChangeValue={(s: ISidefile) => {
                            this.sidefile = s
                            this.stringFile = JSON.stringify(s, null, 2)
                            this.update()
                        }} />
                </EuiFlexGroup>
            ),
        };
    }

    buildTabs() {
        const tabView = this.buildTabView()
        const tabScript = this.buildTabScript();
        return [tabView, tabScript]
    }

    update() {
        const { update } = this.state
        this.setState({ update: !update })
    }

    async onClickSelectFile(files: FileList | null) {
        if ((files !== null) && (files !== undefined)) {
            if (files.length > 0) {
                const txt = await files[0].text()
                const of: ISidefile = JSON.parse(txt)
                this.sidefile = of
                this.stringFile = txt

            } else {
                this.sidefile = this.sideFileEmpty;
                this.stringFile = JSON.stringify(this.sideFileEmpty, null, 2)
            }
        }
        this.update()
    }

    render() {
        this.tabs = this.buildTabs()
        const renderTabs = () => {
            return this.tabs.map((tab, index) => (
                <EuiTab
                    key={index}
                    onClick={() => { this.idxTabClick = tab.id; this.update() }}
                    isSelected={tab.id === this.idxTabClick}
                    disabled={tab.disabled}
                    prepend={tab.prepend}
                    append={tab.append}
                >
                    {tab.name}
                </EuiTab>
            ));
        };
        let tabSel: React.ReactNode = false
        for (let i = 0; i < this.tabs.length; i++) {
            if (this.tabs[i].id === this.idxTabClick) {
                tabSel = this.tabs[i].content
                break;
            }
        }
        const principalItem =
            <EuiFlexGroup responsive={false} direction='column' style={{ height: '100%' }} gutterSize='none'>
                <EuiFlexItem grow={false} style={{ marginTop: "8px", marginBottom: "0px" }}>
                    <EuiFlexGroup responsive={false} direction="row" gutterSize='s'>
                        <EuiFlexItem grow={true}>
                            <EuiFilePicker
                                fullWidth
                                id={htmlIdGenerator('filePicker')()}
                                disabled={!this.props.editMode}
                                multiple={false}
                                initialPromptText="Select or drag and drop file"
                                onChange={(files: FileList | null) => { this.onClickSelectFile(files) }}
                                display="default"
                                accept=".side"
                            />
                        </EuiFlexItem>
                    </EuiFlexGroup>
                </EuiFlexItem>
                <EuiFlexItem grow={false}>
                    <EuiTabs>{renderTabs()}</EuiTabs>
                </EuiFlexItem>
                <EuiFlexItem style={{ minHeight: personalStyle.minHeight, minWidth: personalStyle.minWidth }}>
                    {tabSel}
                </EuiFlexItem>
            </EuiFlexGroup>
        return (
            <EuiPanel style={{ height: "100%", paddingBottom: 4, paddingTop: 2 }}>
                <EuiFlexGroup responsive={false} direction="column" style={{ height: "100%" }} gutterSize='none'>
                    <EuiFlexItem style={{ height: "100%" }}>
                        {principalItem}
                    </EuiFlexItem>
                    <EuiFlexItem grow={false} css={{ marginBottom: "15px !important" }}>
                        <EuiFlexGroup responsive={false} gutterSize='s' wrap>
                            {this.state.editMode ?
                                <EuiFlexItem grow={false}>
                                    <EuiButton iconType='lockOpen' onClick={() => { this.onClickGVE() }} />
                                </EuiFlexItem>
                                :
                                <EuiFlexItem grow={false}>
                                    <EuiButton iconType='lock' onClick={() => { this.onClickGEE() }} />
                                </EuiFlexItem>
                            }
                            <EuiFlexItem grow={false}>
                                <EuiButton disabled={!this.state.editMode} onClick={() => { this.onClickApply() }}>Apply</EuiButton>
                            </EuiFlexItem>
                            <EuiFlexItem grow={false}>
                                <EuiButton onClick={() => { this.onClickCancel() }}>Cancel</EuiButton>
                            </EuiFlexItem>
                        </EuiFlexGroup>
                    </EuiFlexItem>
                </EuiFlexGroup>
            </EuiPanel >
        );
    }
}
