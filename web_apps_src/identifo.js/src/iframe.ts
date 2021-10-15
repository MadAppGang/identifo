const Iframe = {
  create(): HTMLIFrameElement {
    const iframe = document.createElement('iframe');
    iframe.style.display = 'none';
    document.body.appendChild(iframe);
    return iframe;
  },

  remove(iframe: HTMLIFrameElement): void {
    setTimeout(() => {
      if (document.body.contains(iframe)) {
        document.body.removeChild(iframe);
      }
    }, 0);
  },

  captureMessage(iframe: HTMLIFrameElement, src: string): Promise<string> {
    return new Promise((resolve, reject) => {
      const handleMessage = (event: MessageEvent<{ error: string; accessToken: string }>) => {
        if (event.data.error) reject(event.data.error);

        resolve(event.data.accessToken);

        window.removeEventListener('message', handleMessage);
      };

      window.addEventListener('message', handleMessage, false);
      // eslint-disable-next-line no-param-reassign
      iframe.src = src;
    });
  },
};

export default Iframe;
