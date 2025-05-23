export const decode = (base64String) => {
  const binaryString = atob(base64String);
  const byteArray = new Uint8Array(binaryString.length);
  for (let i = 0; i < binaryString.length; i++) {
    byteArray[i] = binaryString.charCodeAt(i);
  }
  return byteArray;
};

export const encode = (buffer)=> {
    let binary = "";
    const bytes = new Uint8Array(buffer);
    const len = bytes.byteLength;
    for(let i=0; i<len;i++){
        binary+= String.fromCharCode(bytes[i]);
    }
    return btoa(binary);
}