import {Toast} from "@elastic/eui/src/components/toast/global_toast_list";
import {EuiGlobalToastList} from "@elastic/eui";
import React from "react";

interface Props {
    show?: boolean
}

interface IComponentState {
    update: boolean
}

export class Toaster extends React.Component<Props, IComponentState> {
    private static readonly container: Toast[] = [];
    private static toastId: number = 0;
    private static instance: Toaster | null

    constructor (props: Props) {
        super(props);
        this.state = { update: false}
        Toaster.instance = this
    }

    public static cId() : string {
        return "Toaster"
    }

    public static get(): Toast[] {
        const r : Toast[] = []
        Toaster.container.forEach((t) => {
            r.push({...t})
        })
        return r;
    }

    public static add(severity: string, iconType: string, title: string, text: string): string {
        Toaster.toastId++;
        const id = Toaster.toastId.toString();
        // "primary" | "success" | "warning" | "danger"

        const color = severity == "success" ? "success" : "danger"

        const toast : Toast = {
            id: id,
            text: text,
            //toastLifeTimeMs: 5000,
            title: title,
            color: color,
            iconType: iconType,
        }
        Toaster.container.push(toast);

        Toaster.update()

        return id;
    }

    private static update() {
        const self = Toaster.instance
        if (self != null) {
            self.setState({update: !self.state.update})
        } else {
            console.log('MISSING TOASTER INSTANCE!')
        }
    }

    public static remove(id: string) {
        //this.container.filter((toast) => toast.id !== id)
        //return;

        for (let i = 0; i < Toaster.container.length; ++i) {
            if (Toaster.container[i].id === id) {
                Toaster.container.splice(i, 1);
                Toaster.update()
                return;
            }
        }
    }

    render() {
        return (
            <EuiGlobalToastList
                toasts={Toaster.get()}
                dismissToast={(removedToast) => {
                    Toaster.remove(removedToast.id)
                }}
                toastLifeTimeMs={5000}
            />
        )
    }
}
