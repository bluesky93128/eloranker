declare module '*.vue' {
  import Vue from 'vue';
  export default Vue;
}

declare module 'free-google-image-search' {
  export default class GoogleImageSearch {
    static searchImage(query: string): Promise<string[]>;
    static googleGetMobile(images: NodeList): string[];
    static googleGetDesktop(images: NodeList): string[];
  }
}
