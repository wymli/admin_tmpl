import React from 'react';
import ReactDOM from 'react-dom';

import { Button, Card } from '@arco-design/web-react'

import MocacoEditor from '@monaco-editor/react';

export function Editor({
    title = "",
    value = "",
    setValue = null,
    lang = "json",
    extra = <></>,
    readOnly = false
}) {
    return (
        <Card title={title} bodyStyle={{ padding: 0 }} extra={extra}>
            <div style={{ border: "1px solid var(--color-neutral-3)", borderRadius: "4px", padding: "5px", }}>
                <MocacoEditor
                    height={"50vh"}
                    defaultLanguage={lang}
                    value={value}
                    theme="Custom1" // light | vs-dark
                    onChange={(value, event) => {
                        if (setValue) {
                            setValue(value)
                        }
                    }}
                    options={{ readOnly: readOnly }}
                    beforeMount={(monaco) => {
                        monaco.editor.defineTheme("Custom1", {
                            base: "vs",
                            inherit: true,
                            colors: {
                                "editor.background": "#fffff1",
                            },
                            rules: []
                        });
                        monaco.editor.setTheme("Custom1");
                    }}
                />
            </div>
        </Card>

    );
}

