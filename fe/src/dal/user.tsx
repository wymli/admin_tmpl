import { doRequest } from '@/utils/axios'

function GetUser({
    query = null,
    data = null,
    setLoadingFn = null,
    okFn = null,
    errFn = null,
}) {
    doRequest({
        method: "GET",
        url: "/api/v1/user",
        query: query,
        data: data,
        setLoadingFn: setLoadingFn,
        okFn: okFn,
        errFn: errFn
    })
}
