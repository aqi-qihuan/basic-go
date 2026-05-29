import '@wangeditor/editor/dist/css/style.css'
import { useState, useEffect } from 'react'
import { Editor, Toolbar } from '@wangeditor/editor-for-react'
import { IDomEditor, IEditorConfig, IToolbarConfig } from '@wangeditor/editor'

interface Props {
    html: any,
    setHtmlFn: any
}

function WangEditor(props: Props) {
    // editor 实例
    const [editor, setEditor] = useState<IDomEditor | null>(null)   // TS 语法
    // const [editor, setEditor] = useState(null)                   // JS 语法

    // 编辑器内容
    const {html, setHtmlFn} = props


    // 模拟 ajax 请求，异步设置 html
    // useEffect(() => {
    //     setTimeout(() => {
    //         setHtml('<p>hello world</p>')
    //     }, 1500)
    // }, [])

    // 工具栏配置
    const toolbarConfig: Partial<IToolbarConfig> = { }  // TS 语法
    // const toolbarConfig = { }                        // JS 语法

    // 编辑器配置
    const editorConfig: Partial<IEditorConfig> = {    // TS 语法
        // const editorConfig = {                         // JS 语法
        placeholder: '请输入内容...',
    }

    // 及时销毁 editor ，重要！
    useEffect(() => {
        return () => {
            if (editor == null) return
            editor.destroy()
            setEditor(null)
        }
    }, [editor])

    return (
        <>
            <div style={{ 
                border: '1px solid rgba(240, 192, 96, 0.3)', 
                zIndex: 100,
                borderRadius: 10,
                overflow: 'hidden',
                boxShadow: '0 0 10px rgba(240, 192, 96, 0.1)',
            }}>
                <Toolbar
                    editor={editor}
                    defaultConfig={toolbarConfig}
                    mode="default"
                    style={{ 
                        borderBottom: '1px solid rgba(240, 192, 96, 0.3)',
                        background: 'rgba(19, 21, 32, 0.9)',
                    }}
                />
                <Editor
                    defaultConfig={editorConfig}
                    value={html}
                    onCreated={setEditor}
                    onChange={editor => setHtmlFn(editor.getHtml())}
                    mode="default"
                    style={{ 
                        height: '500px', 
                        overflowY: 'hidden',
                        background: 'rgba(19, 21, 32, 0.8)',
                        color: '#F5F0E8',
                    }}
                />
            </div>
        </>
    )
}
export default WangEditor
