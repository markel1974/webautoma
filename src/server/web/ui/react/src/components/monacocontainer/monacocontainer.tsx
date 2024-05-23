import React from 'react';
import MonacoEditor, { monaco } from 'react-monaco-editor';
import {
  EuiButtonIcon,
  EuiResizeObserver,
} from '@elastic/eui';
import { EuiButtonGroupOptionProps } from "@elastic/eui/src/components/button/button_group/button_group";

export interface IBorderRadius {
  bottomLeft: string | number
  bottomRight: string | number
  topLeft: string | number
  topRight: string | number
}

interface Props {
  show: boolean
  editMode: boolean
  height?: number
  width?: number
  language?: string,
  fullEditor?: boolean
  value?: string
  onChangeScript?: (script: string) => void;
  border: IBorderRadius
  hasShadow: boolean
}

interface IComponentState {
  dialogId: string;
  dialogProp: string;
  endpoint: string;
  options: string;
  codeChanged: boolean;
  update: boolean;
  editorFullScreen: boolean;
}

const _themeId = 'myCustomTheme'

export class MonacoContainer extends React.Component<Props, IComponentState> {
  private code: string;
  private editor: monaco.editor.IStandaloneCodeEditor | null = null
  private height: number = 0
  private width: number = 0
  private borderRadius: IBorderRadius
  private boxShadow: string = "rgba(0, 0, 0, 0.1) 0px 3px 4px 2px, rgba(0, 0, 0, 0.09) 0px 5px 8px 3px, rgba(0, 0, 0, 0.08) 0px 7px 12px 4px, rgba(0, 0, 0, 0.07) 0px 15px 15px -1px"

  constructor(props: Props) {
    super(props);
    this.code = '';
    if (props.value) this.code = props.value;
    this.state = {
      codeChanged: false,
      dialogId: '',
      dialogProp: '',
      endpoint: '',
      options: '',
      update: false,
      editorFullScreen: false,
    };
    this.borderRadius = props.border
  }

  componentDidMount() {
  }

  componentDidUpdate() {
  }

  componentWillUnmount() {
  }

  editorDidMount(editor: monaco.editor.IStandaloneCodeEditor, monacoT: typeof monaco) {
    this.setupTheme(monacoT)
    this.editor = editor
  }

  //see https://github.com/brijeshb42/monaco-themes/tree/master/themes
  setupTheme(monacoT: typeof monaco) {
    //.monaco-mouse-cursor-text NON DEVE ESSERE DEFINITO!!!!!
    monacoT.editor.defineTheme(_themeId, {
      "base": "vs-dark",
      "inherit": true,
      "rules": [],
      "colors": {
        "editorCursor.foreground": "#0FFFFF",
        "editor.selectionBackground": "#834C5ECC",
      }
    })
    monacoT.editor.setTheme(_themeId)
  }

  startFullScreen() {
    const self = this
    this.setState({ editorFullScreen: true }, () => {
      setTimeout(() => {
        self.setState({ editorFullScreen: false })
      }, 250)
    });
  }

  createButtons() {
    const buttons = ['New:indexEdit', "Delete:eraser", "Save:save", "Options:package"];
    const groupButtons: Array<EuiButtonGroupOptionProps> = [];
    buttons.forEach((x) => {
      const z = x.split(':');
      groupButtons.push({ id: z[0], label: z[0], iconSide: 'left', iconType: z[1] })
    });
    return groupButtons
  }

  updateOptions(code: string): string {
    let opt = '';
    if (code) {
      const regexpOption = /^#options:\s*(.+)$/m;
      const match = code.match(regexpOption);
      if (match && match.length > 1) {
        opt = match[1];
      }
    }
    return opt;
  }

  private doLayoutFullScreen() {
    const rect = document.documentElement.getBoundingClientRect()
    this.doLayout(rect.width, rect.height)
  }

  private doLayoutThis() {
    this.doLayout(this.width, this.height)
  }

  private doLayout(width: number, height: number) {
    this.editor?.layout({ width: width, height: height } as monaco.editor.IDimension)
  }

  private getEditorKey() {
    return "divMonacoEditorSideEditor";
  }

