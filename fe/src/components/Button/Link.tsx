import React from 'react';
import { Modal, Button, Link as ArcoLink } from '@arco-design/web-react';
import { IconFullscreen, IconLaunch } from '@arco-design/web-react/icon'
import { Link as RouterLink } from "react-router-dom"


export function Link({ href, refresh = false, newTab = false, children }) {
    if (refresh || newTab) {
        return <ArcoLink href={href} target={newTab ? '_blank' : ''}>{children}</ArcoLink>
    }
    return <RouterLink to={href}>{children}</RouterLink>
}

export function LinkButton({ href, text = "", refresh = false, newTab = false, btnType = "primary" }) {
    // @ts-ignore
    return <Link href={href} refresh={refresh} newTab={newTab}><Button type={btnType} icon={newTab ? <IconLaunch /> : <></>}>{text}</Button></Link>
}


export function LinkText({ href, refresh = false, newTab = false, text = "" }) {
    return <Link href={href} refresh={refresh} newTab={newTab}><span style={{ textDecoration: "underline", color: "rgb(22, 93, 255)" }}>{text}</span></Link>
}

