import base64 from "base-64";


export const encode = (str)=> {
    return base64.encode(str)
}