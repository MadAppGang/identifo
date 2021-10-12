import {
  setAssetPath,
  defineCustomElements,
  JSX as LocalJSX,
} from "@identifo/web-element";
import { HTMLAttributes } from "react";

type StencilToReact<T> = {
  [P in keyof T]?: T[P] &
    Omit<HTMLAttributes<Element>, "className"> & {
      class?: string;
    };
};

declare global {
  export namespace JSX {
    interface IntrinsicElements
      extends StencilToReact<LocalJSX.IntrinsicElements> {}
  }
}

// You need to copy node_modules/@identifo/web-element/assets to public folder and setAssetPath
setAssetPath(window.location.origin + "/identifo/");
defineCustomElements(window);
