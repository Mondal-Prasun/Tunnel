export const decode = (str) => {
  const parts = str.split(",");
  const base64Data = parts.length > 1 ? parts[1] : parts[0];
  const contentType =
    parts.length > 1 ? parts[0].split(":")[1].split(";")[0] : "";

  const byteCharacters = atob(base64Data);
  const byteNumbers = new Array(byteCharacters.length);
  for (let i = 0; i < byteCharacters.length; i++) {
    byteNumbers[i] = byteCharacters.charCodeAt(i);
  }
  const byteArray = new Uint8Array(byteNumbers);

  return new Blob([byteArray], { type: contentType });
};
