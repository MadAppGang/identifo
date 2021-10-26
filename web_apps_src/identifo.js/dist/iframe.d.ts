declare const Iframe: {
    create(): HTMLIFrameElement;
    remove(iframe: HTMLIFrameElement): void;
    captureMessage(iframe: HTMLIFrameElement, src: string): Promise<string>;
};
export default Iframe;