  private setFullscreen() {
    const self = this
    const container = document.getElementById(this.getEditorKey())
    if (!container) {
      console.log("monacoContainer: missing container")
      return
    }
    let elem = container.getElementsByClassName("react-monaco-editor-container")[0];
    if (!elem) {
      console.log("monacoContainer: missing element")
      return
    }
    elem.onfullscreenchange = ((ev: Event) => {
      if (!document.fullscreenElement) {
        setTimeout(() => { self.doLayoutThis() }, 0)
      } else {
        self.doLayoutFullScreen()
      }
    })
    elem.requestFullscreen().catch((err) => {
      alert(`Error attempting to enable fullscreen mode: ${err.message} (${err.name})`);
    });
  }

  onCodeChange(x: string | undefined) {
    const { codeChanged } = this.state;
    this.code = x ? x : '';
    if (this.props.onChangeScript) {
      this.props.onChangeScript(this.code);
    }
    if (!codeChanged) {
      this.setState({ codeChanged: true })
    }
  }

  /*
  renderEditor() {
    const language = this.props.language ? this.props.language : "javascript"
    return (
      <EuiResizeObserver onResize={(dimension: { height: number; width: number; }) => { this.height = dimension.height; this.width = dimension.width; this.doLayoutThis() }}>
        {(resizeRef: React.LegacyRef<HTMLDivElement> | undefined) => (
          <div id={this.getEditorKey()} ref={resizeRef}
            style={{
              boxShadow: this.props.hasShadow ? this.boxShadow : 'none',
              borderTopLeftRadius: this.borderRadius.topLeft,
              borderTopRightRadius: this.borderRadius.topRight,
              borderBottomLeftRadius: this.borderRadius.bottomLeft,
              borderBottomRightRadius: this.borderRadius.bottomRight,
              position: 'relative', width: '100%', height: '100%', overflow: "hidden"
            }}>
            <div style={{ position: 'absolute', width: '100%', height: '100%', overflow: "hidden" }}>
              <EuiButtonIcon aria-label={"fullscreen"} color="danger" iconType="fullScreen" size="s" iconSize="s" title="toggle fullscreen" onClick={() => { this.startFullScreen() }} style={{ zIndex: 1000, position: "absolute", top: 0, right: 0 }}></EuiButtonIcon>
              <MonacoEditor
                language={language}
                defaultValue={this.code}
                value={this.code}
                //theme="vs-dark"
                //theme="hc-black"
                //theme={_themeId}
                options={{
                  selectOnLineNumbers: true,
                  automaticLayout: true,
                  quickSuggestions: true,
                  quickSuggestionsDelay: 500,
                  readOnly: !this.props.editMode,
                  minimap: { enabled: true }
                }}
                onChange={(x: string | undefined) => this.onCodeChange(x)}
                editorDidMount={(editor: any, monacoT: any) => { this.editorDidMount(editor, monacoT) }}
              />
            </div>
          </div>)}
      </EuiResizeObserver>
    )
  }
  */


  renderEditor() {
    const language = this.props.language ? this.props.language : "javascript"
    return (
              <div id={this.getEditorKey()}
                   style={{
                     boxShadow: this.props.hasShadow ? this.boxShadow : 'none',
                     borderTopLeftRadius: this.borderRadius.topLeft,
                     borderTopRightRadius: this.borderRadius.topRight,
                     borderBottomLeftRadius: this.borderRadius.bottomLeft,
                     borderBottomRightRadius: this.borderRadius.bottomRight,
                     position: 'relative', width: '100%', height: '100%', overflow: "hidden"
                   }}>
                <div style={{ position: 'absolute', width: '100%', height: '100%', overflow: "hidden" }}>
                  <EuiButtonIcon aria-label={"fullscreen"} color="danger" iconType="fullScreen" size="s" iconSize="s" title="toggle fullscreen" onClick={() => { this.startFullScreen() }} style={{ zIndex: 1000, position: "absolute", top: 0, right: 0 }}></EuiButtonIcon>
                  <MonacoEditor
                      language={language}
                      defaultValue={this.code}
                      value={this.code}
                      //theme="vs-dark"
                      //theme="hc-black"
                      //theme={_themeId}
                      options={{
                        selectOnLineNumbers: true,
                        automaticLayout: true,
                        quickSuggestions: true,
                        quickSuggestionsDelay: 500,
                        readOnly: !this.props.editMode,
                        minimap: { enabled: true }
                      }}
                      onChange={(x: string | undefined) => this.onCodeChange(x)}
                      editorDidMount={(editor: any, monacoT: any) => { this.editorDidMount(editor, monacoT) }}
                  />
                </div>
              </div>
    )
  }

  render() {
    if (this.state.editorFullScreen) { this.setFullscreen() }
    return this.renderEditor()
  }
}
