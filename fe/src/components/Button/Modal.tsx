import React from 'react';
import { Modal, Button } from '@arco-design/web-react';
import { IconFullscreen } from '@arco-design/web-react/icon'
import { Editor } from '@/components/Editor'

export function ModalButton({
    btnType = "primary",
    btnSize = "mini",
    btnBody = "Button",
    btnIcon = <></>,
    modalTitle = "Title",
    modalBody = <></>,
    modalWidth = null,
    modalHeight = null,
}) {
    const [visible, setVisible] = React.useState(false);
    return (
        <div>
            <Button
                // @ts-ignore
                type={btnType} size={btnSize}
                icon={btnIcon}
                onClick={() => { setVisible(!visible) }}>
                {btnBody}
            </Button>
            <Modal
                title={modalTitle}
                visible={visible}
                onOk={() => setVisible(false)}
                onCancel={() => setVisible(false)}
                autoFocus={false}
                focusLock={true}
                footer={null}
                style={{ height: modalHeight, width: modalWidth }}
            >
                {modalBody}
            </Modal>
        </div>
    )
}

export function ModalButtonPreview({ modalTitle = "Title", modalBody = <></>, modalTextBody = null }) {
    if (modalTextBody) {
        // use 'modalBody' instead if u want to ctrl Editor
        modalBody = <Editor readOnly={true} value={modalTextBody}></Editor>
    }

    return <ModalButton
        btnIcon={<IconFullscreen />}
        btnBody="FullScreen"
        modalTitle={modalTitle}
        modalBody={modalBody}
        modalWidth="80vw"
        modalHeight="80vh">
    </ModalButton>
}   