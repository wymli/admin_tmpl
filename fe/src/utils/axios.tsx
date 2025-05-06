import { Message } from '@arco-design/web-react';
import axios from 'axios'

export function doRequest({
    method,
    url,
    query = null,
    data = null,
    setLoadingFn = null,
    okFn = null,
    errFn = null
}) {
    if (setLoadingFn) {
        setLoadingFn(true)
    }

    axios({
        method: method,
        url: url,
        data: data,
        params: query,
    }).then(res => {
        const { data, message, code } = res.data
        if (code != 0) {
            Message.error(`failed to ${method} '${url}', err=${message}`)
            if (errFn) {
                errFn()
            }
        }

        if (okFn) {
            okFn(data)
        }

    }).catch(err => {
        console.log(err);

        Message.error(`failed to ${method} '${url}', err=${err}`)
        if (errFn) {
            errFn()
        }

    }).finally(() => {
        if (setLoadingFn) {
            setLoadingFn(false)
        }
    })
}

